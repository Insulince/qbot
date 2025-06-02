package qbot

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"

	"github.com/Insulince/qbot/internal/config"
	"github.com/Insulince/qbot/internal/models"
)

// QBot struct holds the Discord session and state.
type QBot struct {
	discordBotToken    string
	enterTimeout       time.Duration
	fullTimeout        time.Duration
	warnThreshold      time.Duration
	errorChannelId     string
	notificationRoleId string

	store  Store
	guilds Guilds

	session *discordgo.Session

	queue       []QueueItem
	queueMutex  sync.Mutex
	currentUser *QueueItem
}

type Store interface {
	InsertTournament(name, shortName string) error
	GetLatestTournament() (*models.Tournament, error)
	GetTournamentByShortName(shortName string) (*models.Tournament, error)
	ListTournaments(limit, offset int) ([]*models.Tournament, error)
	CountTournaments() (int, error)

	InsertTournamentEntry(guildId string, tournamentId int64, userId, username, displayName string, waves int) error
	GetTournamentEntries(guildId string, tournamentId int64) ([]*models.TournamentEntry, error)
	GetLatestTournamentEntries(guildId string) ([]*models.TournamentEntry, error)
	GetTournamentWinner(guildId string, tournamentId int64, maxWaves int64) (*models.TournamentEntry, error)
	GetTournamentStats(guildId string, tournamentId int64) (entrants int, maxWaves *int64, averageWaves *float64, _ error)

	io.Closer
}

// New initializes and returns a new QBot instance.
func New(cfg config.Config, s Store) (*QBot, error) {
	q := new(QBot)

	q.discordBotToken = cfg.DiscordBotToken
	q.enterTimeout = time.Duration(cfg.EnterTimeoutMinutes) * time.Minute
	q.fullTimeout = time.Duration(cfg.FullTimeoutMinutes) * time.Minute
	q.warnThreshold = time.Duration(cfg.WarnThresholdMinutes) * time.Minute
	q.errorChannelId = cfg.ErrorChannelId
	q.notificationRoleId = cfg.NotificationRoleId

	q.guilds = make(map[string]Guild, len(cfg.Guilds))
	for id, cg := range cfg.Guilds {
		var g Guild
		g.Name = cg.Name
		g.AnnouncementChannelId = cg.AnnouncementChannelId
		q.guilds[id] = g
	}

	q.store = s

	session, err := q.newSession()
	if err != nil {
		return nil, errors.Wrap(err, "new session")
	}
	q.session = session

	return q, nil
}

func MustNew(cfg config.Config, s Store) *QBot {
	return jmust.Must[*QBot](New, cfg, s)[0]
}

func (q *QBot) Run(ctx context.Context) chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		if err := q.run(ctx); err != nil {
			errs <- errors.Wrap(err, "run")
		}
	}()

	return errs
}

func (q *QBot) run(ctx context.Context) error {
	if err := q.session.Open(); err != nil {
		return errors.Wrap(err, "opening session")
	}

	defer func() {
		if v := recover(); v != nil {
			q.reportError(fmt.Errorf("panic: %v", v))
			panic(v)
		}
	}()

	q.Go(q.timeoutChecker)
	q.Go(q.startScheduler)

	fmt.Println("Q is running")

	select {
	case <-ctx.Done():
	}

	return nil
}

func (q *QBot) newSession() (*discordgo.Session, error) {
	auth := fmt.Sprintf("Bot %s", q.discordBotToken)
	session, err := discordgo.New(auth)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Discord session")
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	session.AddHandler(q.messageHandler)

	return session, nil
}

// Close shuts down the bot.
func (q *QBot) Close() error {
	if err := q.session.Close(); err != nil {
		return errors.Wrap(err, "session close")
	}

	if err := q.store.Close(); err != nil {
		return errors.Wrap(err, "store close")
	}

	return nil
}

func (q *QBot) MustClose() {
	jmust.Must[any](q.Close)
}
