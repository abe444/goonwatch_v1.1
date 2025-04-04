package handlers

import (
	"fmt"
	"goonwatch/database"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isModerator(s *discordgo.Session, guildID, userID string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}
		if role.Permissions&discordgo.PermissionManageRoles != 0 || role.Permissions&discordgo.PermissionAdministrator != 0 {
			return true
		}
	}

	return false
}

func HandleSetDaysCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m.GuildID, m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !setdays <user> <days>")
		return
	}

	userMention := args[1]
	days := args[2]

	userID := strings.TrimPrefix(strings.TrimSuffix(userMention, ">"), "<@")
	if userID == "" {
		s.ChannelMessageSend(m.ChannelID, "Invalid user mention.")
		return
	}

	daysInt := 0
	_, err := fmt.Sscan(days, &daysInt)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of days.")
		return
	}

	err = database.UpdateStreak(userID, daysInt)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error updating streak.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set <@%s>'s streak to %d days.", userID, daysInt))
}

func HandleDeleteDaysCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m.GuildID, m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !deletedays <user> <days>")
		return
	}

	userMention := args[1]
	days := args[2]

	userID := strings.TrimPrefix(strings.TrimSuffix(userMention, ">"), "<@")
	if userID == "" {
		s.ChannelMessageSend(m.ChannelID, "Invalid user mention.")
		return
	}

	daysInt := 0
	_, err := fmt.Sscan(days, &daysInt)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid number of days.")
		return
	}

	currentStreak, err := database.GetUserStreak(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching user streak.")
		return
	}

	newStreak := currentStreak - daysInt
	if newStreak < 0 {
		newStreak = 0
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

	err = database.UpdateStreak(userID, newStreak)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error updating streak.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Subtracted %d days from <@%s>'s streak. New streak: %d days.", daysInt, userID, newStreak))
}
