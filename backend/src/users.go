package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"profilsactifs/models"
)

func Me(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	var user models.User
	if err := DB.First(&user, userID).Error; err != nil {
		http.Error(res, "user not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(user)
}

func listUsers(res http.ResponseWriter, req *http.Request) {
	var users []models.User
	DB.Find(&users)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(users)
}

func getUser(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(res, "invalid id", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := DB.First(&user, id).Error; err != nil {
		http.Error(res, "user not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(user)
}
