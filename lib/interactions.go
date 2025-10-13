package lib

import (
	"errors"
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
	ID          Snowflake `json:"id"`
	FileName    string    `json:"filename"`
	Description string    `json:"description,omitempty"`
	ContentType string    `json:"content_type,omitempty"`
	URL         string    `json:"url"`
}

type AllowedMentions struct {
	Parse []string `json:"parse,omitempty"`
}

func (bot *Bot) commandInteractionHandler(w http.ResponseWriter, interaction CommandInteraction) {
	interaction, command, available := bot.handleInteraction(interaction)
	if !available {
		w.Header().Add("Content-Type", CONTENT_TYPE_JSON)
		if _, err := w.Write(bodyUnknownCommandResponse); err != nil {
			http.Error(w, "internal server error - failed to write response", http.StatusInternalServerError)
			LogError("Failed to write unknown command response", "error", err, "interaction_id", interaction.ID, "user_id", interaction.User.ID)
			return
		}
		LogError("Command unavailable", "command", interaction.Data.Name, "guild", interaction.GuildID)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	interaction.Bot = bot

	LogInfo("Command executed", "command", command.Name)

	command.CommandHandler(&interaction)
}

func (bot *Bot) handleInteraction(interaction CommandInteraction) (CommandInteraction, Command, bool) {
	if len(interaction.Data.Options) > 0 && interaction.Data.Options[0].Type == SUB_OPTION_TYPE {
		finalName := interaction.Data.Name + "@" + interaction.Data.Options[0].Name
		subCommand, available := bot.commands.Get(finalName)
		if available {
			interaction.Data.Name, interaction.Data.Options = finalName, interaction.Data.Options[0].Options
		}
		return interaction, subCommand, available
	}

	command, available := bot.commands.Get(interaction.Data.Name)
	return interaction, command, available
}

func (interaction CommandInteraction) Defer(ephemeral bool) {
	var flags MessageFlags = 0

	if ephemeral {
		flags = MESSAGE_FLAG_EPHEMERAL
	}

	if err := interaction.Bot.makeHttpRequestToDiscord(http.MethodPost, "/interactions/"+interaction.ID.String()+"/"+interaction.Token+"/callback", ResponseMessage{
		Type: DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE_RESPONSE_TYPE,
		Data: &ResponseMessageData{
			Flags: flags,
		},
	}, nil, false); err != nil {
		LogError("Failed to defer interaction", "error", err, "interaction_id", interaction.ID, "user_id", interaction.User.ID)
		return
	}
}

func (interaction CommandInteraction) SendReply(reply ResponseMessageData, ephemeral bool, files []DiscordFile) {
	if ephemeral && reply.Flags == 0 {
		reply.Flags = MESSAGE_FLAG_EPHEMERAL
	}

	if err := interaction.Bot.makeHttpRequestToDiscord(http.MethodPost, "/interactions/"+interaction.ID.String()+"/"+interaction.Token+"/callback", ResponseMessage{
		Type: CHANNEL_MESSAGE_WITH_SOURCE_RESPONSE_TYPE,
		Data: &reply,
	}, files, false); err != nil {
		LogError("Failed to send interaction reply", "error", err, "interaction_id", interaction.ID, "user_id", interaction.User.ID)
		return
	}
}

func (interaction CommandInteraction) SendSimpleReply(content string, ephemeral bool) {
	interaction.SendReply(ResponseMessageData{
		Content: content,
	}, ephemeral, nil)
}

func (interaction CommandInteraction) EditReply(reply ResponseMessageData, ephemeral bool, files []DiscordFile) {
	if ephemeral && reply.Flags == 0 {
		reply.Flags = MESSAGE_FLAG_EPHEMERAL
	}

	if err := interaction.Bot.makeHttpRequestToDiscord(http.MethodPatch, "/webhooks/"+interaction.Bot.ApplicationID.String()+"/"+interaction.Token+"/messages/@original", reply, files, false); err != nil {
		LogError("Failed to edit interaction reply", "error", err, "interaction_id", interaction.ID, "user_id", interaction.User.ID)
		return
	}
}

func (interaction CommandInteraction) GetIntOptionValue(name string, fallback int) (int, error) {
	options := interaction.Data.Options
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

func (interaction CommandInteraction) GetFloatOptionValue(name string, fallback float64) (float64, error) {
	options := interaction.Data.Options
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

func (interaction CommandInteraction) GetStringOptionValue(name string, fallback string) (string, error) {
	options := interaction.Data.Options
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

func (interaction CommandInteraction) GetBoolOptionValue(name string, fallback bool) (bool, error) {
	options := interaction.Data.Options
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

func (interaction CommandInteraction) GetAttachmentOptionId(name string, fallback string) (string, error) {
	options := interaction.Data.Options
	if len(options) == 0 {
		return fallback, errors.New("no options provided")
	}

	for _, option := range options {
		if option.Name == name {
			return option.Value.(string), nil
		}
	}

	return "", nil
}

func (interaction CommandInteraction) GetUser() *User {
	if interaction.User != nil {
		return interaction.User
	}
	if interaction.Member != nil {
		return interaction.Member.User
	}
	return nil
}
