package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
)

func ListGuildEmojis(getDto dto.ListGuildEmojisDto, token string) ([]structs.Emoji, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/emojis"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var emojis []structs.Emoji
	err = json.Unmarshal(resp, &emojis)
	if err != nil {
		return nil, err
	}

	return emojis, nil
}

func GetGuildEmoji(getDto dto.GetGuildEmojiDto, token string) (*structs.Emoji, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/emojis/" + getDto.EmojiID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var emoji structs.Emoji
	err = json.Unmarshal(resp, &emoji)
	if err != nil {
		return nil, err
	}

	return &emoji, nil
}

func CreateGuildEmoji(createDto dto.CreateGuildEmojiDto, token string) (*structs.Emoji, error) {
	path := "/guilds/" + createDto.GuildID.ToString() + "/emojis"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(createDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	var emoji structs.Emoji
	err = json.Unmarshal(resp, &emoji)
	if err != nil {
		return nil, err
	}

	return &emoji, nil
}

func ModifyGuildEmoji(modifyDto dto.ModifyGuildEmojiDto, token string) (*structs.Emoji, error) {
	path := "/guilds/" + modifyDto.GuildID.ToString() + "/emojis/" + modifyDto.EmojiID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(modifyDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return nil, err
	}

	var emoji structs.Emoji
	err = json.Unmarshal(resp, &emoji)
	if err != nil {
		return nil, err
	}

	return &emoji, nil
}

func DeleteGuildEmoji(deleteDto dto.DeleteGuildEmojiDto, token string) error {
	path := "/guilds/" + deleteDto.GuildID.ToString() + "/emojis/" + deleteDto.EmojiID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

type listApplicationEmojisResponseWrapper struct {
	Items []structs.Emoji `json:"items"`
}

func ListApplicationEmojis(getDto dto.ListApplicationEmojisDto, token string) ([]structs.Emoji, error) {
	path := "/applications/" + getDto.ApplicationID.ToString() + "/emojis"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var respWrapper listApplicationEmojisResponseWrapper
	err = json.Unmarshal(resp, &respWrapper)
	if err != nil {
		return nil, err
	}

	return respWrapper.Items, nil
}

func GetApplicationEmoji(getDto dto.GetApplicationEmojiDto, token string) (*structs.Emoji, error) {
	path := "/applications/" + getDto.ApplicationID.ToString() + "/emojis/" + getDto.EmojiID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var emoji structs.Emoji
	err = json.Unmarshal(resp, &emoji)
	if err != nil {
		return nil, err
	}

	return &emoji, nil
}

func CreateApplicationEmoji(createDto dto.CreateApplicationEmojiDto, token string) (*structs.Emoji, error) {
	path := "/applications/" + createDto.ApplicationID.ToString() + "/emojis"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(createDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	var emoji structs.Emoji
	err = json.Unmarshal(resp, &emoji)
	if err != nil {
		return nil, err
	}

	return &emoji, nil
}

func ModifyApplicationEmoji(patchDto dto.ModifyApplicationEmojiDto, token string) (*structs.Emoji, error) {
	path := "/applications/" + patchDto.ApplicationID.ToString() + "/emojis/" + patchDto.EmojiID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(patchDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return nil, err
	}

	var emoji structs.Emoji
	err = json.Unmarshal(resp, &emoji)
	if err != nil {
		return nil, err
	}

	return &emoji, nil
}

func DeleteApplicationEmoji(deleteDto dto.DeleteApplicationEmojiDto, token string) error {
	path := "/applications/" + deleteDto.ApplicationID.ToString() + "/emojis/" + deleteDto.EmojiID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}
