// API 错误定义
package api

import "errors"

// 定义常用错误
var (
	ErrFundCodeRequired  = errors.New("基金代码不能为空")
	ErrFundCodesRequired = errors.New("基金代码列表不能为空")
	ErrKeywordsRequired  = errors.New("关键词不能为空")
	ErrTooManyFunds      = errors.New("基金数量超过限制")
	ErrInvalidParams     = errors.New("参数无效")
	ErrDataNotFound      = errors.New("数据不存在")
	ErrInternalError     = errors.New("内部服务器错误")
)
