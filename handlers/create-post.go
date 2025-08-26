package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/draw"

	m "this_project_id_285410/models"
)

func CreatePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract post fields from form values
	authorID := r.FormValue("userId")
	hostID := r.FormValue("hostId")
	content := r.FormValue("content")

	if content == "" {
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

	now := int(time.Now().Unix())

	// Check if image is present
	imagePresent := 0
	file, _, imgErr := r.FormFile("image")
	if imgErr == nil && file != nil {
		imagePresent = 1
	}

	result, err := db.Exec("INSERT INTO posts (authorId, hostId, content, imagePresent, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)", post.AuthorID, post.HostID, post.Content, imagePresent, now, now)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	post.ID = int(id)

	if file != nil {
		defer file.Close()
		imgDir := os.Getenv("ROOT_CONTENT_DIR") + "/post-images"
		imgPath := imgDir + "/" + fmt.Sprintf("%d.jpg", post.ID)
		if err := saveImage(imgPath, file); err != nil {
			http.Error(w, "Image save error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(post)

}

func saveImage(path string, file io.Reader) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create image directory: %w", err)
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return fmt.Errorf("could not read image file: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return fmt.Errorf("file is not a valid image: %w", err)
	}

	// Resize longer side to 500px, retain proportions
	origBounds := img.Bounds()
	ow, oh := origBounds.Dx(), origBounds.Dy()
	var nw, nh int
	if ow > oh {
		nw = 500
		nh = int(float64(oh) * (500.0 / float64(ow)))
	} else {
		nh = 500
		nw = int(float64(ow) * (500.0 / float64(oh)))
	}

	resized := image.NewRGBA(image.Rect(0, 0, nw, nh))
	draw.CatmullRom.Scale(resized, resized.Bounds(), img, origBounds, draw.Over, nil)

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create image file: %w", err)
	}
	defer out.Close()
	err = jpeg.Encode(out, resized, &jpeg.Options{Quality: 90})
	if err != nil {
		return fmt.Errorf("could not encode image: %w", err)
	}
	return nil
}
