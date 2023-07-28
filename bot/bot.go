package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

type Config struct {
	BotApiKey string
	MongoURI  string
	Admin     int64
}

func NewConfig() Config {
	admin, err := strconv.ParseInt(os.Getenv("ADMIN"), 10, 64)
	if err != nil || admin == 0 {
		panic("ADMIN value is not valid")
	}

	return Config{
		BotApiKey: os.Getenv("BOT_API_KEY"),
		MongoURI:  os.Getenv("MONGO_URI"),
		Admin:     admin,
	}
}

type PinBot struct {
	sync.Mutex

	api *tgbotapi.BotAPI
	db  *mongo.Database

	cmdViews map[string]ViewFunc
	admins   []int64
}

func NewPinBot(api *tgbotapi.BotAPI, db *mongo.Database) *PinBot {
	return &PinBot{
		api:      api,
		db:       db,
		cmdViews: make(map[string]ViewFunc),
		admins:   make([]int64, 0),
	}
}

func (b *PinBot) Run(ctx context.Context) error {
	go b.RunFetchPost(ctx)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = UPDATE_TIMEOUT
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.handleUpdate(ctx, update)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *PinBot) RunFetchPost(ctx context.Context) error {
	ticker := time.NewTicker(20 * time.Second) // TODO

	for {
		select {
		case <-ticker.C:
			err := b.runAutopost(ctx)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (b *PinBot) AnswerMsg(update tgbotapi.Update, format string, args ...any) {
	text := fmt.Sprintf(format, args...)

	if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text)); err != nil {
		log.Printf("failed to send error message: %v", err)
	}
}

func (b *PinBot) SendMsg(id int64, format string, args ...any) {
	text := fmt.Sprintf(format, args...)

	if _, err := b.api.Send(tgbotapi.NewMessage(id, text)); err != nil {
		log.Printf("failed to send error message: %v", err)
	}
}

func (b *PinBot) RegisterCmdView(cmd string, view ViewFunc, adminOnly bool) {
	if adminOnly {
		b.cmdViews[cmd] = b.AdminsOnly(view)
	} else {
		b.cmdViews[cmd] = view
	}
}

func (b *PinBot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if (update.Message == nil || !update.Message.IsCommand()) && update.CallbackQuery == nil {
		return
	}

	var view ViewFunc

	cmd := update.Message.Command()
	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		return
	}
	view = cmdView

	if err := view(ctx, b.api, update); err != nil {
		log.Printf("failed to execute view: %v", err)

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error")); err != nil {
			log.Printf("failed to send error message: %v", err)
		}
	}
}

func (b *PinBot) hasPostAccess(channelId int64) bool {
	chatMemberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: channelId,
			UserID: b.api.Self.ID,
		},
	}
	chatMember, err := b.api.GetChatMember(chatMemberConfig)
	if err != nil {
		return false
	}

	if !chatMember.CanPostMessages {
		return false
	}
	return true
}
