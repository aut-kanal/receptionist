package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/receptionist/configuration"
	"gitlab.com/kanalbot/receptionist/mq"
	"gitlab.com/kanalbot/receptionist/ui/text"
)

var (
	bot *miyanbor.Bot
)

func StartBot() {
	botDebug := configuration.GetInstance().GetBool("bot.telegram.debug")
	botToken := configuration.GetInstance().GetString("bot.telegram.token")
	botSessionTimeout := configuration.GetInstance().GetInt("bot.telegram.session-timeout")
	botUpdaterTimeout := configuration.GetInstance().GetInt("bot.telegram.updater-timeout")

	var err error
	bot, err = miyanbor.NewBot(botToken, botDebug, botSessionTimeout)
	if err != nil {
		logrus.WithError(err).Fatalf("can't init bot")
	}
	logrus.Infof("telegram bot initialized completely")

	mq.SubscribeAcceptedMsgs(acceptedMessageHandler)
	mq.SubscribeRejectedMsgs(rejectedMessageHandler)

	logrus.Infof("===================================")

	setCallbacks(bot)
	bot.StartUpdater(0, botUpdaterTimeout)
}

func setCallbacks(bot *miyanbor.Bot) {
	bot.SetSessionStartCallbackHandler(sessionStartHandler)
	bot.SetFallbackCallbackHandler(unknownMessageHandler)

	bot.AddCommandHandler(text.CancelCommandRegex, cancelCommandHandler)
	bot.AddCommandHandler(text.NewMessageCommandRegex, newMessageCommandHandler)
	bot.AddCommandHandler(text.KanalCommandRegex, kanalCommandHandler)
	bot.AddCommandHandler(text.FeedbackCommandRegex, feedbackCommandHandler)
	bot.AddCommandHandler(text.HelpCommandRegex, helpCommandHandler)
}
