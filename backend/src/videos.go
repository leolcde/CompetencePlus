package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"profilsactifs/models"
	"strings"
	"time"
)

const maxUploadSize = 100 << 20 // 100 MB

func uploadFile(res http.ResponseWriter, req *http.Request, userID uint) {
	req.Body = http.MaxBytesReader(res, req.Body, maxUploadSize)
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(res, "file too big (max 100 MB)", http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(res, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	extension := strings.ToLower(filepath.Ext(header.Filename))
	mimeType := header.Header.Get("Content-Type")

	knownExt := extension == ".mp4" || extension == ".webm" || extension == ".mov" || extension == ".m4v" || extension == ".ogg"
	if !knownExt && !strings.HasPrefix(mimeType, "video/") {
		http.Error(res, "unsupported format (video files only)", http.StatusBadRequest)
		return
	}
	if extension == "" {
		switch mimeType {
		case "video/webm":
			extension = ".webm"
		case "video/quicktime":
			extension = ".mov"
		default:
			extension = ".mp4"
		}
	}

	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	dest, err := os.Create(filepath.Join("uploads", name))
	if err != nil {
		http.Error(res, "server error", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		http.Error(res, "write error", http.StatusInternalServerError)
		return
	}

	video := models.Video{URL: "/uploads/" + name, UserID: userID}
	if err := DB.Create(&video).Error; err != nil {
		http.Error(res, "db error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(video)
}

func uploadLink(res http.ResponseWriter, req *http.Request, userID uint) {
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.URL == "" {
		http.Error(res, "missing url", http.StatusBadRequest)
		return
	}

	video := models.Video{URL: body.URL, UserID: userID}
	if err := DB.Create(&video).Error; err != nil {
		http.Error(res, "db error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(video)
}

func getMyVideo(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	var video models.Video
	DB.Where("user_id = ?", userID).Order("created_at DESC, id DESC").First(&video)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(video)
}

func getUserVideo(res http.ResponseWriter, req *http.Request) {
	var video models.Video
	DB.Where("user_id = ?", req.PathValue("id")).Order("created_at DESC, id DESC").First(&video)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(video)
}

func uploadVideo(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)
	contentType := req.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		uploadLink(res, req, userID)
	case strings.HasPrefix(contentType, "multipart/form-data"):
		uploadFile(res, req, userID)
	}
}
