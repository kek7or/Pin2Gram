package dbtypes

type Channel struct {
	ChannelId int64
	Boards    []string
	Cron      string
}
