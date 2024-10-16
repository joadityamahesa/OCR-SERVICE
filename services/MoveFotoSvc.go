package services

import (
	"OCR-SERVICE/models"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// submit leads
// wso2 hit api move
// move dr tempLeads ke survey
// return fileid & filepath ke wso2
// wso2 akan save data file ke trx upload file

func MoveFotoSvc(c *gin.Context, req models.MoveFotoReq) (res models.Respons) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	log.Println("START PROSES Move FOTO ->", "idNo : ", req.IdNo, " filepath : ", req.FilePath, ", renameFile :", req.RenameFile)

	// koneksi sftp
	configsrc := &ssh.ClientConfig{
		User:            os.Getenv("USER_SFTP"),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nil,
		Auth:            []ssh.AuthMethod{ssh.Password(os.Getenv("PWD_SFTP"))},
		Timeout:         5 * time.Second,
	}

	clientsrc, err := ssh.Dial("tcp", os.Getenv("HOST_SFTP"), configsrc)
	if err != nil {

		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Internal Server Error",
			ResponseTimestamp: "",
			Errors:            "Failed Dial SFTP Server" + err.Error(),
			Data:              nil,
		}

		c.JSON(http.StatusInternalServerError, res)
		return res
	}
	defer clientsrc.Close()

	sftpsrc, errClient := sftp.NewClient(clientsrc)
	if errClient != nil {

		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Internal Server Error",
			ResponseTimestamp: timeStr,
			Errors:            "Error Connection",
			Data:              nil,
		}

		c.JSON(http.StatusInternalServerError, res)
		return res
	}
	defer sftpsrc.Close()

	fmt.Println("konek ke sftp")
	// destFile := "/home/esta_dev/" + "HO1/"
	sourceFile := os.Getenv("SourcePath") + req.IdNo + ".jpeg"

	srcFileNasabah, errpsg := sftpsrc.Open(sourceFile)

	if errpsg != nil {
		res = models.Respons{
			ResponseCode:      "404",
			ResponseMessage:   "Not Found",
			ResponseTimestamp: timeStr,
			Errors:            "File Not Found",
			Data:              nil,
		}
		c.JSON(http.StatusNotFound, res)
		return res

	}
	defer srcFileNasabah.Close()

	// upload ke sftp
	// filename := "/home/esta_dev/HO2/" + req.Namafile
	aa := os.Getenv("DestPath") + req.FilePath
	sftpsrc.MkdirAll(aa)
	filename := aa + req.RenameFile + ".jpeg"
	remoteFile, err := sftpsrc.Create(filename)
	if err != nil {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Error Create File",
			ResponseTimestamp: timeStr,
			Errors:            "Error Create File" + err.Error(),
			Data:              nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return res

	}
	defer remoteFile.Close()

	_, errs := io.Copy(remoteFile, srcFileNasabah)
	if errs != nil {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Error Copy File",
			ResponseTimestamp: timeStr,
			Errors:            "Error Copy File",
			Data:              nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return res
	}

	defer func() {
		srcFileNasabah.Close()
		if err := sftpsrc.Remove(sourceFile); err != nil {
			return
		}

	}()

	// dekrip req filepath
	dataEnkrip := MD5Hash(filename)

	data := models.MoveFotoRes{
		FileId:   dataEnkrip,
		FilePath: filename,
	}

	res = models.Respons{
		ResponseCode:      "200",
		ResponseMessage:   "Success",
		ResponseTimestamp: timeStr,
		Errors:            "",
		Data:              data,
	}

	c.JSON(http.StatusOK, res)

	return res
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}
