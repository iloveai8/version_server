package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "有效版本号-1.0.0",
			input:   "1.0.0",
			wantErr: false,
		},
		{
			name:    "有效版本号-1.2.3",
			input:   "1.2.3",
			wantErr: false,
		},
		{
			name:    "有效版本号-1.2.3.4",
			input:   "1.2.3.4",
			wantErr: false,
		},
		{
			name:    "有效版本号-10.20.30",
			input:   "10.20.30",
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:    "无效版本号-1.0",
			input:   "1.0",
			wantErr: true,
		},
		{
			name:    "无效版本号-abc",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "无效版本号-1.0.0.0.0",
			input:   "1.0.0.0.0",
			wantErr: true,
		},
		{
			name:    "无效版本号-v1.0.0",
			input:   "v1.0.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateVersionFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{
			name:   "有效格式-1.0.0",
			input:  "1.0.0",
			expect: true,
		},
		{
			name:   "有效格式-1.2.3.4",
			input:  "1.2.3.4",
			expect: true,
		},
		{
			name:   "无效格式-1.0",
			input:  "1.0",
			expect: false,
		},
		{
			name:   "无效格式-abc",
			input:  "abc",
			expect: false,
		},
		{
			name:   "空字符串",
			input:  "",
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateVersionFormat(tt.input)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "有效的HTTP URL",
			input:   "http://example.com",
			wantErr: false,
		},
		{
			name:    "有效的HTTPS URL",
			input:   "https://example.com",
			wantErr: false,
		},
		{
			name:    "有效的带路径URL",
			input:   "https://example.com/path/to/resource",
			wantErr: false,
		},
		{
			name:    "有效的带端口URL",
			input:   "http://example.com:8080",
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:    "缺少协议",
			input:   "example.com",
			wantErr: true,
		},
		{
			name:    "无效协议",
			input:   "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "缺少主机",
			input:   "http://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHTTPURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "有效的HTTP URL",
			input:   "http://example.com",
			wantErr: false,
		},
		{
			name:    "有效的HTTPS URL",
			input:   "https://example.com",
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:    "缺少协议",
			input:   "example.com",
			wantErr: true,
		},
		{
			name:    "FTP协议",
			input:   "ftp://example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHTTPURL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEnv(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "开发环境",
			input:   "dev",
			wantErr: false,
		},
		{
			name:    "生产环境",
			input:   "pro",
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:    "无效环境",
			input:   "test",
			wantErr: true,
		},
		{
			name:    "无效环境-大写",
			input:   "DEV",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnv(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "有效的IP地址",
			input:   "192.168.1.1",
			wantErr: false,
		},
		{
			name:    "有效的IP地址-10段",
			input:   "10.0.0.1",
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:    "无效的IP地址",
			input:   "192.168.1",
			wantErr: true,
		},
		{
			name:    "无效的IP地址-字符",
			input:   "abc.def.ghi.jkl",
			wantErr: true,
		},
		{
			name:    "无效的IP地址-超出范围",
			input:   "256.256.256.256",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIP(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePageAndPageSize(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		wantErr     bool
	}{
		{
			name:     "有效参数",
			page:     1,
			pageSize: 10,
			wantErr:  false,
		},
		{
			name:     "页码为0",
			page:     0,
			pageSize: 10,
			wantErr:  true,
		},
		{
			name:     "页码为负数",
			page:     -1,
			pageSize: 10,
			wantErr:  true,
		},
		{
			name:     "每页数量为0",
			page:     1,
			pageSize: 0,
			wantErr:  true,
		},
		{
			name:     "每页数量超过100",
			page:     1,
			pageSize: 101,
			wantErr:  true,
		},
		{
			name:     "每页数量为100（边界值）",
			page:     1,
			pageSize: 100,
			wantErr:  false,
		},
		{
			name:     "每页数量为1（边界值）",
			page:     1,
			pageSize: 1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePageAndPageSize(tt.page, tt.pageSize)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
