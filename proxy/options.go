package proxy

import (
	"net/url"
	"strconv"
	"time"
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

func WithPageNumber(pageNumber int) option {
	return func(val url.Values) {
		val.Set("page_no", strconv.Itoa(pageNumber))
	}
}

func WithPageIndex(pageIndex int) option {
	return func(val url.Values) {
		val.Set("page_index", strconv.Itoa(pageIndex))
	}
}

func WithPageSize(pageSize int) option {
	return func(val url.Values) {
		val.Set("page_size", strconv.Itoa(pageSize))
	}
}

func WithFields(fields string) option {
	return func(val url.Values) {
		val.Set("fields", fields)
	}
}

func WithStartCreated(t time.Time) option {
	return func(val url.Values) {
		val.Set("start_created", t.Format("2006-01-02 15:04:05"))
	}
}

func WithEndCreated(t time.Time) option {
	return func(val url.Values) {
		val.Set("end_created", t.Format("2006-01-02 15:04:05"))
	}
}

func WithStartModified(t time.Time) option {
	return func(val url.Values) {
		val.Set("start_modified", t.Format("2006-01-02 15:04:05"))
	}
}

func WithEndModified(t time.Time) option {
	return func(val url.Values) {
		val.Set("end_modified", t.Format("2006-01-02 15:04:05"))
	}
}

func WithTId(tid string) option {
	return func(val url.Values) {
		val.Set("tid", tid)
	}
}

func WithIncludeOAId(includeOAId bool) option {
	return func(val url.Values) {
		include := "false"
		if includeOAId {
			include = "true"
		}
		val.Set("include_oaid", include)
	}
}
