package dto

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/util"
)

type GetGlobalApplicationCommandsDto struct {
	WithLocalizations *bool `json:"with_localizations,omitempty"`
}

type BulkOverwriteGlobalApplicationCommandsDto struct {
	Commands []CreateGlobalApplicationCommandDto `json:"commands"`
}

type EditApplicationCommandPermissionsDto struct {
	Permissions []structs.ApplicationCommandPermissions `json:"permissions"`
}

type EditGuildApplicationCommandDto interface {
	SetName(name string) error
	SetNameLocalizations(localizations map[string]string) error
	SetDescription(description string) error
	SetDescriptionLocalizations(localizations map[string]string) error
	SetOptions(options []structs.ApplicationCommandOption) error
	SetDefaultMemberPermissions(permissions structs.Bitfield[structs.Permission])
	SetDefaultPermission(defaultPermission bool)
	SetNsfw(nsfw bool)
}

type editGuildApplicationCommandDto struct {
	Name                     *string                             `json:"name,omitempty"`
	NameLocalizations        *map[string]string                  `json:"name_localizations,omitempty"`
	Description              *string                             `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string                  `json:"description_localizations,omitempty"`
	Options                  *[]structs.ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                             `json:"default_member_permissions,omitempty"`
	DefaultPermission        *bool                               `json:"default_permission,omitempty"`
	Nsfw                     *bool                               `json:"nsfw,omitempty"`
}

var _ EditGuildApplicationCommandDto = (*editGuildApplicationCommandDto)(nil)

func (e *editGuildApplicationCommandDto) SetName(name string) error {
	if len(name) < 1 || len(name) > 32 {
		return fmt.Errorf("name must be between 1 and 32 characters")
	}
	e.Name = &name
	return nil
}

func (e *editGuildApplicationCommandDto) SetNameLocalizations(localizations map[string]string) error {
	for key, val := range localizations {
		if structs.FindLocaleByCode(key) == nil {
			return fmt.Errorf("invalid locale: %s", key)
		}

		if len(val) < 1 || len(val) > 32 {
			return fmt.Errorf("name must be between 1 and 32 characters")
		}
	}
	e.NameLocalizations = &localizations
	return nil
}

func (e *editGuildApplicationCommandDto) SetDescription(description string) error {
	if len(description) < 1 || len(description) > 100 {
		return fmt.Errorf("description must be between 1 and 100 characters")
	}
	e.Description = &description
	return nil
}

func (e *editGuildApplicationCommandDto) SetDescriptionLocalizations(localizations map[string]string) error {
	for key, val := range localizations {
		if structs.FindLocaleByCode(key) == nil {
			return fmt.Errorf("invalid locale: %s", key)
		}

		if len(val) < 1 || len(val) > 100 {
			return fmt.Errorf("description must be between 1 and 100 characters")
		}
	}
	e.DescriptionLocalizations = &localizations
	return nil
}

func (e *editGuildApplicationCommandDto) SetOptions(options []structs.ApplicationCommandOption) error {
	if len(options) > 25 {
		return fmt.Errorf("options must be less than or equal to 25")
	}
	e.Options = &options
	return nil
}

func (e *editGuildApplicationCommandDto) SetDefaultMemberPermissions(permissions structs.Bitfield[structs.Permission]) {
	e.DefaultMemberPermissions = util.ToPtr(permissions.ToString())
}

func (e *editGuildApplicationCommandDto) SetDefaultPermission(defaultPermission bool) {
	e.DefaultPermission = util.ToPtr(defaultPermission)
}

func (e *editGuildApplicationCommandDto) SetNsfw(nsfw bool) {
	e.Nsfw = util.ToPtr(nsfw)
}

type CreateGuildApplicationCommandDto interface {
	SetName(name string) error
	SetNameLocalizations(localizations map[string]string) error
	SetDescription(description string) error
	SetDescriptionLocalizations(localizations map[string]string) error
	SetOptions(options []structs.ApplicationCommandOption) error
	SetDefaultMemberPermissions(permissions structs.Bitfield[structs.Permission])
	SetDefaultPermission(defaultPermission bool)
	SetType(commandType structs.ApplicationCommandType)
	SetNsfw(nsfw bool)
}

type createGuildApplicationCommandDto struct {
	Name                     string                              `json:"name"`
	NameLocalizations        *map[string]string                  `json:"name_localizations,omitempty"`
	Description              *string                             `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string                  `json:"description_localizations,omitempty"`
	Options                  *[]structs.ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                             `json:"default_member_permissions,omitempty"`
	DefaultPermission        *bool                               `json:"default_permission,omitempty"`
	Type                     *structs.ApplicationCommandType     `json:"type,omitempty"`
	Nsfw                     *bool                               `json:"nsfw,omitempty"`
}

var _ CreateGuildApplicationCommandDto = (*createGuildApplicationCommandDto)(nil)

func (c *createGuildApplicationCommandDto) SetName(name string) error {
	if len(name) < 1 || len(name) > 32 {
		return fmt.Errorf("name must be between 1 and 32 characters")
	}

	if c.Type != nil && *c.Type == structs.ChatInputCommand {
		if err := validateName(name); err != nil {
			return err
		}
	}

	c.Name = name
	return nil
}

func (c *createGuildApplicationCommandDto) SetNameLocalizations(localizations map[string]string) error {
	for key, val := range localizations {
		if structs.FindLocaleByCode(key) == nil {
			return fmt.Errorf("invalid locale: %s", key)
		}

		if len(val) < 1 || len(val) > 32 {
			return fmt.Errorf("name must be between 1 and 32 characters")
		}

		if c.Type != nil && *c.Type == structs.ChatInputCommand {
			if err := validateName(val); err != nil {
				return err
			}
		}
	}
	c.NameLocalizations = &localizations
	return nil
}

func (c *createGuildApplicationCommandDto) SetDescription(description string) error {
	if len(description) < 1 || len(description) > 100 {
		return fmt.Errorf("description must be between 1 and 100 characters")
	}

	c.Description = &description
	return nil
}

func (c *createGuildApplicationCommandDto) SetDescriptionLocalizations(localizations map[string]string) error {
	for key, val := range localizations {
		if structs.FindLocaleByCode(key) == nil {
			return fmt.Errorf("invalid locale: %s", key)
		}

		if len(val) < 1 || len(val) > 100 {
			return fmt.Errorf("description must be between 1 and 100 characters")
		}
	}
	c.DescriptionLocalizations = &localizations
	return nil
}

func (c *createGuildApplicationCommandDto) SetOptions(options []structs.ApplicationCommandOption) error {
	if len(options) > 25 {
		return fmt.Errorf("options must be less than or equal to 25")
	}

	c.Options = &options
	return nil
}

func (c *createGuildApplicationCommandDto) SetDefaultMemberPermissions(permissions structs.Bitfield[structs.Permission]) {
	c.DefaultMemberPermissions = util.ToPtr(permissions.ToString())
}

func (c *createGuildApplicationCommandDto) SetDefaultPermission(defaultPermission bool) {
	c.DefaultPermission = util.ToPtr(defaultPermission)
}

func (c *createGuildApplicationCommandDto) SetType(commandType structs.ApplicationCommandType) {
	c.Type = util.ToPtr(commandType)
}

func (c *createGuildApplicationCommandDto) SetNsfw(nsfw bool) {
	c.Nsfw = util.ToPtr(nsfw)
}

type CreateGlobalApplicationCommandDto interface {
	SetName(name string) error
	SetNameLocalizations(localizations map[string]string) error
	SetDescription(description string) error
	SetDescriptionLocalizations(localizations map[string]string) error
	SetOptions(options []structs.ApplicationCommandOption) error
	SetDefaultMemberPermissions(permissions structs.Bitfield[structs.Permission])
	SetDmPermission(dmPermission bool)
	SetDefaultPermission(defaultPermission bool)
	SetIntegrationTypes(integrationTypes []structs.IntegrationType)
	SetContexts(contexts []structs.IntegrationContextType)
	SetType(commandType structs.ApplicationCommandType)
	SetNsfw(nsfw bool)
}

type createGlobalApplicationCommandDto struct {
	Name                     string                              `json:"name"`
	NameLocalizations        *map[string]string                  `json:"name_localizations,omitempty"`
	Description              *string                             `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string                  `json:"description_localizations,omitempty"`
	Options                  *[]structs.ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                             `json:"default_member_permissions,omitempty"`
	DmPermission             *bool                               `json:"dm_permission,omitempty"`      // DEPRECATED use Contexts instead
	DefaultPermission        *bool                               `json:"default_permission,omitempty"` // replaced by DefaultMemberPermissions defaults to true
	IntegrationTypes         *[]structs.IntegrationType          `json:"integration_types,omitempty"`
	Contexts                 *[]structs.IntegrationContextType   `json:"contexts,omitempty"`
	Type                     *structs.ApplicationCommandType     `json:"type,omitempty"` // defaults to 1 (ChatInputCommand) if nil
	Nsfw                     *bool                               `json:"nsfw,omitempty"`
}

var _ CreateGlobalApplicationCommandDto = (*createGlobalApplicationCommandDto)(nil)

func (c *createGlobalApplicationCommandDto) SetName(name string) error {
	if len(name) < 1 || len(name) > 32 {
		return fmt.Errorf("name must be between 1 and 32 characters")
	}

	if c.Type != nil && *c.Type == structs.ChatInputCommand {
		if err := validateName(name); err != nil {
			return err
		}
	}

	c.Name = name
	return nil
}

func (c *createGlobalApplicationCommandDto) SetNameLocalizations(localizations map[string]string) error {
	for key, val := range localizations {
		if structs.FindLocaleByCode(key) == nil {
			return fmt.Errorf("invalid locale: %s", key)
		}

		if len(val) < 1 || len(val) > 32 {
			return fmt.Errorf("name must be between 1 and 32 characters")
		}

		if c.Type != nil && *c.Type == structs.ChatInputCommand {
			if err := validateName(val); err != nil {
				return err
			}
		}
	}
	c.NameLocalizations = &localizations
	return nil
}

func (c *createGlobalApplicationCommandDto) SetDescription(description string) error {
	if len(description) < 1 || len(description) > 100 {
		return fmt.Errorf("description must be between 1 and 100 characters")
	}

	c.Description = &description
	return nil
}

func (c *createGlobalApplicationCommandDto) SetDescriptionLocalizations(localizations map[string]string) error {
	for key, val := range localizations {
		if structs.FindLocaleByCode(key) == nil {
			return fmt.Errorf("invalid locale: %s", key)
		}

		if len(val) < 1 || len(val) > 100 {
			return fmt.Errorf("description must be between 1 and 100 characters")
		}
	}
	c.DescriptionLocalizations = &localizations
	return nil
}

func (c *createGlobalApplicationCommandDto) SetOptions(options []structs.ApplicationCommandOption) error {
	if len(options) > 25 {
		return fmt.Errorf("options must be less than or equal to 25")
	}

	c.Options = &options
	return nil
}

func (c *createGlobalApplicationCommandDto) SetDefaultMemberPermissions(permissions structs.Bitfield[structs.Permission]) {
	c.DefaultMemberPermissions = util.ToPtr(permissions.ToString())
}

func (c *createGlobalApplicationCommandDto) SetDmPermission(dmPermission bool) {
	c.DmPermission = util.ToPtr(dmPermission)
}

func (c *createGlobalApplicationCommandDto) SetDefaultPermission(defaultPermission bool) {
	c.DefaultPermission = util.ToPtr(defaultPermission)
}

func (c *createGlobalApplicationCommandDto) SetIntegrationTypes(integrationTypes []structs.IntegrationType) {
	c.IntegrationTypes = &integrationTypes
}

func (c *createGlobalApplicationCommandDto) SetContexts(contexts []structs.IntegrationContextType) {
	c.Contexts = &contexts
}

func (c *createGlobalApplicationCommandDto) SetType(commandType structs.ApplicationCommandType) {
	c.Type = util.ToPtr(commandType)
}

func (c *createGlobalApplicationCommandDto) SetNsfw(nsfw bool) {
	c.Nsfw = util.ToPtr(nsfw)
}

// NewGlobalApplicationCommandDto creates a new CreateGlobalApplicationCommandDto with the given name, description, and command type.
// The command type defaults to ChatInputCommand if nil.
// The name must match the regex `^[-_\p{L}\p{N}\p{sc=Deva}\p{sc=Thai}]{1,32}$` if the commandType is ChatInputCommand.
func NewGlobalApplicationCommandDto(name string, commandType *structs.ApplicationCommandType) CreateGlobalApplicationCommandDto {
	if len(name) < 1 || len(name) > 32 {
		log.Fatalf("invalid name: name must be between 1 and 32 characters")
		return nil
	}

	if commandType == nil {
		commandType = util.ToPtr(structs.ChatInputCommand)
	}

	if *commandType == structs.ChatInputCommand {
		if err := validateName(name); err != nil {
			log.Fatalf("invalid name: %s", err)
			return nil
		}
	}

	return &createGlobalApplicationCommandDto{
		Name: name,
		Type: commandType,
	}
}

// NewGuildApplicationCommandDto creates a new CreateGuildApplicationCommandDto with the given name, description, and command type.
// The command type defaults to ChatInputCommand if nil.
// The name must match the regex `^[-_\p{L}\p{N}\p{sc=Deva}\p{sc=Thai}]{1,32}$` if the commandType is ChatInputCommand.
func NewGuildApplicationCommandDto(name string, commandType *structs.ApplicationCommandType) CreateGuildApplicationCommandDto {
	if len(name) < 1 || len(name) > 32 {
		log.Fatalf("invalid name: name must be between 1 and 32 characters")
		return nil
	}

	if commandType == nil {
		commandType = util.ToPtr(structs.ChatInputCommand)
	}

	if *commandType == structs.ChatInputCommand {
		if err := validateName(name); err != nil {
			log.Fatalf("invalid name: %s", err)
			return nil
		}
	}

	return &createGuildApplicationCommandDto{
		Name: name,
		Type: commandType,
	}
}

// validateName is used locally to validate the Name of a CreateGlobalApplicationCommandDto
// it should only be used to validate names for the ChatInputCommand type
func validateName(name string) error {
	if len(name) < 1 || len(name) > 32 {
		return fmt.Errorf("name must be between 1 and 32 characters")
	}

	pattern := `^[-_\p{L}\p{N}\p{Devanagari}\p{Thai}]{1,32}$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	if !re.MatchString(name) {
		return fmt.Errorf("invalid name: %s", name)
	}

	return nil
}
