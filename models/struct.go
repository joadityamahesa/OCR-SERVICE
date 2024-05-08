package models

type BodyReq struct {
	Image1       string `json:"img1Base64"`
	Image2       string `json:"img2Base64"`
	FileName1    string `json:"fileName1"`
	FileName2    string `json:"fileName2"`
	EpService    string `json:"epService"`
	TokenService string `json:"tokenService"`
	Image        string `json:"imgBase64"`
	FileName     string `json:"fileName"`
}

type ServiceResponse struct {
	Code              string      `json:"code"`
	Message           string      `json:"message"`
	Data              interface{} `json:"data"`
	Extra             interface{} `json:"extra"`
	TransactionID     string      `json:"transactionId"`
	PricingStrategy   string      `json:"pricingStrategy"`
	RefId             string      `json:"refId"`
	ResponseTimestamp string      `json:"responseTimestamp"`
}

type CustomError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
