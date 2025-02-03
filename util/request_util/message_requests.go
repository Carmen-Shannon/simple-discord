package request_util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func GetChannelMessages(query dto.GetChannelMessagesDto, token string) ([]structs.Message, error) {
	path := "/channels/" + query.ChannelID.ToString() + "/messages"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	queryParams := util.BuildQueryString(query)
	if queryParams != "" {
		path += queryParams
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var messages []structs.Message
	err = json.Unmarshal(resp, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func GetChannelMessage(idStore dto.GetChannelMessageDto, token string) (*structs.Message, error) {
	path := "/channels/" + idStore.ChannelID.ToString() + "/messages/" + idStore.MessageID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var message structs.Message
	err = json.Unmarshal(resp, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func CreateMessage(reqDto dto.CreateMessageDto, token string) (*structs.Message, error) {
	path := "/channels/" + reqDto.ChannelID.ToString() + "/messages"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	if len(reqDto.Files) > 0 {
		var reqBody bytes.Buffer
		writer := multipart.NewWriter(&reqBody)

		payloadJson, err := json.Marshal(reqDto)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormField("payload_json")
		if err != nil {
			return nil, err
		}
		part.Write(payloadJson)

		for key, fileContent := range reqDto.Files {
			attachmentID := key
			for _, attachment := range reqDto.Attachments {
				if attachment.ID.ToString() == attachmentID {
					part, err := writer.CreateFormFile(fmt.Sprintf("files[%s]", attachmentID), attachment.FileName)
					if err != nil {
						return nil, err
					}
					part.Write(fileContent)
					break
				}
			}
		}

		writer.Close()

		headers["Content-Type"] = writer.FormDataContentType()
		resp, err := HttpRequest("POST", path, headers, reqBody.Bytes())
		if err != nil {
			return nil, err
		}

		var message structs.Message
		err = json.Unmarshal(resp, &message)
		if err != nil {
			return nil, err
		}

		return &message, nil
	}

	body, err := json.Marshal(reqDto)
	if err != nil {
		return nil, err
	}

	headers["Content-Type"] = "application/json"
	resp, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	var message structs.Message
	err = json.Unmarshal(resp, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func CrossPostChannelMessage(idStore dto.GetChannelMessageDto, token string) (*structs.Message, error) {
	path := "/channels/" + idStore.ChannelID.ToString() + "/messages/" + idStore.MessageID.ToString() + "/crosspost"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("POST", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var message structs.Message
	err = json.Unmarshal(resp, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func CreateChannelMessageReaction(reaction dto.CreateReactionDto, token string) error {
	escapedEmoji := util.EncodeStructToURL(reaction.Emoji)
	if escapedEmoji == "" {
		return errors.New("failed to encode emoji")
	}
	path := "/channels/" + reaction.ChannelID.ToString() + "/messages/" + reaction.MessageID.ToString() + "/reactions/" + escapedEmoji + "/@me"

	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("PUT", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMyChannelMessageReaction(reaction dto.CreateReactionDto, token string) error {
	escapedEmoji := util.EncodeStructToURL(reaction.Emoji)
	if escapedEmoji == "" {
		return errors.New("failed to encode emoji")
	}
	path := "/channels/" + reaction.ChannelID.ToString() + "/messages/" + reaction.MessageID.ToString() + "/reactions/" + escapedEmoji + "/@me"

	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserChannelMessageReaction(reaction dto.DeleteUserReactionDto, token string) error {
	escapedEmoji := util.EncodeStructToURL(reaction.Emoji)
	if escapedEmoji == "" {
		return errors.New("failed to encode emoji")
	}
	path := "/channels/" + reaction.ChannelID.ToString() + "/messages/" + reaction.MessageID.ToString() + "/reactions/" + escapedEmoji + "/" + reaction.UserID.ToString()

	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetChannelMessageReactions(store dto.GetReactionsDto, token string) ([]structs.User, error) {
	escapedEmoji := util.EncodeStructToURL(store.Emoji)
	if escapedEmoji == "" {
		return nil, errors.New("failed to encode emoji")
	}
	path := "/channels/" + store.ChannelID.ToString() + "/messages/" + store.MessageID.ToString() + "/reactions/" + escapedEmoji

	query := util.BuildQueryString(store)
	if query != "" {
		path += query
	}

	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var users []structs.User
	err = json.Unmarshal(resp, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteAllChannelMessageReactions(idStore dto.GetChannelMessageDto, token string) error {
	path := "/channels/" + idStore.ChannelID.ToString() + "/messages/" + idStore.MessageID.ToString() + "/reactions"

	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAllChannelMessageReactionsForEmoji(store dto.CreateReactionDto, token string) error {
	escapedEmoji := util.EncodeStructToURL(store.Emoji)
	if escapedEmoji == "" {
		return errors.New("failed to encode emoji")
	}
	path := "/channels/" + store.ChannelID.ToString() + "/messages/" + store.MessageID.ToString() + "/reactions/" + escapedEmoji

	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func EditChannelMessage(messageDto dto.EditMessageDto, token string) (*structs.Message, error) {
	path := "/channels/" + messageDto.ChannelID.ToString() + "/messages/" + messageDto.MessageID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(messageDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return nil, err
	}

	var message structs.Message
	err = json.Unmarshal(resp, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func DeleteChannelMessage(idStore dto.GetChannelMessageDto, token string) error {
	path := "/channels/" + idStore.ChannelID.ToString() + "/messages/" + idStore.MessageID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func BulkDeleteChannelMessages(deleteDto dto.BulkDeleteMessagesDto, token string) error {
	path := "/channels/" + deleteDto.ChannelID.ToString() + "/messages/bulk-delete"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(deleteDto)
	if err != nil {
		return err
	}

	_, err = HttpRequest("POST", path, headers, body)
	if err != nil {
		return err
	}

	return nil
}
