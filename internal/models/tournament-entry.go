package models

type TournamentEntry struct {
	Id           int
	TournamentId int
	UserId       string
	Username     string
	Waves        int
	DisplayName  string
	GuildId      string
}
