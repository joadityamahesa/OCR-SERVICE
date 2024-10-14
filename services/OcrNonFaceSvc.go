package services

import (
	"OCR-SERVICE/constanta"
	"OCR-SERVICE/models"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

//function ocr non face (KTP)

func OcrNonFaceSvc(bodyReq models.BodyReq) models.ServiceResponse {

	res := models.ServiceResponse{}

	id := uuid.New()
	//fmt.Println(id.String())
	refId := id.String()

	// log
	/* today := time.Now().Format("2006-01-02")
	logFile, err := os.OpenFile(today+".log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error processing Log file",
			Data:            nil,
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}

		return res
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.LstdFlags) */

	//currentTime := time.Now()

	//timelog := currentTime.Format("2006-01-02 15:04:05.000")

	//reqLog, _ := json.Marshal(bodyReq)
	//log.Println("START PROSES OCR KTP")
	img := bodyReq.Image
	imgLog := img[:10]

	//log.Println("START PROSES FACE COMPARE")
	log.Println("START PROSES OCR KTP ->", "refId : ", id, " image : ", imgLog, ", fileName :", bodyReq.FileName)
	//log.Println("START PROSES OCR KTP:", string(reqLog))

	// Decode base64 images
	image1Data, err := base64.StdEncoding.DecodeString(bodyReq.Image)
	if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error Decode Image",
			Data:            nil, // models.DataAAI{},
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))

		return res
	}

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the first image
	image1Part, err := writer.CreateFormFile("ocrImage", bodyReq.FileName)
	if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error Processing Image",
			Data:            nil, //models.DataAAI{},
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))
		return res
	}
	io.Copy(image1Part, bytes.NewReader(image1Data))

	// Close the writer
	writer.Close()

	// Make a POST request to another service with the multipart form data
	// Replace "YOUR_ACCESS_TOKEN" with your actual access token
	/* 	url := bodyReq.EpService
	   	req, err := http.NewRequest("POST", url, body)
	   	if err != nil {
	   		res = models.ServiceResponse{
	   			Code:            "REQUEST_ERROR",
	   			Message:         "Error request to external service",
	   			Data:            nil,
	   			Extra:           nil,
	   			TransactionID:   "",
	   			PricingStrategy: "",
	   		}

	   		return res
	   	}
	   	req.Header.Set("X-ACCESS-TOKEN", bodyReq.TokenService)
	   	req.Header.Set("Content-Type", writer.FormDataContentType())

	   	client := &http.Client{}
	   	res, err := client.Do(req)
	   	if err != nil {

	   		// Return the default error response
	   		returnDefaultErrorResponse(c, "Error making request to external service", "REQUEST_FAILED")
	   		return
	   	}
	   	defer res.Body.Close() */

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	respaai := models.ResponseAai{}

	restyClient := resty.New()
	_, errSend := restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", writer.FormDataContentType()).
		SetHeader("X-ACCESS-TOKEN", bodyReq.TokenService).
		SetResult(&respaai).
		SetError(&res).
		SetBody(body.Bytes()).
		Post(bodyReq.EpService)
	//fmt.Println("val resp", restires, restires.StatusCode())
	//log.Println("req to endpoint:", string(reqLog))
	// Get the current time
	currentTime := time.Now()

	// Format the time as per the desired layout
	responseTimeService := currentTime.Format("2006-01-02 15:04:05.000")

	// Format the time as per the desired layout
	respTime := responseTimeService

	if ctx.Err() == context.DeadlineExceeded {
		res = models.ServiceResponse{
			Code:              "REQUEST_TIME_OUT",
			Message:           "koneksi timeout ke OCR service. Silahkan coba beberapa saat lagi",
			Data:              nil, //models.DataAAI{},
			Extra:             nil,
			TransactionID:     "",
			PricingStrategy:   "",
			RefId:             refId,
			ResponseTimestamp: respTime,
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))
		return res
	}

	if errSend != nil {
		res = models.ServiceResponse{
			Code:              "REQUEST_ERROR",
			Message:           "Terkendala Jaringan/Koneksi, Silahkan Coba Beberapa Saat Lagi",
			Data:              nil, //models.DataAAI{},
			Extra:             nil,
			TransactionID:     "",
			PricingStrategy:   "",
			RefId:             refId,
			ResponseTimestamp: respTime,
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))
		return res
	}

	fmt.Println(respaai.Code)

	if respaai.Code != "SUCCESS" {
		res = models.ServiceResponse{
			Code:              "REQUEST_ERROR",
			Message:           "Terkendala Jaringan/Koneksi, Silahkan Coba Beberapa Saat Lagi",
			Data:              nil, //models.DataAAI{},
			Extra:             nil,
			TransactionID:     "",
			PricingStrategy:   "",
			RefId:             refId,
			ResponseTimestamp: respTime,
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))
		return res
	}

	respgoogle := models.ResponseGoogleGeo{}

	reqtogoogle := map[string]string{
		"address": respaai.Data.Address + " KELURAHAN " + respaai.Data.Village + " KECAMATAN " + respaai.Data.District, // The address query
		"key":     constanta.API_KEY_GOOGLE_GEOCODING,                                             // Response format
	}

	restyClient = resty.New()
	_, errs := restyClient.R().
		SetQueryParams(reqtogoogle).
		SetResult(&respgoogle). // This will parse the result into the struct
		Get("https://maps.googleapis.com/maps/api/geocode/json")

	log.Println("Req To Google :", reqtogoogle)

	if errs != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Resp Error From Google",
			Data:            nil, //models.DataAAI{},
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))

		return res
	}

	if respgoogle.Status != "OK" {
		res = models.ServiceResponse{
			Code:            respgoogle.Status,
			Message:         "Resp Error From Google",
			Data:            nil, //models.DataAAI{},
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))

		return res
	}

	res = models.ServiceResponse{
		Code:    respaai.Code,
		Message: respaai.Message,
		Data: &models.DataAAI{
			Address:            respaai.Data.Address,
			BirthPlaceBirthday: respaai.Data.BirthPlaceBirthday,
			BloodType:          respaai.Data.BloodType,
			City:               respaai.Data.City,
			District:           respaai.Data.District,
			ExpiryDate:         respaai.Data.ExpiryDate,
			Gender:             respaai.Data.Gender,
			IDNumber:           respaai.Data.IDNumber,
			MaritalStatus:      respaai.Data.MaritalStatus,
			Name:               respaai.Data.Name,
			Nationality:        respaai.Data.Nationality,
			Occupation:         respaai.Data.Occupation,
			Province:           respaai.Data.Province,
			Religion:           respaai.Data.Religion,
			Rtrw:               respaai.Data.Rtrw,
			Village:            respaai.Data.Village,
			Lat:                fmt.Sprintf("%f", respgoogle.Results[0].Geometry.Location.Lat),
			Lon:                fmt.Sprintf("%f", respgoogle.Results[0].Geometry.Location.Lng),
		},
		Extra:             nil,
		TransactionID:     respaai.TransactionID,
		PricingStrategy:   respaai.PricingStrategy,
		RefId:             refId,
		ResponseTimestamp: respTime,
	}

	res.ResponseTimestamp = respTime
	res.RefId = refId
	// res.Data.BirthPlaceBirthday = "Bandar Lampung Nusantara 21-06-1975"
	resLog, _ := json.Marshal(res)

	log.Println("RESPONSE BODY OCR KTP ->", "refId :", refId, " ,", string(resLog))

	// Write the response from the external service back to the client
	//c.JSON(http.StatusOK, res)

	return res

}

// Function to create base64 string from JSON body
/* func CreateBase64String(bodyReq models.BodyReq) (string, error) {
	// Marshal JSON to bytes
	jsonData, err := json.Marshal(bodyReq)
	if err != nil {
		return "", err
	}

	// Encode bytes to base64 string
	base64String := base64.StdEncoding.EncodeToString(jsonData)

	return base64String, nil
} */
