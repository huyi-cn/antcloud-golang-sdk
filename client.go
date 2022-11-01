package sdk

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	signer "../../antcloud-golang-sdk/sdk/auth"
	errors "../../antcloud-golang-sdk/sdk/errors"
	"../../antcloud-golang-sdk/sdk/request"
	"../../antcloud-golang-sdk/sdk/response"
	"../../antcloud-golang-sdk/sdk/utils"
)

// SDK Client
type Client struct {
	endpoint     string
	accessKey    string
	accessSecret string
	config       *Config
	httpClient   *http.Client
	debug        bool
}

// 构造Client
func NewClient(endpoint string, accessKey string, accessSecret string, config *Config) (client *Client) {
	client = &Client{}
	client.endpoint = endpoint
	client.accessKey = accessKey
	client.accessSecret = accessSecret
	client.httpClient = &http.Client{}
	if config == nil {
		config = NewConfig()
	}
	client.config = config
	// 配置超时时间
	if config.Timeout > 0 {
		client.httpClient.Timeout = config.Timeout
	}
	return client
}

// 通过DTO类调用
func (client *Client) Execute(antCloudRequest request.AntCloudRequest, antCloudResponse response.AntCloudResponse) (err error) {
	// parse antCloudRequest to clientRequest
	clientRequest, err := request.ParseFromAntCloudRequest(antCloudRequest)
	if err != nil {
		return
	}

	// execute client request
	clientResponse, err := client.ExecuteClientRequest(clientRequest)
	if err != nil {
		return
	}

	// parse client response to antCloudResponse
	err = response.ParseFromClientResponse(clientResponse, antCloudResponse)
	if err != nil {
		return
	}
	return
}

// 通过map参数调用
func (client *Client) ExecuteClientRequest(clientRequest *request.ClientRequest) (clientResponse *response.ClientResponse, err error) {
	// 构造参数
	finalParams := client.buildFinalParams(clientRequest)

	// 参数校验
	err = client.validateParams(finalParams)
	if err != nil {
		return
	}

	// 构造HTTP请求
	httpRequest, err := client.buildHttpRequest(finalParams)
	if err != nil {
		return
	}

	if client.config.Debug {
		// Ignore the error
		req, _ := httputil.DumpRequest(httpRequest, true)
		fmt.Printf("Execute HTTP request: %s\n", string(req))
	}

	// 发送HTTP请求
	var httpResponse *http.Response
	for retryTimes := 0; retryTimes <= client.config.MaxRetryTimes; retryTimes++ {
		httpResponse, err = client.httpClient.Do(httpRequest)
		if err != nil {
			if !client.config.AutoRetry {
				return
			} else if retryTimes >= client.config.MaxRetryTimes {
				// retry failed
				transportErrorMsg := fmt.Sprintf("Retry %s times failed", strconv.Itoa(retryTimes+1))
				err = errors.NewClientError(errors.SdkTransportErrorCode, transportErrorMsg, err)
				return
			}
		}
		// retry if status code >= 500 or timeout, will trigger retry
		if client.config.AutoRetry && (err != nil || httpResponse.StatusCode >= http.StatusInternalServerError) {
			// rebuild params
			finalParams = client.buildFinalParams(clientRequest)
			httpRequest, err = client.buildHttpRequest(finalParams)
			if err != nil {
				return
			}
			continue
		}
		break
	}

	if client.config.Debug {
		// Ignore the error
		rsp, _ := httputil.DumpResponse(httpResponse, true)
		fmt.Printf("Receive HTTP response: %s\n", string(rsp))
	}

	// 构造ClientResponse
	clientResponse, err = client.buildClientResponse(httpResponse)
	if err != nil {
		return
	}

	// 校验返回签名
	err = client.checkResponseSign(clientResponse)
	if err != nil {
		return
	}
	return
}

// 构造HTTP请求最后的参数
func (client *Client) buildFinalParams(request *request.ClientRequest) map[string]string {
	finalParams := make(map[string]string)
	// Copy from the original params to the final params
	for key, value := range request.Params {
		finalParams[key] = value
	}
	finalParams["access_key"] = client.accessKey
	finalParams["sign_type"] = "HmacSHA1"
	finalParams["req_msg_id"] = utils.GetUUIDV4()
	finalParams["req_time"] = utils.GetTimeInFormatISO8601()
	finalParams["sdk_version"] = "GOLANG-SDK-1.0.0"
	delete(finalParams, "sign") // 在做签名之前先删掉参数中key名为sign的字段
	finalParams["sign"] = signer.SignParams(finalParams, client.accessSecret)
	return finalParams
}

// 校验参数
func (client *Client) validateParams(params map[string]string) (err error) {
	if params["method"] == "" {
		err = errors.NewClientError(errors.SdkMissingParameterCode, "Request method can't be empty", nil)
	}
	if params["version"] == "" {
		err = errors.NewClientError(errors.SdkMissingParameterCode, "Request version can't be empty", nil)
	}
	return
}

// 构造HttpRequest
func (client *Client) buildHttpRequest(params map[string]string) (httpRequest *http.Request, err error) {
	// 固定为POST方法
	requestMethod := "POST"

	// 初始化httpRequest
	var reader io.Reader
	if params != nil && len(params) > 0 {
		formString := utils.GetUrlFormedMap(params)
		reader = strings.NewReader(formString)
	} else {
		reader = strings.NewReader("")
	}
	httpRequest, err = http.NewRequest(requestMethod, client.endpoint, reader)
	if err != nil {
		return
	}
	httpRequest.Header["Content-Type"] = []string{"application/x-www-form-urlencoded;charset=utf-8"}
	return
}

// 构造ClientResponse
func (client *Client) buildClientResponse(httpResponse *http.Response) (clientResponse *response.ClientResponse, err error) {
	// 读取HttpResponse body
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	// JSON反序列化
	utils.InitJsonParserOnce()
	clientResponse = &response.ClientResponse{}
	err = utils.JsonParser.Unmarshal(body, clientResponse)
	clientResponse.RawResponseContent = body
	return
}

// 校验响应签名
func (client *Client) checkResponseSign(clientResponse *response.ClientResponse) (err error) {
	if clientResponse.IsSuccess() {
		if clientResponse.Sign == "" {
			return errors.NewClientError(errors.SdkBadSignatureCode, `Missing "sign" in response`, nil)
		}
		stringToSign, err := utils.ExtractRespStrToSign(clientResponse.RawResponseContent)
		if err != nil {
			return err
		}
		if signer.Sign(stringToSign, client.accessSecret) != clientResponse.Sign {
			return errors.NewClientError(errors.SdkBadSignatureCode, "Signature is invalid", nil)
		}
	}
	return
}
