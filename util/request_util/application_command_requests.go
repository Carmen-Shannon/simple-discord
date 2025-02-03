package request_util

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func GetGlobalApplicationCommands(dto dto.GetGlobalApplicationCommandsDto, token string) ([]structs.ApplicationCommand, error) {
	path := "/applications/" + dto.ApplicationID + "/commands"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(dto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var commands []structs.ApplicationCommand
	err = json.Unmarshal(resp, &commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func GetGlobalApplicationCommand(dto dto.GetGlobalApplicationCommandDto, token string) (*structs.ApplicationCommand, error) {
	path := "/applications/" + dto.ApplicationID + "/commands/" + dto.CommandID
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var command structs.ApplicationCommand
	err = json.Unmarshal(resp, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func BulkOverwriteGlobalApplicationCommands(dto dto.BulkOverwriteGlobalApplicationCommandsDto, token string) ([]structs.ApplicationCommand, error) {
	path := "/applications/" + dto.ApplicationID + "/commands"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(dto.Commands)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PUT", path, headers, body)
	if err != nil {
		return nil, err
	}

	var commands []structs.ApplicationCommand
	err = json.Unmarshal(resp, &commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func CreateGlobalApplicationCommand(dto dto.CreateGlobalApplicationCommandDto, applicationID string, token string) (*structs.ApplicationCommand, error) {
	path := "/applications/" + applicationID + "/commands"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	var command structs.ApplicationCommand
	err = json.Unmarshal(resp, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func EditGlobalApplicationCommand(dto dto.CreateGlobalApplicationCommandDto, applicationID, commandID, token string) (*structs.ApplicationCommand, error) {
	path := "/applications/" + applicationID + "/commands/" + commandID
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return nil, err
	}

	var command structs.ApplicationCommand
	err = json.Unmarshal(resp, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func DeleteGlobalApplicationCommand(applicationID, commandID, token string) error {
	path := "/applications/" + applicationID + "/commands/" + commandID
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetGuildApplicationCommands(dto dto.GetGlobalApplicationCommandsDto, applicationID, guildID, token string) ([]structs.ApplicationCommand, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(dto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var commands []structs.ApplicationCommand
	err = json.Unmarshal(resp, &commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func CreateGuildApplicationCommand(dto dto.CreateGuildApplicationCommandDto, applicationID, guildID, token string) (*structs.ApplicationCommand, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	var command structs.ApplicationCommand
	err = json.Unmarshal(resp, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func GetGuildApplicationCommand(applicationID, guildID, commandID, token string) (*structs.ApplicationCommand, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands/" + commandID
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var command structs.ApplicationCommand
	err = json.Unmarshal(resp, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func EditGuildApplicationCommand(dto dto.EditGuildApplicationCommandDto, applicationID, guildID, commandID, token string) (*structs.ApplicationCommand, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands/" + commandID
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return nil, err
	}

	var command structs.ApplicationCommand
	err = json.Unmarshal(resp, &command)
	if err != nil {
		return nil, err
	}

	return &command, nil
}

func DeleteGuildApplicationCommand(applicationID, guildID, commandID, token string) error {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands/" + commandID
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func BulkOverwriteGuildApplicationCommands(dto dto.BulkOverwriteGlobalApplicationCommandsDto, applicationID, guildID, token string) ([]structs.ApplicationCommand, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(dto.Commands)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PUT", path, headers, body)
	if err != nil {
		return nil, err
	}

	var commands []structs.ApplicationCommand
	err = json.Unmarshal(resp, &commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func GetGuildApplicationCommandPermissions(applicationID, guildID, token string) ([]structs.ApplicationCommandPermissions, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands/permissions"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var permissions []structs.ApplicationCommandPermissions
	err = json.Unmarshal(resp, &permissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func GetApplicationCommandPermissions(applicationID, guildID, commandID, token string) (*structs.ApplicationCommandPermissions, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands/" + commandID + "/permissions"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var permissions structs.ApplicationCommandPermissions
	err = json.Unmarshal(resp, &permissions)
	if err != nil {
		return nil, err
	}

	return &permissions, nil
}

func EditApplicationCommandPermissions(dto dto.EditApplicationCommandPermissionsDto, applicationID, guildID, token string) (*structs.ApplicationCommandPermissions, error) {
	path := "/applications/" + applicationID + "/guilds/" + guildID + "/commands/permissions"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PUT", path, headers, body)
	if err != nil {
		return nil, err
	}

	var permissions structs.ApplicationCommandPermissions
	err = json.Unmarshal(resp, &permissions)
	if err != nil {
		return nil, err
	}

	return &permissions, nil
}
