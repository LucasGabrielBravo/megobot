package megobot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Megobot struct {
	token    string
	session  *discordgo.Session
	commands []*discordgo.ApplicationCommand
	guildId  string
	timeout  time.Duration
}

func New(token string, guildId string) (*Megobot, error) {
	s, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, fmt.Errorf("error creating session: %s", err)
	}

	s.Identify.Intents = discordgo.IntentsGuilds

	bot := &Megobot{
		token:   token,
		session: s,
		guildId: guildId,
		timeout: 30 * time.Second,
	}

	return bot, nil
}

func (m *Megobot) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func (m *Megobot) AddCommand(name string, description string, options []*discordgo.ApplicationCommandOption, handler func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string) {
	c := &discordgo.ApplicationCommand{
		Name:        name,
		Description: description,
		Options:     options,
	}

	m.commands = append(m.commands, c)

	m.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		if i.ApplicationCommandData().Name == name {
			// Envia resposta de deferral imediatamente (pensando...)
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			if err != nil {
				log.Printf("error sending deferred response: %s", err)
				return
			}

			// Organiza os parametros do comando em um map
			optionsReq := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(optionsReq))
			for _, opt := range optionsReq {
				optionMap[opt.Name] = opt
			}

			// Canal para receber o resultado
			resultChan := make(chan string, 1)

			go func() {
				resultChan <- handler(s, i, optionMap)
			}()

			// Aguarda resultado ou timeout
			var finalResponse string
			select {
			case res := <-resultChan:
				finalResponse = res
			case <-time.After(m.timeout):
				finalResponse = "⏱️ Tempo limite de resposta excedido"
				log.Printf("command '%s' timed out", name)
			}

			if finalResponse != "" {
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &finalResponse,
				})
				if err != nil {
					log.Printf("error editing interaction response: %s", err)
				}
			}
		}
	})
}

func (m *Megobot) loadCommands() {
	_, err := m.session.ApplicationCommandBulkOverwrite(m.session.State.User.ID, m.guildId, m.commands)
	if err != nil {
		log.Println("failed to create commands:", err)
	}
}

func (m *Megobot) clearCommands() {
	_, err := m.session.ApplicationCommandBulkOverwrite(m.session.State.User.ID, m.guildId, []*discordgo.ApplicationCommand{})
	if err != nil {
		log.Println("error cleaning up commands:", err)
	}
}

func (m *Megobot) Start(ctx context.Context) error {
	err := m.session.Open()
	if err != nil {
		log.Println("error opening connection:", err)
		return err
	}
	defer m.session.Close()

	m.loadCommands()
	defer m.clearCommands()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("context cancelled:", ctx.Err())
		return ctx.Err()
	case sig := <-sc:
		log.Println("signal received:", sig)
		return nil
	}
}
