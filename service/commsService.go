package service

import (
	"context"
	"log"

	"github.com/leonardchinonso/lokate-go/config"
	"github.com/leonardchinonso/lokate-go/errors"
	"github.com/leonardchinonso/lokate-go/models/dao"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"github.com/leonardchinonso/lokate-go/utils"
)

type commsService struct {
	contactUsRepository interfaces.ContactUsRepositoryInterface
	smtpUsername        string
	smtpPassword        string
	smtpHost            string
	smtpPort            string
}

// NewCommsService returns an interface for the comms service methods
func NewCommsService(cfg *map[string]string, contactUsRepository interfaces.ContactUsRepositoryInterface) interfaces.CommsServiceInterface {
	username := (*cfg)[config.SmtpUsername]
	password := (*cfg)[config.SmtpPassword]
	host := (*cfg)[config.SmtpHost]
	port := (*cfg)[config.SmtpPort]

	return &commsService{
		contactUsRepository: contactUsRepository,
		smtpUsername:        username,
		smtpPassword:        password,
		smtpHost:            host,
		smtpPort:            port,
	}
}

// SendContactUsEmail sends a message from the contact us to the app email
func (cs *commsService) SendContactUsEmail(ctx context.Context, contactUs *dao.ContactUsDAO) error {
	// make sure userId is not empty
	if contactUs.UserId.IsZero() {
		return errors.ErrBadRequest("invalid user id", nil)
	}

	// send the email to the specified app email
	err := utils.SendSimpleMailSMTP(contactUs.UserEmail, cs.smtpUsername, contactUs.Subject, contactUs.Message, cs.smtpUsername, cs.smtpPassword, cs.smtpHost, cs.smtpPort)
	if err != nil {
		log.Printf("Error sending email as plain text. Error: %v\n", err)
		return errors.ErrInternalServerError("failed to send email as plain text", nil)
	}

	// save the contactUs object to the db
	if err = cs.contactUsRepository.Create(ctx, contactUs); err != nil {
		log.Printf("Error saving contactUs object to the db with userId: %v. Error: %v\n", contactUs.UserId, err)
		return errors.ErrInternalServerError("failed to save contactUs object to db", nil)
	}

	return nil
}
