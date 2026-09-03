package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"jibjob/src/auth"
	utils "jibjob/src/utils"

	"gorm.io/gorm"
)

func handler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "JibJob say Hello !")
}

func health(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		sqlDB, err := gdb.DB()
		if err == nil {
			err = sqlDB.Ping()
		}
		if err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(res).Encode(map[string]string{"status": "ko"})
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(map[string]string{"status": "ok"})
	}
}

func main() {
	gdb, err := utils.Open()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("connexion db impossible: ", err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", health(gdb))
	http.HandleFunc("/auth/register", auth.Register(gdb))
	http.HandleFunc("/auth/login", auth.Login(gdb))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("serv up in http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
