package megobot

import "github.com/bwmarrin/discordgo"

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string

type Command struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     CommandHandler
}

func (c Command) Bind(bot *Megobot) {
	bot.AddCommand(
		c.Name,
		c.Description,
		c.Options,
		c.Handler,
	)
}
