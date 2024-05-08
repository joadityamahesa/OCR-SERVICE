package main

import (
	"OCR-SERVICE/controllers"
	"OCR-SERVICE/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	g.Use(utils.RequestLogger())
	g.POST("/ocr/nonFace", controllers.OcrNonFace)
	g.POST("/ocr/face", controllers.OcrFace)
	g.Run(":8181")

}
