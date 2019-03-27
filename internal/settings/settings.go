package settings

import (
	"context"
)

type (
	service struct {
		repository Repository
	}

	Finder interface {
		With(ctx context.Context) Service

		FindByPrefix(prefix string) (kv KV, err error)
	}

	Service interface {
		Finder
	}
)

func NewService(r Repository) Service {
	svc := &service{
		repository: r,
	}

	return svc
}

func (s service) With(ctx context.Context) Service {
	return &service{
		repository: s.repository.With(ctx),
	}
}

func (s service) FindByPrefix(prefix string) (KV, error) {
	if vv, err := s.repository.Find(Filter{Prefix: prefix}); err != nil {
		return nil, err
	} else {
		return vv.KV(), nil
	}
}
