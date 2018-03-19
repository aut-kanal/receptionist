package miyanbor

import (
	cache "github.com/patrickmn/go-cache"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

// AskStringQuestion sends a question to a Chat and sets a callback to retrieve
// user's answer
func (b *Bot) AskStringQuestion(question string, userID int, chatID int64, callback CallbackFunction) {
	// Send question
	questionMsg := telegramAPI.NewMessage(chatID, question)
	b.Send(questionMsg)

	// Set callback function
	var userSession *UserSession
	userSessionInterface, found := usersSessionCache.Get(getUserSessionKey(userID, chatID))
	if !found {
		userSession = &UserSession{
			chatID:  chatID,
			userID:  userID,
			payload: make(map[string]interface{}),
		}
		usersSessionCache.Add(getUserSessionKey(userID, chatID), userSession, cache.DefaultExpiration)
	} else {
		userSession = userSessionInterface.(*UserSession)
	}
	userSession.messageCallback = func(userSession *UserSession, input interface{}) {
		userSession.messageCallback = nil
		callback(userSession, input)
	}
}

// SendStringMessage sends a string to a Chat
func (b *Bot) SendStringMessage(msg string, chatID int64) {
	telegramMsg := telegramAPI.NewMessage(chatID, msg)
	b.Send(telegramMsg)
}
