package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
)

func sessionStartHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.Debugf("new session started")
}

func unknownMessageHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.Debugf("unknown message received, %+v", input)
}

func newMessageCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("newmessage command received")
}

func kanalCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("kanal command received")
}

func feedbackCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("feedback command received")
}

func helpCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("help command received")
}
