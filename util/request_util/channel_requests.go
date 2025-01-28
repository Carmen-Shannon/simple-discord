package requestutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
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

func EditChannelPermissions(editPermissionsDto dto.EditChannelPermissionsDto, token string) error {
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

func GetChannelInvites(getDto dto.GetChannelDto, token string) ([]structs.Invite, error) {
	path := "/channels/" + getDto.ChannelID.ToString() + "/invites"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var invites []structs.Invite
	err = json.Unmarshal(resp, &invites)
	if err != nil {
		return nil, err
	}

	return invites, nil
}

func CreateChannelInvite(createDto dto.CreateChannelInviteDto, token string) (*structs.Invite, error) {
	path := "/channels/" + createDto.ChannelID.ToString() + "/invites"
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

	var invite structs.Invite
	err = json.Unmarshal(resp, &invite)
	if err != nil {
		return nil, err
	}

	return &invite, nil
}

func DeleteChannelPermission(deleteDto dto.DeleteChannelPermissionDto, token string) error {
	path := "/channels/" + deleteDto.ChannelID.ToString() + "/permissions/" + deleteDto.OverwriteID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func FollowAnnouncementChannel(postDto dto.FollowAnnouncementChannelDto, token string) (*structs.FollowedChannel, error) {
	path := "/channels/" + postDto.ChannelID.ToString() + "/followers"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(postDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	var followedChannel structs.FollowedChannel
	err = json.Unmarshal(resp, &followedChannel)
	if err != nil {
		return nil, err
	}

	return &followedChannel, nil
}

func TriggerTypingIndicator(postDto dto.TriggerTypingIndicatorDto, token string) error {
	path := "/channels/" + postDto.ChannelID.ToString() + "/typing"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("POST", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetPinnedMessages(getDto dto.GetChannelDto, token string) ([]structs.Message, error) {
	path := "/channels/" + getDto.ChannelID.ToString() + "/pins"
	headers := map[string]string{
		"Authorization": "Bot " + token,
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

func PinMessage(putDto dto.PinMessageDto, token string) error {
	path := "/channels/" + putDto.ChannelID.ToString() + "/pins/" + putDto.MessageID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("PUT", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func UnpinMessage(deleteDto dto.PinMessageDto, token string) error {
	path := "/channels/" + deleteDto.ChannelID.ToString() + "/pins/" + deleteDto.MessageID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GroupDMAddRecipient(putDto dto.GroupDMAddRecipientDto, token string) error {
	path := "/channels/" + putDto.ChannelID.ToString() + "/recipients/" + putDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	body, err := json.Marshal(putDto)
	if err != nil {
		return err
	}

	_, err = HttpRequest("PUT", path, headers, body)
	if err != nil {
		return err
	}

	return nil
}

func GroupDMRemoveRecipient(deleteDto dto.GroupDMRemoveRecipientDto, token string) error {
	path := "/channels/" + deleteDto.ChannelID.ToString() + "/recipients/" + deleteDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func StartThreadFromMessage(postDto dto.StartThreadFromMessageDto, token string) (*structs.Channel, error) {
	path := "/channels/" + postDto.ChannelID.ToString() + "/messages/" + postDto.MessageID.ToString() + "/threads"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(postDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
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

func StartThreadWithoutMessage(postDto dto.StartThreadWithoutMessageDto, token string) (*structs.Channel, error) {
	path := "/channels/" + postDto.ChannelID.ToString() + "/threads"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(postDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
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

func StartThreadInForumOrMediaChannel(postDto dto.StartThreadInForumOrMediaChannelDto, token string) (*structs.Channel, error) {
	path := "/channels/" + postDto.ChannelID.ToString() + "/threads"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	// if request includes File uploads, use multipart/form-data and marshal the dto into a payload_json field
	if len(postDto.Files) > 0 {
		var reqBody bytes.Buffer
		writer := multipart.NewWriter(&reqBody)

		payloadJson, err := json.Marshal(postDto)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormField("payload_json")
		if err != nil {
			return nil, err
		}
		part.Write(payloadJson)

		for key, fileContent := range postDto.Files {
			attachmentID := key
			for _, attachment := range postDto.Attachments {
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

		var channel structs.Channel
		err = json.Unmarshal(resp, &channel)
		if err != nil {
			return nil, err
		}

		return &channel, nil
	}

	body, err := json.Marshal(postDto)
	if err != nil {
		return nil, err
	}

	headers["Content-Type"] = "application/json"
	resp, err := HttpRequest("POST", path, headers, body)
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

func JoinThread(putDto dto.GetChannelDto, token string) error {
	path := "/channels/" + putDto.ChannelID.ToString() + "/thread-members/@me"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("PUT", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func AddThreadMember(putDto dto.GroupDMRemoveRecipientDto, token string) error {
	path := "/channels/" + putDto.ChannelID.ToString() + "/thread-members/" + putDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("PUT", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func LeaveThread(deleteDto dto.GetChannelDto, token string) error {
	path := "/channels/" + deleteDto.ChannelID.ToString() + "/thread-members/@me"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func RemoveThreadMember(deleteDto dto.GroupDMRemoveRecipientDto, token string) error {
	path := "/channels/" + deleteDto.ChannelID.ToString() + "/thread-members/" + deleteDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetThreadMember(getDto dto.GetThreadMemberDto, token string) (*structs.ThreadMember, error) {
	path := "/channels/" + getDto.ChannelID.ToString() + "/thread-members/" + getDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var threadMember structs.ThreadMember
	err = json.Unmarshal(resp, &threadMember)
	if err != nil {
		return nil, err
	}

	return &threadMember, nil
}

func ListThreadMembers(getDto dto.ListThreadMembersDto, token string) ([]structs.ThreadMember, error) {
	path := "/channels/" + getDto.ChannelID.ToString() + "/thread-members"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var threadMembers []structs.ThreadMember
	err = json.Unmarshal(resp, &threadMembers)
	if err != nil {
		return nil, err
	}

	return threadMembers, nil
}

func ListPublicArchivedThreads(getDto dto.ListPublicArchivedThreadsDto, token string) (*ListPublicArchivedThreadsResponse, error) {
	path := "/channels/" + getDto.ChannelID.ToString() + "/threads/archived/public"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var response ListPublicArchivedThreadsResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type ListPublicArchivedThreadsResponse struct {
	Threads []structs.Channel      `json:"threads"`
	Members []structs.ThreadMember `json:"members"`
	HasMore bool                   `json:"has_more"`
}

func ListPrivateArchivedThreads(getDto dto.ListPublicArchivedThreadsDto, token string) (*ListPublicArchivedThreadsResponse, error) {
	path := "/channels/" + getDto.ChannelID.ToString() + "/threads/archived/private"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var response ListPublicArchivedThreadsResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func ListJoinedPrivateArchivedThreads(getDto dto.ListPublicArchivedThreadsDto, token string) (*ListPublicArchivedThreadsResponse, error) {
	path := "/channels/" + getDto.ChannelID.ToString() + "/users/@me/threads/archived/private"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var response ListPublicArchivedThreadsResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
