package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/receptionist/mq"
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

	if isMedia(msg) {
		encodedMsg, err := encodeBinary(msg)
		if err != nil {
			logrus.Error(err)
			return
		}
		userSession.Payload["msg"] = encodedMsg

		userSession.Payload["media_url"], err = getMediaURL(msg)
		if err != nil {
			logrus.WithError(err).Error("can't get media's URL")
			bot.SendStringMessage(text.MsgNewMessageError, userSession.ChatID)
			return
		}

		bot.AskStringQuestion(text.MsgNewMessageCaptionDialog,
			userSession.UserID, userSession.ChatID, newMessageCaptionHandler)
	} else {
		// Publish msg to message queue
		encodedMsg, err := encodeBinary(msg)
		if err != nil {
			logrus.Error(err)
			return
		}
		err = mq.PublishMsg(&amqp.Publishing{
			ContentType: "application/x-binary",
			Body:        []byte(encodedMsg),
		})
		if err != nil {
			// Send error report
			logrus.WithError(err).Error("can't publish msg to message queue")
			bot.SendStringMessage(text.MsgNewMessageError, userSession.ChatID)
		} else {
			// Send success report
			bot.SendStringMessage(text.MsgNewMessageSuccessful, userSession.ChatID)
		}
	}
}

func newMessageCaptionHandler(userSession *miyanbor.UserSession, input interface{}) {
	captionMsg, ok := input.(*telegramAPI.Message)
	if !ok {
		logrus.Errorln("can't cast input caption to telegram message")
		return
	}

	lastMsg := &telegramAPI.Message{}
	err := decodeBinary(userSession.Payload["msg"].(string), lastMsg)
	if err != nil {
		logrus.Error(err)
		return
	}
	lastMsg.Caption = captionMsg.Text
	delete(userSession.Payload, "msg")

	// Publish msg to message queue
	encodedMsg, err := encodeBinary(lastMsg)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = mq.PublishMsg(&amqp.Publishing{
		ContentType: "application/x-binary",
		Body:        []byte(encodedMsg),
	})
	if err != nil {
		// Send error report
		logrus.WithError(err).Error("can't publish msg to message queue")
		bot.SendStringMessage(text.MsgNewMessageError, userSession.ChatID)
		return
	}

	fileID, _ := getMediaFileID(lastMsg)
	mediaURL, _ := getMediaURL(lastMsg)
	media := Media{
		FileID: fileID,
		URL:    mediaURL,
	}

	// Publish media to URLs queue
	encodedMsg, err = encodeBinary(media)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = mq.PublishMediaURL(&amqp.Publishing{
		ContentType: "application/x-binary",
		Body:        []byte(encodedMsg),
	})
	if err != nil {
		// Send error report
		logrus.WithError(err).Error("can't publish media to media queue")
		bot.SendStringMessage(text.MsgNewMessageError, userSession.ChatID)
		return
	}

	// Send success report
	bot.SendStringMessage(text.MsgNewMessageSuccessful, userSession.ChatID)
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
