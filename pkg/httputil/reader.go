// Package httputil 提供HTTP工具函数
//
// 该包实现了：
// - 请求体读取（带大小限制）
// - 请求ID中间件
// - HTTP相关工具函数
//
// 使用示例：
//   body, err := httputil.ReadRequestBody(r)
//   if err != nil {
//       return err
//   }
package httputil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	// MaxBodySize 默认最大请求体大小（1MB）
	MaxBodySize = 1 << 20
)

// ReadRequestBody 安全地读取HTTP请求体
//
// 使用默认大小限制（1MB）读取HTTP请求体
//
// 参数:
//   r: HTTP请求对象
//
// 返回:
//   []byte: 请求体内容
//   error: 读取失败或请求体过大时返回错误
//
// 示例:
//   body, err := httputil.ReadRequestBody(r)
func ReadRequestBody(r *http.Request) ([]byte, error) {
	return ReadRequestBodyWithLimit(r, MaxBodySize)
}

// ReadRequestBodyWithLimit 读取HTTP请求体（带大小限制）
//
// 使用指定大小限制读取HTTP请求体，防止内存耗尽攻击
//
// 参数:
//   r: HTTP请求对象
//   maxSize: 最大允许的请求体大小（字节）
//
// 返回:
//   []byte: 请求体内容
//   error: 读取失败或请求体过大时返回错误
//
// 示例:
//   body, err := httputil.ReadRequestBodyWithLimit(r, 2<<20) // 2MB
func ReadRequestBodyWithLimit(r *http.Request, maxSize int64) ([]byte, error) {
	if r == nil {
		return nil, errors.New("请求对象为空")
	}

	if r.Body == nil {
		return nil, errors.New("请求体为空")
	}

	// 使用io.LimitReader限制读取大小
	limitedReader := io.LimitReader(r.Body, maxSize+1) // +1用于检测是否超过限制

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("读取请求体失败: %w", err)
	}

	// 检查是否超过大小限制
	if int64(len(body)) > maxSize {
		return nil, fmt.Errorf("请求体过大，最大允许%d字节", maxSize)
	}

	return body, nil
}
