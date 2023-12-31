package bot

import (
	"context"
	"fmt"
	"math/rand"
	"pinterest-tg-autopost/dbtypes"
	"pinterest-tg-autopost/pinterest"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *PinBot) runAutopost(ctx context.Context) error {
	channelsColl := b.db.Collection(CHANNELS_COLLECTION)
	postsColl := b.db.Collection(POSTS_COLLECTION)

	cur, err := channelsColl.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to get channels from database")
	}

	var channels []dbtypes.Channel
	if err := cur.All(ctx, &channels); err != nil {
		return fmt.Errorf("failed to get channels from database")
	}

	for _, channel := range channels {
		pins := make([]pinterest.Pin, 0)

		for _, board := range channel.Boards {
			pins_, err := pinterest.GetPinsFromBoard(board)
			if err != nil {
				return err
			}

			pins = append(pins, pins_...)
		}

		filter := bson.M{"channelId": channel.ChannelId}
		cur, err := postsColl.Find(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to find posts from database")
		}

		var posts []dbtypes.Post
		if err := cur.All(ctx, &posts); err != nil {
			return fmt.Errorf("failed to get posts from database")
		}

		pinsToPost := NonImplicationPins(pins, posts)

		if len(pinsToPost) == 0 {
			continue
		}
		randomPin := pinsToPost[rand.Int63n(int64(len(pinsToPost)))]

		photo := tgbotapi.NewPhoto(channel.ChannelId, tgbotapi.FileURL(randomPin.Images.I564x.URL))
		sentMsg, err := b.api.Send(photo)
		if err != nil {
			return fmt.Errorf("failed to send error message: %v", err)
		}

		post := bson.M{
			"channelId": channel.ChannelId,
			"msgId":     sentMsg.MessageID,
			"pinId":     randomPin.ID,
			"time":      time.Now().Unix()}

		_, err = postsColl.InsertOne(ctx, post)
		if err != nil {
			return fmt.Errorf("failed to insert post to database")
		}
	}

	return nil
}
