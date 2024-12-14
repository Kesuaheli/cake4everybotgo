package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

// handleSubcommandUpdate handles the functionality of the update subcommand
func (cmd Chat) handleSubcommandUpdate() {
	cmd.ReplyDeferedHidden()
	players, err := cmd.getPlayers()
	if err != nil {
		cmd.ReplyError()
		return
	}

	var failedToSend string
	for _, p := range players {
		var DMChannel *discordgo.Channel
		DMChannel, err = cmd.Session.UserChannelCreate(p.User.ID)
		if err != nil {
			log.Printf("ERROR: could not create DM channel for user %s: %+v", p.User.ID, err)
			failedToSend += "\n- " + p.Mention()
			continue
		}

		if p.MessageID == "" {
			var msg *discordgo.Message
			msg, err = cmd.Session.ChannelMessageSendComplex(DMChannel.ID, cmd.inviteMessage(p))
			if err != nil {
				log.Printf("ERROR: could not send invite message for %s: %+v", p.DisplayName(), err)
				failedToSend += "\n- " + p.Mention()
				continue
			}
			p.MessageID = msg.ID
		} else {
			_, err = cmd.Session.ChannelMessageEditComplex(util.MessageComplexEdit(cmd.inviteMessage(p), DMChannel.ID, p.MessageID))
			if err != nil {
				log.Printf("ERROR: could not update bot message for %s '%s/%s': %+v", p.DisplayName(), cmd.Interaction.ChannelID, p.MessageID, err)
				failedToSend += "\n- " + p.Mention()
				return
			}
		}
	}

	if failedToSend != "" {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.cmd.update.error"), failedToSend)
		return
	}
	cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.cmd.update.success"), len(players))
}
