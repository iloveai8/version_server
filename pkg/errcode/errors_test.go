package errcode

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := &AppError{
		Code:    ErrCodeInvalidParams,
		Message: "参数错误",
	}
	assert.Equal(t, "参数错误", err.Error())
}

func TestAppError_WithDetails(t *testing.T) {
	baseErr := &AppError{
		Code:    ErrCodeVersionNotFound,
		Message: "版本不存在",
	}
	details := map[string]interface{}{
		"version": "1.0.0",
	}
	err := baseErr.WithDetails(details)

	assert.Equal(t, ErrCodeVersionNotFound, err.Code)
	assert.Equal(t, "版本不存在", err.Message)
	assert.Equal(t, "1.0.0", err.Details["version"])

	// 原始错误不应被修改
	assert.Nil(t, baseErr.Details)
}

func TestAppError_WithMessage(t *testing.T) {
	baseErr := &AppError{
		Code:    ErrCodeInternalError,
		Message: "内部错误",
	}
	err := baseErr.WithMessage("自定义错误消息")

	assert.Equal(t, ErrCodeInternalError, err.Code)
	assert.Equal(t, "自定义错误消息", err.Message)

	// 原始错误不应被修改
	assert.Equal(t, "内部错误", baseErr.Message)
}

func TestIsAppError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "AppError",
			err:  ErrInvalidParams,
			want: true,
		},
		{
			name: "标准error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsAppError(tt.err))
		})
	}
}

func TestAsAppError(t *testing.T) {
	// 测试AppError
	appErr := ErrVersionNotFound
	result, ok := AsAppError(appErr)
	assert.True(t, ok)
	assert.Equal(t, appErr, result)

	// 测试标准error
	stdErr := errors.New("standard error")
	result, ok = AsAppError(stdErr)
	assert.False(t, ok)
	assert.Nil(t, result)

	// 测试nil
	result, ok = AsAppError(nil)
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		code     string
		expected int
	}{
		{ErrCodeInvalidParams, 400},
		{ErrCodeVersionNotFound, 404},
		{ErrCodeInternalError, 500},
		{ErrCodeRedisConnError, 503},
		{ErrCodeTimeout, 504},
		{"UNKNOWN_CODE", 500}, // 默认500
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetHTTPStatus(tt.code))
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{ErrCodeTimeout, true},
		{ErrCodeRedisError, true},
		{ErrCodeRedisConnError, true},
		{ErrCodeInvalidParams, false},
		{ErrCodeVersionNotFound, false},
		{ErrCodeInternalError, false},
		{"UNKNOWN_CODE", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsRetryable(tt.code))
		})
	}
}

func TestGetRetryAfter(t *testing.T) {
	tests := []struct {
		code     string
		expected int
	}{
		{ErrCodeTimeout, 3},
		{ErrCodeRedisError, 1},
		{ErrCodeRedisConnError, 5},
		{ErrCodeInvalidParams, 0},
		{ErrCodeVersionNotFound, 0},
		{"UNKNOWN_CODE", 0},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetRetryAfter(tt.code))
		})
	}
}

func TestNewParamsError(t *testing.T) {
	msg := "自定义参数错误"
	err := NewParamsError(msg)

	assert.Equal(t, ErrCodeInvalidParams, err.Code)
	assert.Equal(t, msg, err.Message)
	assert.Equal(t, 400, err.HTTPStatus)
	assert.False(t, err.Retryable)
}

func TestNewVersionNotFoundError(t *testing.T) {
	version := "1.2.0"
	err := NewVersionNotFoundError(version)

	assert.Equal(t, ErrCodeVersionNotFound, err.Code)
	assert.Contains(t, err.Message, version)
	assert.Equal(t, 404, err.HTTPStatus)
	assert.Equal(t, version, err.Details["version"])
	assert.False(t, err.Retryable)
}

func TestNewInvalidVersionFormatError(t *testing.T) {
	version := "invalid"
	err := NewInvalidVersionFormatError(version)

	assert.Equal(t, ErrCodeInvalidParams, err.Code)
	assert.Contains(t, err.Message, "版本号格式无效")
	assert.Equal(t, 400, err.HTTPStatus)
	assert.Equal(t, "vsn", err.Details["field"])
	assert.Equal(t, version, err.Details["value"])
	assert.False(t, err.Retryable)
}

func TestNewRedisError(t *testing.T) {
	msg := "Redis连接失败"
	err := NewRedisError(msg)

	assert.Equal(t, ErrCodeRedisError, err.Code)
	assert.Equal(t, msg, err.Message)
	assert.Equal(t, 500, err.HTTPStatus)
	assert.True(t, err.Retryable)
	assert.Equal(t, 1, err.RetryAfter)
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          *AppError
		code         string
		httpStatus   int
		retryable    bool
		retryAfter   int
	}{
		{
			name:       "ErrInvalidParams",
			err:        ErrInvalidParams,
			code:       ErrCodeInvalidParams,
			httpStatus: 400,
			retryable:  false,
		},
		{
			name:       "ErrVersionNotFound",
			err:        ErrVersionNotFound,
			code:       ErrCodeVersionNotFound,
			httpStatus: 404,
			retryable:  false,
		},
		{
			name:       "ErrAccessDenied",
			err:        ErrAccessDenied,
			code:       ErrCodeAccessDenied,
			httpStatus: 403,
			retryable:  false,
		},
		{
			name:       "ErrAuthenticationFailed",
			err:        ErrAuthenticationFailed,
			code:       ErrCodeAuthenticationFailed,
			httpStatus: 401,
			retryable:  false,
		},
		{
			name:       "ErrGMConfigNotFound",
			err:        ErrGMConfigNotFound,
			code:       ErrCodeGMConfigNotFound,
			httpStatus: 404,
			retryable:  false,
		},
		{
			name:       "ErrIPNotFound",
			err:        ErrIPNotFound,
			code:       ErrCodeIPNotFound,
			httpStatus: 404,
			retryable:  false,
		},
		{
			name:       "ErrInternalError",
			err:        ErrInternalError,
			code:       ErrCodeInternalError,
			httpStatus: 500,
			retryable:  false,
		},
		{
			name:       "ErrRedisError",
			err:        ErrRedisError,
			code:       ErrCodeRedisError,
			httpStatus: 500,
			retryable:  true,
			retryAfter: 1,
		},
		{
			name:       "ErrRedisTimeout",
			err:        ErrRedisTimeout,
			code:       ErrCodeRedisTimeout,
			httpStatus: 504,
			retryable:  true,
			retryAfter: 3,
		},
		{
			name:       "ErrTimeout",
			err:        ErrTimeout,
			code:       ErrCodeTimeout,
			httpStatus: 504,
			retryable:  true,
			retryAfter: 3,
		},
		{
			name:       "ErrJSONError",
			err:        ErrJSONError,
			code:       ErrCodeJSONError,
			httpStatus: 500,
			retryable:  false,
		},
		{
			name:       "ErrRedisConnError",
			err:        ErrRedisConnError,
			code:       ErrCodeRedisConnError,
			httpStatus: 503,
			retryable:  true,
			retryAfter: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.err.Code)
			assert.Equal(t, tt.httpStatus, tt.err.HTTPStatus)
			assert.Equal(t, tt.retryable, tt.err.Retryable)
			if tt.retryable {
				assert.Equal(t, tt.retryAfter, tt.err.RetryAfter)
			}
		})
	}
}

func TestHTTPStatusMapping(t *testing.T) {
	// 验证所有预定义错误的HTTP状态码与映射表一致
	predefinedErrors := []*AppError{
		ErrInvalidParams,
		ErrVersionNotFound,
		ErrInvalidIPFormat,
		ErrInvalidURLFormat,
		ErrInvalidEnv,
		ErrRequestBodyTooLarge,
		ErrAccessDenied,
		ErrAuthenticationFailed,
		ErrGMConfigNotFound,
		ErrIPNotFound,
		ErrInternalError,
		ErrRequestBodyReadError,
		ErrRedisError,
		ErrRedisTimeout,
		ErrTimeout,
		ErrJSONError,
		ErrRedisConnError,
	}

	for _, err := range predefinedErrors {
		t.Run(err.Code, func(t *testing.T) {
			mappedStatus := GetHTTPStatus(err.Code)
			assert.Equal(t, err.HTTPStatus, mappedStatus,
				"HTTP status mismatch for %s: predefined=%d, mapped=%d",
				err.Code, err.HTTPStatus, mappedStatus)
		})
	}
}
