package request

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequst(t *testing.T) {
	appKey := "appkey"
	appSecret := "appSecret"
	sessionKey := "session"
	req := NewRequest(appKey, appSecret, sessionKey)
	assert.Equal(t, appKey, req.config["app_key"])
	assert.Equal(t, sessionKey, req.config["session"])
	assert.Equal(t, "2.0", req.config["v"])
	assert.Equal(t, "json", req.config["format"])
	assert.Equal(t, "true", req.config["simplify"])

	version := "3.0"
	format := "xml"
	isSimplified := false
	req = NewRequest(appKey, appSecret, sessionKey, WithVersion(version), WithFormat(format), WithSimplify(isSimplified))
	assert.Equal(t, appKey, req.config["app_key"])
	assert.Equal(t, sessionKey, req.config["session"])
	assert.Equal(t, version, req.config["v"])
	assert.Equal(t, format, req.config["format"])
	assert.Equal(t, "false", req.config["simplify"])
}

func TestGetAll(t *testing.T) {
	parseTotal := func(val Values) int {
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

	req := NewRequest("", "", "")
	vals, err := req.GetAll("", url.Values{}, "", "", parseTotal)
	assert.Nil(t, err)
	assert.Equal(t, int(math.Ceil(5924.0/100)), len(vals))
	result := []interface{}{}
	for _, val := range vals {
		if pageData, err := val.GetResult(Values{}); err == nil {
			if pageVal, ok := pageData.(*Values); ok {
				tradesVal := pageVal.Get("trades")
				if trades, ok := tradesVal.([]interface{}); ok {
					result = append(result, trades...)
				} else {
					fmt.Printf("tradesVal: %+v", tradesVal)
					fmt.Printf("type tradesVal: %+v", reflect.TypeOf(tradesVal))
				}
			} else {
				fmt.Printf("pageData: %+v", pageData)
				fmt.Printf("type pageData: %+v", reflect.TypeOf(pageData))
			}
		}
	}

	assert.Equal(t, 5924, len(result))
}
