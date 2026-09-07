package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"profilsactifs/src/auth"
	"profilsactifs/src/profils"
	"profilsactifs/src/quiz"
	utils "profilsactifs/src/utils"

	"gorm.io/gorm"
)

func handler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "ProfilsActifs say Hello !")
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

	router := http.NewServeMux()
	router.HandleFunc("/", handler)
	router.HandleFunc("/health", health(gdb))
	router.HandleFunc("/auth/register", auth.Register(gdb))
	router.HandleFunc("/auth/login", auth.Login(gdb))
	router.HandleFunc("/auth/me", auth.Me(gdb))
	router.HandleFunc("/profils", profils.List(gdb))
	router.HandleFunc("/profils/{id}", profils.Detail(gdb))
	router.HandleFunc("/quiz", quiz.List(gdb))
	router.HandleFunc("/quiz/start", quiz.Start(gdb))
	router.HandleFunc("/quiz/answer", quiz.SaveAnswer(gdb))
	router.HandleFunc("/quiz/valider", quiz.Validate(gdb))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("serv up in http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
