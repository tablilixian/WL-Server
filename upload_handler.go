package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 32MB max memory
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve file from posted form-data
	// Debugging: log request headers and form fields
	log.Printf("Content-Type header: %s", r.Header.Get("Content-Type"))
	log.Printf("Available form keys: %v", r.PostForm)

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("FormFile error: %v", err)
		http.Error(w, "Missing file in form-data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 获取并验证gameId
	gameId := r.FormValue("gameId")
	if gameId == "" {
		http.Error(w, "Missing gameId parameter", http.StatusBadRequest)
		return
	}

	// 验证文件扩展名
	if filepath.Ext(header.Filename) != ".zip" {
		http.Error(w, "Only ZIP files are allowed", http.StatusBadRequest)
		return
	}

	// 创建游戏ID子目录
	subDir := filepath.Join(uploadDir, gameId)
	if err := os.MkdirAll(subDir, 0755); err != nil {
		http.Error(w, "Failed to create game directory", http.StatusInternalServerError)
		return
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// 创建目标文件路径
	dstPath := filepath.Join(subDir, header.Filename)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Failed to create destination file", http.StatusInternalServerError)
		return
	}
	defer dstFile.Close()

	// Copy file content
	if _, err := io.Copy(dstFile, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully: " + header.Filename))
}
