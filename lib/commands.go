package lib

import (
	"net/http"
	"strconv"
	"strings"
)

type Command struct {
	Type        CommandType     `json:"type,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     []CommandOption `json:"options,omitempty"`

	CommandHandler func(interaction *CommandInteraction) `json:"-"`
}

type CommandOption struct {
	Type        OptionType      `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Required    bool            `json:"required,omitempty"`
	MinValue    float64         `json:"min_value,omitempty"`
	MaxValue    float64         `json:"max_value,omitempty"`
	MinLength   uint32          `json:"min_length,omitempty"`
	MaxLength   uint32          `json:"max_length,omitempty"`
	Options     []CommandOption `json:"options,omitempty"`
}

type CommandType uint8

const (
	COMMAND_TYPE_CHAT_INPUT CommandType = iota + 1
	COMMAND_TYPE_USER
	COMMAND_TYPE_MESSAGE
)

type CommandInteraction struct {
	ID            Snowflake              `json:"id"`
	ApplicationID Snowflake              `json:"application_id"`
	Type          InteractionType        `json:"type"`
	Data          CommandInteractionData `json:"data"`
	GuildID       Snowflake              `json:"guild_id,omitempty"`
	ChannelID     Snowflake              `json:"channel_id,omitempty"`
	Member        *Member                `json:"member,omitempty"`
	User          *User                  `json:"user,omitempty"`
	Token         string                 `json:"token"`

	Bot *Bot `json:"-"`
}

type CommandInteractionData struct {
	ID       Snowflake                  `json:"id,omitempty"`
	Name     string                     `json:"name"`
	Type     CommandType                `json:"type"`
	Resolved *InteractionDataResolved   `json:"resolved,omitempty"`
	Options  []CommandInteractionOption `json:"options,omitempty"`
	GuildID  Snowflake                  `json:"guild_id,omitempty"`
	TargetID Snowflake                  `json:"target_id,omitempty"`
}

type Member struct {
	User            *User  `json:"user,omitempty"`
	Nickname        string `json:"nick,omitempty"`
	GuildAvatarHash string `json:"avatar,omitempty"`
}

func (member Member) GuildAvatarURL(guildId string, userId string) string {
	if member.GuildAvatarHash == "" {
		return ""
	}

	if guildId == "" || userId == "" {
		return ""
	}

	return DISCORD_CDN_URL + "/guilds/" + guildId + "/users/" + userId + "/avatars/" + member.GuildAvatarHash + ".png"
}

type User struct {
	ID         Snowflake `json:"id"`
	Username   string    `json:"username"`
	GlobalName string    `json:"global_name,omitempty"`
	AvatarHash string    `json:"avatar,omitempty"`
}

func (user User) AvatarURL() string {
	if user.AvatarHash == "" {
		return DISCORD_CDN_URL + "/embed/avatars/" + strconv.FormatUint(uint64(user.ID>>22)%6, 10) + ".png"
	}

	return DISCORD_CDN_URL + "/avatars/" + user.ID.String() + "/" + user.AvatarHash + ".png"
}

type CommandInteractionOption struct {
	Name    string                     `json:"name"`
	Type    OptionType                 `json:"type"`
	Value   any                        `json:"value,omitempty"`
	Options []CommandInteractionOption `json:"options,omitempty"`
	Focused bool                       `json:"focused,omitempty"`
}

type InteractionDataResolved struct {
	Users       map[Snowflake]*User       `json:"users,omitempty"`
	Members     map[Snowflake]*Member     `json:"members,omitempty"`
	Messages    map[Snowflake]*Message    `json:"messages,omitempty"`
	Attachments map[Snowflake]*Attachment `json:"attachments,omitempty"`
}

func (data *InteractionDataResolved) String() string {
	return "Users: " + strconv.Itoa(len(data.Users)) +
		", Members: " + strconv.Itoa(len(data.Members)) +
		", Messages: " + strconv.Itoa(len(data.Messages)) +
		", Attachments: " + strconv.Itoa(len(data.Attachments))
}

type Message struct {
	ID     Snowflake `json:"id"`
	Author *User     `json:"author,omitempty"`
}

type OptionType uint8

const (
	SUB_OPTION_TYPE OptionType = iota + 1
	_                          // OPTION_SUB_COMMAND_GROUP (not supported)
	STRING_OPTION_TYPE
	INTEGER_OPTION_TYPE
	BOOLEAN_OPTION_TYPE
	USER_OPTION_TYPE
	CHANNEL_OPTION_TYPE
	ROLE_OPTION_TYPE
	MENTIONABLE_OPTION_TYPE
	NUMBER_OPTION_TYPE
	ATTACHMENT_OPTION_TYPE
)

func (bot *Bot) RegisterCommand(cmd Command) {
	if bot.commands.Has(cmd.Name) {
		LogError("bot already has registered \"" + cmd.Name + "\" slash command (name is already in use)")
		return
	}

	if cmd.Type == 0 {
		cmd.Type = COMMAND_TYPE_CHAT_INPUT
	}

	bot.commands.Set(cmd.Name, cmd)

}

func (bot *Bot) RegisterSubCommand(subCommand Command, parentCommandName string) {
	if !bot.commands.Has(parentCommandName) {
		LogError("missing \"" + parentCommandName + "\" slash command in registry (parent command needs to be registered in bot before adding subcommands)")
		return
	}

	finalName := parentCommandName + "@" + subCommand.Name
	if bot.commands.Has(finalName) {
		LogError("bot already has registered \"" + finalName + "\" slash command (name for subcommand is already in use)")
		return
	}

	if subCommand.Type == 0 {
		subCommand.Type = COMMAND_TYPE_CHAT_INPUT
	}

	bot.commands.Set(finalName, subCommand)
}

func (bot *Bot) SyncCommandsWithDiscord(guildIDs []Snowflake) error {
	commands := parseCommandsForDiscordAPI(bot.commands)

	url := "/applications/" + bot.ApplicationID.String() + "/commands"
	if len(guildIDs) != 0 {
		url = "/applications/" + bot.ApplicationID.String() + "/guilds/" + guildIDs[0].String() + "/commands"
	}

	err := bot.makeHttpRequestToDiscord(http.MethodPut, url, commands, nil, true)

	return err
}

func parseCommandsForDiscordAPI(commands *SharedMap[string, Command]) []Command {
	commands.mu.RLock()

	cmdTree := make(map[string]Command, len(commands.cache))
	parsedCommands := make([]Command, 0, len(commands.cache))

	for fullName, command := range commands.cache {
		if !strings.Contains(fullName, "@") {
			cmdTree[fullName] = command
		}
	}

	for fullName, command := range commands.cache {
		if strings.Contains(fullName, "@") {
			names := strings.SplitN(fullName, "@", 2)
			parentName := names[0]
			parentCommand, exists := cmdTree[parentName]
			if !exists {
				LogError("parent command \"" + parentName + "\" not found in command tree while parsing commands for Discord API")
				continue
			}

			parentCommand.Options = append(parentCommand.Options, CommandOption{
				Name:        command.Name,
				Description: command.Description,
				Type:        SUB_OPTION_TYPE,
				Options:     command.Options,
			})
			cmdTree[parentName] = parentCommand
		}
	}

	commands.mu.RUnlock()
	for _, command := range cmdTree {
		parsedCommands = append(parsedCommands, command)
	}

	return parsedCommands
}
