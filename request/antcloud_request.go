package request

type AntCloudRequest interface {
	GetMethod() string
	SetMethod(method string)
	GetVersion() string
	SetVersion(version string)
	GetProductInstanceId() string
	SetProductInstanceId(productInstanceId string)
}

type BaseAntCloudRequest struct {
	Method            string `json:"method"`
	Version           string `json:"version"`
	ProductInstanceId string `json:"product_instance_id"`
}

func (antCloudRequest *BaseAntCloudRequest) GetMethod() string {
	return antCloudRequest.Method
}

func (antCloudRequest *BaseAntCloudRequest) SetMethod(method string) {
	antCloudRequest.Method = method
}

func (antCloudRequest *BaseAntCloudRequest) GetVersion() string {
	return antCloudRequest.Version
}

func (antCloudRequest *BaseAntCloudRequest) SetVersion(version string) {
	antCloudRequest.Version = version
}

func (antCloudRequest *BaseAntCloudRequest) GetProductInstanceId() string {
	return antCloudRequest.ProductInstanceId
}

func (antCloudRequest *BaseAntCloudRequest) SetProductInstanceId(productInstanceId string) {
	antCloudRequest.ProductInstanceId = productInstanceId
}
