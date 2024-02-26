package v1

import (
	"context"
	"encoding/json"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type DirCreator interface {
	CreateDir(ctx context.Context, request entity.DirCreateRequest) (response entity.DirCreateResponse, err error)
}

type DirGetter interface {
	GetDirectory(ctx context.Context, request entity.DirGetRequest) (response entity.DirGetResponse, err error)
}

type DirLister interface {
	ListDir(ctx context.Context, request entity.DirListRequest) (response entity.DirListResponse, err error)
}

type directoryRoutes struct {
	creator DirCreator
	getter  DirGetter
	lister  DirLister
	l       *slog.Logger
}

func newDirectoryRoutes(router chi.Router, creator DirCreator, getter DirGetter, lister DirLister, l *slog.Logger) {
	r := &directoryRoutes{
		creator: creator,
		getter:  getter,
		lister:  lister,
		l:       l,
	}

	router.Route("/dir", func(router chi.Router) {
		router.Post("/create", r.createDir)
		router.Post("/get", r.getDir)
		router.Post("/list", r.listDir)
	})
}

func (r *directoryRoutes) createDir(writer http.ResponseWriter, request *http.Request) {
	const op = "http.v1.directoryRoutes.createDir()"
	logger := r.l.With(slog.String("operation", op))

	var model entity.DirCreateRequest

	err := json.NewDecoder(request.Body).Decode(&model)
	if err != nil {
		logger.Error("unable to decode request body", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusBadRequest)
		return
	}

	err = model.Validate()
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := r.creator.CreateDir(request.Context(), model)
	if err != nil {
		logger.Error("unable to create dir", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		logger.Error("unable to encode response", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

}

func (r *directoryRoutes) getDir(writer http.ResponseWriter, request *http.Request) {
	const op = "http.v1.directoryRoutes.getDir()"
	logger := r.l.With(slog.String("operation", op))

	var model entity.DirGetRequest

	err := json.NewDecoder(request.Body).Decode(&model)
	if err != nil {
		logger.Error("unable to decode request body", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusBadRequest)
		return
	}

	err = model.Validate()
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := r.getter.GetDirectory(request.Context(), model)
	if err != nil {
		logger.Error("unable to get dir", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		logger.Error("unable to encode response", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(writer, request, resp.DirPath)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func (r *directoryRoutes) listDir(writer http.ResponseWriter, request *http.Request) {
	const op = "http.v1.directoryRoutes.listDir()"
	logger := r.l.With(slog.String("operation", op))

	var model entity.DirListRequest

	err := json.NewDecoder(request.Body).Decode(&model)
	if err != nil {
		logger.Error("unable to decode request body", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusBadRequest)
		return
	}

	err = model.Validate()
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := r.lister.ListDir(request.Context(), model)
	if err != nil {
		logger.Error("unable to list dir", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		logger.Error("unable to encode response", sl.Err(err))
		http.Error(writer, errInternalError.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
