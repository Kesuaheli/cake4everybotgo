package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

)

func (cmd MsgCmd) handler() {
	const emojiName = "👍"

	msg := cmd.data.Resolved.Messages[cmd.data.TargetID]
	if len(msg.Reactions) == 0 {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.no_reactions"), emojiName)
		return
	}
	var reaction *discordgo.MessageReactions
	for _, r := range msg.Reactions {
		if r.Emoji.Name != emojiName {
			continue
		}
		reaction = r
		break
	}

	if reaction == nil {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.no_reactions"), emojiName)
		return
	}

	emojiID := reaction.Emoji.ID
	if emojiID == "" {
		emojiID = reaction.Emoji.Name
	}
	users, err := cmd.Session.MessageReactions(msg.ChannelID, msg.ID, emojiID, 100, "", "")
	if err != nil {
		log.Printf("Error on get users: %v\n", err)
		cmd.ReplyError()
		return
	}

	if len(users) < 2 {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.not_enough_reactions"), 2)
		return
	}
}
