package structs

import "encoding/json"

type ApplicationCommand struct {
	ID                       Snowflake                   `json:"id"`
	Type                     ApplicationCommandType      `json:"type"`
	ApplicationID            Snowflake                   `json:"application_id"`
	GuildID                  *Snowflake                  `json:"guild_id,omitempty"`
	Name                     string                      `json:"name"`
	NameLocalizations        map[string]string           `json:"name_localizations,omitempty"`
	Description              string                      `json:"description"`
	DescriptionLocalizations map[string]string           `json:"description_localizations,omitempty"`
	Options                  *[]ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *Bitfield[Permission]       `json:"default_member_permissions,omitempty"`
	DmPermission             *bool                       `json:"dm_permission,omitempty"`      // DEPRECATED use contexts
	DefaultPermission        bool                        `json:"default_permission,omitempty"` // not recommended for use, soon deprecated
	NSFW                     bool                        `json:"nsfw"`
	IntegrationTypes         []IntegrationType           `json:"integration_types,omitempty"`
	Contexts                 []IntegrationContextType    `json:"contexts,omitempty"`
	Version                  Snowflake                   `json:"version"`
}

type ApplicationCommandType int

const (
	ChatInputCommand         ApplicationCommandType = 1
	UserCommand              ApplicationCommandType = 2
	MessageCommand           ApplicationCommandType = 3
	PrimaryEntryPointCommand ApplicationCommandType = 4
)

type ApplicationCommandOption struct {
	Type                     ApplicationCommandOptionType     `json:"type"`
	Name                     string                           `json:"name"`
	NameLocalizations        map[string]string                `json:"name_localizations,omitempty"`
	Description              string                           `json:"description"`
	DescriptionLocalizations map[string]string                `json:"description_localizations,omitempty"`
	Required                 bool                             `json:"required"`
	Choices                  []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options                  []ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes             []ChannelType                    `json:"channel_types,omitempty"`
	MinValue                 float64                          `json:"min_value"`
	MaxValue                 float64                          `json:"max_value"`
	MinLength                int                              `json:"min_length"`
	MaxLength                int                              `json:"max_length"`
	Autocomplete             bool                             `json:"autocomplete"`
}

type ApplicationCommandOptionType int

const (
	SubCommandOptionType      ApplicationCommandOptionType = 1
	SubCommandGroupOptionType ApplicationCommandOptionType = 2
	StringOptionType          ApplicationCommandOptionType = 3
	IntegerOptionType         ApplicationCommandOptionType = 4
	BooleanOptionType         ApplicationCommandOptionType = 5
	UserOptionType            ApplicationCommandOptionType = 6
	ChannelOptionType         ApplicationCommandOptionType = 7
	RoleOptionType            ApplicationCommandOptionType = 8
	MentionableOptionType     ApplicationCommandOptionType = 9
	NumberOptionType          ApplicationCommandOptionType = 10
	AttachmentOptionType      ApplicationCommandOptionType = 11
)

type ApplicationCommandOptionChoice struct {
	Name              string            `json:"name"`
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	Value             interface{}       `json:"value"` // can be a string, int or float64
}

type EntryPointCommandHandlerType int

const (
	AppHandlerType            EntryPointCommandHandlerType = 1
	DiscordLaunchActivityType EntryPointCommandHandlerType = 2
)

func (a *ApplicationCommandOptionChoice) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if name, ok := raw["name"].(string); ok {
		a.Name = name
	}

	if nameLocalizations, ok := raw["name_localizations"].(map[string]string); ok {
		a.NameLocalizations = nameLocalizations
	}

	if value, ok := raw["value"].(string); ok {
		a.Value = value
	} else if value, ok := raw["value"].(int); ok {
		a.Value = value
	} else if value, ok := raw["value"].(float64); ok {
		a.Value = value
	}

	return nil
}

func (a *ApplicationCommandOptionChoice) MarshalJSON() ([]byte, error) {
	raw := map[string]interface{}{
		"name":  a.Name,
		"value": a.Value,
	}

	if len(a.NameLocalizations) > 0 {
		raw["name_localizations"] = a.NameLocalizations
	}

	return json.Marshal(raw)
}
