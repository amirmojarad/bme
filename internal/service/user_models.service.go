package service

import (
	"bme/internal/constants"
	"bme/pkg/errorext"
	"bme/pkg/jwt"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

type UserEntity struct {
	ID             uint
	Username       string
	PhoneNumber    string
	FirstName      string
	LastName       string
	HashedPassword string
}

type AuthRegisterRequest struct {
	AuthClaims
	FirstName   string
	LastName    string
	PhoneNumber string
}

type AuthClaims struct {
	Username string
	Password string
}

type AuthLoginRequest struct {
	AuthClaims
}

type CreateUserRequest struct {
	Username            string
	HashedPassword      string
	FirstName, LastName string
	PhoneNumber         string
}

type FirstUserFilter struct {
	ID             *uint
	Username       *string
	TelegramUserID *uint
}

type AuthUserWithRoleEntity struct {
	UserID uint               `json:"user_id"`
	Role   constants.UserRole `json:"role"`
}

func (f FirstUserFilter) FilterMap() map[string]any {
	filterMap := make(map[string]any)

	if f.ID != nil {
		filterMap["id"] = f.ID
	}

	if f.TelegramUserID != nil {
		filterMap["telegram_user_id"] = f.TelegramUserID
	}

	if f.Username != nil {
		filterMap["username"] = f.Username
	}

	return filterMap
}

func (req AuthClaims) validate() error {
	if req.Username == "" {
		return errorext.NewValidation(errors.New("username can not be empty"), errorext.ErrValidation)
	}

	if req.Password == "" {
		return errorext.NewValidation(errors.New("password can not be empty"), errorext.ErrValidation)
	}

	return nil
}

func (entity UserEntity) UserClaims() jwt.UserClaims {
	return jwt.UserClaims{
		UserID:   entity.ID,
		Username: entity.Username,
	}
}

func (entity UserEntity) String() string {
	if entity.FirstName != "" && entity.LastName != "" {
		return fmt.Sprintf("%s %s", entity.FirstName, entity.LastName)
	}

	return entity.Username
}

func (entity UserEntity) isEmpty() bool {
	return entity.ID == 0
}

func (entity *AuthUserWithRoleEntity) key() string {
	return strconv.Itoa(int(entity.UserID))
}

func (entity *AuthUserWithRoleEntity) isAdmin() bool {
	return entity.Role == constants.UserRoleAdmin
}

func (entity *AuthUserWithRoleEntity) valueMap() map[string]any {
	return map[string]any{
		"role": entity.Role.String(),
	}
}

func (entity *AuthUserWithRoleEntity) value() ([]byte, error) {
	valueMap := entity.valueMap()

	valueMapAsBytes, err := json.Marshal(&valueMap)
	if err != nil {
		return nil, err
	}

	return valueMapAsBytes, nil
}

func (entity *AuthUserWithRoleEntity) createFromBytes(data []byte) error {
	return json.Unmarshal(data, entity)
}

func (entity *AuthUserWithRoleEntity) IsAdmin() bool {
	return entity.Role == constants.UserRoleAdmin
}
