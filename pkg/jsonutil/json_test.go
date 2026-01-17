package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "正常对象",
			input: map[string]interface{}{
				"name": "test",
				"value": 123,
			},
			wantErr: false,
		},
		{
			name:    "空对象",
			input:   struct{}{},
			wantErr: false,
		},
		{
			name:    "nil对象",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "字符串",
			input:   "hello",
			wantErr: false,
		},
		{
			name:    "数字",
			input:   123,
			wantErr: false,
		},
		{
			name:    "布尔值",
			input:   true,
			wantErr: false,
		},
		{
			name: "切片",
			input: []int{1, 2, 3},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "正常对象",
			input: map[string]string{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:    "嵌套对象",
			input:   map[string]interface{}{"user": map[string]string{"name": "test"}},
			wantErr: false,
		},
		{
			name:    "nil对象",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalIndent(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		data    []byte
		target  interface{}
		wantErr bool
	}{
		{
			name: "正常JSON",
			data: []byte(`{"name":"test","value":123}`),
			target: &TestStruct{},
			wantErr: false,
		},
		{
			name:    "空数据",
			data:    []byte{},
			target:  &TestStruct{},
			wantErr: true,
		},
		{
			name:    "nil数据",
			data:    nil,
			target:  &TestStruct{},
			wantErr: true,
		},
		{
			name:    "nil目标对象",
			data:    []byte(`{"name":"test"}`),
			target:  nil,
			wantErr: true,
		},
		{
			name:    "无效JSON",
			data:    []byte(`{invalid}`),
			target:  &TestStruct{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.data, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMarshalToString(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "正常对象",
			input:   map[string]string{"key": "value"},
			wantErr: false,
		},
		{
			name:    "nil对象",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalToString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestUnmarshalFromString(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		str     string
		target  interface{}
		wantErr bool
	}{
		{
			name:    "正常JSON字符串",
			str:     `{"name":"test"}`,
			target:  &TestStruct{},
			wantErr: false,
		},
		{
			name:    "空字符串",
			str:     "",
			target:  &TestStruct{},
			wantErr: true,
		},
		{
			name:    "无效JSON字符串",
			str:     `{invalid}`,
			target:  &TestStruct{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalFromString(tt.str, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	original := map[string]interface{}{
		"name":  "test",
		"value": float64(123), // JSON数字会被解析为float64
		"flag":  true,
	}

	// 序列化
	data, err := Marshal(original)
	assert.NoError(t, err)

	// 反序列化
	var result map[string]interface{}
	err = Unmarshal(data, &result)
	assert.NoError(t, err)

	// 验证
	assert.Equal(t, original["name"], result["name"])
	assert.Equal(t, original["value"], result["value"])
	assert.Equal(t, original["flag"], result["flag"])
}

func TestMarshalIndent_Format(t *testing.T) {
	input := map[string]string{"key": "value"}

	// 普通序列化
	compact, _ := Marshal(input)

	// 格式化序列化
	indented, _ := MarshalIndent(input)

	// 格式化的版本应该更长
	assert.Greater(t, len(indented), len(compact))
}
