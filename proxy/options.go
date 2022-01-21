package proxy

import (
	"net/url"
	"strconv"
)

type option func(val url.Values)

func WithOrderCode(orderCode string) option {
	return func(val url.Values) {
		val.Set("order_code", orderCode)
	}
}

func WithOrderCodeType(orderCodeType string) option {
	return func(val url.Values) {
		val.Set("order_code_type", orderCodeType)
	}
}

func WithPageNumber(pageNum int) option {
	return func(val url.Values) {
		val.Set("page_no", strconv.Itoa(pageNum))
	}
}

func WithFields(fields string) option {
	return func(val url.Values) {
		val.Set("fields", fields)
	}
}
