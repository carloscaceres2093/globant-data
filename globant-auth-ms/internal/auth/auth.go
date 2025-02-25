package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"globant-auth-ms/local-lib/web/middleware/log"

	"github.com/google/uuid"
)

type Error struct {
	Message        string
	AdditionalInfo string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.AdditionalInfo)
}

var (
	ErrFileType       = errors.New("file type error")
	ErrFormatModel    = errors.New("error getting columnsm model")
	ErrParseFile      = errors.New("error parsing file")
	ErrNameDuplicated = errors.New("name duplicated")
	ErrInvalidToken   = errors.New("invalid token")
	ErrInvalidUser    = errors.New("invalid user")
)

type Store interface {
	CreateUser(user User) (User, error)
	GetToken(userCode string) (User, error)
}

type AuthService struct {
	store Store
	salt  string
}

func NewAuthService(store Store, salt string) *AuthService {
	return &AuthService{
		store: store,
		salt:  salt,
	}
}

func (e *AuthService) CreateUser(authReq AuthRequest) (AuthResponse, error) {
	var (
		response User
		err      error
	)

	token := generateRandomToken()

	tokenHash := hashToken(token, saltBase64(e.salt))
	userModel := User{
		UserName:  authReq.UserName,
		TokenHash: tokenHash,
		Active:    true,
	}

	err = retry(3, func() (err error) {
		genUUID := uuid.New().String()
		userModel.UserCode = genUUID
		response, err = e.store.CreateUser(userModel)
		return
	})
	if err != nil {

		return AuthResponse{}, err
	}

	return AuthResponse{
		UserName: response.UserName,
		Token:    token,
		UserCode: response.UserCode,
	}, nil
}
func (e *AuthService) ValidateToken(userCode string, token string) (AuthResponse, error) {
	user, err := e.store.GetToken(userCode)
	if err != nil {
		return AuthResponse{}, err
	}
	tokenHash := hashToken(token, saltBase64(e.salt))
	if user.TokenHash != tokenHash {
		return AuthResponse{}, ErrInvalidToken
	}
	return AuthResponse{
		UserName: user.UserName,
		Token:    token,
		UserCode: user.UserCode,
	}, nil
}

func generateRandomToken() string {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		panic("failed to generate random token")
	}
	return base64.URLEncoding.EncodeToString(tokenBytes)
}

func retry(attempts int, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Err(errors.New(fmt.Sprint("retrying after error: ", err.Error())))
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}
func hashToken(token, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(token + salt))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func saltBase64(salt string) string {
	saltBytes := []byte(salt)
	return base64.URLEncoding.EncodeToString(saltBytes)
}
