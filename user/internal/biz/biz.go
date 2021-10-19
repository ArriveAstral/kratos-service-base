package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUseCase, NewAddressUseCase)

const SystemInternalErrorMsg = "系统繁忙,请稍后再试"
const SystemForbiddenErrorMsg = "暂无权限访问"
const SystemBadRequestErrorMsg = "请求参数错误"
const SystemNotfoundErrorMsg = "数据不存在"
