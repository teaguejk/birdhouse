package main

import (
	"net/http"
	"path/filepath"

	"birdhouse.api/internal/images"
)

func (app *application) GetLatestImage(w http.ResponseWriter, r *http.Request) {
	dir := "/Users/jteague/dev/apps/birdhouse/shared/images"

	latestImage, result := images.GetLatestImage(dir)

	if result.Status != http.StatusOK {
		http.Error(w, result.Message, result.Status)
		return
	}

	http.ServeFile(w, r, latestImage)
}

func (app *application) UploadImage(w http.ResponseWriter, r *http.Request) {
	dir := "/Users/jteague/dev/apps/birdhouse/shared/images"

	app.logger.Info("uploading image")

	err := r.ParseMultipartForm(10 << 20) // limit upload size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := filepath.Join(dir, handler.Filename)

	result := images.UploadImage(filename, file)

	app.logger.Info("image uploaded successfully", "filename", handler.Filename)

	w.WriteHeader(result.Status)
	w.Write([]byte(result.Message))
}
