package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/IceWreck/BetterBin/db"
	"github.com/IceWreck/BetterBin/logger"
	"github.com/go-chi/chi"
)

const maxUploadSize = 1024 * 1024 * 10 // 10MB

var errUploadingFile = errors.New("error uploading file")
var errFileTooLarge = errors.New("uploaded file is too big")
var errDropNotFound = errors.New("drop not found")

// NewDropPage - webpage to make a new file drop
func NewDropPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_drop", nil)
}

// UploadFile - upload a new file
func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		logger.Error("uploaded file too big", err)
		renderError(w, r, errFileTooLarge, http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("upload")
	if err != nil {
		logger.Error("error uploading file", err)
		renderError(w, r, errUploadingFile, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// create file id
	fileID := newID(10)
	title := r.PostFormValue("title")
	newFileName := fileID + filepath.Ext(fileHeader.Filename)

	// create and save file
	newFile, err := os.Create("./drops/" + newFileName)
	defer newFile.Close()
	if err != nil {
		logger.Error("error uploading file", err)
		renderError(w, r, errUploadingFile, http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(newFile, file); err != nil {
		logger.Error("error uploading file", err)
		renderError(w, r, errUploadingFile, http.StatusInternalServerError)
		return
	}

	// db entry at the end only if upload succeeded
	err = db.NewDrop(fileID, title, newFileName)
	if err != nil {
		logger.Error("error creating db entry", err)
		renderError(w, r, errUploadingFile, http.StatusInternalServerError)
		return
	}

	logger.Info("uploaded", newFileName, "size", fileHeader.Size)
	renderSuccess(w, r, fileID)
}

// ViewDrop - preview and download file drop
func ViewDrop(w http.ResponseWriter, r *http.Request) {
	dropID := chi.URLParam(r, "dropID")
	drop, err := db.GetDrop(dropID)
	if err != nil {
		renderTemplate(w, "drop_not_found", nil)
		return
	}
	// preview type for formats that the browser can display natively
	switch filepath.Ext(drop.FileName) {
	case ".jpg", ".jpeg", ".webp", ".png", ".gif":
		drop.Preview = "image"
		logger.Debug("preview image")
	case ".mp4", ".webm":
		drop.Preview = "video"
		logger.Debug("preview video")
	default:
		drop.Preview = "none"
		logger.Debug("no preview")

	}

	renderTemplate(w, "view_drop", drop)
}
