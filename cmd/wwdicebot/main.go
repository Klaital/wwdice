package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/klaital/wwdice/pkg/characters"
	"github.com/klaital/wwdice/pkg/dice"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var token string
var buffer = make([][]byte, 0)

func main() {
	// Read bot token from env var first, but allow a commandline flag to override it
	token = os.Getenv("DISCORD_BOT_TOKEN")
	if len(token) == 0 {
		flag.StringVar(&token, "t", "", "Discord bot token")
		flag.Parse()
	}

	if token == "" {
		fmt.Println("No token provided. Please run: wwdicebot -t <bot token>")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	rand.Seed(time.Now().UnixNano())

	// Register readyHandler as a callback for the readyHandler events.
	dg.AddHandler(readyHandler)

	// Register messageCreateHandler as a callback for the messageCreateHandler events.
	dg.AddHandler(messageCreateHandler)

	// We need information about messages
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		panic(err)
	}

	// Load an HTTP Handler to respond to k8s healthchecks and also roll dice if someone wants
	http.HandleFunc("/healthz", healthCheckHandler)
	http.HandleFunc("/roll", rollDiceHandler)
	http.ListenAndServe(":8080", nil)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("WW Dice Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusOK)
}

func rollDiceHandler(w http.ResponseWriter, r *http.Request)  {
	diceCfg, err := dice.ParseDiceString("3d6")
	if err != nil {
		log.WithError(err).Error("Failed to parse string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rolls := diceCfg.RollDice()
	results := diceCfg.CalculateResults(rolls)
	if results.Botch {
		_, err = fmt.Fprintf(w, "BOTCH %+v", rolls)
		if err != nil {
			log.WithError(err).Error("Failed to write botch results")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		_, err = fmt.Fprintf(w, "%d successes %v", results.Successes, rolls)
		if err != nil {
			log.WithError(err).Error("Failed to write success results")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// This function will be called (due to AddHandler above) when the bot receives
// the "readyHandler" event from Discord.
func readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateStatus(0, "!roll")
}


// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if the message is "!roll"
	if strings.HasPrefix(m.Content, "!roll") {
		// Parse the dice command
		diceCfg, err := dice.ParseDiceString(strings.Replace(m.Content, "!roll", "", 1))
		if err != nil {
			log.WithError(err).WithField("cmd", m.Content).Error("Failed to parse dice string")
			return
		}
		// Roll the dice
		rolls := diceCfg.RollDice()
		results := diceCfg.CalculateResults(rolls)
		if results.Botch {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("BOTCH! %v", rolls))
			if err != nil {
				log.WithError(err).Error("Failed to send message")
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d successes %v", results.Successes, rolls))
			if err != nil {
				log.WithError(err).Error("Failed to send message")
			}
		}
	}
}


func registerCharacter(id string, c characters.Character) error {
	return errors.New("not implemented yet")
}
func loadPlayerCharacter(id string) (characters.Character, error) {
	return nil, errors.New("not implemented yet")
}
func loadCharacter() (characters.Character, error) {
	return nil, errors.New("not implemented yet")
}
