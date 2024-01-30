package ws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	go_logger "github.com/pefish/go-logger"
	"github.com/pefish/go-okx"
	"github.com/pefish/go-okx/events"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

// ClientWs is the websocket api client
//
// https://www.okex.com/docs-v5/en/#websocket-api
type ClientWs struct {
	Cancel        context.CancelFunc
	ErrChan       chan *events.Error
	SubscribeChan chan *events.Subscribe
	UnsubscribeCh chan *events.Unsubscribe
	LoginChan     chan *events.Login
	SuccessChan   chan *events.Success
	url           map[bool]okex.BaseURL // need or not login -> url
	apiKey        string
	secretKey     []byte
	passphrase    string
	AuthRequested *time.Time
	Authorized    bool
	Private       *Private
	Public        *Public
	Trade         *Trade
	ctx           context.Context
	logger        go_logger.InterfaceLogger
}

const (
	redialTick = 2 * time.Second
	writeWait  = 3 * time.Second
	pongWait   = 30 * time.Second
	PingPeriod = (pongWait * 8) / 10
)

// NewClient returns a pointer to a fresh ClientWs
func NewClient(
	ctx context.Context,
	apiKey,
	secretKey,
	passphrase string,
	url map[bool]okex.BaseURL,
) *ClientWs {
	ctx, cancel := context.WithCancel(ctx)
	c := &ClientWs{
		apiKey:     apiKey,
		secretKey:  []byte(secretKey),
		passphrase: passphrase,
		ctx:        ctx,
		Cancel:     cancel,
		url:        url,
	}
	c.Private = NewPrivate(c)
	c.Public = NewPublic(c)
	c.Trade = NewTrade(c)
	return c
}

func (c *ClientWs) SetLogger(logger go_logger.InterfaceLogger) *ClientWs {
	c.logger = logger
	return c
}

// ReConnect into the server
//
// https://www.okex.com/docs-v5/en/#websocket-api-connect
func (c *ClientWs) ReConnect(needLogin bool, sendData []byte) error {
	senderChan := make(chan []byte, 3)
	sendErrChan := make(chan error)

	err := c.dial(senderChan, needLogin, sendErrChan)
	if err != nil {
		return err
	}
	c.logger.InfoF("Connect success.")
	senderChan <- sendData

	go func() {
		for {
			select {
			case err := <-sendErrChan:
				c.logger.ErrorF("Send error <%+v>, reconnect...\n", err)
				err = c.dial(senderChan, needLogin, sendErrChan)
				if err != nil {
					c.logger.ErrorF("Reconnect failed. %+v\n", err)
					return
				}
				c.logger.InfoF("Connect success.")
				senderChan <- sendData
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Login
//
// https://www.okex.com/docs-v5/en/#websocket-api-login
func (c *ClientWs) Login() error {
	if c.Authorized {
		return nil
	}
	if c.AuthRequested != nil && time.Since(*c.AuthRequested).Seconds() < 30 {
		return nil
	}
	now := time.Now()
	c.AuthRequested = &now
	method := http.MethodGet
	path := "/users/self/verify"
	ts, sign := c.sign(method, path)
	args := []map[string]string{
		{
			"apiKey":     c.apiKey,
			"passphrase": c.passphrase,
			"timestamp":  ts,
			"sign":       sign,
		},
	}
	return c.Send(true, okex.LoginOperation, args)
}

// Subscribe
// Users can choose to subscribe to one or more channels, and the total length of multiple channels cannot exceed 4096 bytes.
//
// https://www.okex.com/docs-v5/en/#websocket-api-subscribe
func (c *ClientWs) Subscribe(needLogin bool, args []map[string]string) error {
	return c.Send(needLogin, okex.SubscribeOperation, args)
}

// Unsubscribe into channel(s)
//
// https://www.okex.com/docs-v5/en/#websocket-api-unsubscribe
func (c *ClientWs) Unsubscribe(needLogin bool, args []map[string]string) error {
	return c.Send(needLogin, okex.UnsubscribeOperation, args)
}

// Send message through either connections
func (c *ClientWs) Send(needLogin bool, op okex.Operation, args []map[string]string) error {
	j, err := json.Marshal(map[string]interface{}{
		"op":   op,
		"args": args,
	})
	if err != nil {
		return err
	}
	err = c.ReConnect(needLogin, j)
	if err != nil {
		return err
	}
	if needLogin {
		err = c.WaitForAuthorization()
		if err != nil {
			return err
		}
	}

	return nil
}

// WaitForAuthorization waits for the auth response and try to log in if it was needed
func (c *ClientWs) WaitForAuthorization() error {
	if c.Authorized {
		return nil
	}
	if err := c.Login(); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Millisecond * 300)
	defer ticker.Stop()
	for range ticker.C {
		if c.Authorized {
			return nil
		}
	}
	return nil
}

func (c *ClientWs) dial(
	senderChan chan []byte,
	needLogin bool,
	sendErrChan chan<- error,
) error {
	conn, res, err := websocket.DefaultDialer.Dial(string(c.url[needLogin]), nil)
	if err != nil {
		var statusCode int
		if res != nil {
			statusCode = res.StatusCode
		}
		return errors.New(fmt.Sprintf("error %d: %w", statusCode, err))
	}
	defer res.Body.Close()
	go func() {
		err := c.receiver(conn)
		if err != nil {
			c.logger.ErrorF("Receiver error: %v\n", err)
		}
	}()
	go func() {
		err := c.sender(conn, senderChan)
		if err != nil {
			sendErrChan <- err
			c.logger.ErrorF("Sender error: %v\n", err)
		}
	}()
	return nil
}
func (c *ClientWs) sender(
	conn *websocket.Conn,
	senderChan chan []byte,
) error {
	ticker := time.NewTicker(time.Millisecond * 300)
	defer ticker.Stop()
	for {
		select {
		case data := <-senderChan:
			err := conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return err
			}
			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}
			if _, err = w.Write(data); err != nil {
				return err
			}
			if err := w.Close(); err != nil {
				return err
			}
		case <-ticker.C:
			senderChan <- []byte("ping")
		case <-c.ctx.Done():
			return errors.New("operation cancelled: sender")
		}
	}
}
func (c *ClientWs) receiver(conn *websocket.Conn) error {
	for {
		select {
		case <-c.ctx.Done():
			return errors.New("operation cancelled: receiver")
		default:
			err := conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				return err
			}
			mt, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return conn.Close()
				}
				return err
			}
			if mt == websocket.TextMessage && string(data) != "pong" {
				e := &events.Basic{}
				if err := json.Unmarshal(data, &e); err != nil {
					return err
				}
				go func() {
					c.process(data, e)
				}()
			}
		}
	}
}
func (c *ClientWs) sign(method, path string) (string, string) {
	t := time.Now().UTC().Unix()
	ts := fmt.Sprint(t)
	s := ts + method + path
	p := []byte(s)
	h := hmac.New(sha256.New, c.secretKey)
	h.Write(p)
	return ts, base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// TODO: break each case into a separate function
func (c *ClientWs) process(data []byte, e *events.Basic) bool {
	switch e.Event {
	case "error":
		e := events.Error{}
		_ = json.Unmarshal(data, &e)
		go func() {
			c.ErrChan <- &e
		}()
		return true
	case "subscribe":
		e := events.Subscribe{}
		_ = json.Unmarshal(data, &e)
		go func() {
			if c.SubscribeChan != nil {
				c.SubscribeChan <- &e
			}
		}()
		return true
	case "unsubscribe":
		e := events.Unsubscribe{}
		_ = json.Unmarshal(data, &e)
		go func() {
			if c.UnsubscribeCh != nil {
				c.UnsubscribeCh <- &e
			}
		}()
		return true
	case "login":
		if time.Since(*c.AuthRequested).Seconds() > 30 {
			c.AuthRequested = nil
			_ = c.Login()
			break
		}
		c.Authorized = true
		e := events.Login{}
		_ = json.Unmarshal(data, &e)
		go func() {
			if c.LoginChan != nil {
				c.LoginChan <- &e
			}
		}()
		return true
	}
	if c.Private.Process(data, e) {
		return true
	}
	if c.Public.Process(data, e) {
		return true
	}
	if e.ID != "" {
		if e.Code != 0 {
			ee := *e
			ee.Event = "error"
			return c.process(data, &ee)
		}
		e := events.Success{}
		_ = json.Unmarshal(data, &e)
		go func() {
			if c.SuccessChan != nil {
				c.SuccessChan <- &e
			}
		}()
		return true
	}
	return false
}
