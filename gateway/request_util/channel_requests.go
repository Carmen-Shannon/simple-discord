package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
)

func GetChannel(getDto dto.GetChannelDto, token string) (*structs.Channel, error) {
	path := "/channels/" + getDto.ChannelID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var channel structs.Channel
	err = json.Unmarshal(resp, &channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func ModifyChannel(updates dto.UpdateChannelDto, token string) (*structs.Channel, error) {
	path := "/channels/" + updates.ChannelID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(updates)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return nil, err
	}

	var channel structs.Channel
	err = json.Unmarshal(resp, &channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func DeleteChannel(deleteDto dto.GetChannelDto, token string) (*structs.Channel, error) {
	path := "/channels/" + deleteDto.ChannelID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var channel structs.Channel
	err = json.Unmarshal(resp, &channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func EditChannelPermissions(editPermissionsDto dto.EditChannelPermissionsDto, token string) (error) {
	path := "/channels/" + editPermissionsDto.ChannelID.ToString() + "/permissions/" + editPermissionsDto.OverwriteID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(editPermissionsDto)
	if err != nil {
		return err
	}

	_, err = HttpRequest("PUT", path, headers, body)
	if err != nil {
		return err
	}

	return nil
}