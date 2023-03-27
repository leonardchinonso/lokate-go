package interfaces

import (
	"context"
	"github.com/leonardchinonso/lokate-go/models/dao"
)

// AboutRepositoryInterface defines the methods for the about repository
type AboutRepositoryInterface interface {
	GetDetails(ctx context.Context, about *dao.About) error
}

// AboutServiceInterface defines the methods for the about action
type AboutServiceInterface interface {
	Details(ctx context.Context) (*dao.About, error)
}
