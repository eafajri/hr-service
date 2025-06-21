package usecase_test

import (
	"errors"
	"testing"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserUseCase_GetUserByUsernaname(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(userRepository *mocks.UserRepository)
		wantErr  error
		wantRes  entity.User
	}{
		{
			name: "success",
			mockFunc: func(userRepository *mocks.UserRepository) {
				userRepository.On("GetUserByUsername", mock.Anything).Return(entity.User{ID: 12, Username: "username"}, nil)
			},
			wantErr: nil,
			wantRes: entity.User{
				ID:       12,
				Username: "username",
			},
		},
		{
			name: "error",
			mockFunc: func(userRepository *mocks.UserRepository) {
				userRepository.On("GetUserByUsername", mock.Anything).Return(entity.User{}, errors.New("error"))
			},
			wantErr: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := mocks.NewUserRepository(t)

			tt.mockFunc(userRepository)

			usecase := usecase.NewUserUseCase(userRepository)
			res, err := usecase.GetUserByUsernaname("")
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.wantRes, res)
				assert.NoError(t, err)
			}
		})
	}
}
