package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type customResponseWriter struct {
	gin.ResponseWriter
	body *string
}

func (crw *customResponseWriter) Write(data []byte) (int, error) {
	if crw.body == nil {
		body := string(data)
		crw.body = &body
	} else {
		*crw.body += string(data)
	}
	return crw.ResponseWriter.Write(data)
}

func (crw *customResponseWriter) Status() int {
	return crw.ResponseWriter.Status()
}

func captureRequestBody(body io.ReadCloser) []byte {
	requestBody, _ := ioutil.ReadAll(body)
	return requestBody
}

func logRequestResponse(c *gin.Context, responseBody *string, requestBody []byte) {

	log.Println("Request Method: ", c.Request.Method)
	log.Println("Request URL: ", c.Request.URL)

	// log.Printf("Request Headers: %v", c.Request.Header)

	// Log request body (if present)
	//log.Printf("Request Body: %s", string(requestBody))

	// Log response details
	log.Printf("Response Status: %d", c.Writer.Status())

	// Log response body (if present)
	if c.Writer.Status() == 200 {
		if responseBody != nil {
			log.Printf("Response Body: %s\n\n", *responseBody)
		}
	} else {
		if responseBody != nil {
			log.Printf("Response Body: %s\n\n", *responseBody)
		}
	}

}

func createLogFile() (*os.File, error) {
	// Create a directory for logs if it doesn't exist
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	// Generate log file name based on the current date
	logFileName := fmt.Sprintf("%s/ocr_service_log_%s.log", logDir, time.Now().Format("2006-01-02"))

	// Open or create the log file
	return os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logFile, err := createLogFile()
		if err != nil {
			log.Printf("Error creating log file: %v", err)
			c.Next()
			return
		}
		defer logFile.Close()

		// Set the log output to the file
		log.SetOutput(logFile)

		startTime := time.Now()

		// Use the customResponseWriter to capture the response body
		crw := &customResponseWriter{c.Writer, nil}
		c.Writer = crw

		// Create a tee reader to duplicate the request body for logging
		requestBody := captureRequestBody(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

		// Log the request details
		log.Printf("[%s] %s %s\n", startTime.Format(time.RFC3339), c.Request.Method, c.Request.URL)

		// Process the request
		c.Next()

		endTime := time.Now()
		duration := endTime.Sub(startTime)

		status := crw.Status()
		if status >= http.StatusOK && status < http.StatusMultipleChoices {
			log.Printf("[%s] %s %s %d (%v)\n", endTime.Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, status, duration)
		} else {
			log.Printf("[%s] %s %s %d (%v)\n", endTime.Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, status, duration)
		}
		// Log response details
		logRequestResponse(c, crw.body, requestBody)
	}
}
