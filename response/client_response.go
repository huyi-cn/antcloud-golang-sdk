package response

const (
	SuccessResultCode = "OK"
)

type ClientResponse struct {
	RawResponseContent []byte
	Response           BaseAntCloudResponse `json:"response"`
	Sign               string               `json:"sign"`
}

func (clientResponse *ClientResponse) IsSuccess() bool {
	return clientResponse.Response.IsSuccess()
}
