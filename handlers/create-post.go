package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"

	m "this_project_id_285410/models"
)

func CreatePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		log.Println("Invalid form data:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract post fields from form values
	authorID := r.FormValue("userId")
	hostID := r.FormValue("hostId")
	content := r.FormValue("content")

	if content == "" {
		log.Println("Missing content")
		http.Error(w, "Missing content", http.StatusBadRequest)
		return
	}

	var post m.Post
	post.Content = content

	if authorID != "" {
		// Assume valid int string from form
		post.AuthorID = 0
		fmt.Sscanf(authorID, "%d", &post.AuthorID)
	}
	if hostID != "" {
		post.HostID = 0
		fmt.Sscanf(hostID, "%d", &post.HostID)
	}

	now := time.Now()

	// Check if image is present
	imagePresent := false
	file, _, imgErr := r.FormFile("image")
	if imgErr == nil && file != nil {
		imagePresent = true
	}

	err = db.QueryRow(
		"INSERT INTO posts (author_id, host_id, content, image_present, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		post.AuthorID, post.HostID, post.Content, imagePresent, now, now,
	).Scan(&post.ID)
	if err != nil {
		log.Println("DB error:", err)
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if file != nil {
		defer file.Close()
		imgDir := os.Getenv("ROOT_CONTENT_DIR") + "/post-images"
		imgPath := imgDir + "/" + fmt.Sprintf("%d.jpg", post.ID)
		if err := saveImage(imgPath, file); err != nil {
			log.Println("Image save error:", err)
			http.Error(w, "Image save error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(post)

}

func saveImage(path string, file io.Reader) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Println("could not create image directory:", err)
		return fmt.Errorf("could not create image directory: %w", err)
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		log.Println("could not read image file:", err)
		return fmt.Errorf("could not read image file: %w", err)
	}

	img, err := imaging.Decode(bytes.NewReader(buf.Bytes()), imaging.AutoOrientation(true))
	if err != nil {
		log.Println("file is not a valid image:", err)
		return fmt.Errorf("file is not a valid image: %w", err)
	}

	resized := imaging.Fit(img, 500, 500, imaging.Lanczos)

	out, err := os.Create(path)
	if err != nil {
		log.Println("could not create image file:", err)
		return fmt.Errorf("could not create image file: %w", err)
	}
	defer out.Close()
	err = imaging.Encode(out, resized, imaging.JPEG)
	if err != nil {
		log.Println("could not encode image:", err)
		return fmt.Errorf("could not encode image: %w", err)
	}
	return nil
}
