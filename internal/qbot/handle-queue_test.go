package qbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_QBot_formatPlayersAhead(t *testing.T) {
	makeQueue := func(userIDs ...string) []QueueItem {
		items := make([]QueueItem, len(userIDs))
		for i, id := range userIDs {
			items[i] = QueueItem{UserID: id}
		}
		return items
	}

	t.Run("position 0 returns empty string", func(t *testing.T) {
		q := &QBot{queue: makeQueue("user1", "user2")}
		require.Equal(t, "", q.formatPlayersAhead(0))
	})

	t.Run("position 1 returns first player", func(t *testing.T) {
		q := &QBot{queue: makeQueue("user1", "user2", "user3")}
		require.Equal(t, "<@user1>", q.formatPlayersAhead(1))
	})

	t.Run("position 3 returns three players", func(t *testing.T) {
		q := &QBot{queue: makeQueue("user1", "user2", "user3", "user4")}
		require.Equal(t, "<@user1>, <@user2>, <@user3>", q.formatPlayersAhead(3))
	})

	t.Run("caps display at 10 and shows overflow count", func(t *testing.T) {
		// 15 users in queue, position 14 → position capped at 10, overflow = 15-10 = 5
		ids := make([]string, 15)
		for i := range ids {
			ids[i] = "user" + string(rune('A'+i))
		}
		q := &QBot{queue: makeQueue(ids...)}
		result := q.formatPlayersAhead(14)
		require.Contains(t, result, "_and 5 more..._")
	})

	t.Run("exactly 10 in queue shows no overflow", func(t *testing.T) {
		// 10 users in queue, len(q.queue) == limit so overflow condition is false
		ids := make([]string, 10)
		for i := range ids {
			ids[i] = "user" + string(rune('A'+i))
		}
		q := &QBot{queue: makeQueue(ids...)}
		result := q.formatPlayersAhead(10)
		require.NotContains(t, result, "more")
	})
}
