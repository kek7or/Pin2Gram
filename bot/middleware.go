package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *PinBot) AdminsOnly(view ViewFunc) ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		b.Mutex.Lock()
		for _, admin := range b.admins {
			if admin == update.SentFrom().ID {
				return view(ctx, bot, update)
			}
		}
		b.Mutex.Unlock()

		_, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "NO_ACCESS"))

		return err
	}
}
