package proxy

import (
	"fmt"
	"log"
	"net/url"

	"github.com/jy01095902/topapi/request"
)

type Trade request.Values

func (trade Trade) TId() string {
	if tid, ok := trade["tid"].(string); ok {
		return tid
	}

	return ""
}

func (trade Trade) Status() string {
	if tid, ok := trade["status"].(string); ok {
		return tid
	}

	return ""
}

type TradeProxy struct {
	req request.APIRequest
}

func NewTradeProxy(config Config) TradeProxy {
	proxy := TradeProxy{
		req: request.NewRequest(config.AppKey, config.AppSecret, config.SessionKey),
	}

	return proxy
}

func (proxy TradeProxy) ListSoldTrades(opts ...option) ([]Trade, error) {
	params := url.Values{}
	params.Add("page_no", "1")
	params.Add("page_size", "100")
	for i := range opts {
		opts[i](params)
	}

	parseTotal := func(val request.Values) int {
		type response struct {
			Total int `json:"total_results"`
		}

		res, err := val.GetResult(response{})
		if err != nil {
			return 0
		}

		r, ok := res.(*response)
		if !ok {
			return 0
		}

		return r.Total
	}

	resp, err := proxy.req.GetAll("taobao.trades.sold.get", params, "page_no", "page_size", parseTotal)
	if err != nil {
		return nil, err
	}

	type response struct {
		Trades []Trade `json:"trades"`
	}
	result := []Trade{}
	for _, res := range resp {
		vals, err := res.GetResult(response{})
		if err != nil {
			log.Printf("%s: parse result error", request.ErrTOPAPIBizError.Error())

			continue
		}

		r, ok := vals.(*response)
		if !ok {
			log.Printf("%s: parse result error", request.ErrTOPAPIBizError.Error())

			continue
		}

		result = append(result, r.Trades...)
	}

	return result, nil
}

func (proxy TradeProxy) ListIncrementTrades(opts ...option) ([]Trade, error) {
	params := url.Values{}
	params.Add("page_no", "1")
	params.Add("page_size", "100")
	for i := range opts {
		opts[i](params)
	}

	parseTotal := func(val request.Values) int {
		type response struct {
			Total int `json:"total_results"`
		}

		res, err := val.GetResult(response{})
		if err != nil {
			return 0
		}

		r, ok := res.(*response)
		if !ok {
			return 0
		}

		return r.Total
	}

	resp, err := proxy.req.GetAll("taobao.trades.sold.increment.get", params, "page_no", "page_size", parseTotal)
	if err != nil {
		return nil, err
	}

	type response struct {
		Trades []Trade `json:"trades"`
	}
	result := []Trade{}
	for _, res := range resp {
		vals, err := res.GetResult(response{})
		if err != nil {
			log.Printf("%s: parse result error", request.ErrTOPAPIBizError.Error())

			continue
		}

		r, ok := vals.(*response)
		if !ok {
			log.Printf("%s: parse result error", request.ErrTOPAPIBizError.Error())

			continue
		}

		result = append(result, r.Trades...)
	}

	return result, nil
}

func (proxy TradeProxy) GetFullinfoTrade(opts ...option) (Trade, error) {
	params := url.Values{}
	for i := range opts {
		opts[i](params)
	}

	resp, err := proxy.req.Get("taobao.trade.fullinfo.get", params)
	if err != nil {
		return nil, err
	}

	type response struct {
		Trade Trade `json:"trade"`
	}
	vals, err := resp.GetResult(response{})
	if err != nil {
		return nil, err
	}

	trade, ok := vals.(*response)
	if !ok {
		return nil, fmt.Errorf("%w: result is not trade response", request.ErrTOPAPIBizError)
	}

	return trade.Trade, nil
}
