package services

import (
	"fmt"

	"resqiar.com-server/constants"
	"resqiar.com-server/entities"
	"resqiar.com-server/inputs"
	"resqiar.com-server/repositories"
)

type UserService interface {
	GetUsernameList() ([]string, error)
	RegisterUser(profile *entities.GooglePayload) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
	FindUserByID(userID string) (*entities.SafeUser, error)
	FindUserByUsername(username string) (*entities.SafeUser, error)
	CheckUsernameExist(username string) bool
	UpdateUser(payload *inputs.UpdateUserInput, userID string) error
}

type UserServiceImpl struct {
	UtilService UtilService
	Repository  repositories.UserRepository
}

func (service *UserServiceImpl) GetUsernameList() ([]string, error) {
	result, err := service.Repository.GetUsernameList()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *UserServiceImpl) RegisterUser(profile *entities.GooglePayload) (*entities.User, error) {
	// format the given name from the provider
	formattedName := service.UtilService.FormatUsername(profile.GivenName)

	// concatenate formatted name with the nano id
	formattedName = fmt.Sprintf("%s_%s", formattedName, service.UtilService.GenerateRandomID(7))

	newUser := &entities.User{
		Username:   formattedName,
		Email:      profile.Email,
		Provider:   constants.Google,
		ProviderID: profile.SUB,
		PictureURL: profile.Picture,
	}

	result, err := service.Repository.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *UserServiceImpl) FindUserByEmail(email string) (*entities.User, error) {
	result, err := service.Repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *UserServiceImpl) FindUserByID(userID string) (*entities.SafeUser, error) {
	safeUser, err := service.Repository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return safeUser, nil
}

func (service *UserServiceImpl) FindUserByUsername(username string) (*entities.SafeUser, error) {
	safeUser, err := service.Repository.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return safeUser, nil
}

func (service *UserServiceImpl) CheckUsernameExist(username string) bool {
	exist, _ := service.Repository.FindByUsername(username)
	if exist != nil {
		return true
	}

	return false
}

func (service *UserServiceImpl) UpdateUser(payload *inputs.UpdateUserInput, userID string) error {
	user, err := service.Repository.FindByID(userID)
	if err != nil {
		return err
	}

	if err := service.Repository.UpdateUser(user.ID, payload); err != nil {
		return err
	}

	return nil
}
