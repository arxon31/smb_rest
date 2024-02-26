package v1

import (
	"git.spbec-mining.ru/arxon31/sambaMW/internal/service/webAPI/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func NewRouter(
	router chi.Router,
	l *slog.Logger,
	uploadFile usecase.FileSaveUsecase,
	downloadFile usecase.FileGetUsecase,
	listDirectory usecase.DirectoryListUsecase,
	getDirectory usecase.DirectoryGetUsecase,
	createDirectory usecase.DirectoryCreateUsecase) {

	router.Use(middleware.Logger, middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	router.Route("/api/v1", func(router chi.Router) {
		newDirectoryRoutes(router, createDirectory, getDirectory, listDirectory, l)
		newFilesRoutes(router, uploadFile, downloadFile, l)
	})
}
