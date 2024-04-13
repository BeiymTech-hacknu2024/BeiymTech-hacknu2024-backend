package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/controllers"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/handlers"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/postgres"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/routes"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/pkg"
	"github.com/gorilla/sessions"
)

func main() {
	DB, err := postgres.ConnectDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	lg, err := pkg.NewLogger()

	if err != nil {
		log.Fatalf("error in creating logger: %v", err)

	}
	userc := controllers.NewUserController(DB, lg)
	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))
	userh := handlers.NewUserHandler(userc, store, lg)

  
	router := routes.NewRouter(userh)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
