package dbtypes

type Channel struct {
	ChannelId int64
	Boards    []string
	Cron      string
}

type Post struct {
	ChannelId int64
	MsgId     int64
	PinId     string
	Time      int64
}
