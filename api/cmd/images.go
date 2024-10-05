package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

func (app *application) GetLatestImage(w http.ResponseWriter, r *http.Request) {
	dir := "/Users/jteague/dev/apps/birdhouse/shared/images"

	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		http.Error(w, "No images found", http.StatusNotFound)
		return
	}

	// sort.Slice(files, func(i, j int) bool {
	// 	infoI, errI := files[i].Info()
	// 	infoJ, errJ := files[j].Info()
	// 	if errI != nil || errJ != nil {
	// 		http.Error(w, "Unable to get file info", http.StatusInternalServerError)
	// 		return false
	// 	}
	// 	return infoI.ModTime().After(infoJ.ModTime())
	// })

	slices.SortFunc(files, func(a, b os.DirEntry) int {
		infoA, errA := a.Info()
		infoB, errB := b.Info()
		if errA != nil || errB != nil {
			http.Error(w, "Unable to get file info", http.StatusInternalServerError)
			return 0
		}
		if infoA.ModTime().After(infoB.ModTime()) {
			return -1
		}
		if infoA.ModTime().Before(infoB.ModTime()) {
			return 1
		}
		return 0
	})

	latestImage := filepath.Join(dir, files[0].Name())

	http.ServeFile(w, r, latestImage)
}

func (app *application) UploadImage(w http.ResponseWriter, r *http.Request) {
	dir := "/Users/jteague/dev/apps/birdhouse/shared/images"

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

	dst, err := os.Create(filepath.Join(dir, handler.Filename))
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Image uploaded successfully"))
}
