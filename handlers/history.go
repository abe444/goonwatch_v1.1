package handlers

import (
	"fmt"
	"goonwatch/database"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleHistoryCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	history, err := database.GetStreakHistory(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching your streak history.")
		return
	}

	var historyMessage string
	if len(history) == 0 {
		historyMessage = "You have no streak history yet."
	} else {
		historyMessage = "**Your Streak History:**\n"
		for _, entry := range history {
			streak := entry["streak"].(int)
			recordedAt := entry["recordedAt"].(time.Time).Format("2006-01-02 15:04:05")
			historyMessage += fmt.Sprintf("- **%d days** on %s\n", streak, recordedAt)
		}
	}

	s.ChannelMessageSend(m.ChannelID, historyMessage)
}
