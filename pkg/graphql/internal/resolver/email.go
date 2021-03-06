package resolver

import (
	"context"
	"strconv"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func NewEmail(service iface.Service) *Email {
	return &Email{
		service: service,
	}
}

type Email struct {
	service iface.Service
}

func (r *Email) User(ctx context.Context, e *entity.Email) (*entity.User, error) {
	u, err := r.service.GetUserByEmail(ctx, e.Address)
	if err == nil {
		return entity.NewUser(u), nil
	}
	return nil, Wrap(ctx, err, "fail to get user by email")
}

func (r *Email) Email(ctx context.Context, rawEmailID string) (*entity.Email, error) {
	emailID, err := strconv.ParseInt(rawEmailID, 10, 64)
	if err != nil || emailID == 0 {
		return nil, iface.ErrInvalidID
	}

	emails, err := r.service.FilterEmails(ctx, iface.FilterEmails{EmailID: emailID})
	if err != nil {
		return nil, Wrap(ctx, err, "fail to filter emails")
	}

	if len(emails) == 0 {
		return nil, iface.ErrNotFound
	}

	return entity.NewEmail(emails[0]), nil
}
