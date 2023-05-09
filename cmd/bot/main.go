package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	expirationTime = 5 * time.Minute
)

type UserStorage struct {
	values map[int64][2]string
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		values: make(map[int64][2]string),
	}
}

func (us *UserStorage) Set(userID int64, args string) error {
	values := strings.Split(args, " ")
	if len(values) != 2 {
		return fmt.Errorf("write login and password please :)")
	}
	us.values[userID] = [2]string{values[0], values[1]}
	time.AfterFunc(expirationTime, func() {
		delete(us.values, userID)
	})
	return nil
}

func (us *UserStorage) Get(userID int64) ([2]string, bool) {
	value, ok := us.values[userID]
	return value, ok
}

func (us *UserStorage) Delete(userID int64) {
	delete(us.values, userID)
}

func main() {
	bot, err := tgbotapi.NewBotAPI("6023164955:AAHgDhUulJMkdlpDxKQzxrvn0c_ZNNQxDHM")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	userStorage := NewUserStorage()

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			userID := update.Message.From.ID

			switch update.Message.Command() {
			case "set":
				args := update.Message.CommandArguments()
				err := userStorage.Set(userID, args)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Login and password set")
					bot.Send(msg)
				}
			case "get":
				value, ok := userStorage.Get(userID)
				if ! ok {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Login and password not found")
					bot.Send(msg)
				}else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Login: "+value[0]+"\nPassword: "+value[1])
					bot.Send(msg)
				}
			case "del":
				userStorage.Delete(userID)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Paswword deleted")
				bot.Send(msg)

			case "list":
				var sb strings.Builder
				for userID, value := range userStorage.values {
					 sb.WriteString(strconv.FormatInt(userID, 10))
					 sb.WriteString(": ")
					 sb.WriteString(value[1])
					 sb.WriteString("\n")
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, sb.String())
				bot.Send(msg)
				
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
				bot.Send(msg)
			}
		}
	}
}

