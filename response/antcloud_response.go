package response

import "../../../antcloud-golang-sdk/sdk/utils"

type AntCloudResponse interface {
	GetResultMsg() string
	SetResultMsg(resultMsg string)
	GetReqMsgId() string
	SetReqMsgId(reqMsgId string)
	GetResultCode() string
	SetResultCode(resultCode string)
	IsSuccess() bool
}

type BaseAntCloudResponse struct {
	ResultMsg  string `json:"result_msg"`
	ReqMsgId   string `json:"req_msg_id"`
	ResultCode string `json:"result_code"`
}

func (antCloudResponse *BaseAntCloudResponse) GetResultMsg() string {
	return antCloudResponse.ResultMsg
}

func (antCloudResponse *BaseAntCloudResponse) SetResultMsg(resultMsg string) {
	antCloudResponse.ResultMsg = resultMsg
}

func (antCloudResponse *BaseAntCloudResponse) GetReqMsgId() string {
	return antCloudResponse.ReqMsgId
}

func (antCloudResponse *BaseAntCloudResponse) SetReqMsgId(reqMsgId string) {
	antCloudResponse.ReqMsgId = reqMsgId
}

func (antCloudResponse *BaseAntCloudResponse) GetResultCode() string {
	return antCloudResponse.ResultCode
}

func (antCloudResponse *BaseAntCloudResponse) SetResultCode(resultCode string) {
	antCloudResponse.ResultCode = resultCode
}

func (antCloudResponse *BaseAntCloudResponse) IsSuccess() bool {
	return antCloudResponse.ResultCode == SuccessResultCode
}

func ParseFromClientResponse(clientResponse *ClientResponse, antCloudResponse AntCloudResponse) (err error) {
	utils.InitJsonParserOnce()
	respStr, err := utils.ExtractRespStr(clientResponse.RawResponseContent)
	if err != nil {
		return
	}
	err = utils.JsonParser.UnmarshalFromString(respStr, antCloudResponse)
	if err != nil {
		return
	}
	return
}
