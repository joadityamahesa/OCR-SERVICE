package main

import (
	"OCR-SERVICE/controllers"
	"OCR-SERVICE/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	g := gin.Default()

	g.Use(utils.RequestLogger())
	g.POST("/ocr/nonFace", controllers.OcrNonFace)
	g.POST("/ocr/face", controllers.OcrFace) // latest

	g.POST("/v2/ocr/nonFace", controllers.OcrNonFaceV2) // latest
	g.POST("/v2/move/foto", controllers.MoveFoto)

	g.Run(":8181")

}
