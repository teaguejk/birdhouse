package images

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

type Result struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func GetLatestImage(dir string) (string, Result) {
	var result Result

	files, err := os.ReadDir(dir)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = "Unable to read directory"

		return "", result
	}

	if len(files) == 0 {
		result.Status = http.StatusNotFound
		result.Message = "No images found"

		return "", result
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

	result.Status = http.StatusOK
	result.Message = "found"

	latestImage := filepath.Join(dir, files[0].Name())

	return latestImage, result
}

func UploadImage(name string, file multipart.File) Result {
	var result Result

	dst, err := os.Create(name)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = "Unable to create file"

		return result
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = "Unable to save file"

		return result
	}

	result.Status = http.StatusCreated
	result.Message = "File uploaded successfully"

	return result
}
