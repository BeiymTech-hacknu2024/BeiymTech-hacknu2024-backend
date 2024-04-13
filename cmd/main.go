package main

import (
	"log"
  "io"
  "os"

  "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/postgres"
)

func main() {
  lg, err := NewLogger()
  if err != nil {
    log.Fatalf("error in creating logger: %v", err)
  }
  DB, err := postgres.ConnectDB()
  if err != nil {
    log.Fatalf("Unable to connect to database: %v\n", err)
  }

}

func init() {
	logrus.SetFormatter(logrus.StandardLogger().Formatter)
	logrus.SetReportCaller(true)
}

func NewLogger() (*logrus.Logger, error) {
	f, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout, f),
		Level: logrus.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
	logger.SetReportCaller(true)
	return logger, nil
}
