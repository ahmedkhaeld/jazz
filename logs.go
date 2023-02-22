package jazz

import (
	"io"
	"log"
)

//createLoggers create two loggers instance as the terminal is the destination
func (j *Jazz) createLoggers(out io.Writer) (*log.Logger, *log.Logger) {
	infoLog := log.New(out, "\nINFO:\t", log.Ldate|log.Ltime)
	errLog := log.New(out, "\nERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errLog
}
