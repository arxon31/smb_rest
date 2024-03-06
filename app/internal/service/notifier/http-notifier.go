package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
	"net/http"
)

type Notifier struct {
	client   *http.Client
	endpoint string
	logger   *slog.Logger
}

func NewNotifier(endpoint string, logger *slog.Logger) *Notifier {
	client := &http.Client{}

	return &Notifier{
		client:   client,
		endpoint: endpoint,
		logger:   logger,
	}
}

func (n *Notifier) Notify(ctx context.Context, dirs []string) error {
	err := n.send(dirs)
	if err != nil {
		n.logger.Error("failed to notify", sl.Err(err))
		return err
	}
	return nil
}

func (n *Notifier) send(dirs []string) error {
	dirsModel := entity.DirNotify{Dirs: dirs}
	dirsJSON, err := json.Marshal(dirsModel)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(dirsJSON)

	req, err := http.NewRequest(http.MethodPost, n.endpoint, buf)
	if err != nil {
		return err
	}

	_, err = n.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
