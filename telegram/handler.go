package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/receptionist/ui/text"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

func sessionStartHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.Debugf("new session started")

	updateUserInfo(userSession)
}

func unknownMessageHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.Debugf("unknown message received, %+v", input)
}

func newMessageCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("newmessage command received")

	bot.AskStringQuestion(text.MsgNewMessageDialog, userSession.UserID, userSession.ChatID, newMessageContentHandler)
}

func newMessageContentHandler(userSession *miyanbor.UserSession, input interface{}) {
	msg, ok := input.(*telegramAPI.Message)
	if !ok {
		logrus.Errorln("can't cast input content to telegram message")
		return
	}

	if msg.Audio != nil || msg.Photo != nil || msg.Video != nil || msg.Document != nil {
		bot.AskStringQuestion(text.MsgNewMessageCaptionDialog,
			userSession.UserID, userSession.ChatID, newMessageCaptionHandler)
	} else {
		bot.SendStringMessage(text.MsgNewMessageSuccessful, userSession.ChatID)
		// TODO: Send content to admins
	}
}

func newMessageCaptionHandler(userSession *miyanbor.UserSession, input interface{}) {
	_, ok := input.(*telegramAPI.Message)
	if !ok {
		logrus.Errorln("can't cast input caption to telegram message")
		return
	}
	bot.SendStringMessage(text.MsgNewMessageSuccessful, userSession.ChatID)
	// TODO: Send content to admins
}

func kanalCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("kanal command received")
	bot.SendStringMessage(text.MsgKanalLink, userSession.ChatID)
}

func feedbackCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("feedback command received")
}

func helpCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("help command received")
	bot.SendStringMessage(text.MsgHelp, userSession.ChatID)
}
