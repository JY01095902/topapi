package proxy

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/jy01095902/topapi/request"
)

type Trade request.Values

type tradeProxy struct {
	req request.APIRequest
}

func NewTradeProxy(config Config) tradeProxy {
	proxy := tradeProxy{
		req: request.NewRequest(config.AppKey, config.AppSecret, config.SessionKey),
	}

	return proxy
}

/*
	淘宝只能保证最近三个月内数据的正确性，所以此API查询的订单范围为从调用日期往前推三个月的0点开始到调用日期前一天的24点结束。
	调用日期：2022年01月19日11:20:24
	开始时间：2021年10月19日00:00:00
	结束时间：2022年01月18日23:59:59
*/
func (proxy tradeProxy) ListBaseTrades() ([]Trade, error) {
	now := time.Now()
	firstday := now.AddDate(0, -3, 0)
	start := time.Date(firstday.Year(), firstday.Month(), firstday.Day(), 0, 0, 0, 0, now.Location())
	yesterday := now.AddDate(0, 0, -1)
	end := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, now.Location())

	params := url.Values{}
	params.Add("start_created", start.Format("2006-01-02 15:04:05"))
	params.Add("end_created", end.Format("2006-01-02 15:04:05"))
	params.Add("page_no", "1")
	params.Add("page_size", "100")
	params.Add("fields", "total_results,tid,created")

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

	resp, err := proxy.req.GetAll("taobao.trades.sold.get", params, parseTotal)
	if err != nil {
		return nil, err
	}

	result := []Trade{}
	for _, res := range resp {
		if vals, ok := res.Get("trades").([]interface{}); ok {
			for i := range vals {
				if val, ok := vals[i].(map[string]interface{}); ok {
					result = append(result, Trade(request.Values(val)))
				} else {

					fmt.Printf("vals[i]type: %+v \n", reflect.TypeOf(vals[i]))
				}
			}
		} else {

			fmt.Printf("res.Get('trades') type: %+v \n", reflect.TypeOf(res.Get("trades")))
		}
	}

	return result, nil
}
