package services

import (
	"fmt"
	"github.com/vivasoft-ltd/go-ems/domain"
	"github.com/vivasoft-ltd/golang-course-utils/logger"
)

type Mail struct {
	userRepo domain.UserRepository
}

func NewMailService(userRepo domain.UserRepository) *Mail {
	return &Mail{
		userRepo: userRepo,
	}
}
func (m *Mail) SendInvitationEmail(userIds []int) error {
	users, err := m.userRepo.ReadUsers(userIds)
	if err != nil {
		return err
	}
	for _, user := range users {
		logger.Info(fmt.Sprintf("Send Invitation Mail to %s", user.Email))
	}

	return nil
}
