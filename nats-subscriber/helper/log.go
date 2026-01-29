package helper

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetOutput(gin.DefaultWriter)
}
