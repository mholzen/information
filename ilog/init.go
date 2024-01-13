package ilog

import (
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func Init() {
	logFile, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	spew.Config.Indent = "\t"
}
