package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestGetUUIDV4(t *testing.T) {
	uuid := GetUUIDV4()
	assert.Equal(t, 32, len(uuid))
	assert.NotEqual(t, GetUUIDV4(), GetUUIDV4())
}

func TestGetTimeInFormatISO8601(t *testing.T) {
	s := GetTimeInFormatISO8601()
	assert.Equal(t, 20, len(s))
	// 2006-01-02T15:04:05Z
	re := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`)
	assert.True(t, re.MatchString(s))
}

func TestGetUrlFormedMap(t *testing.T) {
	m := make(map[string]string)
	m["key"] = "value"
	s := GetUrlFormedMap(m)
	assert.Equal(t, "key=value", s)
	m["key2"] = "http://domain/?key=value&key2=value2"
	s2 := GetUrlFormedMap(m)
	assert.Equal(t, "key=value&key2=http%3A%2F%2Fdomain%2F%3Fkey%3Dvalue%26key2%3Dvalue2", s2)
}

func TestGetTimeInFormatISO8601WithTZData(t *testing.T) {
	TZData = []byte(`"GMT"`)
	LoadLocationFromTZData = func(name string, data []byte) (location *time.Location, e error) {
		if strings.Contains(string(data), name) {
			location, _ = time.LoadLocation(name)
		}
		e = fmt.Errorf("here is a error in test")
		return location, e
	}
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	s := GetTimeInFormatISO8601()
	assert.Equal(t, 20, len(s))
	// 2006-01-02T15:04:05Z
	re := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`)
	assert.True(t, re.MatchString(s))
}

func TestExtractRespStrToSign(t *testing.T) {
	// 标准格式
	result1, err := ExtractRespStrToSign([]byte(`{"response":{"a":1,"c":"hello"}, "sign":"abcde"}`))
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,"c":"hello"}`, result1)

	// response中带较多空格
	result2, err := ExtractRespStrToSign([]byte(`{"response":{"a":1,    "c":"hello"}, "sign":"abcde"}`))
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,    "c":"hello"}`, result2)

	// response中带回车
	result3, err := ExtractRespStrToSign([]byte(`{"response":{"a":1,
"c":"hello"}, "sign":"abcde"}`))
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,
"c":"hello"}`, result3)

	//
	// response和sign换个位置
	//

	// 标准格式
	result4, err := ExtractRespStrToSign([]byte(`{"sign":"abcde", "response":{"a":1,"c":"hello"}}`))
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,"c":"hello"}`, result4)

	// response中带较多空格
	result5, err := ExtractRespStrToSign([]byte(`{"sign":"abcde", "response":{"a":1,    "c":"hello"}}`))
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,    "c":"hello"}`, result5)

	// response中带回车
	result6, err := ExtractRespStrToSign([]byte(`{"sign":"abcde", "response":{"a":1,
"c":"hello"}}`))
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,
"c":"hello"}`, result6)
}
