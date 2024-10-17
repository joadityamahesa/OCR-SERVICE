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
	IdNo         string `json:"idNo"`
	FlagCheckOcr string `json:"flagCheckOcr"`
	FlagBcg      string `json:"flagBcg"`
}

type ServiceResponseFace struct {
	Code              string      `json:"code"`
	Message           string      `json:"message"`
	Data              interface{} `json:"data"`
	Extra             interface{} `json:"extra"`
	TransactionID     string      `json:"transactionId"`
	PricingStrategy   string      `json:"pricingStrategy"`
	RefId             string      `json:"refId"`
	ResponseTimestamp string      `json:"responseTimestamp"`
}

type ServiceResponse struct {
	Code              string      `json:"code"`
	Message           string      `json:"message"`
	Data              *DataAAI    `json:"data"`
	Extra             interface{} `json:"extra"`
	TransactionID     string      `json:"transactionId"`
	PricingStrategy   string      `json:"pricingStrategy"`
	RefId             string      `json:"refId"`
	ResponseTimestamp string      `json:"responseTimestamp"`
}

type DataAAI struct {
	Address            string      `json:"address"`
	BirthPlaceBirthday string      `json:"birthPlaceBirthday"`
	BloodType          interface{} `json:"bloodType"`
	City               string      `json:"city"`
	District           string      `json:"district"`
	ExpiryDate         string      `json:"expiryDate"`
	Gender             string      `json:"gender"`
	IDNumber           string      `json:"idNumber"`
	MaritalStatus      string      `json:"maritalStatus"`
	Name               string      `json:"name"`
	Nationality        string      `json:"nationality"`
	Occupation         string      `json:"occupation"`
	Province           string      `json:"province"`
	Religion           string      `json:"religion"`
	Rtrw               string      `json:"rtrw"`
	Village            string      `json:"village"`
	Lat                string      `json:"latitude"`
	Lon                string      `json:"longitude"`
}

type CustomError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
