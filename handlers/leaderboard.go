package handlers

import (
	"fmt"
	"goonwatch/database"

	"github.com/bwmarrin/discordgo"
)

func HandleLeaderboardCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	users, err := database.GetLeaderboard()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching leaderboard.")
		return
	}

	var leaderboard string
	for _, user := range users {
		member, err := s.GuildMember(m.GuildID, user.ID)
		if err != nil {
			fmt.Println("Error fetching member:", err)
			continue
		}
		leaderboard += fmt.Sprintf("%s: %d days\n", member.User.Username, user.Streak)
	}

	s.ChannelMessageSend(m.ChannelID, "**Leaderboard**\n"+leaderboard)
}

func HandleStatusCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	streak, err := database.GetUserStreak(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching your status.")
		return
	}

	var status string
	if streak == 0 {
		status = "relapsed"
	} else {
		status = fmt.Sprintf("%d days goon free!", streak)
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, your status: %s", userID, status))
}

func HandleHallOfGoonersCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	db := database.GetDB()
	rows, err := db.Query(`
		SELECT h.user_id, h.last_streak, h.best_streak, h.relapse_count, u.streak 
		FROM hall_of_gooners h 
		JOIN users u ON h.user_id = u.id 
		ORDER BY h.relapse_count DESC 
		LIMIT 10
	`)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching Hall of Gooners.")
		return
	}
	defer rows.Close()

	message := "🤡 **HALL OF GOONERS** 🤡\n"
	for rows.Next() {
		var userID string
		var lastStreak, bestStreak, relapseCount, currentStreak int
		err := rows.Scan(&userID, &lastStreak, &bestStreak, &relapseCount, &currentStreak)
		if err != nil {
			continue
		}

		member, err := s.GuildMember(m.GuildID, userID)
		if err != nil {
			fmt.Println("Error fetching member:", err)
			continue
		}

		message += fmt.Sprintf("%s - %d relapses | Best: %d days | Current: %d days\n",
			member.User.Username, relapseCount, bestStreak, currentStreak)
	}

	s.ChannelMessageSend(m.ChannelID, message)
}
