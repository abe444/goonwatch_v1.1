package handlers

import (
	"fmt"
	"goonwatch/database"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func assignRankRole(s *discordgo.Session, guildID, userID, rank string) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		fmt.Println("Error fetching roles:", err)
		return
	}

	var roleID string
	for _, role := range roles {
		if role.Name == rank {
			roleID = role.ID
			break
		}
	}

	if roleID != "" {
		err := s.GuildMemberRoleAdd(guildID, userID, roleID)
		if err != nil {
			fmt.Println("Error assigning role:", err)
		}
	}
}

func HandleLogCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !goon <yes/no>")
		return
	}

	userID := m.Author.ID
	response := strings.ToLower(args[1])

	if response == "no" {
		canLog, err := database.CanLogAgain(userID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error checking log cooldown.")
			return
		}

		if !canLog {
			timeRemaining, err := database.GetTimeUntilNextLog(userID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error checking cooldown time.")
				return
			}

			hours := int(timeRemaining.Hours())
			minutes := int(timeRemaining.Minutes()) % 60
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"You need to wait %d hours and %d minutes before logging your next successful day.",
				hours, minutes))
			return
		}
	}

	currentStreak, err := database.GetUserStreak(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching your streak.")
		return
	}

	var newStreak int
	if response == "no" {
		newStreak = currentStreak + 1

		_, err = database.GetDB().Exec(`
			INSERT INTO users (id, streak) 
			VALUES ($1, $2) 
			ON CONFLICT (id) DO UPDATE 
			SET streak = $2
		`, userID, newStreak)
		if err != nil {
			fmt.Println("Error ensuring user exists:", err)
			s.ChannelMessageSend(m.ChannelID, "Error updating streak.")
			return
		}

		err = database.UpdateStreak(userID, newStreak)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error updating streak.")
			return
		}

		err = database.AddStreakHistory(userID, newStreak)
		if err != nil {
			fmt.Println("Error recording streak history:", err)
		}

		goonerStatus, err := database.GetGoonerStatus(userID)
		if err != nil {
			fmt.Println("Error checking gooner status:", err)
		}
		if goonerStatus != nil && newStreak > goonerStatus.LastStreak {
			err = database.RemoveFromHallOfGooners(userID)
			if err != nil {
				fmt.Println("Error removing from hall of gooners:", err)
			}
			s.ChannelMessageSend(m.ChannelID, "Congratulations! You've surpassed your previous streak and escaped the Hall of Gooners!")
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Streak updated! Current streak: %d days", newStreak))
	} else if response == "yes" {
		newStreak = 0

		_, err = database.GetDB().Exec(`
			INSERT INTO users (id, streak) 
			VALUES ($1, $2) 
			ON CONFLICT (id) DO UPDATE 
			SET streak = $2
		`, userID, newStreak)
		if err != nil {
			fmt.Println("Error ensuring user exists:", err)
			s.ChannelMessageSend(m.ChannelID, "Error updating streak.")
			return
		}

		err = database.AddToHallOfGooners(userID, currentStreak)
		if err != nil {
			fmt.Println("Error adding to hall of gooners:", err)
		}

		err = database.UpdateStreak(userID, newStreak)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error updating streak.")
			return
		}

		err = database.AddStreakHistory(userID, newStreak)
		if err != nil {
			fmt.Println("Error recording streak history:", err)
		}

		s.ChannelMessageSend(m.ChannelID, "You GOONED. Relapse count added.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Invalid response. Usage: !goon <yes/no>")
		return
	}

	currentRank := getRank(currentStreak)
	newRank := getRank(newStreak)
	if newRank != currentRank {
		err := database.UpdateUserRank(userID, newRank)
		if err != nil {
			fmt.Println("Error updating user rank:", err)
		}

		err = database.AddRankHistory(userID, newRank)
		if err != nil {
			fmt.Println("Error adding rank history:", err)
		}

		if newRank != "unranked" {
			assignRankRole(s, m.GuildID, userID, newRank)

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, congratulations! You achieved the **%s** rank!", userID, newRank))
		}
	}

	err = database.UpdateLastLogTime(userID)
	if err != nil {
		fmt.Println("Error updating last log time:", err)
	}
}
