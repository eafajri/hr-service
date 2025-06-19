package usecase

import "github.com/eafajri/hr-service.git/module/employee/internal/entity"

//go:generate mockery --name UserUseCase --output ./mocks
type UserUseCase interface {
	GetUserByID(userID int64) (entity.User, error)
}

type UserUseCaseImpl struct {
	userRepository UserRepository
}

func NewProfileUseCase(userRepository UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userRepository: userRepository,
	}
}

func (u *UserUseCaseImpl) GetUserByID(userID int64) (entity.User, error) {
	user, err := u.userRepository.GetUserByID(userID)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
