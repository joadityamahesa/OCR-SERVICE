package controllers

import (
	"OCR-SERVICE/models"
	"OCR-SERVICE/services"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MoveFoto(c *gin.Context) {
	req := models.MoveFotoReq{} // struct berisi entity req dari FE
	// res := models.Respons{}                  //struct berisi response yang akan dikembalikan ke FE
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	err := json.Unmarshal(dataBodyReq, &req)

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, dataBodyReq); err != nil {
		return
	}

	if err != nil {
		res := models.Respons{
			ResponseCode:      "400",
			ResponseMessage:   "Error, Unmarshall body Request",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Data:              nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	services.MoveFotoSvc(c, req)

}
