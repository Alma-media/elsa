package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Alma-media/elsa/pipe"
)

type Storage interface {
	Load(context.Context) (pipe.Pipe, error)
	Save(context.Context, pipe.Pipe) error
}

type Manager interface {
	Apply(context.Context, pipe.Pipe) (<-chan struct{}, error)
}

type Handler struct {
	mu sync.Mutex

	storage Storage
	manager Manager

	await  <-chan struct{}
	cancel func()
}

func restore(storage Storage, manager Manager) (func(), <-chan struct{}, error) {
	ctx, cancel := context.WithCancel(context.Background())

	routes, err := storage.Load(ctx)
	if err != nil {
		cancel()

		return nil, nil, err
	}

	await, err := manager.Apply(ctx, routes)
	if err != nil {
		cancel()

		return nil, nil, err
	}

	return cancel, await, err
}

func NewHandler(storage Storage, manager Manager) (*Handler, error) {
	cancel, await, err := restore(storage, manager)
	if err != nil {
		return nil, err
	}

	return &Handler{
		storage: storage,
		manager: manager,

		await:  await,
		cancel: cancel,
	}, nil
}

func (h *Handler) Stop() {
	h.cancel()
	<-h.await
}

func (h *Handler) ApplyHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	defer r.Body.Close()

	var pipe pipe.Pipe

	if err := json.NewDecoder(r.Body).Decode(&pipe); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	h.Stop()

	ctx, cancel := context.WithCancel(context.Background())

	await, err := h.manager.Apply(ctx, pipe)
	if err != nil {
		cancel()
		<-await

		h.cancel, h.await, _ = restore(h.storage, h.manager)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := h.storage.Save(r.Context(), pipe); err != nil {
		cancel()
		<-await

		h.cancel, h.await, _ = restore(h.storage, h.manager)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	h.cancel, h.await = cancel, await
}