// Package jsonutil 提供JSON序列化和反序列化工具函数
//
// 该包是对标准库encoding/json的封装，提供了更友好的错误处理
//
// 使用示例：
//   data, err := jsonutil.Marshal(obj)
//   if err != nil {
//       return err
//   }
package jsonutil

import (
	"encoding/json"
	"fmt"
)

// Marshal 序列化对象为JSON（带错误检查）
//
// 将对象序列化为JSON字节数组，包含参数验证
//
// 参数:
//   v: 待序列化的对象
//
// 返回:
//   []byte: JSON字节数组
//   error: 对象为空或序列化失败时返回错误
//
// 示例:
//   data, err := jsonutil.Marshal(obj)
func Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("序列化对象不能为空")
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %w", err)
	}

	return data, nil
}

// MarshalIndent 序列化对象为格式化的JSON（带错误检查）
//
// 将对象序列化为格式化的JSON字节数组（带缩进），包含参数验证
//
// 参数:
//   v: 待序列化的对象
//
// 返回:
//   []byte: 格式化的JSON字节数组
//   error: 对象为空或序列化失败时返回错误
//
// 示例:
//   data, err := jsonutil.MarshalIndent(obj)
func MarshalIndent(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("序列化对象不能为空")
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("JSON格式化序列化失败: %w", err)
	}

	return data, nil
}

// Unmarshal 反序列化JSON为对象（带错误检查）
//
// 将JSON字节数组反序列化为对象，包含参数验证
//
// 参数:
//   data: JSON字节数组
//   v: 目标对象（必须是指针类型）
//
// 返回:
//   error: 数据为空、目标对象为空或反序列化失败时返回错误
//
// 示例:
//   err := jsonutil.Unmarshal(data, &obj)
func Unmarshal(data []byte, v interface{}) error {
	if data == nil || len(data) == 0 {
		return fmt.Errorf("JSON数据不能为空")
	}

	if v == nil {
		return fmt.Errorf("目标对象不能为空")
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("JSON反序列化失败: %w", err)
	}

	return nil
}

// MarshalToString 序列化对象为JSON字符串
//
// 将对象序列化为JSON字符串
//
// 参数:
//   v: 待序列化的对象
//
// 返回:
//   string: JSON字符串
//   error: 序列化失败时返回错误
//
// 示例:
//   str, err := jsonutil.MarshalToString(obj)
func MarshalToString(v interface{}) (string, error) {
	data, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalFromString 从JSON字符串反序列化
//
// 将JSON字符串反序列化为对象
//
// 参数:
//   str: JSON字符串
//   v: 目标对象（必须是指针类型）
//
// 返回:
//   error: 反序列化失败时返回错误
//
// 示例:
//   err := jsonutil.UnmarshalFromString(`{"name":"test"}`, &obj)
func UnmarshalFromString(str string, v interface{}) error {
	return Unmarshal([]byte(str), v)
}
