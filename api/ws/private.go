package ws

import (
	"encoding/json"
	"github.com/pefish/go-okx"
	"github.com/pefish/go-okx/events"
	"github.com/pefish/go-okx/events/private"
	requests "github.com/pefish/go-okx/requests/ws/private"
)

// Private
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel
type Private struct {
	*ClientWs
	AccountCh            chan *private.Account
	PositionCh           chan *private.Position
	BalanceAndPositionCh chan *private.BalanceAndPosition
	OrderCh              chan *private.Order
}

// NewPrivate returns a pointer to a fresh Private
func NewPrivate(c *ClientWs) *Private {
	return &Private{ClientWs: c}
}

// Account
// Retrieve account information. Data will be pushed when triggered by events such as placing/canceling order, and will also be pushed in regular interval according to subscription granularity.
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-account-channel
func (c *Private) Account(req []requests.Account, ch ...chan *private.Account) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "account"
	}
	if len(ch) > 0 {
		c.AccountCh = ch[0]
	}
	return c.Subscribe(true, m)
}

// UAccount
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-account-channel
func (c *Private) UAccount(req []requests.Account, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "account"
	}
	if len(rCh) > 0 && rCh[0] {
		c.AccountCh = nil
	}
	return c.Unsubscribe(true, m)
}

// Position
// Retrieve position information. Initial snapshot will be pushed according to subscription granularity. Data will be pushed when triggered by events such as placing/canceling order, and will also be pushed in regular interval according to subscription granularity.
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-positions-channel
func (c *Private) Position(req []requests.Position, ch ...chan *private.Position) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "positions"
	}
	if len(ch) > 0 {
		c.PositionCh = ch[0]
	}
	return c.Subscribe(true, m)
}

// UPosition
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-positions-channel
func (c *Private) UPosition(req []requests.Position, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "positions"
	}
	if len(rCh) > 0 && rCh[0] {
		c.PositionCh = nil
	}
	return c.Unsubscribe(true, m)
}

// BalanceAndPosition
// Retrieve account balance and position information. Data will be pushed when triggered by events such as filled order, funding transfer.
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-balance-and-position-channel
func (c *Private) BalanceAndPosition(ch ...chan *private.BalanceAndPosition) error {
	m := []map[string]string{
		{
			"channel": "balance_and_position",
		},
	}
	if len(ch) > 0 {
		c.BalanceAndPositionCh = ch[0]
	}
	return c.Subscribe(true, m)
}

// UBalanceAndPosition unsubscribes a position channel
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-balance-and-position-channel
func (c *Private) UBalanceAndPosition(rCh ...bool) error {
	m := []map[string]string{
		{
			"channel": "balance_and_position",
		},
	}
	if len(rCh) > 0 && rCh[0] {
		c.BalanceAndPositionCh = nil
	}
	return c.Unsubscribe(true, m)
}

// Order
// Retrieve position information. Initial snapshot will be pushed according to subscription granularity. Data will be pushed when triggered by events such as placing/canceling order, and will also be pushed in regular interval according to subscription granularity.
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-order-channel
func (c *Private) Order(req []requests.Order, ch ...chan *private.Order) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "orders"
	}
	if len(ch) > 0 {
		c.OrderCh = ch[0]
	}
	return c.Subscribe(true, m)
}

// UOrder
//
// https://www.okex.com/docs-v5/en/#websocket-api-private-channel-order-channel
func (c *Private) UOrder(req []requests.Order, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "orders"
	}
	if len(rCh) > 0 && rCh[0] {
		c.OrderCh = nil
	}
	return c.Unsubscribe(true, m)
}

func (c *Private) Process(data []byte, e *events.Basic) bool {
	if e.Event == "" && e.Arg != nil && e.Data != nil && len(e.Data) > 0 {
		ch, ok := e.Arg.Get("channel")
		if !ok {
			return false
		}
		switch ch {
		case "account":
			e := private.Account{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.AccountCh != nil {
				c.AccountCh <- &e
			}
			return true
		case "positions":
			e := private.Position{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.PositionCh != nil {
				c.PositionCh <- &e
			}
			return true
		case "balance_and_position":
			e := private.BalanceAndPosition{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.BalanceAndPositionCh != nil {
				c.BalanceAndPositionCh <- &e
			}
			return true
		case "orders":
			e := private.Order{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.OrderCh != nil {
				c.OrderCh <- &e
			}
			return true
		}
	}
	return false
}
