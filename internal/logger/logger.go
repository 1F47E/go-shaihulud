package logger

// import logrus
import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	log := initLogger()
	return &Logger{log}
}

func initLogger() *logrus.Logger {

	log := logrus.New()

	var format logrus.TextFormatter
	format.ForceColors = true
	format.DisableTimestamp = true
	log.Out = os.Stdout
	log.SetFormatter(&format)

	if os.Getenv("DEBUG") == "1" {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}
