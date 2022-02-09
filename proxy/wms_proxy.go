package proxy

import (
	"log"
	"net/url"

	"github.com/jy01095902/topapi/request"
)

type SerialNumberInfo request.Values

func (sninfo SerialNumberInfo) ItemId() string {
	if id, ok := sninfo["item_id"].(string); ok {
		return id
	}

	return ""
}

func (sninfo SerialNumberInfo) ItemCode() string {
	if code, ok := sninfo["item_code"].(string); ok {
		return code
	}

	return ""
}

func (sninfo SerialNumberInfo) SerialNumber() string {
	if sn, ok := sninfo["sn_code"].(string); ok {
		return sn
	}

	return ""
}

type WMSProxy struct {
	req request.APIRequest
}

func NewWMSProxy(config Config) WMSProxy {
	proxy := WMSProxy{
		req: request.NewRequest(config.AppKey, config.AppSecret, config.SessionKey),
	}

	return proxy
}

func (proxy WMSProxy) ListSerialNumberInfos(opts ...option) ([]SerialNumberInfo, error) {
	params := url.Values{}
	params.Add("page_index", "1")
	for i := range opts {
		opts[i](params)
	}

	parseTotal := func(val request.Values) int {
		type response struct {
			Result struct {
				Total int `json:"total_count"`
			} `json:"result"`
		}

		res, err := val.GetResult(response{})
		if err != nil {
			return 0
		}

		r, ok := res.(*response)
		if !ok {
			return 0
		}

		// log.Printf("parse total val: %+v", val)
		// log.Printf("parse total r: %+v", r)
		return r.Result.Total
	}

	resp, err := proxy.req.GetAll("taobao.wlb.wms.sn.info.query", params, "page_index", "page_size", parseTotal)
	if err != nil {
		return nil, err
	}

	type response struct {
		Result struct {
			List []struct {
				SerialNumberInfo SerialNumberInfo `json:"sn_info"`
			} `json:"sn_info_list"`
		} `json:"result"`
	}
	result := []SerialNumberInfo{}
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

		for i := range r.Result.List {
			result = append(result, r.Result.List[i].SerialNumberInfo)
		}
	}

	return result, nil
}
