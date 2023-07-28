package main

import (
	"context"
	"log"

	"pinterest-tg-autopost/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	config := bot.NewConfig()

	api, err := tgbotapi.NewBotAPI(config.BotApiKey)
	if err != nil {
		log.Panic(err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Panic(err)
	}

	db := mongoClient.Database(bot.DATABASE_NAME)

	pinbot := bot.NewPinBot(api, db)

	pinbot.RegisterCmdView("addChannel", pinbot.ViewCmdAddChannel(), false)
	pinbot.RegisterCmdView("removeChannel", pinbot.ViewCmdRemoveChannels(), false)
	pinbot.RegisterCmdView("channels", pinbot.ViewCmdChannels(), false)
	pinbot.RegisterCmdView("addBoards", pinbot.ViewCmdAddBoards(), false)
	pinbot.RegisterCmdView("removeBoards", pinbot.ViewCmdRemoveBoards(), false)

	pinbot.Run(context.TODO())
}
