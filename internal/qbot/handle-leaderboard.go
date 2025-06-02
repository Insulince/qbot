package qbot

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

// Handle !leaderboard
func (q *QBot) handleLeaderboard(cmd Cmd, final bool) error {
	g := q.guilds.MustGet(cmd.GuildId)
	channelId := g.AnnouncementChannelId
	if cmd.Message != nil {
		channelId = cmd.Message.ChannelID
	}

	tournamentEntries, err := q.store.GetLatestTournamentEntries(cmd.GuildId)
	if err != nil {
		return errors.Wrap(err, "getting latest tournament entries")
	}

	leaderboardMsg := fmt.Sprintf("🏆 **Latest Tournament Leaderboard** 🏆\n")
	var entries []string
	var lastPlaceDisplayName string
	for i, tournamentEntry := range tournamentEntries {
		entry := fmt.Sprintf("%d. **<@%s>** - Wave %d", i, tournamentEntry.UserId, tournamentEntry.Waves)
		entries = append(entries, entry)

		lastPlaceDisplayName = tournamentEntry.DisplayName // Keep track of the last user
	}

	entriesMsg := "_No entries yet._"
	if len(entries) > 0 {
		entriesMsg = strings.Join(entries, "\n")
	}
	leaderboardMsg += entriesMsg

	// Send the leaderboard message first
	q.mustPostWithoutTags(channelId, leaderboardMsg)

	if final {
		q.congratulateLoser(channelId, lastPlaceDisplayName)
	}

	return nil
}

// Send a local image file with the username added directly to the image
func (q *QBot) congratulateLoser(channelID, lastPlaceDisplayName string) {
	const templatePath = "/app/assets/celebrate.png"
	const outputPath = "/tmp/last_place_meme.png"

	// Open the template image
	templateFile, err := os.Open(templatePath)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Could not load template image.")
		return
	}
	defer jmust.MustClose(templateFile)

	// Decode the PNG image
	img, err := png.Decode(templateFile)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Could not decode template image.")
		return
	}

	// Create a new RGBA image
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Load font
	fontBytes, err := os.ReadFile("/app/assets/impact.ttf")
	if err != nil {
		// Fallback to system font if custom font not available
		fontBytes, err = os.ReadFile("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf")
		if err != nil {
			q.mustPost(channelID, "❌ Error: Could not load font.")
			return
		}
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Could not parse font.")
		return
	}

	// Calculate text width to right-align it
	fontSize := 48.0
	opts := truetype.Options{
		Size: fontSize,
		DPI:  72,
	}
	face := truetype.NewFace(f, &opts)

	// Get text width for right alignment
	text := lastPlaceDisplayName
	width := font.MeasureString(face, text).Ceil()

	// Position the text in the lower right corner
	padding := 20 // Padding from the right edge
	x := bounds.Max.X - width - padding
	y := bounds.Max.Y - 80 // Position near the bottom

	// Draw text with black outline/shadow effect
	// First draw the black outline/shadow by drawing the text multiple times with slight offsets
	outlineC := freetype.NewContext()
	outlineC.SetDPI(72)
	outlineC.SetFont(f)
	outlineC.SetFontSize(fontSize)
	outlineC.SetClip(bounds)
	outlineC.SetDst(rgba)
	outlineC.SetSrc(image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 255})) // Black outline

	// Draw the outline by drawing the text multiple times with slight offsets
	offsets := []struct{ dx, dy int }{{-2, 0}, {2, 0}, {0, -2}, {0, 2}, {-2, -2}, {2, -2}, {-2, 2}, {2, 2}}
	for _, offset := range offsets {
		pt := freetype.Pt(x+offset.dx, y+offset.dy)
		if _, err = outlineC.DrawString(text, pt); err != nil {
			q.mustPost(channelID, "❌ Error: Could not draw text outline on image.")
			return
		}
	}

	// Now draw the white text on top
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(bounds)
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255})) // White text

	// Draw the main text
	pt := freetype.Pt(x, y)
	if _, err = c.DrawString(text, pt); err != nil {
		q.mustPost(channelID, "❌ Error: Could not draw text on image.")
		return
	}

	// Save the modified image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Could not create output file.")
		return
	}

	if err := png.Encode(outputFile, rgba); err != nil {
		q.mustPost(channelID, "❌ Error: Could not encode output image.")
		return
	}

	// Reopen the file for reading
	jmust.MustClose(outputFile)
	file, err := os.Open(outputPath)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Could not open modified image.")
		return
	}
	defer jmust.MustClose(file)

	// Create a message with the modified image
	message := &discordgo.MessageSend{
		Content: fmt.Sprintf("😂 **Congrats %s!** 😂", lastPlaceDisplayName),
		Files: []*discordgo.File{
			{
				Name:   "last_place_meme.png",
				Reader: file,
			},
		},
	}

	// Send the message with the image
	_, err = q.session.ChannelMessageSendComplex(channelID, message)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Failed to send image.")
	}

	// Clean up the temporary file
	if err := os.Remove(outputPath); err != nil {
		q.mustPost(channelID, fmt.Sprintf("❌ Error: Failed to delete temporary image from %q.", outputPath))
		return
	}
}
