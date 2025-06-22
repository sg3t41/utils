package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string
	Email     string
	Name      string
	Password  string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewUser(email, name string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) UpdateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	u.Version++
	return nil
}

func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	u.Email = email
	u.UpdatedAt = time.Now()
	u.Version++
	return nil
}

func (u *User) UpdatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	u.Password = password
	u.UpdatedAt = time.Now()
	u.Version++
	return nil
}

func (u *User) Update(name, email *string) error {
	var updated bool
	
	if name != nil && *name != u.Name {
		if *name == "" {
			return errors.New("name cannot be empty")
		}
		u.Name = *name
		updated = true
	}
	
	if email != nil && *email != u.Email {
		if *email == "" {
			return errors.New("email cannot be empty")
		}
		u.Email = *email
		updated = true
	}
	
	if updated {
		u.UpdatedAt = time.Now()
		u.Version++
	}
	
	return nil
}

func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
	u.Version++
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}