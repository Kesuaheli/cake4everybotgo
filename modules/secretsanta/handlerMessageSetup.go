package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

)

func (cmd MsgCmd) handler() {
	joinEmoji := util.GetConfigEmoji("secretsanta")
	joinEmojiID := joinEmoji.ID
	if joinEmojiID == "" {
		joinEmojiID = joinEmoji.Name
	}

	msg := cmd.data.Resolved.Messages[cmd.data.TargetID]
	if len(msg.Reactions) == 0 {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.no_reactions"), joinEmojiID)
		return
	}
	var hasReaction bool
	for _, r := range msg.Reactions {
		if !util.CompareEmoji(r.Emoji, joinEmoji) {
			continue
		}
		hasReaction = true
		break
	}

	if !hasReaction {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.no_reactions"), joinEmojiID)
		return
	}

	users, err := cmd.Session.MessageReactions(msg.ChannelID, msg.ID, joinEmojiID, 100, "", "")
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
