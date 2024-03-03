package v1

import (
	"context"
	"encoding/json"
	"errors"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"path"
)

const (
	maxFileSize = 10 * 1024 * 1024
	fileKey     = "content"
	pathKey     = "filepath"
)

type Saver interface {
	SaveFile(ctx context.Context, request entity.FileSaveRequest) (response entity.FileSaveResponse, err error)
}

type Getter interface {
	GetFile(ctx context.Context, request entity.FileGetRequest) (response entity.FileGetResponse, err error)
}

type filesRoutes struct {
	saver  Saver
	getter Getter
	l      *slog.Logger
}

func newFilesRoutes(router chi.Router, fileSaver Saver, fileGetter Getter, l *slog.Logger) {
	r := &filesRoutes{
		saver:  fileSaver,
		getter: fileGetter,
		l:      l,
	}

	router.Route("/file", func(router chi.Router) {
		router.Post("/get", r.getFile)
		router.Post("/put", r.putFile)
	})
}

func (r *filesRoutes) getFile(w http.ResponseWriter, req *http.Request) {
	const op = "http.v1.files.getFile()"
	logger := r.l.With(slog.String("operation", op))

	// Собираем запрос и валидируем его
	var model entity.FileGetRequest
	err := json.NewDecoder(req.Body).Decode(&model)
	if err != nil {
		logger.Error("unable to decode request body", sl.Err(err))
		http.Error(w, errInternalError.Error(), http.StatusBadRequest)
		return
	}
	err = model.Validate()
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := r.getter.GetFile(req.Context(), model)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Error("file not found", sl.Err(err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logger.Error("unable to get file", sl.Err(err))
		http.Error(w, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("unable to encode response", sl.Err(err))
		http.Error(w, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	http.ServeFile(w, req, resp.FilePath)
}

func (r *filesRoutes) putFile(w http.ResponseWriter, req *http.Request) {
	const op = "http.v1.files.putFile()"

	logger := r.l.With(slog.String("operation", op))

	// Устанавливаем максимальный размер файла
	err := req.ParseMultipartForm(maxFileSize)
	if err != nil {
		logger.Error("unable to parse multipart form", sl.Err(err))
		http.Error(w, errFileSizeIsTooBig.Error(), http.StatusBadRequest)
		return
	}

	// Читаем по ключу файл из формы
	file, fileHandler, err := req.FormFile(fileKey)
	if err != nil {
		logger.Error("unable to get file from form", sl.Err(err), slog.String("key", fileKey))
		http.Error(w, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	// Читаем по ключу путь к папке из формы
	filePath := req.FormValue(pathKey)
	if filePath == "" {
		logger.Error("unable to get file path from form", sl.Err(err), slog.String("key", pathKey))
		http.Error(w, errBadRequest.Error(), http.StatusBadRequest)
		return
	}

	// Читаем содержимое файла и помещаем в структуру модели
	b := make([]byte, fileHandler.Size)
	_, err = file.Read(b)
	if err != nil {
		logger.Error("unable to read file", sl.Err(err))
		http.Error(w, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	var model = entity.FileSaveRequest{
		FilePath: path.Join(filePath, fileHandler.Filename),
		Content:  b,
	}

	// Валидируем структуру модели
	err = model.Validate()
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отдаем юзкейсу модель для сохранения, в ответ получаем модель с полным путём к файлу
	resp, err := r.saver.SaveFile(req.Context(), model)
	if err != nil {
		logger.Error("unable to save file", sl.Err(err))
		http.Error(w, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	// Выставляем заголовок и отдаеем ответ с полным путём к сохраненному файлу
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("unable to encode response", sl.Err(err))
		http.Error(w, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
