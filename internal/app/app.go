package app

import (
	"fmt"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/config"
	v1 "git.spbec-mining.ru/arxon31/sambaMW/internal/controller/http/v1"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/service/webAPI/usecase"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/httpserver"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/samba"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {

	l := logger.New(cfg.App.Name, cfg.App.Env)

	l.Info("application starting", slog.String("version", cfg.App.Version))

	smbClient, err := samba.New(l, cfg.SS.Host, cfg.SS.Port, cfg.SS.User, cfg.SS.Password, cfg.SS.ShareName, cfg.SS.ConnectionPoolSize, cfg.App.TmpDirectoryPath, cfg.App.TmpFilePath)
	if err != nil {
		l.Error("can not create samba client", sl.Err(err))
		os.Exit(1)
	}

	saveFileUseCase := usecase.NewFileSaveUsecase(smbClient, nil, l)
	downloadFileUseCase := usecase.NewFileGetUsecase(smbClient, l)
	listDirectoryUseCase := usecase.NewDirectoryListUsecase(smbClient, l)
	downloadDirectoryUseCase := usecase.NewDirectoryGetUsecase(smbClient, nil, l)
	createDirectoryUseCase := usecase.NewDirectoryCreateUsecase(smbClient, nil, l)

	router := chi.NewRouter()
	v1.NewRouter(router, l, saveFileUseCase, downloadFileUseCase, listDirectoryUseCase, downloadDirectoryUseCase, createDirectoryUseCase)

	server := httpserver.New(router, httpserver.Addr(cfg.HTTP.Host, cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app.Run() signal: ", slog.String("signal", s.String()))
	case err = <-server.Notify():
		l.Error("http server error", sl.Err(fmt.Errorf("app.Run(): httpServer.Notify: %w", err)))
	}

	err = server.Shutdown()
	if err != nil {
		l.Error("http server shutdown error", sl.Err(err))
	}

}
