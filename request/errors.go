package request

import "errors"

var (
	ErrCallTOPAPIFailed = errors.New("call Taobao Open Platform API failed")
	ErrTOPAPIBizError   = errors.New("call Taobao Open Platform API biz error")
)
