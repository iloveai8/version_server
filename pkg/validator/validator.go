// Package validator 提供输入验证功能
//
// 该包实现了：
// - 版本号格式验证（X.Y.Z或X.Y.Z.N）
// - URL格式验证
// - 环境参数验证
// - IP地址格式验证
// - 分页参数验证
//
// 使用示例：
//   if err := validator.ValidateVersion("1.2.3"); err != nil {
//       return err
//   }
package validator

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
)

var (
	// 版本号正则表达式: X.Y.Z 或 X.Y.Z.N
	// 例如: 1.0.0, 1.2.3, 1.2.3.4
	versionRegex = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:\.(\d+))?$`)

	// URL正则表达式
	urlRegex = regexp.MustCompile(`^https?://.+`)
)

// ValidateVersion 验证版本号格式
//
// 检查版本号是否符合X.Y.Z或X.Y.Z.N格式，其中X、Y、Z、N为非负整数
//
// 参数:
//   vsn: 版本号字符串，如"1.0.0"或"1.2.3.4"
//
// 返回:
//   error: 版本号格式无效时返回错误，格式正确返回nil
//
// 示例:
//   err := validator.ValidateVersion("1.2.3")
//   if err != nil {
//       return err
//   }
func ValidateVersion(vsn string) error {
	if vsn == "" {
		return fmt.Errorf("版本号不能为空")
	}

	if !versionRegex.MatchString(vsn) {
		return fmt.Errorf("版本号格式无效，应为X.Y.Z或X.Y.Z.N格式")
	}

	return nil
}

// ValidateVersionFormat 验证版本号格式（返回bool）
//
// 检查版本号格式是否有效，返回布尔值
//
// 参数:
//   vsn: 版本号字符串
//
// 返回:
//   bool: 格式有效返回true，否则返回false
//
// 示例:
//   if validator.ValidateVersionFormat("1.2.3") {
//       // 版本号格式正确
//   }
func ValidateVersionFormat(vsn string) bool {
	return versionRegex.MatchString(vsn)
}

// ValidateURL 验证URL格式
//
// 使用标准库验证URL格式，检查协议和主机
//
// 参数:
//   rawURL: URL字符串
//
// 返回:
//   error: URL格式无效时返回错误，格式正确返回nil
//
// 示例:
//   err := validator.ValidateURL("https://example.com")
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL不能为空")
	}

	// 使用标准库验证
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("URL解析失败: %w", err)
	}

	// 检查协议
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL必须包含协议（http://或https://）")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL协议必须是http或https")
	}

	// 检查主机
	if parsedURL.Host == "" {
		return fmt.Errorf("URL必须包含主机地址")
	}

	return nil
}

// ValidateHTTPURL 验证HTTP/HTTPS URL
//
// 快速验证URL是否以http://或https://开头
//
// 参数:
//   rawURL: URL字符串
//
// 返回:
//   error: URL不以http://或https://开头时返回错误
//
// 示例:
//   err := validator.ValidateHTTPURL("https://example.com")
func ValidateHTTPURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL不能为空")
	}

	if !urlRegex.MatchString(rawURL) {
		return fmt.Errorf("URL必须以http://或https://开头")
	}

	return nil
}

// ValidateEnv 验证环境参数
//
// 检查环境参数是否为允许的值（dev或pro）
//
// 参数:
//   env: 环境参数字符串
//
// 返回:
//   error: 环境参数无效时返回错误
//
// 示例:
//   err := validator.ValidateEnv("dev")
func ValidateEnv(env string) error {
	if env == "" {
		return fmt.Errorf("环境参数不能为空")
	}

	validEnvs := map[string]bool{
		"dev": true,
		"pro": true,
	}

	if !validEnvs[env] {
		return fmt.Errorf("环境参数无效，必须是dev或pro")
	}

	return nil
}

// ValidateIP 验证IP地址格式
//
// 使用标准库验证IP地址格式（支持IPv4和IPv6）
//
// 参数:
//   ip: IP地址字符串
//
// 返回:
//   error: IP地址格式无效时返回错误
//
// 示例:
//   err := validator.ValidateIP("192.168.1.1")
func ValidateIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP地址不能为空")
	}

	// 使用net.ParseIP进行严格验证
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("IP地址格式无效")
	}

	return nil
}

// ValidateIPOrCIDR 验证IP地址或CIDR格式
//
// 支持单个IP地址（IPv4/IPv6）和CIDR网段格式
//
// 参数:
//   ip: IP地址或CIDR字符串
//
// 返回:
//   error: 格式无效时返回错误
//
// 示例:
//   err := validator.ValidateIPOrCIDR("192.168.1.1")      // 单个IP
//   err := validator.ValidateIPOrCIDR("10.0.0.0/24")      // CIDR网段
func ValidateIPOrCIDR(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP地址不能为空")
	}

	// 先尝试作为单个IP验证
	if net.ParseIP(ip) != nil {
		return nil
	}

	// 尝试作为CIDR验证
	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		return fmt.Errorf("IP地址或CIDR格式无效")
	}

	return nil
}

// ValidateCIDR 验证CIDR格式
//
// 专门验证CIDR网段格式
//
// 参数:
//   cidr: CIDR字符串
//
// 返回:
//   error: CIDR格式无效时返回错误
//
// 示例:
//   err := validator.ValidateCIDR("10.0.0.0/24")
func ValidateCIDR(cidr string) error {
	if cidr == "" {
		return fmt.Errorf("CIDR不能为空")
	}

	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("CIDR格式无效")
	}

	return nil
}

// ValidatePageAndPageSize 验证分页参数
//
// 检查分页参数是否合法（页码>0，每页数量1-100）
//
// 参数:
//   page: 页码
//   pageSize: 每页数量
//
// 返回:
//   error: 分页参数无效时返回错误
//
// 示例:
//   err := validator.ValidatePageAndPageSize(1, 20)
func ValidatePageAndPageSize(page, pageSize int) error {
	if page < 1 {
		return fmt.Errorf("页码必须大于0")
	}

	if pageSize < 1 || pageSize > 100 {
		return fmt.Errorf("每页数量必须在1-100之间")
	}

	return nil
}
