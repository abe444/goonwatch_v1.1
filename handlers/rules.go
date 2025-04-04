package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleRulesCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "📜 GoonWatchBot Rules 📜",
		Description: "Just don't fucking goon bruh. \n",
		Color:       0x00ff00, // Green color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "24 Hour cooldown",
				Value:  "You must wait 1 day before logging a successful day after a relapse.",
				Inline: false,
			},
			{
				Name:   "Wet Dreams, sex, succubi attacking DO NOT COUNT..maybe. ",
				Value:  "Porn is out the window, but what you additionally abstain to is your business.",
				Inline: false,
			},
			{
				Name:   "Hall Of Gooners",
				Value:  "Relapse count stays with you forever.\n A successful day logged will kick you out of the Hall of Gooners.",
				Inline: false,
			},
		},
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println("Error sending rules embed:", err)
	}
}
