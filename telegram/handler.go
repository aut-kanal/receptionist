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

func sessionStartHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.WithField("user", *userSession).Debugf("new session started")

	updateUserInfo(userSession)
}

func unknownMessageHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.WithField("user", *userSession).Debugf("unknown message received")
	// bot.SendStringMessage(text.MsgUnknownInput, userSession.ChatID)
	replyMsg := telegramAPI.NewMessage(userSession.ChatID, text.MsgUnknownInput)
	replyMsg.ReplyMarkup = telegramAPI.NewRemoveKeyboard(false)
	bot.Send(replyMsg)
}

func cancelCommandHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.Debug("cancel command")
	userSession.ResetSession()
	bot.SendStringMessage(text.MsgCanceledSuccessfully, userSession.ChatID)
}

func newMessageCommandHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.WithField("user", *userSession).Debugf("newmessage command received")

	bot.AskStringQuestion(text.MsgNewMessageDialog, userSession.UserID, userSession.ChatID, newMessageContentHandler)
}

func newMessageContentHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	telegramMsg := update.(*telegramAPI.Update).Message
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

func newMessageCaptionHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	lastMsg := &Message{}

	// Add caption to lastMsg
	captionMsg := update.(*telegramAPI.Update).Message
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

func kanalCommandHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.Debugf("kanal command received")
	bot.SendStringMessage(text.MsgKanalLink, userSession.ChatID)
}

func feedbackCommandHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.Debugf("feedback command received, user: %v", *userSession)
	bot.AskStringQuestion(text.MsgFeedback, userSession.UserID, userSession.ChatID, feedbackMessageHandler)
}

func feedbackMessageHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	adminChatID := configuration.GetInstance().GetInt64("admin-chatid")
	feedbackMsg := update.(*telegramAPI.Update).Message
	feedbackForward := telegramAPI.NewForward(adminChatID, feedbackMsg.Chat.ID, feedbackMsg.MessageID)
	go bot.Send(feedbackForward)

	successfulSent := telegramAPI.NewMessage(userSession.ChatID, text.MsgFeedbackSent)
	bot.Send(successfulSent)
}

func helpCommandHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.Debugf("help command received")
	bot.SendStringMessage(text.MsgHelp, userSession.ChatID)
}
