package errcode

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSetRetryHeaders(t *testing.T) {
	tests := []struct {
		name       string
		err        *AppError
		wantHeader string
	}{
		{
			name: "可重试错误",
			err: &AppError{
				Code:       ErrCodeTimeout,
				HTTPStatus: 504,
				Retryable:  true,
				RetryAfter: 3,
			},
			wantHeader: "3",
		},
		{
			name: "不可重试错误",
			err: &AppError{
				Code:       ErrCodeInvalidParams,
				HTTPStatus: 400,
				Retryable:  false,
			},
			wantHeader: "",
		},
		{
			name: "RetryAfter为0",
			err: &AppError{
				Code:       ErrCodeRedisError,
				HTTPStatus: 500,
				Retryable:  true,
				RetryAfter: 0,
			},
			wantHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			SetRetryHeaders(c, tt.err)

			header := c.Writer.Header().Get("Retry-After")
			assert.Equal(t, tt.wantHeader, header)
		})
	}
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedBody string
		checkRetry   bool
		retryAfter   string
	}{
		{
			name:         "AppError - 参数错误",
			err:          ErrInvalidParams,
			expectedCode: 400,
			expectedBody: `{"error":{"code":"INVALID_PARAMS","message":"参数格式错误"}}`,
			checkRetry:   false,
		},
		{
			name:         "AppError - 超时（可重试）",
			err:          ErrTimeout,
			expectedCode: 504,
			expectedBody: `{"error":{"code":"REQUEST_TIMEOUT","message":"请求超时"}}`,
			checkRetry:   true,
			retryAfter:   "3",
		},
		{
			name:         "AppError - Redis连接错误（可重试）",
			err:          ErrRedisConnError,
			expectedCode: 503,
			expectedBody: `{"error":{"code":"REDIS_CONNECTION_ERROR","message":"Redis连接失败"}}`,
			checkRetry:   true,
			retryAfter:   "5",
		},
		{
			name:         "AppError - 版本不存在",
			err:          ErrVersionNotFound,
			expectedCode: 404,
			expectedBody: `{"error":{"code":"VERSION_NOT_FOUND","message":"版本不存在"}}`,
			checkRetry:   false,
		},
		{
			name:         "标准error",
			err:          assert.AnError,
			expectedCode: 500,
			expectedBody: `{"error":{"code":"INTERNAL_ERROR","message":"服务器内部错误"}}`,
			checkRetry:   false,
		},
		{
			name:         "nil error",
			err:          nil,
			expectedCode: 500,
			expectedBody: `{"error":{"code":"INTERNAL_ERROR","message":"未知错误"}}`,
			checkRetry:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			RespondWithError(c, tt.err)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			if tt.checkRetry {
				retryAfter := w.Header().Get("Retry-After")
				assert.Equal(t, tt.retryAfter, retryAfter, "Retry-After header mismatch")
			}
		})
	}
}

func TestRespondWithSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"key": "value"}
	RespondWithSuccess(c, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"data":{"key":"value"}}`, w.Body.String())
}

func TestRespondWithSuccessAndStatus(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"key": "value"}
	RespondWithSuccessAndStatus(c, http.StatusCreated, data)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"data":{"key":"value"}}`, w.Body.String())
}

func TestRespondWithCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"id": "123"}
	RespondWithCreated(c, data)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"data":{"id":"123"}}`, w.Body.String())
}

func TestRespondWithNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	RespondWithNoContent(c)

	// gin的c.Status()会设置状态码，但在测试中需要检查writer的状态
	// 由于c.Status()不会自动写入响应，我们检查Header中的状态
	status := c.Writer.Status()
	assert.Equal(t, http.StatusNoContent, status)
	assert.Empty(t, w.Body.String())
}

func TestRespondWithBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	message := "无效的参数"
	RespondWithBadRequest(c, message)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":{"code":"INVALID_PARAMS","message":"无效的参数"}}`, w.Body.String())
}

func TestRespondWithNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	message := "资源不存在"
	RespondWithNotFound(c, message)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":{"code":"VERSION_NOT_FOUND","message":"资源不存在"}}`, w.Body.String())
}

func TestRespondWithInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	message := "服务器出错了"
	RespondWithInternalError(c, message)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":{"code":"INTERNAL_ERROR","message":"服务器出错了"}}`, w.Body.String())
}

func TestRespondWithTimeout(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	message := "请求超时"
	RespondWithTimeout(c, message)

	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
	assert.JSONEq(t, `{"error":{"code":"REQUEST_TIMEOUT","message":"请求超时"}}`, w.Body.String())
}

func TestErrorResponse_Struct(t *testing.T) {
	errInfo := &ErrorInfo{
		Code:    ErrCodeInvalidParams,
		Message: "参数错误",
		Details: map[string]interface{}{
			"field": "vsn",
		},
	}
	resp := ErrorResponse{Error: errInfo}

	assert.Equal(t, errInfo, resp.Error)
	assert.Equal(t, ErrCodeInvalidParams, resp.Error.Code)
	assert.Equal(t, "参数错误", resp.Error.Message)
	assert.Equal(t, "vsn", resp.Error.Details["field"])
}

func TestSuccessResponse_Struct(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := SuccessResponse{Data: data}

	assert.Equal(t, data, resp.Data)
	assert.Equal(t, "value", resp.Data.(map[string]string)["key"])
}
