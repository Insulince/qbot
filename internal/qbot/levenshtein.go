package qbot

// knownCommands is the canonical list used for typo suggestions.
// Keep this in sync with the switch statement in commands.go.
var knownCommands = []string{
	"!clearwave", "!commands", "!dev", "!enter", "!forcewave",
	"!full", "!help", "!history", "!leaderboard", "!leave",
	"!moretime", "!owned", "!position", "!progress", "!queue",
	"!remove", "!reset", "!shame", "!skip", "!submitwave",
	"!version", "!view",
}

// suggestCommand returns the closest known command to the query if it is within
// the suggestion threshold, or an empty string if nothing is close enough.
func suggestCommand(query string) string {
	const maxDist = 3

	best := ""
	bestDist := maxDist + 1

	for _, cmd := range knownCommands {
		d := levenshtein(query, cmd)
		if d < bestDist {
			bestDist = d
			best = cmd
		}
	}

	// Also require the distance to be less than half the query length to avoid
	// suggesting commands for very short unrelated inputs (e.g. "!zz" → "!q").
	if bestDist > len([]rune(query))/2 {
		return ""
	}

	return best
}

// levenshtein returns the edit distance between strings a and b.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)

	dp := make([]int, lb+1)
	for j := range dp {
		dp[j] = j
	}

	for i := 1; i <= la; i++ {
		prev := dp[0]
		dp[0] = i
		for j := 1; j <= lb; j++ {
			tmp := dp[j]
			if ra[i-1] == rb[j-1] {
				dp[j] = prev
			} else {
				dp[j] = 1 + min(prev, dp[j], dp[j-1])
			}
			prev = tmp
		}
	}

	return dp[lb]
}
