package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"pinterest-tg-autopost/pinterest"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *PinBot) ViewCmdAddChannel() ViewFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args := strings.Split(update.Message.CommandArguments(), " ")

		if !AreAddChannelArgsValid(args) {
			b.AnswerMsg(update, "unable to parse arguments")
			return nil
		}

		channelId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			b.AnswerMsg(update, "failed to convert channelId to integer")
			return nil
		}

		coll := b.db.Collection(CHANNELS_COLLECTION)

		res := coll.FindOne(ctx, bson.M{"channelId": channelId})
		if res.Err() != mongo.ErrNoDocuments {
			b.AnswerMsg(update, "'%d' already exists in database", channelId)
			return nil
		}

		boards := args[1:]
		boards, err = handleBoards(boards)
		if err != nil {
			b.AnswerMsg(update, err.Error())
		}

		if hasPostAccess := b.hasPostAccess(channelId); !hasPostAccess {
			b.AnswerMsg(update, "no access to post messages in '%d' channel", channelId)
			return nil
		}

		_, err = coll.InsertOne(ctx, bson.M{"channelId": channelId, "cron": nil, "boards": boards})
		if err != nil {
			b.AnswerMsg(update, "failed to add '%d' channel to database", channelId)
		}

		b.AnswerMsg(update, "channel has been successfully added to database")

		return nil
	}
}

func handleBoards(boards []string) ([]string, error) {
	validBoards := 0

	boards = TrimElements(boards, "/")
	boards = RemoveDuplicate(boards)

	for _, board := range boards {
		uri := fmt.Sprintf(PIDGETS_URI, board)
		resp, err := http.Get(uri)
		if err != nil {
			return nil, fmt.Errorf("failed to access '%s' board", board)
		}
		defer resp.Body.Close()

		var decoded pinterest.GetApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
			return nil, fmt.Errorf("failed to decode '%s' board info to json", board)
		}

		if decoded.Status != "success" {
			return nil, fmt.Errorf("failed to get info about '%s' board", board)
		}

		validBoards++
	}
	if validBoards == 0 {
		return nil, fmt.Errorf("zero valid boards")
	}

	return boards, nil
}
