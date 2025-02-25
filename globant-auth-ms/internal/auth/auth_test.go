package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type StoreMock struct {
	mock.Mock
}
type SaltMock struct {
	mock.Mock
}

func (s *StoreMock) CreateUser(_ User) (User, error) {
	args := s.Called()
	return args.Get(0).(User), args.Error(1)
}
func (s *StoreMock) GetToken(_ string) (User, error) {
	args := s.Called()
	return args.Get(0).(User), args.Error(1)
}
func TestService_CreateUser(t *testing.T) {
	serviceReq := AuthRequest{
		UserName: "test",
	}
	dbResponse := User{
		UserCode:  "test",
		UserName:  "test",
		TokenHash: "test",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	responseService := AuthResponse{
		UserName: dbResponse.UserName,
		Token:    dbResponse.TokenHash,
		UserCode: dbResponse.UserCode,
	}
	var tests = []struct {
		name              string
		store             *StoreMock
		satl              string
		request           AuthRequest
		expectedRepsponse AuthResponse
		expectedError     error
	}{
		{
			name: "Ok - CreateUser",
			store: func() *StoreMock {
				m := StoreMock{}
				m.On("CreateUser", mock.Anything).Return(dbResponse, nil)
				return &m
			}(),
			satl:              "test",
			request:           serviceReq,
			expectedRepsponse: responseService,
		},
		{
			name: "Fail - Saving",
			store: func() *StoreMock {
				m := StoreMock{}
				m.On("CreateUser", mock.Anything).Return(dbResponse, errors.New("random error"))
				return &m
			}(),
			satl:              "test",
			request:           serviceReq,
			expectedRepsponse: responseService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAuthService(tt.store, tt.satl)
			result, err := service.CreateUser(tt.request)
			require.Equal(t, tt.expectedRepsponse, result)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestService_ValidateToken(t *testing.T) {
	userCode := "test"
	token := "test"
	tokenWrong := "test_wrong"
	dbResponse := User{
		UserCode:  "test",
		UserName:  "test",
		TokenHash: "test",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	responseService := AuthResponse{
		UserName: dbResponse.UserName,
		Token:    token,
		UserCode: dbResponse.UserCode,
	}
	var tests = []struct {
		name              string
		store             *StoreMock
		satl              string
		userCode          string
		token             string
		expectedRepsponse AuthResponse
		expectedError     error
	}{
		{
			name: "Ok - Validate",
			store: func() *StoreMock {
				m := StoreMock{}
				m.On("GetToken", mock.Anything).Return(dbResponse, nil)
				return &m
			}(),
			userCode:          userCode,
			token:             token,
			satl:              "test",
			expectedRepsponse: responseService,
		},
		{
			name: "Fail - Get Data DB",
			store: func() *StoreMock {
				m := StoreMock{}
				m.On("GetToken", mock.Anything).Return(User{}, errors.New("random error"))
				return &m
			}(),
			userCode:          userCode,
			token:             token,
			satl:              "test",
			expectedRepsponse: responseService,
			expectedError:     errors.New("OK"),
		},
		{
			name: "Fail - TokenValid",
			store: func() *StoreMock {
				m := StoreMock{}
				m.On("GetToken", mock.Anything).Return(dbResponse, nil)
				return &m
			}(),
			userCode:          userCode,
			token:             tokenWrong,
			satl:              "test",
			expectedRepsponse: responseService,
			expectedError:     errors.New("OK"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAuthService(tt.store, tt.satl)
			result, err := service.ValidateToken(tt.userCode, tt.token)
			require.Equal(t, tt.expectedRepsponse, result)
			require.Equal(t, tt.expectedError, err)
		})
	}
}
