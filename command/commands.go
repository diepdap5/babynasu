package command

import "github.com/bwmarrin/discordgo"

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate) error
}
type (
	HandlerSig = func(s *discordgo.Session, i *discordgo.InteractionCreate) (discordgo.InteractionResponse, error)
)

var Commands = map[string]Command{
	"ping": {
		Command: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Pong!",
		},
		Handler: Ping,
	},
}
