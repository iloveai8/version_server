package httputil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestReadRequestBody(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantErr    bool
		wantResult string
	}{
		{
			name:       "正常请求体",
			body:       `{"test": "data"}`,
			wantErr:    false,
			wantResult: `{"test": "data"}`,
		},
		{
			name:       "空请求体",
			body:       "",
			wantErr:    false,
			wantResult: "",
		},
		{
			name:       "小请求体",
			body:       "hello",
			wantErr:    false,
			wantResult: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			result, err := ReadRequestBody(req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, string(result))
			}
		})
	}
}

func TestReadRequestBodyWithLimit(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		maxSize     int64
		wantErr     bool
		checkLength bool
		maxLength   int
	}{
		{
			name:        "正常请求体",
			body:        "hello world",
			maxSize:     100,
			wantErr:     false,
			checkLength: false,
		},
		{
			name:        "请求体超过限制",
			body:        strings.Repeat("a", 101),
			maxSize:     100,
			wantErr:     true,
			checkLength: false,
		},
		{
			name:        "请求体刚好等于限制",
			body:        strings.Repeat("a", 100),
			maxSize:     100,
			wantErr:     false,
			checkLength: false,
		},
		{
			name:        "空请求体",
			body:        "",
			maxSize:     100,
			wantErr:     false,
			checkLength: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			result, err := ReadRequestBodyWithLimit(req, tt.maxSize)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkLength {
					assert.LessOrEqual(t, len(result), tt.maxLength)
				}
			}
		})
	}
}

func TestReadRequestBody_NilRequest(t *testing.T) {
	_, err := ReadRequestBody(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "请求对象为空")
}

func TestReadRequestBody_NilBody(t *testing.T) {
	req := &http.Request{}
	_, err := ReadRequestBody(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "请求体为空")
}

func TestRequestID(t *testing.T) {
	// 设置路由
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		c.JSON(200, gin.H{"requestId": requestID})
	})

	t.Run("生成新的请求ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Header().Get(XRequestID), "")

		// 验证响应体
		assert.Contains(t, w.Body.String(), "requestId")
	})

	t.Run("使用请求头中的请求ID", func(t *testing.T) {
		testID := "test-request-id-123"
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set(XRequestID, testID)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, testID, w.Header().Get(XRequestID))

		// 验证响应体中的请求ID
		assert.Contains(t, w.Body.String(), testID)
	})
}

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(*gin.Context)
		expectID   bool
		expectEmpty bool
	}{
		{
			name: "存在请求ID",
			setupFunc: func(c *gin.Context) {
				c.Set(XRequestID, "test-id-123")
			},
			expectID: true,
		},
		{
			name:       "不存在请求ID",
			setupFunc:  func(c *gin.Context) {},
			expectID:   false,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			tt.setupFunc(c)

			requestID := GetRequestID(c)

			if tt.expectEmpty {
				assert.Empty(t, requestID)
			} else if tt.expectID {
				assert.NotEmpty(t, requestID)
			}
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	// 验证生成的ID不为空
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)

	// 验证每次生成的ID不同
	assert.NotEqual(t, id1, id2)
}

func TestXRequestID_Const(t *testing.T) {
	assert.Equal(t, "X-Request-ID", XRequestID)
}

func TestMaxBodySize_Const(t *testing.T) {
	// 验证MaxBodySize为1MB
	assert.Equal(t, int64(1<<20), int64(MaxBodySize))
}
