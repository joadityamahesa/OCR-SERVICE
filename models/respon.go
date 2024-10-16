package models

type ResponseAai struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
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
	} `json:"data"`
	Extra             interface{} `json:"extra"`
	TransactionID     string      `json:"transactionId"`
	PricingStrategy   string      `json:"pricingStrategy"`
	RefID             string      `json:"refId"`
	ResponseTimestamp string      `json:"responseTimestamp"`
}

type ResponseGoogleGeo struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID  string `json:"place_id"`
		PlusCode struct {
			CompoundCode string `json:"compound_code"`
			GlobalCode   string `json:"global_code"`
		} `json:"plus_code"`
		Types []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

type Respons struct {
	ResponseCode      string      `json:"responseCode"`
	ResponseMessage   string      `json:"responseMessage"`
	ResponseTimestamp string      `json:"responseTimestamp"`
	Errors            string      `json:"errors"`
	Data              interface{} `json:"data"`
}
