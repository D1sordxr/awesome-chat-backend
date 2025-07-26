package fast

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

type MessageSaveAndBroadcastUseCase struct {
	broadcastURL string
	saveURL      string
	client       *http.Client
	timeout      struct {
		broadcast time.Duration
		save      time.Duration
	}
}

func NewMessageSaveAndBroadcastUseCase(broadcastURL, saveURL string) *MessageSaveAndBroadcastUseCase {
	return &MessageSaveAndBroadcastUseCase{
		broadcastURL: broadcastURL,
		saveURL:      saveURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  false,
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: 20,
			},
		},
		timeout: struct {
			broadcast time.Duration
			save      time.Duration
		}{
			broadcast: 3 * time.Second,
			save:      5 * time.Second,
		},
	}
}

func (uc *MessageSaveAndBroadcastUseCase) Execute(ctx context.Context, payload []byte) error {
	broadcastBuff := bytes.NewBuffer(payload)
	saveBuff := bytes.NewBuffer(payload)

	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(2)
	go uc.broadcastMessage(ctx, broadcastBuff, errChan, &wg)
	go uc.saveMessage(ctx, saveBuff, errChan, &wg)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	errs := make([]error, 0, 2)
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (uc *MessageSaveAndBroadcastUseCase) broadcastMessage(
	ctx context.Context,
	msgBuff *bytes.Buffer,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(ctx, uc.timeout.broadcast)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uc.broadcastURL, msgBuff)
	if err != nil {
		errChan <- err
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := uc.client.Do(req)
	if err != nil {
		errChan <- err
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		errChan <- errors.New("broadcast failed with status: " + resp.Status)
		return
	}

	errChan <- nil
}

func (uc *MessageSaveAndBroadcastUseCase) saveMessage(
	ctx context.Context,
	msgBuff *bytes.Buffer,
	errChan chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(ctx, uc.timeout.save)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uc.saveURL, msgBuff)
	if err != nil {
		errChan <- err
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := uc.client.Do(req)
	if err != nil {
		errChan <- err
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		errChan <- errors.New("save failed with status: " + resp.Status)
		return
	}

	errChan <- nil
}
