package interfaces

import (
	"context"
	"github.com/leonardchinonso/lokate-go/models/dao"
)

// ContactUsRepositoryInterface defines the methods for the contact us repository
type ContactUsRepositoryInterface interface {
	Create(ctx context.Context, contactUs *dao.ContactUsDAO) error
}

type ContactUsServiceInterface interface {
	SendContactUsEmail(ctx context.Context, contactUs *dao.ContactUsDAO) error
}
