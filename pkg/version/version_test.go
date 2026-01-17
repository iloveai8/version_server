package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantMajor int
		wantMinor int
		wantPatch int
		wantBuild int
	}{
		{
			name:      "有效版本号-三位",
			input:     "1.2.3",
			wantErr:   false,
			wantMajor: 1,
			wantMinor: 2,
			wantPatch: 3,
			wantBuild: 0,
		},
		{
			name:      "有效版本号-四位",
			input:     "1.2.3.4",
			wantErr:   false,
			wantMajor: 1,
			wantMinor: 2,
			wantPatch: 3,
			wantBuild: 4,
		},
		{
			name:      "大版本号",
			input:     "10.20.30",
			wantErr:   false,
			wantMajor: 10,
			wantMinor: 20,
			wantPatch: 30,
			wantBuild: 0,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:    "无效格式-两位",
			input:   "1.2",
			wantErr: true,
		},
		{
			name:    "无效格式-字符",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "无效格式-带前缀",
			input:   "v1.2.3",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantMajor, result.Major)
				assert.Equal(t, tt.wantMinor, result.Minor)
				assert.Equal(t, tt.wantPatch, result.Patch)
				assert.Equal(t, tt.wantBuild, result.Build)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name     string
		version  *Version
		expected string
	}{
		{
			name:     "三位版本号",
			version:  &Version{Major: 1, Minor: 2, Patch: 3},
			expected: "1.2.3",
		},
		{
			name:     "四位版本号",
			version:  &Version{Major: 1, Minor: 2, Patch: 3, Build: 4},
			expected: "1.2.3.4",
		},
		{
			name:     "Build为0",
			version:  &Version{Major: 1, Minor: 2, Patch: 3, Build: 0},
			expected: "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Version
		v2       *Version
		expected int
	}{
		{
			name:     "v1大于v2-Major",
			v1:       &Version{Major: 2, Minor: 0, Patch: 0},
			v2:       &Version{Major: 1, Minor: 0, Patch: 0},
			expected: 1,
		},
		{
			name:     "v1小于v2-Major",
			v1:       &Version{Major: 1, Minor: 0, Patch: 0},
			v2:       &Version{Major: 2, Minor: 0, Patch: 0},
			expected: -1,
		},
		{
			name:     "v1大于v2-Minor",
			v1:       &Version{Major: 1, Minor: 2, Patch: 0},
			v2:       &Version{Major: 1, Minor: 1, Patch: 0},
			expected: 1,
		},
		{
			name:     "v1大于v2-Patch",
			v1:       &Version{Major: 1, Minor: 2, Patch: 3},
			v2:       &Version{Major: 1, Minor: 2, Patch: 2},
			expected: 1,
		},
		{
			name:     "v1大于v2-Build",
			v1:       &Version{Major: 1, Minor: 2, Patch: 3, Build: 4},
			v2:       &Version{Major: 1, Minor: 2, Patch: 3, Build: 3},
			expected: 1,
		},
		{
			name:     "相等",
			v1:       &Version{Major: 1, Minor: 2, Patch: 3},
			v2:       &Version{Major: 1, Minor: 2, Patch: 3},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Compare(tt.v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersion_GreaterThan(t *testing.T) {
	v1 := &Version{Major: 2, Minor: 0, Patch: 0}
	v2 := &Version{Major: 1, Minor: 0, Patch: 0}

	assert.True(t, v1.GreaterThan(v2))
	assert.False(t, v2.GreaterThan(v1))
	assert.False(t, v1.GreaterThan(v1))
}

func TestVersion_LessThan(t *testing.T) {
	v1 := &Version{Major: 1, Minor: 0, Patch: 0}
	v2 := &Version{Major: 2, Minor: 0, Patch: 0}

	assert.True(t, v1.LessThan(v2))
	assert.False(t, v2.LessThan(v1))
	assert.False(t, v1.LessThan(v1))
}

func TestVersion_Equal(t *testing.T) {
	v1 := &Version{Major: 1, Minor: 2, Patch: 3}
	v2 := &Version{Major: 1, Minor: 2, Patch: 3}
	v3 := &Version{Major: 1, Minor: 2, Patch: 4}

	assert.True(t, v1.Equal(v2))
	assert.False(t, v1.Equal(v3))
	assert.True(t, v1.Equal(v1))
}

func TestVersion_GreaterThanOrEqual(t *testing.T) {
	v1 := &Version{Major: 2, Minor: 0, Patch: 0}
	v2 := &Version{Major: 1, Minor: 0, Patch: 0}
	v3 := &Version{Major: 2, Minor: 0, Patch: 0}

	assert.True(t, v1.GreaterThanOrEqual(v2))
	assert.True(t, v1.GreaterThanOrEqual(v3))
	assert.False(t, v2.GreaterThanOrEqual(v1))
}

func TestVersion_LessThanOrEqual(t *testing.T) {
	v1 := &Version{Major: 1, Minor: 0, Patch: 0}
	v2 := &Version{Major: 2, Minor: 0, Patch: 0}
	v3 := &Version{Major: 1, Minor: 0, Patch: 0}

	assert.True(t, v1.LessThanOrEqual(v2))
	assert.True(t, v1.LessThanOrEqual(v3))
	assert.False(t, v2.LessThanOrEqual(v1))
}

func TestMustParse(t *testing.T) {
	// 正常解析
	v := MustParse("1.2.3")
	assert.NotNil(t, v)
	assert.Equal(t, 1, v.Major)
	assert.Equal(t, 2, v.Minor)
	assert.Equal(t, 3, v.Patch)

	// panic情况
	assert.Panics(t, func() {
		MustParse("invalid")
	})
}

func TestCompareStrings(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
		wantErr  bool
	}{
		{
			name:     "v1大于v2",
			v1:       "2.0.0",
			v2:       "1.0.0",
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "v1小于v2",
			v1:       "1.0.0",
			v2:       "2.0.0",
			expected: -1,
			wantErr:  false,
		},
		{
			name:     "相等",
			v1:       "1.2.3",
			v2:       "1.2.3",
			expected: 0,
			wantErr:  false,
		},
		{
			name:    "v1无效",
			v1:      "invalid",
			v2:      "1.0.0",
			wantErr: true,
		},
		{
			name:    "v2无效",
			v1:      "1.0.0",
			v2:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareStrings(tt.v1, tt.v2)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{
			name:   "有效版本号",
			input:  "1.2.3",
			expect: true,
		},
		{
			name:   "有效四位版本号",
			input:  "1.2.3.4",
			expect: true,
		},
		{
			name:   "无效版本号",
			input:  "1.2",
			expect: false,
		},
		{
			name:   "空字符串",
			input:  "",
			expect: false,
		},
		{
			name:   "带v前缀",
			input:  "v1.2.3",
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValid(tt.input)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "带v前缀（小写）",
			input:    "v1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "带V前缀（大写）",
			input:    "V1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "无前缀",
			input:    "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "混合大小写",
			input:    "V1.2.3",
			expected: "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
