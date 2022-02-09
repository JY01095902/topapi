package request

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func sortParams(params url.Values) []string {
	rawParamSlice := []string{}
	for k := range params {
		val := params.Get(k)
		rawParamSlice = append(rawParamSlice, k+val)
	}
	sort.Strings(rawParamSlice)

	return rawParamSlice
}

func encryptSHA256(content, secret string) []byte {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	_, err := h.Write([]byte(content))
	if err != nil {
		return nil
	}

	return h.Sum(nil)
}

func signBySHA256(params url.Values, appSecret string) string {
	paramString := strings.Join(sortParams(params), "")
	signature := encryptSHA256(paramString, appSecret)

	return fmt.Sprintf("%X", signature)
}

func signByMD5(params url.Values, appSecret string) string {
	paramString := strings.Join(sortParams(params), "")
	signature := md5.Sum([]byte(appSecret + paramString + appSecret))

	return fmt.Sprintf("%X", signature)
}

func makeQueryString(params url.Values, appSecret, signMethod string) string {
	paramSlice := []string{}
	for k := range params {
		val := params.Get(k)
		paramSlice = append(paramSlice, k+"="+url.QueryEscape(val))
	}

	sign := createSignFunc(signMethod)
	signature := sign(params, appSecret)
	paramSlice = append(paramSlice, "sign="+signature)

	return strings.Join(paramSlice, "&")
}

type sign func(params url.Values, appSecret string) string

func createSignFunc(signMethod string) sign {
	fns := map[string]sign{
		"hmac-sha256": signBySHA256,
		"md5":         signByMD5,
	}

	if fn, exist := fns[signMethod]; exist {
		return fn
	}

	return signBySHA256
}
