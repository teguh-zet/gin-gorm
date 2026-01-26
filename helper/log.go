package helper

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupLogOutput() {
	dt := time.Now()
	nmFile := dt.Format("01-02-2006 15-04-05")

	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir("./logs", 0755)
	}
	f, _ := os.Create("./logs/" + nmFile + ".log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}
