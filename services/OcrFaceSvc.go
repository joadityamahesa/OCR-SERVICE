package services

import (
	"OCR-SERVICE/models"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

//function ocr non face (KTP)

func OcrFaceSvc(bodyReq models.BodyReq) models.ServiceResponse {

	res := models.ServiceResponse{}
	id := uuid.New()
	//fmt.Println(id.String())
	refId := id.String()

	// log
	//today := time.Now().Format("2006-01-02")
	//logFile, err := os.OpenFile(today+".log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	/* if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error processing Log file",
			Data:            nil,
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}

		return res
	} */
	//defer logFile.Close()
	//log.SetOutput(logFile)
	//log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.LstdFlags)

	//currentTime := time.Now()

	//timelog := currentTime.Format("2006-01-02 15:04:05.000")

	//reqLog, _ := json.Marshal(bodyReq)
	img1 := bodyReq.Image1
	img1Log := img1[:10]
	img2 := bodyReq.Image2
	img2Log := img2[:10]
	//log.Println("START PROSES FACE COMPARE")
	log.Println("START PROSES FACE COMPARE ->", "refId : ", id, ", image 1 :", img1Log, " ,", "image 2 : ",
		img2Log, ", fileName 1 :", bodyReq.FileName1, ", fileName 2 :", bodyReq.FileName2)

	// Decode base64 images
	image1Data, err := base64.StdEncoding.DecodeString(bodyReq.Image1)
	if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error Decode Image",
			Data:            nil,
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))

		return res
	}

	// Decode base64 images
	image2Data, err := base64.StdEncoding.DecodeString(bodyReq.Image2)
	if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error Decode Image",
			Data:            nil,
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
	image1Part, err := writer.CreateFormFile("firstImage", bodyReq.FileName1)
	if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error Processing Image",
			Data:            nil,
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))
		return res
	}
	io.Copy(image1Part, bytes.NewReader(image1Data))

	// Add the second image
	image2Part, err := writer.CreateFormFile("secondImage", bodyReq.FileName2)
	if err != nil {
		res = models.ServiceResponse{
			Code:            "REQUEST_ERROR",
			Message:         "Error Processing Image",
			Data:            nil,
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)

		log.Println("RESPONSE BODY:", string(resLog))
		return res
	}
	io.Copy(image2Part, bytes.NewReader(image2Data))

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

	//restyClient.SetTimeout(1 * time.Millisecond)
	//var tm int
	//var tm = 1

	//ctx, cancel := context.WithTimeout(context.Background(), tm*time.Millisecond*time.Millisecond)
	//defer cancel()

	//restyClient.SetTimeout(1 * time.Microsecond)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	restyClient := resty.New()
	_, errSend := restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", writer.FormDataContentType()).
		SetHeader("X-ACCESS-TOKEN", bodyReq.TokenService).
		SetResult(&res).
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

	log.Println(respTime)

	if ctx.Err() == context.DeadlineExceeded {
		res = models.ServiceResponse{
			Code:              "REQUEST_TIME_OUT",
			Message:           "koneksi timeout ke OCR service. Silahkan coba beberapa saat lagi",
			Data:              nil,
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
			Data:              nil,
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

	if res.Code == "" {
		res = models.ServiceResponse{
			Code:              "REQUEST_ERROR",
			Message:           "Terkendala Jaringan/Koneksi, Silahkan Coba Beberapa Saat Lagi",
			Data:              nil,
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

	res.ResponseTimestamp = respTime
	res.RefId = refId
	resLog, _ := json.Marshal(res)

	log.Println("RESPONSE BODY FACE COMPARE ->", "refId :", refId, " ,", string(resLog))

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
