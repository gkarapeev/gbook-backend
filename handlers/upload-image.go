package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

// UploadImageHandler handles image uploads, resizes and crops to 200x200, and saves as JPEG.
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse multipart form"))
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not get file from form"))
		return
	}
	defer file.Close()

	userId := r.FormValue("userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing userId field in form"))
		return
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read file"))
		return
	}

	img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File is not a valid image: " + err.Error()))
		return
	}

	if format != "jpeg" && format != "png" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only PNG and JPEG images are allowed"))
		return
	}

	// Resize: fit the smaller side to 200, keep proportions
	origBounds := img.Bounds()
	ow, oh := origBounds.Dx(), origBounds.Dy()
	var nw, nh int
	if ow < oh {
		nw = 200
		nh = int(float64(oh) * (200.0 / float64(ow)))
	} else {
		nh = 200
		nw = int(float64(ow) * (200.0 / float64(oh)))
	}

	resized := image.NewRGBA(image.Rect(0, 0, nw, nh))
	draw.CatmullRom.Scale(resized, resized.Bounds(), img, origBounds, draw.Over, nil)

	// Center crop to 200x200
	x0 := (nw - 200) / 2
	y0 := (nh - 200) / 2
	cropped := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(cropped, cropped.Bounds(), resized, image.Point{X: x0, Y: y0}, draw.Src)

	// Save as JPEG
	destDir := os.Getenv("ROOT_CONTENT_DIR") + "/avatars"
	if destDir == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ROOT CONTENT DIR not set"))
		return
	}
	os.MkdirAll(destDir, 0755)

	var ext string
	if format == "png" {
		ext = ".png"
	} else {
		ext = ".jpg"
	}

	filename := userId + ext
	filepath := filepath.Join(destDir, filename)

	out, err := os.Create(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not save file:" + err.Error()))
		return
	}
	defer out.Close()

	if format == "png" {
		err = png.Encode(out, cropped)
	} else {
		err = jpeg.Encode(out, cropped, &jpeg.Options{Quality: 90})
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not encode image: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Image uploaded and saved as %s", filename)))
}
