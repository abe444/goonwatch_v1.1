package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func HandleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpEmbed := &discordgo.MessageEmbed{
		Title: "GoonWatch Bot Commands",
		Color: 0x00ff00, // Green color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Daily Logging",
				Value: "`!goon no` - Log a successful day (24h cooldown)\n" +
					"`!goon yes` - Log a relapse (no cooldown)\n" +
					"`!status` - Check your current streak and status",
				Inline: false,
			},
			{
				Name: "Progress & Rankings",
				Value: "`!rank` - Check your current rank\n" +
					"`!ranks` - View all available ranks and requirements\n" +
					"`!history` - View your streak history",
				Inline: false,
			},
			{
				Name: "Leaderboards",
				Value: "`!leaderboard` - View the top streaks\n" +
					"`!hallofgooners` - View the Hall of Gooners (relapse leaderboard)",
				Inline: false,
			},
			{
				Name: "Moderator Commands",
				Value: "`!setdays <@user> <days>` - Set a user's streak days\n" +
					"`!deletedays <@user> <days>` - Remove days from a user's streak",
				Inline: false,
			},
			{
				Name: "Other",
				Value: "`!help` - Show this help message\n" +
					"`!test` - Check if the bot is online",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "NEVER GOON",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, helpEmbed)
}
