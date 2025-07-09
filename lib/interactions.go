package lib

import (
	"errors"
	"log/slog"
	"net/http"
)

// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object
type ResponseMessage struct {
	Type ResponseType         `json:"type"`
	Data *ResponseMessageData `json:"data,omitempty"`
}

// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-messages
type ResponseMessageData struct {
	Content         string           `json:"content,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlags     `json:"flags,omitempty"`
	Attachments     []Attachment     `json:"attachments,omitempty"`
}

type Attachment struct {
	ID           Snowflake `json:"id"`
	FileName     string    `json:"filename"`
	Description  string    `json:"description,omitempty"`
	ContentType  string    `json:"content_type,omitempty"`
	URL          string    `json:"url"`
}

type AllowedMentions struct {
	Parse       []string `json:"parse,omitempty"`
}


func (bot *Bot) commandInteractionHandler(w http.ResponseWriter, interaction CommandInteraction) {
	itx, command, available := bot.handleInteraction(interaction)
	if !available {
		w.Header().Add("Content-Type", CONTENT_TYPE_JSON)
		w.Write(bodyUnknownCommandResponse)
		slog.Error("Command unavailable", "command", itx.Data.Name, "guild", itx.GuildID, "user", itx.User.ID)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	itx.Bot = bot

	slog.Info("Command executed", "command", command.Name)

	command.CommandHandler(&itx)
}

func (bot *Bot) handleInteraction(itx CommandInteraction) (CommandInteraction, Command, bool) {
	if len(itx.Data.Options) > 0 && itx.Data.Options[0].Type == SUB_OPTION_TYPE {
		finalName := itx.Data.Name + "@" + itx.Data.Options[0].Name
		subCommand, available := bot.commands.Get(finalName)
		if available {
			if itx.Member != nil {
				itx.Member.GuildID = itx.GuildID
			}

			itx.Data.Name, itx.Data.Options = finalName, itx.Data.Options[0].Options
		}
		return itx, subCommand, available
	}

	if itx.Member != nil {
		itx.Member.GuildID = itx.GuildID
	}

	command, available := bot.commands.Get(itx.Data.Name)
	return itx, command, available
}

func (itx CommandInteraction) Defer(ephemeral bool) error {
	var flags MessageFlags = 0

	if ephemeral {
		flags = MESSAGE_FLAG_EPHEMERAL
	}

	err := itx.Bot.makeHttpRequestToDiscord(http.MethodPost, "/interactions/"+itx.ID.String()+"/"+itx.Token+"/callback", ResponseMessage{
		Type: DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE_RESPONSE_TYPE,
		Data: &ResponseMessageData{
			Flags: flags,
		},
	}, nil, false)

	return err
}

func (itx CommandInteraction) SendReply(reply ResponseMessageData, ephemeral bool, files []DiscordFile) error {
	if ephemeral && reply.Flags == 0 {
		reply.Flags = MESSAGE_FLAG_EPHEMERAL
	}

	err := itx.Bot.makeHttpRequestToDiscord(http.MethodPost, "/interactions/"+itx.ID.String()+"/"+itx.Token+"/callback", ResponseMessage{
		Type: CHANNEL_MESSAGE_WITH_SOURCE_RESPONSE_TYPE,
		Data: &reply,
	}, files, false)

	return err
}

func (itx CommandInteraction) SendSimpleReply(content string, ephemeral bool) error {
	return itx.SendReply(ResponseMessageData{
		Content: content,
	}, ephemeral, nil)
}

func (itx CommandInteraction) EditReply(reply ResponseMessageData, ephemeral bool, files []DiscordFile) error {
	if ephemeral && reply.Flags == 0 {
		reply.Flags = MESSAGE_FLAG_EPHEMERAL
	}

	err := itx.Bot.makeHttpRequestToDiscord(http.MethodPatch, "/webhooks/"+itx.Bot.ApplicationID.String()+"/"+itx.Token+"/messages/@original", reply, files, false)

	return err
}

func (itx CommandInteraction) GetIntOptionValue(name string, fallback int) (int, error) {
	options := itx.Data.Options
	if len(options) == 0 {
		return 0, errors.New("no options provided")
	}

	for _, option := range options {
		if option.Name == name {
			return int(option.Value.(float64)), nil
		}
	}

	return fallback, nil
}

func (itx CommandInteraction) GetFloatOptionValue(name string, fallback float64) (float64, error) {
	options := itx.Data.Options
	if len(options) == 0 {
		return 0, errors.New("no options provided")
	}
	for _, option := range options {
		if option.Name == name {
			return option.Value.(float64), nil
		}
	}
	return fallback, nil
}

func (itx CommandInteraction) GetStringOptionValue(name string, fallback string) (string, error) {
	options := itx.Data.Options
	if len(options) == 0 {
		return "", errors.New("no options provided")
	}

	for _, option := range options {
		if option.Name == name {
			return option.Value.(string), nil
		}
	}

	return fallback, nil
}

func (itx CommandInteraction) GetBoolOptionValue(name string, fallback bool) (bool, error) {
	options := itx.Data.Options
	if len(options) == 0 {
		return false, errors.New("no options provided")
	}

	for _, option := range options {
		if option.Name == name {
			return option.Value.(bool), nil
		}
	}

	return fallback, nil
}