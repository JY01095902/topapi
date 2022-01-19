package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jy01095902/gokits/elves"
)

type config map[string]string

// default APIRequest confit
func newConfig(appKey, sessionKey string) config {
	config := map[string]string{
		"app_key":     appKey,
		"session":     sessionKey,
		"v":           "2.0",
		"format":      "json",
		"simplify":    "true",
		"sign_method": "hmac-sha256",
	}

	return config
}

type option func(config config)

func WithVersion(version string) option {
	return func(config config) {
		config["v"] = version
	}
}

func WithFormat(format string) option {
	return func(config config) {
		config["format"] = format
	}
}

func WithSimplify(isSimplified bool) option {
	return func(config config) {
		val := "false"
		if isSimplified {
			val = "true"
		}
		config["simplify"] = val
	}
}

type APIRequest struct {
	baseURL   string
	appSecret string
	config    config
}

func NewRequest(appKey, appSecret, sessionKey string, options ...option) APIRequest {
	req := APIRequest{
		baseURL:   "https://eco.taobao.com/router/rest",
		appSecret: appSecret,
		config:    newConfig(appKey, sessionKey),
	}

	for _, opt := range options {
		opt(req.config)
	}

	return req
}

func (req APIRequest) checkAPIResponse(val Values) error {
	type ErrResponse struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"msg"`
		} `json:"error_response"`
	}

	result, err := val.GetResult(ErrResponse{})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrTOPAPIBizError, err.Error())
	}

	resp, ok := result.(*ErrResponse)
	if !ok {
		return fmt.Errorf("%w: result is not ErrResponse", ErrTOPAPIBizError)
	}

	if resp.Error.Code != 0 {
		return fmt.Errorf("%w: %s", ErrTOPAPIBizError, resp.Error.Message)
	}

	return nil
}

func (req APIRequest) execute(r *resty.Request, method, url string) (Values, error) {
	resp, err := r.
		SetResult(Values{}).
		Execute(method, url)

	if err != nil {
		return nil, fmt.Errorf("%w error: %s", ErrCallTOPAPIFailed, err.Error())
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%w error: %s", ErrCallTOPAPIFailed, resp.String())
	}

	var result Values
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return Values{}, fmt.Errorf("%w: type of result is not Values", ErrTOPAPIBizError)
	}

	if err := req.checkAPIResponse(result); err != nil {
		return Values{}, err
	}

	return result, err
}

func (req APIRequest) Get(apiName string, params url.Values) (Values, error) {
	mergedParams := url.Values{}
	mergedParams.Add("method", apiName)
	mergedParams.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	for k := range params {
		mergedParams.Add(k, params.Get(k))
	}
	for k, v := range req.config {
		mergedParams.Add(k, v)
	}
	r := resty.New().R().
		EnableTrace()

	url := req.baseURL + "?" + makeQueryString(mergedParams, req.appSecret)

	return req.execute(r, resty.MethodGet, url)
}

// 从结果中解析总数据量，因为postAll方法不知道返回值的格式
type parseTotalFunc func(vals Values) int

func (req APIRequest) GetMock(apiName string, params url.Values) (Values, error) {
	ms := time.Now().Nanosecond() / 1000000
	if ms > 100 && ms < 500 {
		return nil, errors.New("访问太频繁")
	}

	total := 5924
	data := make([]Values, total)
	for i := 0; i < total; i++ {
		data[i] = Values{
			"id": i + 1,
		}
	}
	pageNumber, err := strconv.Atoi(params.Get("page_no"))
	if err != nil {
		return nil, err
	}
	pageSize, err := strconv.Atoi(params.Get("page_size"))
	if err != nil {
		return nil, err
	}

	start := (pageNumber - 1) * pageSize
	end := pageNumber * pageSize
	if end > len(data) {
		end = len(data)
	}
	return Values{
		"total_results": total,
		"trades":        data[start:end],
	}, nil
}
func (req APIRequest) GetAll(apiName string, params url.Values, parseTotal parseTotalFunc) ([]Values, error) {
	copyParams := func(params url.Values) url.Values {
		result := url.Values{}
		for k := range params {
			result.Add(k, params.Get(k))
		}

		return result
	}

	// 每页返回的结果，客户端自己解析里边的内容
	result := []Values{} // 每一条是一页的数据
	pageNumber, pageSize := 1, 100
	copiedParams := copyParams(params)
	copiedParams.Set("page_no", strconv.Itoa(pageNumber))
	copiedParams.Set("page_size", strconv.Itoa(pageSize))

	resp, err := req.Get(apiName, copiedParams)
	if err != nil {
		return nil, err
	}
	// 先把第一页的结果加进去，之后从第二页开始查询
	result = append(result, resp)

	total := parseTotal(resp)
	pageCnt := int(math.Ceil(float64(total) / float64(pageSize)))
	if pageCnt == 0 {
		return nil, errors.New("not found")
	}

	valsChan := make(chan Values)
	var wgRes sync.WaitGroup
	wgRes.Add(1)
	go func() {
		defer wgRes.Done()

		for vals := range valsChan {
			result = append(result, vals)
			// log.Printf("result len: %+v", len(result))
		}
	}()

	pool, err := elves.NewPool(10)
	if err != nil {
		return nil, err
	}

	failedChan := make(chan url.Values)
	var wgFailed sync.WaitGroup
	wgFailed.Add(1)
	go func() {
		defer wgFailed.Done()

		for params := range failedChan {
		Retry:
			// time.Sleep(100 * time.Millisecond)
			log.Printf("retry body: %+v", params)
			if resp, err := req.Get(apiName, params); err == nil {
				valsChan <- resp
				log.Printf("retry  ok: %s", "")
			} else {
				log.Printf("call api retry failed, error: %s", err.Error())

				goto Retry
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(pageCnt - 1)
	for pageNum := 2; pageNum <= pageCnt; pageNum++ {
		pool.Execute(func() {
			defer wg.Done()

			copiedParams := copyParams(params)
			copiedParams.Set("page_no", strconv.Itoa(pageNum))
			copiedParams.Set("page_size", strconv.Itoa(pageSize))
			// log.Printf("pageNumber: %s, pageSize: %s", copiedParams.Get("page_no"), copiedParams.Get("page_size"))
			if resp, err := req.Get(apiName, copiedParams); err == nil {
				valsChan <- resp
			} else {
				log.Printf("call api failed, error: %s", err.Error())

				failedChan <- copiedParams
			}
		})

		// time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	close(failedChan)
	pool.Destroy()

	wgFailed.Wait()
	close(valsChan)

	wgRes.Wait()

	return result, nil
}
