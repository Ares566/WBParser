package repository

import (
	"WBParser/internal/domain"
	"context"
)

type Product interface {
	Add(ctx context.Context, products []domain.ProductCard) error
	Update(ctx context.Context, product domain.ProductCard) error
}
