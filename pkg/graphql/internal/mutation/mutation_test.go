package mutation

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)

	m := NewMutation(service)

	ctx := context.TODO()

	// succeed
	{
		name := "name"

		service.EXPECT().AddUser(ctx, name).Return(int64(1), nil)

		u, err := m.AddUser(ctx, entity.AddUserInput{
			Name: name,
		})
		assert.Nil(t, err)
		assert.NotNil(t, u)
	}

	// fails if service fails
	{
		name := "name"

		service.EXPECT().AddUser(ctx, name).Return(int64(0), fmt.Errorf("opz"))

		u, err := m.AddUser(ctx, entity.AddUserInput{
			Name: name,
		})
		assert.Equal(t, err.Error(), "service failed")
		assert.Nil(t, u)
	}
}

func TestAddEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockService(ctrl)

	m := NewMutation(service)

	ctx := context.TODO()

	// succeed
	{
		address := "email@email.com"
		userID := int64(12)

		service.EXPECT().AddEmail(ctx, userID, address).Return(int64(1), nil)

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  strconv.FormatInt(userID, 10),
			Address: address,
		})
		assert.Nil(t, err)
		assert.Equal(t, u.Email.ID, "1")
	}

	// fails if userID is invalid
	{
		address := "email@email.com"
		userID := "0"

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Equal(t, err.Error(), "input: invalid userID")
		assert.Nil(t, u)
	}

	// fails if email is invalid
	{
		address := "email"
		userID := "1"

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Equal(t, err.Error(), "input: invalid email address")
		assert.Nil(t, u)
	}

	// fails if service fails with duplicated
	{
		address := "email@email.com"
		userID := int64(12)

		service.EXPECT().AddEmail(ctx, userID, address).Return(int64(0), iface.ErrAlreadyExists)

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  strconv.FormatInt(userID, 10),
			Address: address,
		})
		assert.Equal(t, err.Error(), fmt.Sprintf("input: %v", iface.ErrAlreadyExists))
		assert.Nil(t, u)
	}

	// fails if service fails
	{
		address := "email@email.com"
		userID := int64(12)

		service.EXPECT().AddEmail(ctx, userID, address).Return(int64(0), fmt.Errorf("opz"))

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  strconv.FormatInt(userID, 10),
			Address: address,
		})
		assert.Equal(t, err.Error(), "service failed")
		assert.Nil(t, u)
	}
}
