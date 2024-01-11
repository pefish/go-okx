package ws

import (
	"encoding/json"
	"fmt"
	"github.com/pefish/go-okx"
	"github.com/pefish/go-okx/events"
	"github.com/pefish/go-okx/events/public"
	requests "github.com/pefish/go-okx/requests/ws/public"
	"strings"
)

// Public
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels
type Public struct {
	*ClientWs
	InstrumentsCh                    chan *public.Instruments
	TickersCh                        chan *public.Tickers
	OpenInterestCh                   chan *public.OpenInterest
	CandlesticksCh                   chan *public.Candlesticks
	TradesCh                         chan *public.Trades
	EstimatedDeliveryExercisePriceCh chan *public.EstimatedDeliveryExercisePrice
	MarkPriceCh                      chan *public.MarkPrice
	MarkPriceCandlesticksCh          chan *public.MarkPriceCandlesticks
	PriceLimitCh                     chan *public.PriceLimit
	OrderBookCh                      chan *public.OrderBook
	OptionSummaryCh                  chan *public.OptionSummary
	FundingRateCh                    chan *public.FundingRate
	IndexCandlesticksCh              chan *public.IndexCandlesticks
	IndexTickersCh                   chan *public.IndexTickers
	LiquidationOrdersCh              chan *public.LiquidationOrders
}

// NewPublic returns a pointer to a fresh Public
func NewPublic(c *ClientWs) *Public {
	return &Public{ClientWs: c}
}

// Instruments
// The full instrument list will be pushed for the first time after subscription. Subsequently, the instruments will be pushed if there's any change to the instrument’s state (such as delivery of FUTURES, exercise of OPTION, listing of new contracts / trading pairs, trading suspension, etc.).
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-instruments-channel
func (c *Public) Instruments(req []requests.Instruments, ch ...chan *public.Instruments) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "instruments"
	}
	if len(ch) > 0 {
		c.InstrumentsCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UInstruments
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-instruments-channel
func (c *Public) UInstruments(req []requests.Instruments, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "instruments"
	}
	if len(rCh) > 0 && rCh[0] {
		c.InstrumentsCh = nil
	}
	return c.Unsubscribe(false, m)
}

// Tickers
// Retrieve the last traded price, bid price, ask price and 24-hour trading volume of instruments. Data will be pushed every 100 ms.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-tickers-channel
func (c *Public) Tickers(req []requests.Tickers, ch ...chan *public.Tickers) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "tickers"
	}
	if len(ch) > 0 {
		c.TickersCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UTickers
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-tickers-channel
func (c *Public) UTickers(req []requests.Tickers, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "tickers"
	}
	if len(rCh) > 0 && rCh[0] {
		c.TickersCh = nil
	}
	return c.Unsubscribe(false, m)
}

// OpenInterest
// Retrieve the open interest. Data will by pushed every 3 seconds.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-open-interest-channel
func (c *Public) OpenInterest(req []requests.OpenInterest, ch ...chan *public.OpenInterest) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "open-interest"
	}
	if len(ch) > 0 {
		c.OpenInterestCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UOpenInterest
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-open-interest-channel
func (c *Public) UOpenInterest(req []requests.OpenInterest, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "open-interest"
	}
	if len(rCh) > 0 && rCh[0] {
		c.OpenInterestCh = nil
	}
	return c.Unsubscribe(false, m)
}

// Candlesticks
// Retrieve the open interest. Data will by pushed every 3 seconds.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-candlesticks-channel
func (c *Public) Candlesticks(req []requests.Candlesticks, ch ...chan *public.Candlesticks) error {
	if len(ch) > 0 {
		c.CandlesticksCh = ch[0]
	}
	return c.Subscribe(false, okex.StructSlice2MapSlice(req))
}

// UCandlesticks
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-candlesticks-channel
func (c *Public) UCandlesticks(req []requests.Candlesticks, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	if len(rCh) > 0 && rCh[0] {
		c.CandlesticksCh = nil
	}
	return c.Unsubscribe(false, m)
}

// Trades
// Retrieve the recent trades data. Data will be pushed whenever there is a trade.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-trades-channel
func (c *Public) Trades(req []requests.Trades, ch ...chan *public.Trades) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "trades"
	}
	if len(ch) > 0 {
		c.TradesCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UTrades
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-trades-channel
func (c *Public) UTrades(req []requests.Trades, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "trades"
	}
	if len(rCh) > 0 && rCh[0] {
		c.TradesCh = nil
	}
	return c.Unsubscribe(false, m)
}

// EstimatedDeliveryExercisePrice
// Retrieve the estimated delivery/exercise price of FUTURES contracts and OPTION.
//
// Only the estimated delivery/exercise price will be pushed an hour before delivery/exercise, and will be pushed if there is any price change.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-estimated-delivery-exercise-price-channel
func (c *Public) EstimatedDeliveryExercisePrice(req []requests.EstimatedDeliveryExercisePrice, ch ...chan *public.EstimatedDeliveryExercisePrice) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "estimated-price"
	}
	if len(ch) > 0 {
		c.EstimatedDeliveryExercisePriceCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UEstimatedDeliveryExercisePrice
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-estimated-delivery-exercise-price-channel
func (c *Public) UEstimatedDeliveryExercisePrice(req []requests.EstimatedDeliveryExercisePrice, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "estimated-price"
	}
	if len(rCh) > 0 && rCh[0] {
		c.EstimatedDeliveryExercisePriceCh = nil
	}
	return c.Unsubscribe(false, m)
}

// MarkPrice
// Retrieve the mark price. Data will be pushed every 200 ms when the mark price changes, and will be pushed every 10 seconds when the mark price does not change.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-mark-price-channel
func (c *Public) MarkPrice(req []requests.MarkPrice, ch ...chan *public.MarkPrice) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "mark-price"
	}
	if len(ch) > 0 {
		c.MarkPriceCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UMarkPrice
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-mark-price-channel
func (c *Public) UMarkPrice(req []requests.MarkPrice, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "mark-price"
	}
	if len(rCh) > 0 && rCh[0] {
		c.MarkPriceCh = nil
	}
	return c.Unsubscribe(false, m)
}

// MarkPriceCandlesticks
// Retrieve the candlesticks data of the mark price. Data will be pushed every 500 ms.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-mark-price-candlesticks-channel
func (c *Public) MarkPriceCandlesticks(req []requests.MarkPriceCandlesticks, ch ...chan *public.MarkPriceCandlesticks) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "mark-price-" + m[i]["channel"]
	}
	if len(ch) > 0 {
		c.MarkPriceCandlesticksCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UMarkPriceCandlesticks
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-mark-price-candlesticks-channel
func (c *Public) UMarkPriceCandlesticks(req []requests.MarkPriceCandlesticks, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "mark-price-" + m[i]["channel"]
	}
	if len(rCh) > 0 && rCh[0] {
		c.MarkPriceCandlesticksCh = nil
	}
	return c.Unsubscribe(false, m)
}

// PriceLimit
// Retrieve the maximum buy price and minimum sell price of the instrument. Data will be pushed every 5 seconds when there are changes in limits, and will not be pushed when there is no changes on limit.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-price-limit-channel
func (c *Public) PriceLimit(req []requests.PriceLimit, ch ...chan *public.PriceLimit) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "price-limit"
	}
	if len(ch) > 0 {
		c.PriceLimitCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UPriceLimit
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-price-limit-channel
func (c *Public) UPriceLimit(req []requests.PriceLimit, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "price-limit"
	}
	if len(rCh) > 0 && rCh[0] {
		c.PriceLimitCh = nil
	}
	return c.Unsubscribe(false, m)
}

// OrderBook
// Retrieve order book data.
//
// Use books for 400 depth levels, book5 for 5 depth levels, books50-l2-tbt tick-by-tick 50 depth levels, and books-l2-tbt for tick-by-tick 400 depth levels.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-order-book-channel
func (c *Public) OrderBook(req []requests.OrderBook, ch ...chan *public.OrderBook) error {
	m := okex.StructSlice2MapSlice(req)
	if len(ch) > 0 {
		c.OrderBookCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UOrderBook
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-order-book-channel
func (c *Public) UOrderBook(req []requests.OrderBook, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	if len(rCh) > 0 && rCh[0] {
		c.OrderBookCh = nil
	}
	return c.Unsubscribe(false, m)
}

// OPTIONSummary
// Retrieve detailed pricing information of all OPTION contracts. Data will be pushed at once.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-option-summary-channel
func (c *Public) OPTIONSummary(req []requests.OPTIONSummary, ch ...chan *public.OptionSummary) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "opt-summary"
	}
	if len(ch) > 0 {
		c.OptionSummaryCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UOPTIONSummary
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-option-summary-channel
func (c *Public) UOPTIONSummary(req []requests.OPTIONSummary, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "opt-summary"
	}
	if len(rCh) > 0 && rCh[0] {
		c.OptionSummaryCh = nil
	}
	return c.Unsubscribe(false, m)
}

// FundingRate
// Retrieve funding rate. Data will be pushed every minute.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-funding-rate-channel
func (c *Public) FundingRate(req []requests.FundingRate, ch ...chan *public.FundingRate) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "funding-rate"
	}
	if len(ch) > 0 {
		c.FundingRateCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UFundingRate
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-funding-rate-channel
func (c *Public) UFundingRate(req []requests.FundingRate, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "funding-rate"
	}
	if len(rCh) > 0 && rCh[0] {
		c.FundingRateCh = nil
	}
	return c.Unsubscribe(false, m)
}

// IndexCandlesticks
// Retrieve the candlesticks data of the index. Data will be pushed every 500 ms.
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-index-candlesticks-channel
func (c *Public) IndexCandlesticks(req []requests.IndexCandlesticks, ch ...chan *public.IndexCandlesticks) error {
	m := okex.StructSlice2MapSlice(req)
	if len(ch) > 0 {
		c.IndexCandlesticksCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UIndexCandlesticks
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-index-candlesticks-channel
func (c *Public) UIndexCandlesticks(req []requests.IndexCandlesticks, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	if len(rCh) > 0 && rCh[0] {
		c.IndexCandlesticksCh = nil
	}
	return c.Unsubscribe(false, m)
}

// IndexTickers
// Retrieve index tickers data
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-index-tickers-channel
func (c *Public) IndexTickers(req []requests.IndexTickers, ch ...chan *public.IndexTickers) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "index-tickers"
	}
	if len(ch) > 0 {
		c.IndexTickersCh = ch[0]
	}
	return c.Subscribe(false, m)
}

// UIndexTickers
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-index-tickers-channel
func (c *Public) UIndexTickers(req []requests.IndexTickers, rCh ...bool) error {
	m := okex.StructSlice2MapSlice(req)
	for i, _ := range m {
		m[i]["channel"] = "index-tickers"
	}
	if len(rCh) > 0 && rCh[0] {
		c.IndexTickersCh = nil
	}
	return c.Unsubscribe(false, m)
}

func (c *Public) Process(data []byte, e *events.Basic) bool {
	if e.Event == "" && e.Arg != nil && e.Data != nil && len(e.Data) > 0 {
		ch, ok := e.Arg.Get("channel")
		if !ok {
			return false
		}
		switch ch {
		case "instruments":
			e := public.Instruments{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.InstrumentsCh != nil {
				c.InstrumentsCh <- &e
			}
			return true
		case "tickers":
			e := public.Tickers{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.TickersCh != nil {
				c.TickersCh <- &e
			}
			return true
		case "open-interest":
			e := public.OpenInterest{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.OpenInterestCh != nil {
				c.OpenInterestCh <- &e
			}
			return true
		case "trades":
			e := public.Trades{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.TradesCh != nil {
				c.TradesCh <- &e
			}
			return true
		case "estimated-price":
			e := public.EstimatedDeliveryExercisePrice{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.EstimatedDeliveryExercisePriceCh != nil {
				c.EstimatedDeliveryExercisePriceCh <- &e
			}
			return true
		case "mark-price":
			e := public.MarkPrice{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.MarkPriceCh != nil {
				c.MarkPriceCh <- &e
			}
			return true
		case "price-limit":
			e := public.PriceLimit{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.PriceLimitCh != nil {
				c.PriceLimitCh <- &e
			}
			return true
		case "opt-summary":
			e := public.OptionSummary{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.OptionSummaryCh != nil {
				c.OptionSummaryCh <- &e
			}
			return true
		case "funding-rate":
			e := public.FundingRate{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.FundingRateCh != nil {
				c.FundingRateCh <- &e
			}
			return true
		case "index-tickers":
			e := public.IndexTickers{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.IndexTickersCh != nil {
				c.IndexTickersCh <- &e
			}
			return true
		case "liquidation-orders":
			e := public.LiquidationOrders{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return false
			}
			if c.LiquidationOrdersCh != nil {
				c.LiquidationOrdersCh <- &e
			}
			return true
		default:
			// special cases
			// market price candlestick channel
			chName := fmt.Sprint(ch)
			// market price channels
			if strings.Contains(chName, "mark-price-candle") {
				e := public.MarkPriceCandlesticks{}
				err := json.Unmarshal(data, &e)
				if err != nil {
					return false
				}
				if c.MarkPriceCandlesticksCh != nil {
					c.MarkPriceCandlesticksCh <- &e
				}
				return true
			}
			// index chandlestick channels
			if strings.Contains(chName, "index-candle") {
				e := public.IndexCandlesticks{}
				err := json.Unmarshal(data, &e)
				if err != nil {
					return false
				}
				if c.IndexCandlesticksCh != nil {
					c.IndexCandlesticksCh <- &e
				}
				return true
			}
			// candlestick channels
			if strings.Contains(chName, "candle") {
				e := public.Candlesticks{}
				err := json.Unmarshal(data, &e)
				if err != nil {
					return false
				}
				if c.CandlesticksCh != nil {
					c.CandlesticksCh <- &e
				}
				return true
			}
			// order book channels
			if strings.Contains(chName, "books") {
				e := public.OrderBook{}
				err := json.Unmarshal(data, &e)
				if err != nil {
					return false
				}
				if c.OrderBookCh != nil {
					c.OrderBookCh <- &e
				}
				return true
			}
		}
	}
	return false
}
