package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"
)

// AddUser create a new user in the database
func (s *Storage) AddUser(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	now := time.Now()
	return Insert(ctx, tx, "INSERT INTO users (name, created, updated) VALUES (?, ?, ?)", name, now, now)
}

// DeleteUser remove an user from the database
func (s *Storage) DeleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	return Delete(ctx, tx, "DELETE FROM users WHERE id = ?", userID)
}

// FilterUsersID retrieve usersID from the database for a given filter
func (s *Storage) FilterUsersID(ctx context.Context, filter iface.FilterUsers) ([]int64, error) {
	limit := iface.FilterUsersDefaultLimit
	if filter.Limit != 0 {
		limit = filter.Limit
	}

	var args []interface{}
	var query string

	if len(filter.Email) != 0 {
		query = "SELECT u.id FROM users u INNER JOIN emails e ON(e.user_id = u.id) WHERE e.address = ?"
		args = append(args, filter.Email)
	} else {
		query = "SELECT id FROM users LIMIT ?"
		args = append(args, limit)
	}

	rows, err := Select(ctx, s.sql, scanInt, query, args...)
	if err != nil {
		return nil, err
	}

	IDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row != nil {
			IDs = append(IDs, row.(int64))
		}
	}

	return IDs, nil
}

// FetchUsers retrieve users from the database
func (s *Storage) FetchUsers(ctx context.Context, IDs ...int64) ([]*entity.User, error) {
	if len(IDs) == 0 {
		return make([]*entity.User, 0), nil
	}

	query := fmt.Sprintf(
		"SELECT id, name, created, updated "+
			"FROM users WHERE id IN (%s)",
		strings.Repeat("?,", len(IDs))[0:len(IDs)*2-1])

	args := make([]interface{}, 0, len(IDs))
	for _, ID := range IDs {
		args = append(args, ID)
	}
	rows, err := Select(ctx, s.sql, scanUser, query, args...)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, row.(*entity.User))
	}

	return users, nil
}

func scanUser(sc func(dest ...interface{}) error) (interface{}, error) {
	var id int64
	var name string
	var created time.Time
	var updated time.Time

	err := sc(&id, &name, &created, &updated)
	if err != nil {
		return nil, errors.New("could not scan user").SetParent(err)
	}

	return &entity.User{
		ID:      id,
		Name:    name,
		Created: created,
		Updated: updated,
	}, nil
}
