package users

import (
	"context"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/neghi-go/auth-sdk/utils"
	"github.com/neghi-go/database"
)

var (
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9_%+-]+(\.[a-zA-Z0-9_%+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)
)

type User struct {
	store  database.Model[UserModel]
	hasher *utils.Hasher
}

// UserModel represents the user model in the database
type UserModel struct {
	ID        uuid.UUID `json:"id" db:"id,index,unique"`
	Aud       string    `json:"aud" db:"audience"`
	Role      string    `json:"role" db:"role"`
	IsSSOUser bool      `json:"-" db:"is_sso_user"`

	Email           string    `json:"email" db:"email,index.unique"`
	EmailVerified   bool      `json:"email_verified" db:"email_verified"`
	EmailVerifiedAt time.Time `json:"email_verified_at" db:"email_verified_at"`

	EncryptedPassword      string    `json:"-" db:"encrypted_password"`
	EncryptedPasswordSalt  string    `json:"-" db:"encrypted_password_salt"`
	PasswordRecoveryToken  string    `json:"-" db:"password_recovery_token"`
	PasswordRecoverySentAt time.Time `json:"password_recovery_sent_at" db:"password_recovery_sent_at"`

	LastLogin   int       `json:"last_login" db:"last_login"`
	BannedUntil time.Time `json:"banned_until" db:"banned_until"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func New(store database.Model[UserModel]) *User {
	cfg := &User{
		store:  store,
		hasher: utils.NewHasher(),
	}
	return cfg
}

func (u *User) CreateUser(ctx context.Context, email string, opts ...UserModelCreateOptions) (*UserModel, error) {
	if ok := EmailRegex.MatchString(email); !ok {
		return nil, ErrInvalidEmail
	}
	user := &UserModel{
		Email:     email,
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	for _, opt := range opts {
		if err := opt(user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
func (u *User) StoreUser(ctx context.Context, user *UserModel) error {
	if err := u.store.WithContext(ctx).Save(*user); err != nil {
		return err
	}
	return nil
}
func (u *User) RetrieveUser(ctx context.Context, email string) (*UserModel, error) {
	user, err := u.store.WithContext(ctx).Query(database.WithFilter("email", email)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) RetrieveUsers(ctx context.Context) ([]*UserModel, error) {
	user, err := u.store.WithContext(ctx).Query().All()
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *User) UpdateUser(ctx context.Context, user *UserModel) error {
	filter := database.WithFilter("email", user.Email)
	if err := u.store.WithContext(ctx).Query(filter).Update(*user); err != nil {
		return err
	}
	return nil
}
func (u *User) DeleteUser(ctx context.Context, user *UserModel) error {
	filter := database.WithFilter("email", user.Email)
	if err := u.store.WithContext(ctx).Query(filter).Delete(); err != nil {
		return err
	}
	return nil
}
func (u *User) ValidateUserPassword(ctx context.Context, user *UserModel, password string) error {
	return u.hasher.Compare(user.EncryptedPassword, password, user.EncryptedPasswordSalt)
}

func (u *UserModel) Update(opts ...UserModelUpdateOptions) error {
	for _, opt := range opts {
		if err := opt(u); err != nil {
			return err
		}
	}
	return nil
}
