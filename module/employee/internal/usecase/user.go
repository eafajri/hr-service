package usecase

import "github.com/eafajri/hr-service.git/module/employee/internal/entity"

//go:generate mockery --name UserUseCase --output ./mocks
type UserUseCase interface {
	GetUserByUsernaname(username string) (entity.User, error)
}

type UserUseCaseImpl struct {
	userRepository UserRepository
}

func NewUserUseCase(userRepository UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userRepository: userRepository,
	}
}

func (u *UserUseCaseImpl) GetUserByUsernaname(username string) (entity.User, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
