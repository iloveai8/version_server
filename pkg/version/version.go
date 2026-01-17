// Package version 提供版本号解析和比较功能
//
// 该包实现了：
// - 版本号解析（X.Y.Z或X.Y.Z.N格式）
// - 版本号比较（大于、小于、等于）
// - 版本号标准化
// - 版本号验证
//
// 使用示例：
//   v1, _ := version.Parse("1.2.3")
//   v2, _ := version.Parse("1.2.4")
//   if v1.LessThan(v2) {
//       fmt.Println("v1 < v2")
//   }
package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 版本号正则表达式: X.Y.Z 或 X.Y.Z.N
var versionRegex = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:\.(\d+))?$`)

// Version 版本号类型
//
// 表示一个语义化版本号，包含Major、Minor、Patch和可选的Build
type Version struct {
	Major int // 主版本号
	Minor int // 次版本号
	Patch int // 补丁版本号
	Build int // 可选的第4位版本号
}

// Parse 解析版本号字符串
//
// 将版本号字符串解析为Version对象
//
// 参数:
//   versionStr: 版本号字符串（格式：X.Y.Z或X.Y.Z.N）
//
// 返回:
//   *Version: 版本号对象
//   error: 版本号为空或格式无效时返回错误
//
// 示例:
//   v, err := version.Parse("1.2.3")
//   if err != nil {
//       return err
//   }
func Parse(versionStr string) (*Version, error) {
	if versionStr == "" {
		return nil, fmt.Errorf("版本号不能为空")
	}

	matches := versionRegex.FindStringSubmatch(versionStr)
	if matches == nil {
		return nil, fmt.Errorf("版本号格式无效: %s", versionStr)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	build := 0
	if matches[4] != "" {
		build, _ = strconv.Atoi(matches[4])
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		Build: build,
	}, nil
}

// String 返回版本号字符串
//
// 将Version对象转换为字符串格式
//
// 返回:
//   string: 版本号字符串（格式：X.Y.Z或X.Y.Z.N）
//
// 示例:
//   v := version.MustParse("1.2.3.4")
//   fmt.Println(v.String()) // "1.2.3.4"
func (v *Version) String() string {
	if v.Build > 0 {
		return fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Build)
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare 比较两个版本号
//
// 按照Major > Minor > Patch > Build的顺序比较版本号大小
//
// 参数:
//   other: 要比较的另一个版本号
//
// 返回:
//   int: -1表示v<other, 0表示v==other, 1表示v>other
//
// 示例:
//   v1 := version.MustParse("2.0.0")
//   v2 := version.MustParse("1.9.9")
//   result := v1.Compare(v2)  // 返回1，表示v1>v2
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	if v.Build != other.Build {
		if v.Build < other.Build {
			return -1
		}
		return 1
	}

	return 0
}

// GreaterThan 判断当前版本是否大于另一个版本
//
// 参数:
//   other: 要比较的另一个版本号
//
// 返回:
//   bool: 如果v>other返回true，否则返回false
//
// 示例:
//   v1 := version.MustParse("2.0.0")
//   v2 := version.MustParse("1.0.0")
//   if v1.GreaterThan(v2) {
//       fmt.Println("v1 > v2")
//   }
func (v *Version) GreaterThan(other *Version) bool {
	return v.Compare(other) > 0
}

// LessThan 判断当前版本是否小于另一个版本
//
// 参数:
//   other: 要比较的另一个版本号
//
// 返回:
//   bool: 如果v<other返回true，否则返回false
//
// 示例:
//   v1 := version.MustParse("1.0.0")
//   v2 := version.MustParse("2.0.0")
//   if v1.LessThan(v2) {
//       fmt.Println("v1 < v2")
//   }
func (v *Version) LessThan(other *Version) bool {
	return v.Compare(other) < 0
}

// Equal 判断当前版本是否等于另一个版本
//
// 参数:
//   other: 要比较的另一个版本号
//
// 返回:
//   bool: 如果v==other返回true，否则返回false
//
// 示例:
//   v1 := version.MustParse("1.0.0")
//   v2 := version.MustParse("1.0.0")
//   if v1.Equal(v2) {
//       fmt.Println("v1 == v2")
//   }
func (v *Version) Equal(other *Version) bool {
	return v.Compare(other) == 0
}

// GreaterThanOrEqual 判断当前版本是否大于或等于另一个版本
//
// 参数:
//   other: 要比较的另一个版本号
//
// 返回:
//   bool: 如果v>=other返回true，否则返回false
//
// 示例:
//   v1 := version.MustParse("2.0.0")
//   v2 := version.MustParse("1.0.0")
//   if v1.GreaterThanOrEqual(v2) {
//       fmt.Println("v1 >= v2")
//   }
func (v *Version) GreaterThanOrEqual(other *Version) bool {
	return v.Compare(other) >= 0
}

// LessThanOrEqual 判断当前版本是否小于或等于另一个版本
//
// 参数:
//   other: 要比较的另一个版本号
//
// 返回:
//   bool: 如果v<=other返回true，否则返回false
//
// 示例:
//   v1 := version.MustParse("1.0.0")
//   v2 := version.MustParse("2.0.0")
//   if v1.LessThanOrEqual(v2) {
//       fmt.Println("v1 <= v2")
//   }
func (v *Version) LessThanOrEqual(other *Version) bool {
	return v.Compare(other) <= 0
}

// MustParse 解析版本号，失败则panic
//
// 解析版本号字符串，如果失败则panic
// 适用于确定版本号格式正确的场景
//
// 参数:
//   versionStr: 版本号字符串
//
// 返回:
//   *Version: 版本号对象
//
// Panic:
// 当版本号格式无效时panic
//
// 示例:
//   v := version.MustParse("1.2.3")
func MustParse(versionStr string) *Version {
	v, err := Parse(versionStr)
	if err != nil {
		panic(err)
	}
	return v
}

// CompareStrings 比较两个版本号字符串
//
// 直接比较两个版本号字符串，无需先解析
//
// 参数:
//   v1: 版本号字符串1
//   v2: 版本号字符串2
//
// 返回:
//   int: -1表示v1<v2, 0表示v1==v2, 1表示v1>v2
//   error: 版本号格式无效时返回错误
//
// 示例:
//   result, err := version.CompareStrings("1.2.3", "1.2.4")
func CompareStrings(v1, v2 string) (int, error) {
	version1, err := Parse(v1)
	if err != nil {
		return 0, err
	}

	version2, err := Parse(v2)
	if err != nil {
		return 0, err
	}

	return version1.Compare(version2), nil
}

// IsValid 验证版本号字符串是否有效
//
// 检查版本号字符串是否符合格式要求
//
// 参数:
//   versionStr: 版本号字符串
//
// 返回:
//   bool: 格式有效返回true，否则返回false
//
// 示例:
//   if version.IsValid("1.2.3") {
//       fmt.Println("版本号格式正确")
//   }
func IsValid(versionStr string) bool {
	return versionRegex.MatchString(versionStr)
}

// Normalize 标准化版本号字符串（去除v前缀等）
//
// 标准化版本号字符串，去除常见的前缀和后缀
//
// 参数:
//   versionStr: 版本号字符串
//
// 返回:
//   string: 标准化后的版本号字符串
//
// 标准化操作:
// - 去除v或V前缀
// - 转换为小写
//
// 示例:
//   normalized := version.Normalize("v1.2.3")  // 返回"1.2.3"
func Normalize(versionStr string) string {
	// 去除v前缀
	versionStr = strings.TrimPrefix(versionStr, "v")
	// 去除V前缀
	versionStr = strings.TrimPrefix(versionStr, "V")
	// 转换为小写
	versionStr = strings.ToLower(versionStr)
	return versionStr
}
