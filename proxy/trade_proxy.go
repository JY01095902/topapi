package proxy

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/jy01095902/topapi/request"
)

type Trade request.Values

func (trade Trade) TId() string {
	if tid, ok := trade["tid"].(string); ok {
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

/*
	淘宝只能保证最近三个月内数据的正确性，所以此API查询的订单范围为从调用日期往前推三个月的0点开始到调用日期前一天的24点结束。
	调用日期：2022年01月19日11:20:24
	开始时间：2021年10月19日00:00:00
	结束时间：2022年01月18日23:59:59
*/
func (proxy TradeProxy) ListBaseTrades() ([]Trade, error) {
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

func (proxy TradeProxy) ListIncrementTrades(start, end time.Time) ([]Trade, error) {
	params := url.Values{}
	params.Add("start_modified", start.Format("2006-01-02 15:04:05"))
	params.Add("end_modified", end.Format("2006-01-02 15:04:05"))
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

	resp, err := proxy.req.GetAll("taobao.trades.sold.increment.get", params, parseTotal)
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

func (proxy TradeProxy) GetFullinfoTrade(tid string) (Trade, error) {
	fields := []string{
		"payment",
		"post_fee",
		"receiver_name",
		"receiver_state",
		"receiver_address",
		"receiver_zip",
		"receiver_mobile",
		"receiver_phone",
		"consign_time",
		"received_payment",
		"promotion_details",
		"id",
		"gift_item_name",
		"gift_item_id",
		"gift_item_num",
		"has_post_fee",
		"promotion_id",
		"promotion_name",
		"promotion_desc",
		"receiver_country",
		"receiver_town",
		"tid",
		"num_iid",
		"status",
		"title",
		"type",
		"price",
		"discount_fee",
		"total_fee",
		"created",
		"pay_time",
		"buyer_cod_fee",
		"modified",
		"end_time",
		"nr_outer_iid",
		"outer_iid",
		"buyer_nick",
		"credit_card_fee",
		"has_yfx",
		"yfx_fee",
		"step_trade_status",
		"step_paid_fee",
		"shipping_type",
		"adjust_fee",
		"trade_from",
		"service_orders",
		"receiver_city",
		"receiver_district",
		"orders",
		"delivery_time",
		"collect_time",
		"dispatch_time",
		"sign_time",
		"delivery_cps",
		"refund_status",
		"oaid",
		"cid",
		"estimate_con_time",
		"oid",
		"item_oid",
		"service_id",
		"sku_id",
		"item_meal_id",
		"item_meal_name",
		"num",
		"outer_sku_id",
		"order_from",
		"refund_id",
		"is_service_order",
		"bind_oids_all_status",
		"logistics_company",
		"invoice_no",
		"divide_order_fee",
		"part_mjz_discount",
		"store_code",
		"md_fee",
		"customization",
		"inv_type",
		"is_sh_ship",
		"shipper",
		"f_type",
		"f_status",
		"f_term",
		"assembly_rela",
		"assembly_price",
		"assembly_item",
	}
	params := url.Values{}
	params.Add("tid", tid)
	params.Add("include_oaid", "true")
	params.Add("fields", strings.Join(fields, ","))

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
