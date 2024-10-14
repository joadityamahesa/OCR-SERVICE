package controllers

import (
	"OCR-SERVICE/models"
	"OCR-SERVICE/services"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//function ocr non face (KTP)

func OcrFace(c *gin.Context) {

	req := models.BodyReq{}
	//res := models.ServiceResponse{}
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	err := json.Unmarshal(dataBodyReq, &req)

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, dataBodyReq); err != nil {
		fmt.Println(err)
	}

	if err != nil {
		res := models.ServiceResponseFace{
			Code:            "INVALID_JSON",
			Message:         "Error Unmarshall JSON",
			Data:            nil,
			Extra:           nil,
			TransactionID:   "",
			PricingStrategy: "",
		}

		c.JSON(200, res)

		respLog, _ := json.Marshal(res)
		log.Println("response body:", string(respLog))
		return
	}

	svc := services.OcrFaceSvc(req)

	c.JSON(http.StatusOK, svc)
}
