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
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// cek flagcheckocr 1 or 0
// if 1 (ocr on) hit aai get address,kel,kec
// hit ke google get lat lon
// upload ke sftp folder tempLeads
// return data dari aai + lat lon
// if 0 (ocr off) ga hit aai & hit google
// upload ke sftp folder tempLeads

func OcrNonFaceV2Svc(bodyReq models.BodyReq) models.ServiceResponse {

	res := models.ServiceResponse{}

	id := uuid.New()
	//fmt.Println(id.String())
	refId := id.String()

	img := bodyReq.Image
	imgLog := img[:10]

	//log.Println("START PROSES FACE COMPARE")
	log.Println("START PROSES OCR KTP ->", "refId : ", id, " image : ", imgLog, ", fileName :", bodyReq.FileName, ", flagCheckOcr :", bodyReq.FlagCheckOcr, ", idNo :", bodyReq.IdNo)
	//log.Println("START PROSES OCR KTP:", string(reqLog))

	// koneksi sftp source path
	configsrc := &ssh.ClientConfig{
		User:            os.Getenv("USER_SFTP"),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nil,
		Auth:            []ssh.AuthMethod{ssh.Password(os.Getenv("PWD_SFTP"))},
		Timeout:         15 * time.Second,
	}

	clientsrc, err := ssh.Dial("tcp", os.Getenv("HOST_SFTP"), configsrc)
	if err != nil {

		res = models.ServiceResponse{
			Code:            "500",
			Message:         "Failed Dial SFTP",
			Data:            nil, //models.DataAAI{},
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}

		resLog, _ := json.Marshal(res)
		log.Println("RESPONSE BODY:", string(resLog))

		return res
	}

	sftpsrc, errC := sftp.NewClient(clientsrc)
	if errC != nil {
		res = models.ServiceResponse{
			Code:            "500",
			Message:         "Failed Dial New Client",
			Data:            nil, //models.DataAAI{},
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}
		resLog, _ := json.Marshal(res)
		log.Println("RESPONSE BODY:", string(resLog))
		return res
	}
	defer sftpsrc.Close()
	// end koneksi sftp source
	fmt.Println("Successfully connected to ssh server A.")

	// 1. cek checkOcr
	if bodyReq.FlagCheckOcr == "1" {
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

		writer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()

		respaai := models.ResponseAai{}

		// 2. hit ke aai
		restyClient := resty.New()
		_, errSend := restyClient.R().
			SetContext(ctx).
			SetHeader("Content-Type", writer.FormDataContentType()).
			SetHeader("X-ACCESS-TOKEN", bodyReq.TokenService).
			SetResult(&respaai).
			SetError(&res).
			SetBody(body.Bytes()).
			Post(bodyReq.EpService)

		currentTime := time.Now()

		responseTimeService := currentTime.Format("2006-01-02 15:04:05.000")

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
			"key":     constanta.API_KEY_GOOGLE_GEOCODING,                                                                  // Response format
		}
		// 3. hit ke google
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

		log.Println("Res From Google :", respgoogle.Results[0].Geometry.Location.Lat, respgoogle.Results[0].Geometry.Location.Lng)

		// 4. cek flagBcg jika == "B" -> Upload ke sftp
		if bodyReq.FlagBcg == "B" {
			imgBytes, _ := base64.StdEncoding.DecodeString(bodyReq.Image)
			folderPath := os.Getenv("SourcePath")

			sftpsrc.Mkdir(folderPath)

			remoteFile, errCreate := sftpsrc.Create(folderPath + respaai.Data.IDNumber + ".jpeg")
			if errCreate != nil {
				res = models.ServiceResponse{
					Code:            respgoogle.Status,
					Message:         "Failed Dial New Client",
					Data:            nil, //models.DataAAI{},
					Extra:           nil,
					TransactionID:   "",
					PricingStrategy: "",
				}
				resLog, _ := json.Marshal(res)
				log.Println("RESPONSE BODY:", string(resLog))
				return res
			}

			_, err = remoteFile.Write(imgBytes)
			if err != nil {
				return res
			}
		}

		lat := strconv.FormatFloat(respgoogle.Results[0].Geometry.Location.Lat, 'f', -1, 64)
		lng := strconv.FormatFloat(respgoogle.Results[0].Geometry.Location.Lng, 'f', -1, 64)

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
				Lat:                lat, //fmt.Sprintf("%f", respgoogle.Results[0].Geometry.Location.Lat),
				Lon:                lng, //fmt.Sprintf("%f", respgoogle.Results[0].Geometry.Location.Lng),
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
	} else {
		imgBytes, _ := base64.StdEncoding.DecodeString(bodyReq.Image)
		folderPath := os.Getenv("SourcePath")

		sftpsrc.Mkdir(folderPath)

		remoteFile, errCreate := sftpsrc.Create(folderPath + bodyReq.IdNo + ".jpeg")
		if errCreate != nil {
			res = models.ServiceResponse{
				Code:            "500",
				Message:         "Failed Dial New Client",
				Data:            nil, //models.DataAAI{},
				Extra:           nil,
				TransactionID:   "",
				PricingStrategy: "",
			}
			resLog, _ := json.Marshal(res)
			log.Println("RESPONSE BODY:", string(resLog))
			return res
		}

		_, err = remoteFile.Write(imgBytes)
		if err != nil {
			return res
		}

		res = models.ServiceResponse{
			Code:              "SUCCESS",
			Message:           "OK",
			Data:              nil, //&models.DataAAI{},
			Extra:             nil,
			TransactionID:     "",
			PricingStrategy:   "",
			RefId:             refId,
			ResponseTimestamp: "",
		}

		return res
	}

}
