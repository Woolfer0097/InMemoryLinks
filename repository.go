package inmemorylinks

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrNotFound   = errors.New("url not found")
	ErrCodeExists = errors.New("url code already exists")
)

type Url struct {
	ID          uuid.UUID
	UrlCode     string
	OriginalUrl string
}

type UrlRepository struct {
	mu     sync.RWMutex
	byCode map[string]*Url
}

func NewUrlRepository() *UrlRepository {
	return &UrlRepository{
		byCode: make(map[string]*Url),
	}
}

func (r *UrlRepository) Create(ctx context.Context, url *Url) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if url == nil {
		return errors.New("url cannot be nil")
	}
	if url.UrlCode == "" {
		return errors.New("url code cannot be empty")
	}
	if _, exists := r.byCode[url.UrlCode]; exists {
		return ErrCodeExists
	}
	stored := *url
	r.byCode[url.UrlCode] = &stored
	return nil
}

func (r *UrlRepository) GetByCode(ctx context.Context, code string) (*Url, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	url, ok := r.byCode[code]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *url
	return &cp, nil
}
