package handlers

import (
	"fmt"
	"goonwatch/database"
	"sort"

	"github.com/bwmarrin/discordgo"
)

// Define the ranks and their corresponding streak thresholds
var ranks = map[string]int{
	"Noob r/nofap lurker":       5,
	"Intermediate chudlord":     10,
	"Self-improvement junkie":   20,
	"Twitter biohacker":         25,
	"Nonchalant practitioner":   45,
	"Normie repellent":          55,
	"Enlightened":               60,
	"Restless Sage":             75,
	"Monk Mode":                 80,
	"Goonlord Patrol - LVL 90+": 90,
}

// getRank returns the user's rank based on their streak.
func getRank(streak int) string {
	if streak == 0 {
		return "unranked"
	}

	var highestRank string
	var highestThreshold int

	for rank, threshold := range ranks {
		if streak >= threshold && threshold >= highestThreshold {
			highestRank = rank
			highestThreshold = threshold
		}
	}

	if highestRank == "" {
		return "unranked"
	}
	return highestRank
}

func HandleRankCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	streak, err := database.GetUserStreak(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching your rank.")
		return
	}

	rank := getRank(streak)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, your rank: **%s** (%d days)", userID, rank, streak))

	assignRankRole(s, m.GuildID, userID, rank)
}

func HandleRanksCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	type rankInfo struct {
		name      string
		threshold int
	}

	ranksList := make([]rankInfo, 0, len(ranks))
	for rank, threshold := range ranks {
		ranksList = append(ranksList, rankInfo{rank, threshold})
	}

	sort.Slice(ranksList, func(i, j int) bool {
		return ranksList[i].threshold < ranksList[j].threshold
	})

	var message string
	message = "**Available Ranks (Days Required):**\n"
	for _, rank := range ranksList {
		message += fmt.Sprintf("- **%s**: %d+ days\n", rank.name, rank.threshold)
	}

	s.ChannelMessageSend(m.ChannelID, message)
}
