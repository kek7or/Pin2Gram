package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"pinterest-tg-autopost/dbtypes"
	"pinterest-tg-autopost/pinterest"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *PinBot) ViewCmdAddBoards() ViewFunc {
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

		var channel dbtypes.Channel
		filter := bson.M{"channelId": channelId}

		res := coll.FindOne(ctx, filter)
		if res.Err() == mongo.ErrNoDocuments {
			b.AnswerMsg(update, "'%d' does not exist in database", channelId)
			return nil
		}

		if err := res.Decode(&channel); err != nil {
			b.AnswerMsg(update, "failed to decode channel")
			return nil
		}

		toInsert := make([]interface{}, 0)

		userBoards := args[1:]
		userBoards, err = handleBoards(userBoards)
		if err != nil {
			b.AnswerMsg(update, err.Error())
		}

		for _, userBoard := range userBoards {
			for _, dboard := range channel.Boards {
				if userBoard == dboard {
					b.AnswerMsg(update, "board '%s' already exists in database", userBoard)
				} else {
					toInsert = append(toInsert, userBoard)
				}
			}
		}

		if len(toInsert) == 0 {
			b.AnswerMsg(update, "no boards to add")
			return nil
		}

		updatedBoards := bson.D{{"$push", bson.D{{"boards", bson.D{{"$each", toInsert}}}}}}
		if _, err := coll.UpdateOne(ctx, filter, updatedBoards); err != nil {
			b.AnswerMsg(update, "failed to update document")
			return nil
		}

		b.AnswerMsg(update, "boards has been successfully updated")

		return nil
	}
}

func (b *PinBot) ViewCmdRemoveBoards() ViewFunc {
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

		var channel dbtypes.Channel
		filter := bson.M{"channelId": channelId}

		res := coll.FindOne(ctx, filter)
		if res.Err() == mongo.ErrNoDocuments {
			b.AnswerMsg(update, "'%d' does not exist in database", channelId)
			return nil
		}

		if err := res.Decode(&channel); err != nil {
			b.AnswerMsg(update, "failed to decode channel")
			return nil
		}

		toDelete := make([]interface{}, 0)

		userBoards := args[1:]
		userBoards, err = handleBoards(userBoards)
		if err != nil {
			b.AnswerMsg(update, err.Error())
		}

		for _, userBoard := range userBoards {
			for _, dboard := range channel.Boards {
				if userBoard == dboard {
					toDelete = append(toDelete, userBoard)
				}
			}
		}

		if len(toDelete) == 0 {
			b.AnswerMsg(update, "no boards to delete")
			return nil
		}

		updatedBoards := bson.D{{"$pull", bson.D{{"boards", bson.D{{"$in", toDelete}}}}}}
		if _, err := coll.UpdateOne(ctx, filter, updatedBoards); err != nil {
			b.AnswerMsg(update, "failed to update document")
			return nil
		}

		b.AnswerMsg(update, "boards has been successfully updated")

		return nil
	}
}

func (b *PinBot) ViewCmdChannels() ViewFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		coll := b.db.Collection(CHANNELS_COLLECTION)

		cur, err := coll.Find(ctx, bson.M{})
		if err != nil {
			b.AnswerMsg(update, "failed to get channels from database")
			return nil
		}

		var channels []dbtypes.Channel
		if err := cur.All(ctx, &channels); err != nil {
			b.AnswerMsg(update, "failed to get channels from database")
			return nil
		}

		text := ""

		for _, c := range channels {
			text += fmt.Sprintf("%v\n", c)
		}

		b.AnswerMsg(update, text)

		return nil
	}
}

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
