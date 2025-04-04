package main

import (
	"fmt"
	"goonwatch/database"
	"goonwatch/handlers"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	database.InitDB()

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set in the .env file")
	}

	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer discord.Close()

	fmt.Println("Bot is running! Press CTRL+C to exit.")
	<-make(chan struct{})
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, "!goon"):
		handlers.HandleLogCommand(s, m)
	case m.Content == "!leaderboard":
		handlers.HandleLeaderboardCommand(s, m)
	case m.Content == "!status":
		handlers.HandleStatusCommand(s, m)
	case m.Content == "!rank":
		handlers.HandleRankCommand(s, m)
	case m.Content == "!ranks":
		handlers.HandleRanksCommand(s, m)
	case m.Content == "!history":
		handlers.HandleHistoryCommand(s, m)
	case strings.HasPrefix(m.Content, "!setdays"):
		handlers.HandleSetDaysCommand(s, m)
	case strings.HasPrefix(m.Content, "!deletedays"):
		handlers.HandleDeleteDaysCommand(s, m)
	case m.Content == "!hallofgooners":
		handlers.HandleHallOfGoonersCommand(s, m)
	case m.Content == "!help":
		handlers.HandleHelpCommand(s, m)
	case m.Content == "!rules":
		handlers.HandleRulesCommand(s, m)
	case m.Content == "!test":
		s.ChannelMessageSend(m.ChannelID, "Systems are operational. GoonWatchBot online.")
	}
}
