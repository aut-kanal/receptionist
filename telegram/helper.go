package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/receptionist/db"
	"gitlab.com/kanalbot/receptionist/db/models"
)

func updateUserInfo(userSession *miyanbor.UserSession) {
	newUser := &models.User{
		UserID:    userSession.User.ID,
		ChatID:    userSession.ChatID,
		FirstName: userSession.User.FirstName,
		LastName:  userSession.User.LastName,
		Username:  userSession.User.UserName,
	}

	user := &models.User{}
	db.GetInstance().Where(&models.User{UserID: newUser.UserID, ChatID: newUser.ChatID}).Find(user)
	if !user.IsEqual(newUser) {
		db.GetInstance().Create(newUser)
		logrus.WithField("user-id", newUser.UserID).Debugf("user info updated on db")
	}
}
