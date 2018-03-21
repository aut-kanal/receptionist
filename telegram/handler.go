package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/receptionist/configuration"
	"gitlab.com/kanalbot/receptionist/mq"
	"gitlab.com/kanalbot/receptionist/ui/text"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

func sessionStartHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.WithField("user", *userSession).Debugf("new session started")

	updateUserInfo(userSession)
}

func unknownMessageHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.WithField("user", *userSession).Debugf("unknown message received, %+v", input)
}

func cancelCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debug("cancel command")
	userSession.ResetSession()
	bot.SendStringMessage(text.MsgCanceledSuccessfully, userSession.ChatID)
}

func newMessageCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.WithField("user", *userSession).Debugf("newmessage command received")

	bot.AskStringQuestion(text.MsgNewMessageDialog, userSession.UserID, userSession.ChatID, newMessageContentHandler)
}

func newMessageContentHandler(userSession *miyanbor.UserSession, input interface{}) {
	telegramMsg, ok := input.(*telegramAPI.Message)
	if !ok {
		logrus.Errorln("can't cast input content to telegram message")
		return
	}
	msg := &Message{
		Message: telegramMsg,
	}

	if isMedia(msg.Message) {
		encodedMsg, err := encodeBinary(msg)
		if err != nil {
			logrus.Error(err)
			return
		}
		userSession.Payload["msg"] = encodedMsg

		userSession.Payload["media_url"], err = getMediaURL(msg.Message)
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
	lastMsg := &Message{}

	// Add caption to lastMsg
	captionMsg, ok := input.(*telegramAPI.Message)
	if !ok {
		logrus.Errorln("can't cast input caption to telegram message")
		return
	}
	err := decodeBinary(userSession.Payload["msg"].(string), lastMsg)
	if err != nil {
		logrus.Error(err)
		return
	}
	lastMsg.Caption = captionMsg.Text
	delete(userSession.Payload, "msg")

	// Add file's URL to lastMsg
	mediaURL, _ := getMediaURL(lastMsg.Message)
	lastMsg.FileURL = mediaURL

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

	// Send success report
	bot.SendStringMessage(text.MsgNewMessageSuccessful, userSession.ChatID)
}

func kanalCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("kanal command received")
	bot.SendStringMessage(text.MsgKanalLink, userSession.ChatID)
}

func feedbackCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("feedback command received, user: %v", *userSession)
	bot.AskStringQuestion(text.MsgFeedback, userSession.UserID, userSession.ChatID, feedbackMessageHandler)
}

func feedbackMessageHandler(userSession *miyanbor.UserSession, input interface{}) {
	adminChatID := configuration.GetInstance().GetInt64("admin-chatid")
	feedbackMsg := input.(*telegramAPI.Message)
	feedbackForward := telegramAPI.NewForward(adminChatID, feedbackMsg.Chat.ID, feedbackMsg.MessageID)
	bot.Send(feedbackForward)
}

func helpCommandHandler(userSession *miyanbor.UserSession, matches interface{}) {
	logrus.Debugf("help command received")
	bot.SendStringMessage(text.MsgHelp, userSession.ChatID)
}
