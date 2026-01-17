// Package errcode 提供统一的错误处理和响应格式
//
// 该包实现了：
// - 自定义错误类型 AppError，支持重试策略
// - HTTP状态码映射
// - 标准化的错误响应格式
// - Retry-After头设置（RFC 7231）
//
// 使用示例：
//   err := errcode.ErrVersionNotFound.WithDetails(map[string]interface{}{
//       "version": "1.2.0",
//   })
//   errcode.RespondWithError(c, err)
package errcode

// AppError 应用错误类型
//
// 实现了error接口，包含错误码、消息、HTTP状态码和重试信息
type AppError struct {
	Code       string                 `json:"code"`                  // 自定义错误码
	Message    string                 `json:"message"`               // 错误消息
	HTTPStatus int                    `json:"-"`                     // HTTP状态码
	Details    map[string]interface{} `json:"details,omitempty"`     // 错误详情

	// 重试相关字段（内部使用，不暴露给客户端）
	Retryable  bool `json:"-"`  // 是否可重试
	RetryAfter int  `json:"-"`  // 建议重试延迟（秒）
}

// Error 实现error接口
//
// 返回错误消息
func (e *AppError) Error() string {
	return e.Message
}

// WithDetails 添加错误详情
//
// 创建一个新的AppError副本并添加详细信息
//
// 参数:
//   details: 错误详情映射
//
// 返回:
//   *AppError: 新的错误实例（包含详情）
//
// 示例:
//   err := errcode.ErrVersionNotFound.WithDetails(map[string]interface{}{
//       "version": "1.2.0",
//       "availableVersions": []string{"1.0.0", "1.1.0"},
//   })
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	err := *e
	err.Details = details
	return &err
}

// WithMessage 添加错误消息
//
// 创建一个新的AppError副本并修改错误消息
//
// 参数:
//   msg: 新的错误消息
//
// 返回:
//   *AppError: 新的错误实例（包含新消息）
//
// 示例:
//   err := errcode.ErrInvalidParams.WithMessage("版本号格式无效")
func (e *AppError) WithMessage(msg string) *AppError {
	err := *e
	err.Message = msg
	return &err
}

// ==================== 错误码定义 ====================

const (
	// 客户端错误 (4xx)

	// ErrCodeInvalidParams 参数格式错误 (400)
	ErrCodeInvalidParams = "INVALID_PARAMS"

	// ErrCodeVersionNotFound 版本不存在 (404)
	ErrCodeVersionNotFound = "VERSION_NOT_FOUND"

	// ErrCodeInvalidIPFormat IP格式无效 (400)
	ErrCodeInvalidIPFormat = "INVALID_IP_FORMAT"

	// ErrCodeInvalidURLFormat URL格式无效 (400)
	ErrCodeInvalidURLFormat = "INVALID_URL_FORMAT"

	// ErrCodeInvalidEnv 环境参数无效 (400)
	ErrCodeInvalidEnv = "INVALID_ENV"

	// ErrCodeRequestBodyTooLarge 请求体过大 (400)
	ErrCodeRequestBodyTooLarge = "REQUEST_BODY_TOO_LARGE"

	// ErrCodeRequestBodyReadError 请求体读取错误 (500)
	ErrCodeRequestBodyReadError = "REQUEST_BODY_READ_ERROR"

	// ErrCodeAccessDenied 访问被拒绝 (403)
	ErrCodeAccessDenied = "ACCESS_DENIED"

	// ErrCodeAuthenticationFailed 认证失败 (401)
	ErrCodeAuthenticationFailed = "AUTHENTICATION_FAILED"

	// ErrCodeGMConfigNotFound GM配置不存在 (404)
	ErrCodeGMConfigNotFound = "GM_CONFIG_NOT_FOUND"

	// ErrCodeIPNotFound IP不在白名单 (404)
	ErrCodeIPNotFound = "IP_NOT_FOUND"

	// ErrCodeDuplicateIP IP重复添加 (400)
	ErrCodeDuplicateIP = "DUPLICATE_IP"

	// 服务器错误 (5xx)

	// ErrCodeInternalError 内部错误 (500)
	ErrCodeInternalError = "INTERNAL_ERROR"

	// ErrCodeRedisError Redis操作失败 (500)
	ErrCodeRedisError = "REDIS_ERROR"

	// ErrCodeRedisTimeout Redis操作超时 (504)
	ErrCodeRedisTimeout = "REDIS_TIMEOUT"

	// ErrCodeTimeout 请求超时 (504)
	ErrCodeTimeout = "REQUEST_TIMEOUT"

	// ErrCodeJSONError JSON解析失败 (500)
	ErrCodeJSONError = "JSON_PARSE_ERROR"

	// ErrCodeRedisConnError Redis连接失败 (503)
	ErrCodeRedisConnError = "REDIS_CONNECTION_ERROR"
)

// ==================== 预定义错误 ====================
//
// 以下变量提供了常用的预定义错误实例，可以直接使用或通过WithDetails/WithMessage方法定制

var (
	// 客户端错误 4xx

	// ErrInvalidParams 参数错误 (400)
	//
	// 用于表示请求参数格式错误或缺失必要参数
	ErrInvalidParams = &AppError{
		Code:       ErrCodeInvalidParams,
		Message:    "参数格式错误",
		HTTPStatus: 400,
		Retryable:  false,
	}

	// ErrVersionNotFound 版本不存在 (404)
	//
	// 用于表示请求的版本号不存在
	ErrVersionNotFound = &AppError{
		Code:       ErrCodeVersionNotFound,
		Message:    "版本不存在",
		HTTPStatus: 404,
		Retryable:  false,
	}

	// ErrInvalidIPFormat IP格式无效 (400)
	//
	// 用于表示IP地址格式不符合要求
	ErrInvalidIPFormat = &AppError{
		Code:       ErrCodeInvalidIPFormat,
		Message:    "IP地址格式无效",
		HTTPStatus: 400,
		Retryable:  false,
	}

	// ErrInvalidURLFormat URL格式无效 (400)
	//
	// 用于表示URL格式不正确，必须以http://或https://开头
	ErrInvalidURLFormat = &AppError{
		Code:       ErrCodeInvalidURLFormat,
		Message:    "URL格式无效，必须以http://或https://开头",
		HTTPStatus: 400,
		Retryable:  false,
	}

	// ErrInvalidEnv 环境参数无效 (400)
	//
	// 用于表示环境参数不在允许的范围内（必须是dev或pro）
	ErrInvalidEnv = &AppError{
		Code:       ErrCodeInvalidEnv,
		Message:    "环境参数无效，必须是dev或pro",
		HTTPStatus: 400,
		Retryable:  false,
	}

	// ErrRequestBodyTooLarge 请求体过大 (400)
	//
	// 用于表示请求体超过最大限制（1MB）
	ErrRequestBodyTooLarge = &AppError{
		Code:       ErrCodeRequestBodyTooLarge,
		Message:    "请求体过大，最大限制1MB",
		HTTPStatus: 400,
		Retryable:  false,
	}

	// ErrRequestBodyReadError 请求体读取错误 (500)
	//
	// 用于表示请求体读取失败
	ErrRequestBodyReadError = &AppError{
		Code:       ErrCodeRequestBodyReadError,
		Message:    "请求体读取失败",
		HTTPStatus: 500,
		Retryable:  false,
	}

	// ErrAccessDenied 访问拒绝 (403)
	//
	// 用于表示访问被拒绝（如IP白名单检查失败）
	ErrAccessDenied = &AppError{
		Code:       ErrCodeAccessDenied,
		Message:    "访问被拒绝",
		HTTPStatus: 403,
		Retryable:  false,
	}

	// ErrAuthenticationFailed 认证失败 (401)
	//
	// 用于表示身份认证失败
	ErrAuthenticationFailed = &AppError{
		Code:       ErrCodeAuthenticationFailed,
		Message:    "认证失败",
		HTTPStatus: 401,
		Retryable:  false,
	}

	// ErrGMConfigNotFound GM配置不存在 (404)
	//
	// 用于表示GM配置不存在
	ErrGMConfigNotFound = &AppError{
		Code:       ErrCodeGMConfigNotFound,
		Message:    "GM配置不存在",
		HTTPStatus: 404,
		Retryable:  false,
	}

	// ErrIPNotFound IP不在白名单 (404)
	//
	// 用于表示IP不在白名单中
	ErrIPNotFound = &AppError{
		Code:       ErrCodeIPNotFound,
		Message:    "IP不在白名单",
		HTTPStatus: 404,
		Retryable:  false,
	}

	// ErrDuplicateIP IP重复添加 (400)
	//
	// 用于表示IP已存在于白名单中
	ErrDuplicateIP = &AppError{
		Code:       ErrCodeDuplicateIP,
		Message:    "IP已存在于白名单中",
		HTTPStatus: 400,
		Retryable:  false,
	}

	// 服务器错误 5xx

	// ErrInternalError 内部错误 (500)
	//
	// 用于表示服务器内部错误
	ErrInternalError = &AppError{
		Code:       ErrCodeInternalError,
		Message:    "服务器内部错误",
		HTTPStatus: 500,
		Retryable:  false,
	}

	// ErrRedisError Redis错误 (500)
	//
	// 用于表示Redis操作失败，可重试（延迟1秒）
	ErrRedisError = &AppError{
		Code:       ErrCodeRedisError,
		Message:    "Redis操作失败",
		HTTPStatus: 500,
		Retryable:  true,
		RetryAfter: 1,
	}

	// ErrRedisTimeout Redis超时 (504)
	//
	// 用于表示Redis操作超时，可重试（延迟3秒）
	ErrRedisTimeout = &AppError{
		Code:       ErrCodeRedisTimeout,
		Message:    "Redis操作超时",
		HTTPStatus: 504,
		Retryable:  true,
		RetryAfter: 3,
	}

	// ErrTimeout 请求超时 (504)
	//
	// 用于表示请求超时，可重试（延迟3秒）
	ErrTimeout = &AppError{
		Code:       ErrCodeTimeout,
		Message:    "请求超时",
		HTTPStatus: 504,
		Retryable:  true,
		RetryAfter: 3,
	}

	// ErrJSONError JSON解析错误 (500)
	//
	// 用于表示JSON解析失败
	ErrJSONError = &AppError{
		Code:       ErrCodeJSONError,
		Message:    "JSON解析失败",
		HTTPStatus: 500,
		Retryable:  false,
	}

	// ErrRedisConnError Redis连接错误 (503)
	//
	// 用于表示Redis连接失败，可重试（延迟5秒）
	ErrRedisConnError = &AppError{
		Code:       ErrCodeRedisConnError,
		Message:    "Redis连接失败",
		HTTPStatus: 503,
		Retryable:  true,
		RetryAfter: 5,
	}
)

// ==================== 辅助函数 ====================

// NewParamsError 创建参数错误
//
// 创建一个自定义消息的参数错误
//
// 参数:
//   msg: 错误消息
//
// 返回:
//   *AppError: 参数错误实例（HTTP 400）
//
// 示例:
//   err := NewParamsError("版本号不能为空")
func NewParamsError(msg string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidParams,
		Message:    msg,
		HTTPStatus: 400,
	}
}

// NewVersionNotFoundError 创建版本不存在错误
//
// 创建一个版本不存在的错误，并自动添加版本详情
//
// 参数:
//   version: 不存在的版本号
//
// 返回:
//   *AppError: 版本不存在错误实例（HTTP 404）
//
// 示例:
//   err := NewVersionNotFoundError("1.2.0")
func NewVersionNotFoundError(version string) *AppError {
	return &AppError{
		Code:       ErrCodeVersionNotFound,
		Message:    "版本 " + version + " 不存在",
		HTTPStatus: 404,
		Details: map[string]interface{}{
			"version": version,
		},
	}
}

// NewInvalidVersionFormatError 创建版本格式错误
//
// 创建一个版本号格式无效的错误，并自动添加字段值详情
//
// 参数:
//   version: 无效的版本号
//
// 返回:
//   *AppError: 版本格式错误实例（HTTP 400）
//
// 示例:
//   err := NewInvalidVersionFormatError("invalid")
func NewInvalidVersionFormatError(version string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidParams,
		Message:    "版本号格式无效，应为X.Y.Z或X.Y.Z.N",
		HTTPStatus: 400,
		Details: map[string]interface{}{
			"field":  "vsn",
			"value":  version,
			"format": "X.Y.Z or X.Y.Z.N",
		},
	}
}

// NewRedisError 创建Redis错误
//
// 创建一个自定义消息的Redis操作失败错误（可重试）
//
// 参数:
//   msg: 错误消息
//
// 返回:
//   *AppError: Redis错误实例（HTTP 500，可重试，延迟1秒）
//
// 示例:
//   err := NewRedisError("Redis连接失败")
func NewRedisError(msg string) *AppError {
	return &AppError{
		Code:       ErrCodeRedisError,
		Message:    msg,
		HTTPStatus: 500,
		Retryable:  true,
		RetryAfter: 1,
	}
}

// IsAppError 判断是否为AppError
//
// 判断错误是否为AppError类型
//
// 参数:
//   err: 待判断的错误
//
// 返回:
//   bool: 如果是AppError返回true，否则返回false
//
// 示例:
//   if errcode.IsAppError(err) {
//       appErr := err.(*errcode.AppError)
//       log.Printf("Error code: %s", appErr.Code)
//   }
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError 转换为AppError
//
// 尝试将error转换为AppError类型
//
// 参数:
//   err: 待转换的错误
//
// 返回:
//   *AppError: AppError实例（如果不是AppError则返回nil）
//   bool: 是否转换成功
//
// 示例:
//   if appErr, ok := errcode.AsAppError(err); ok {
//       log.Printf("Error code: %s", appErr.Code)
//   }
func AsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// GetHTTPStatus 获取错误对应的HTTP状态码
//
// 根据错误码获取对应的HTTP状态码
//
// 参数:
//   code: 错误码
//
// 返回:
//   int: HTTP状态码（如果未找到则返回500）
//
// 示例:
//   status := errcode.GetHTTPStatus("VERSION_NOT_FOUND")  // 返回404
func GetHTTPStatus(code string) int {
	if status, ok := errorCodeToHTTPStatus[code]; ok {
		return status
	}
	return 500 // 默认返回500
}

// IsRetryable 判断错误是否可重试
//
// 根据错误码判断错误是否可重试
//
// 参数:
//   code: 错误码
//
// 返回:
//   bool: 如果可重试返回true，否则返回false
//
// 示例:
//   if errcode.IsRetryable("REDIS_ERROR") {
//       // 重试逻辑
//   }
func IsRetryable(code string) bool {
	return retryableCodes[code]
}

// GetRetryAfter 获取错误建议的重试延迟（秒）
//
// 根据错误码获取建议的重试延迟时间
//
// 参数:
//   code: 错误码
//
// 返回:
//   int: 重试延迟秒数（如果错误不可重试则返回0）
//
// 示例:
//   delay := errcode.GetRetryAfter("REDIS_TIMEOUT")  // 返回3
//   time.Sleep(time.Duration(delay) * time.Second)
func GetRetryAfter(code string) int {
	return retryAfterMap[code]
}

// ==================== HTTP状态码映射表 ====================
//
// 以下映射表定义了错误码到HTTP状态码的对应关系，遵循RFC 7231标准

// errorCodeToHTTPStatus 错误码到HTTP状态码的映射
//
// 该映射表定义了所有错误码对应的HTTP状态码
var errorCodeToHTTPStatus = map[string]int{
	// 400 Bad Request
	ErrCodeInvalidParams:       400,
	ErrCodeInvalidIPFormat:     400,
	ErrCodeInvalidURLFormat:    400,
	ErrCodeInvalidEnv:          400,
	ErrCodeRequestBodyTooLarge: 400,
	ErrCodeDuplicateIP:         400,

	// 401 Unauthorized
	ErrCodeAuthenticationFailed: 401,

	// 403 Forbidden
	ErrCodeAccessDenied: 403,

	// 404 Not Found
	ErrCodeVersionNotFound:  404,
	ErrCodeGMConfigNotFound: 404,
	ErrCodeIPNotFound:       404,

	// 500 Internal Server Error
	ErrCodeInternalError:       500,
	ErrCodeRequestBodyReadError: 500,
	ErrCodeRedisError:          500,
	ErrCodeJSONError:           500,

	// 503 Service Unavailable
	ErrCodeRedisConnError: 503,

	// 504 Gateway Timeout
	ErrCodeTimeout:      504,
	ErrCodeRedisTimeout: 504,
}

// retryableCodes 可重试的错误码
//
// 定义哪些错误码对应的错误是可以重试的
var retryableCodes = map[string]bool{
	ErrCodeRedisError:      true,
	ErrCodeRedisTimeout:    true,
	ErrCodeTimeout:         true,
	ErrCodeRedisConnError:  true,
}

// retryAfterMap 错误码到重试延迟的映射（秒）
//
// 定义每个可重试错误码建议的重试延迟时间（遵循RFC 7231 Retry-After头）
var retryAfterMap = map[string]int{
	ErrCodeRedisError:     1,  // Redis错误1秒后重试
	ErrCodeRedisTimeout:   3,  // Redis超时3秒后重试
	ErrCodeTimeout:        3,  // 超时3秒后重试
	ErrCodeRedisConnError: 5,  // 连接错误5秒后重试
}
