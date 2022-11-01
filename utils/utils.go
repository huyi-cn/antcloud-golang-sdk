package utils

import (
	"fmt"
	"reflect"
	"strconv"

	errors "../../../antcloud-golang-sdk/sdk/errors"

	"encoding/hex"
	"net/url"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var LoadLocationFromTZData func(name string, data []byte) (*time.Location, error) = nil
var TZData []byte = nil

func GetUUIDV4() (uuidHex string) {
	uuidV4, _ := uuid.NewV4()
	uuidHex = hex.EncodeToString(uuidV4.Bytes())
	return
}

func GetGMTLocation() (*time.Location, error) {
	if LoadLocationFromTZData != nil && TZData != nil {
		return LoadLocationFromTZData("GMT", TZData)
	} else {
		return time.LoadLocation("GMT")
	}
}

func GetTimeInFormatISO8601() (timeStr string) {
	gmt, err := GetGMTLocation()

	if err != nil {
		panic(err)
	}
	return time.Now().In(gmt).Format("2006-01-02T15:04:05Z")
}

func GetUrlFormedMap(source map[string]string) (urlEncoded string) {
	urlEncoder := url.Values{}
	for key, value := range source {
		if key != "" && value != "" {
			urlEncoder.Add(key, value)
		}
	}
	urlEncoded = urlEncoder.Encode()
	return
}

// 提取响应结果字符串
func ExtractRespStr(respBody []byte) (respStr string, err error) {
	// init json parser
	InitJsonParserOnce()

	// 先看看直接marshal的结果是否是原始resp body的子串，是的话就直接返回marshal的结果
	responseMap := make(map[string]interface{})
	err = JsonParser.Unmarshal(respBody, &responseMap)
	if err != nil {
		return "", errors.NewClientError(errors.SdkInvalidResponseFormatCode, "Response format is invalid", err)
	}
	respStr, err = JsonParser.MarshalToString(responseMap["response"])
	return
}

// 从响应结果中提取待签名部分
func ExtractRespStrToSign(respBody []byte) (stringToSign string, err error) {
	// init json parser
	InitJsonParserOnce()

	// 先看看直接marshal的结果是否是原始resp body的子串，是的话就直接返回marshal的结果
	respBodyStr := string(respBody)
	responseMap := make(map[string]interface{})
	err = JsonParser.Unmarshal(respBody, &responseMap)
	if err != nil {
		return "", errors.NewClientError(errors.SdkInvalidResponseFormatCode, "Response format is invalid", err)
	}
	respContentJsonStr, err := JsonParser.MarshalToString(responseMap["response"])
	if strings.Contains(respBodyStr, respContentJsonStr) {
		return respContentJsonStr, err
	}

	// 如果不是的话（比如因为resp body中带有一些不必要的空格换行等字符），我们需要自己手动解析出来
	// 首先判断response和sign在json string中哪个排在前面
	respNodeKey := `"response"`
	signNodeKey := `"sign"`
	respFirstOccurIdx := strings.Index(respBodyStr, respNodeKey)
	signFirstOccurIdx := strings.Index(respBodyStr, signNodeKey)
	if respFirstOccurIdx < 0 || signFirstOccurIdx < 0 {
		return "", errors.NewClientError(errors.SdkInvalidResponseFormatCode, `Response body should contains "response" or "sign" key`, nil)
	}

	respNodeStartIdx := respFirstOccurIdx
	extractStartIdx := strings.Index(respBodyStr[respNodeStartIdx:], "{") + respNodeStartIdx
	var extractEndIdx int
	if respFirstOccurIdx < signFirstOccurIdx {
		// response 出现在 sign 前面
		signLastOccurIdx := strings.LastIndex(respBodyStr, signNodeKey)
		extractEndIdx = strings.LastIndex(respBodyStr[0:signLastOccurIdx], "}")
	} else {
		// response 出现在 sign 后面
		extractEndIdx = strings.LastIndex(respBodyStr[0:len(respBodyStr)-1], "}")
	}
	return respBodyStr[extractStartIdx : extractEndIdx+1], err
}

// 将请求对象打平成请求参数
func FlattenToRequestParamsMap(requestObj interface{}) (params map[string]string, err error) {
	// init json
	InitJsonParserOnce()

	// parse to generic map
	requestJsonStr, err := JsonParser.MarshalToString(requestObj)
	if err != nil {
		return
	}
	var genericMap map[string]interface{}
	err = JsonParser.UnmarshalFromString(requestJsonStr, &genericMap)
	if err != nil {
		return
	}
	params = make(map[string]string)
	flatten(&params, "", genericMap)
	return
}

func flatten(result *map[string]string, path string, obj interface{}) {
	if IsInstanceOf(obj, (map[string]interface{})(nil)) {
		for key, value := range obj.(map[string]interface{}) {
			flatten(result, path+"."+key, value)
		}
	} else if IsInstanceOf(obj, ([]interface{})(nil)) {
		for index, element := range obj.([]interface{}) {
			flatten(result, path+"."+strconv.Itoa(index+1), element)
		}
	} else {
		if obj != nil {
			(*result)[path[1:]] = fmt.Sprintf("%v", obj)
		}
	}
}

func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}
