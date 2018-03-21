package telegram

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/receptionist/db"
	"gitlab.com/kanalbot/receptionist/db/models"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
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

func encodeBinary(i interface{}) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		return "", fmt.Errorf("can't encode binary, %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func decodeBinary(enc string, out interface{}) error {
	b64, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return fmt.Errorf("base64 decode failed, %v", err)
	}
	buf := bytes.Buffer{}
	buf.Write(b64)
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(out)
	if err != nil {
		return fmt.Errorf("can't decode, %v", err)
	}
	return nil
}

func isMedia(msg *telegramAPI.Message) bool {
	if msg.Audio != nil || msg.Voice != nil || msg.Photo != nil ||
		msg.Video != nil || msg.Document != nil {
		return true
	}
	return false
}

func getMediaFileID(msg *telegramAPI.Message) (string, error) {
	if msg.Voice != nil {
		return msg.Voice.FileID, nil
	}
	if msg.Audio != nil {
		return msg.Audio.FileID, nil
	}
	if msg.Video != nil {
		return msg.Video.FileID, nil
	}
	if msg.Document != nil {
		return msg.Document.FileID, nil
	}
	if msg.Photo != nil {
		return (*msg.Photo)[len(*msg.Photo)-1].FileID, nil
	}
	return "", errors.New("message doesn't have media")
}

func getMediaURL(msg *telegramAPI.Message) (string, error) {
	fileID, err := getMediaFileID(msg)
	if err != nil {
		return "", err
	}
	return bot.GetFileDirectURL(fileID)
}
