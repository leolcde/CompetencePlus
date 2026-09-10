package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"profilsactifs/db"
)

func root(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "ProfilsActifs")
}

func healthCheck(res http.ResponseWriter, req *http.Request) {
	status := "ok"
	code := http.StatusOK

	sqlDB, err := DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		status = "db unreachable"
		code = http.StatusServiceUnavailable
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	json.NewEncoder(res).Encode(map[string]string{"status": status})
}

func main() {
	var err error
	DB, err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(DB); err != nil {
		log.Fatal(err)
	}

	if err := db.SeedQuestions(DB); err != nil {
		log.Fatal(err)
	}

	if err := db.SeedUsers(DB); err != nil {
		log.Fatal(err)
	}

	if err := db.SeedAdmins(DB); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("uploads", 0o755); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/", root)
	router.HandleFunc("GET /health", healthCheck)

	router.HandleFunc("GET /videos", requireAuth(getMyVideo))
	router.HandleFunc("POST /videos", requireRole("candidate", uploadVideo))
	router.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	router.HandleFunc("GET /questions", listQuestions)
	router.HandleFunc("POST /questionnaire/start", requireRole("candidate", startQuestionnaire))
	router.HandleFunc("POST /questionnaire/answer", requireRole("candidate", submitAnswer))
	router.HandleFunc("POST /questionnaire/validate", requireRole("candidate", validateQuestionnaire))
	router.HandleFunc("GET /badge", requireAuth(getBadge))

	router.HandleFunc("POST /register", register)
	router.HandleFunc("POST /login", login)

	router.HandleFunc("GET /me", requireAuth(Me))
	router.HandleFunc("GET /users", listUsers)
	router.HandleFunc("GET /users/{id}", getUser)
	router.HandleFunc("GET /users/{id}/video", getUserVideo)

	router.HandleFunc("GET /consent", requireAuth(getConsent))
	router.HandleFunc("POST /consent", requireAuth(approuveConsent))
	router.HandleFunc("DELETE /consent", requireAuth(revokeConsent))

	fmt.Println("server up on 8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(router)); err != nil {
		log.Fatal(err)
	}
}
