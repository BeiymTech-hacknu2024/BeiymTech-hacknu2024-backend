package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/controllers"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/handlers"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/postgres"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/routes"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
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

	userc := controllers.NewUserController(DB, lg)
	assignc := controllers.NewAssignmentController(DB, lg, userc)
	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))
	userh := handlers.NewUserHandler(userc, store, lg)

	assignh := handlers.NewAssignmentHandler(assignc, store, lg)

	router := routes.NewRouter(userh, assignh)

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))

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
