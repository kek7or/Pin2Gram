# Pinterest Autopost Bot
Bot for Telegram that fetches Pinterest boards and posts content to a channel.

# Implemented
- Simple commands to manage bot
- Autopost service. It will post any pins from selected boards to channel with interval from `bot/consts.go`

# Setup
TODO

# Commands
- `/addChannel CHANNEL_ID BOARD1 BOARD2 ...`
- `/removeChannels CHANNEL_ID1 CHANNEL_ID2 ...`
- `/channels`
- `/addBoards CHANNEL_ID BOARD1 BOARD2 ...`
- `/removeBoards CHANNEL_ID BOARD1 BOARD2 ...`

# Important note
Bot uses an old public Pinterest endpoint that may be closed at any time. Use at your own risk. No plans to upgrade to a newer version any time soon
`v3/pidgets/boards/{user_id}/{board}/pins/`
