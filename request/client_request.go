package request

import "github.com/huyi-cn/antcloud-golang-sdk/utils"

type ClientRequest struct {
	Params map[string]string
}

func NewClientRequest(params map[string]string) (antCloudRequest *ClientRequest) {
	if len(params) == 0 {
		params = make(map[string]string)
	}
	var request = &ClientRequest{
		Params: params,
	}
	return request
}

func ParseFromAntCloudRequest(antCloudRequest AntCloudRequest) (clientRequest *ClientRequest, err error) {
	params, err := utils.FlattenToRequestParamsMap(antCloudRequest)
	if err != nil {
		return
	}
	clientRequest = &ClientRequest{
		Params: params,
	}
	return
}

func (request *ClientRequest) GetParams() map[string]string {
	return request.Params
}
