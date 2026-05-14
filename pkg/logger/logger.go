package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Debug *log.Logger
	Error *log.Logger
)

func Init(isDebugMode bool) {
	var output io.Writer = os.Stdout
	if !isDebugMode {
		output = io.Discard
	}

	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	Debug = log.New(output, "[DEBUG] ", log.LstdFlags)
	Error = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
}
