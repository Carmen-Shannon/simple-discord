package request_util

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func CreateGuild(postDto dto.CreateGuildDto, token string) (*structs.Guild, error) {
	path := "/guilds"
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

	var guild structs.Guild
	if err := json.Unmarshal(resp, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func GetGuild(getDto dto.GetGuildDto, token string) (*structs.Guild, error) {
	path := "/guilds/" + getDto.GuildID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	query := path + util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", query, headers, nil)
	if err != nil {
		return nil, err
	}

	var guild structs.Guild
	if err := json.Unmarshal(resp, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func GetGuildPreview(getDto dto.GetGuildPreviewDto, token string) (*structs.GuildPreview, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/preview"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var guildPreview structs.GuildPreview
	if err := json.Unmarshal(resp, &guildPreview); err != nil {
		return nil, err
	}

	return &guildPreview, nil
}

func ModifyGuild(patchDto dto.ModifyGuildDto, token string) (*structs.Guild, error) {
	path := "/guilds/" + patchDto.GuildID.ToString()
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

	var guild structs.Guild
	if err := json.Unmarshal(resp, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func DeleteGuild(deleteDto dto.GetGuildPreviewDto, token string) error {
	path := "/guilds/" + deleteDto.GuildID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetGuildChannels(getDto dto.GetGuildPreviewDto, token string) ([]structs.Channel, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/channels"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var channels []structs.Channel
	if err := json.Unmarshal(resp, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

func CreateGuildChannel(postDto dto.CreateGuildChannelDto, token string) (*structs.Channel, error) {
	path := "/guilds/" + postDto.GuildID.ToString() + "/channels"
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
	if err := json.Unmarshal(resp, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func ModifyGuildChannelPositions(patchDto dto.ModifyGuildChannelPositionsDto, token string) error {
	path := "/guilds/" + patchDto.GuildID.ToString() + "/channels"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(patchDto)
	if err != nil {
		return err
	}

	_, err = HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return err
	}

	return nil
}

type listActiveGuildThreadsResponseWrapper struct {
	Threads []structs.Channel      `json:"threads"`
	Members []structs.ThreadMember `json:"members"`
}

func ListActiveGuildThreads(getDto dto.GetGuildPreviewDto, token string) ([]structs.Channel, []structs.ThreadMember, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/threads/active"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, nil, err
	}

	var wrapper listActiveGuildThreadsResponseWrapper
	if err := json.Unmarshal(resp, &wrapper); err != nil {
		return nil, nil, err
	}

	return wrapper.Threads, wrapper.Members, nil
}

func GetGuildMember(getDto dto.GetGuildMemberDto, token string) (*structs.GuildMember, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/members/" + getDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var member structs.GuildMember
	if err := json.Unmarshal(resp, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func ListGuildMembers(getDto dto.ListGuildMembersDto, token string) ([]structs.GuildMember, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/members"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	query := path + util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", query, headers, nil)
	if err != nil {
		return nil, err
	}

	var members []structs.GuildMember
	if err := json.Unmarshal(resp, &members); err != nil {
		return nil, err
	}

	return members, nil
}

func SearchGuildMembers(getDto dto.SearchGuildMembersDto, token string) ([]structs.GuildMember, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/members/search"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	query := path + util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", query, headers, nil)
	if err != nil {
		return nil, err
	}

	var members []structs.GuildMember
	if err := json.Unmarshal(resp, &members); err != nil {
		return nil, err
	}

	return members, nil
}

func AddGuildMember(putDto dto.AddGuildMemberDto, token string) (*structs.GuildMember, error) {
	path := "/guilds/" + putDto.GuildID.ToString() + "/members/" + putDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(putDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PUT", path, headers, body)
	if err != nil {
		return nil, err
	}

	var member structs.GuildMember
	if err := json.Unmarshal(resp, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func ModifyGuildMember(patchDto dto.ModifyGuildMemberDto, token string) (*structs.GuildMember, error) {
	path := "/guilds/" + patchDto.GuildID.ToString() + "/members/" + patchDto.UserID.ToString()
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

	var member structs.GuildMember
	if err := json.Unmarshal(resp, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func ModifyCurrentMember(patchDto dto.ModifyCurrentMemberDto, token string) (*structs.GuildMember, error) {
	path := "/guilds/" + patchDto.GuildID.ToString() + "/members/@me"
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

	var member structs.GuildMember
	if err := json.Unmarshal(resp, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func AddGuildMemberRole(putDto dto.AddGuildMemberRoleDto, token string) error {
	path := "/guilds/" + putDto.GuildID.ToString() + "/members/" + putDto.UserID.ToString() + "/roles/" + putDto.RoleID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("PUT", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func RemoveGuildMemberRole(deleteDto dto.AddGuildMemberRoleDto, token string) error {
	path := "/guilds/" + deleteDto.GuildID.ToString() + "/members/" + deleteDto.UserID.ToString() + "/roles/" + deleteDto.RoleID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func RemoveGuildMember(deleteDto dto.GetGuildMemberDto, token string) error {
	path := "/guilds/" + deleteDto.GuildID.ToString() + "/members/" + deleteDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetGuildBans(getDto dto.GetGuildBansDto, token string) ([]structs.Ban, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/bans"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	query := path + util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", query, headers, nil)
	if err != nil {
		return nil, err
	}

	var bans []structs.Ban
	if err := json.Unmarshal(resp, &bans); err != nil {
		return nil, err
	}

	return bans, nil
}

func GetGuildBan(getDto dto.GetGuildMemberDto, token string) (*structs.Ban, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/bans/" + getDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var ban structs.Ban
	if err := json.Unmarshal(resp, &ban); err != nil {
		return nil, err
	}

	return &ban, nil
}

func CreateGuildBan(putDto dto.CreateGuildBanDto, token string) error {
	path := "/guilds/" + putDto.GuildID.ToString() + "/bans/" + putDto.UserID.ToString()
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

func RemoveGuildBan(deleteDto dto.GetGuildMemberDto, token string) error {
	path := "/guilds/" + deleteDto.GuildID.ToString() + "/bans/" + deleteDto.UserID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

type bulkGuildBanResponseWrapper struct {
	BannedUsers []structs.Snowflake `json:"banned_users"`
	FailedUsers []structs.Snowflake `json:"failed_users"`
}

func BulkGuildBan(postDto dto.BulkGuildBanDto, token string) (bannedUsers []structs.Snowflake, failedUsers []structs.Snowflake, err error) {
	path := "/guilds/" + postDto.GuildID.ToString() + "/bans"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(postDto)
	if err != nil {
		return nil, nil, err
	}

	resp, err := HttpRequest("PUT", path, headers, body)
	if err != nil {
		return nil, nil, err
	}

	var wrapper bulkGuildBanResponseWrapper
	if err := json.Unmarshal(resp, &wrapper); err != nil {
		return nil, nil, err
	}

	return wrapper.BannedUsers, wrapper.FailedUsers, nil
}

func GetGuildRoles(getDto dto.GetGuildPreviewDto, token string) ([]structs.Role, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/roles"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var roles []structs.Role
	if err := json.Unmarshal(resp, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func GetGuildRole(getDto dto.GetGuildRoleDto, token string) (*[]structs.Role, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/roles/" + getDto.RoleID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var roles []structs.Role
	if err := json.Unmarshal(resp, &roles); err != nil {
		return nil, err
	}

	return &roles, nil
}

func CreateGuildRole(postDto dto.CreateGuildRoleDto, token string) (*structs.Role, error) {
	path := "/guilds/" + postDto.GuildID.ToString() + "/roles"
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

	var role structs.Role
	if err := json.Unmarshal(resp, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

func ModifyGuildRolePositions(patchDto dto.ModifyGuildRolePositionsDto, token string) error {
	path := "/guilds/" + patchDto.GuildID.ToString() + "/roles"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(patchDto)
	if err != nil {
		return err
	}

	_, err = HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return err
	}

	return nil
}

func ModifyGuildRole(patchDto dto.ModifyGuildRoleDto, token string) (*structs.Role, error) {
	path := "/guilds/" + patchDto.GuildID.ToString() + "/roles/" + patchDto.RoleID.ToString()
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

	var role structs.Role
	if err := json.Unmarshal(resp, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

func ModifyGuildMFALevel(postDto dto.ModifyGuildMFALevelDto, token string) (*structs.MFALevel, error) {
	path := "/guilds/" + postDto.GuildID.ToString() + "/mfa"
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

	var mfa structs.MFALevel
	if err := json.Unmarshal(resp, &mfa); err != nil {
		return nil, err
	}

	return &mfa, nil
}

func DeleteGuildRole(deleteDto dto.GetGuildRoleDto, token string) error {
	path := "/guilds/" + deleteDto.GuildID.ToString() + "/roles/" + deleteDto.RoleID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

type prunedWrapper struct {
	Pruned int `json:"pruned"`
}

func GetGuildPruneCount(getDto dto.GetGuildPruneCountDto, token string) (*int, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/prune"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	query := path + util.BuildQueryStringDelimitedSlices(getDto)
	resp, err := HttpRequest("GET", query, headers, nil)
	if err != nil {
		return nil, err
	}

	var count prunedWrapper
	if err := json.Unmarshal(resp, &count); err != nil {
		return nil, err
	}

	return &count.Pruned, nil
}

func BeginGuildPrune(postDto dto.BeginGuildPruneDto, token string) (*int, error) {
	path := "/guilds/" + postDto.GuildID.ToString() + "/prune"
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

	var count prunedWrapper
	if err := json.Unmarshal(resp, &count); err != nil {
		return nil, err
	}

	return &count.Pruned, nil
}

func GetGuildVoiceRegions(getDto dto.GetGuildPreviewDto, token string) ([]structs.VoiceRegion, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/regions"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var regions []structs.VoiceRegion
	if err := json.Unmarshal(resp, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

func GetGuildInvites(getDto dto.GetGuildPreviewDto, token string) ([]structs.Invite, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/invites"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var invites []structs.Invite
	if err := json.Unmarshal(resp, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

func GetGuildIntegrations(getDto dto.GetGuildPreviewDto, token string) ([]structs.GuildIntegration, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/integrations"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var integrations []structs.GuildIntegration
	if err := json.Unmarshal(resp, &integrations); err != nil {
		return nil, err
	}

	return integrations, nil
}

func DeleteGuildIntegration(deleteDto dto.DeleteGuildIntegrationDto, token string) error {
	path := "/guilds/" + deleteDto.GuildID.ToString() + "/integrations/" + deleteDto.IntegrationID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetGuildWidgetSettings(getDto dto.GetGuildPreviewDto, token string) (*structs.GuildWidgetSettings, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/widget"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var settings structs.GuildWidgetSettings
	if err := json.Unmarshal(resp, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func ModifyGuildWidget(patchDto dto.ModifyGuildWidgetDto, token string) (*structs.GuildWidgetSettings, error) {
	path := "/guilds/" + patchDto.GuildID.ToString() + "/widget"
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

	var settings structs.GuildWidgetSettings
	if err := json.Unmarshal(resp, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func GetGuildWidget(getDto dto.GetGuildPreviewDto, token string) (*structs.GuildWidget, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/widget.json"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var widget structs.GuildWidget
	if err := json.Unmarshal(resp, &widget); err != nil {
		return nil, err
	}

	return &widget, nil
}

func GetGuildVanityURL(getDto dto.GetGuildPreviewDto, token string) (*structs.Invite, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/vanity-url"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var url structs.Invite
	if err := json.Unmarshal(resp, &url); err != nil {
		return nil, err
	}

	return &url, nil
}

func GetGuildWidgetImage(getDto dto.GetGuildWidgetImageDto, token string) (*string, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/widget.png"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	path += util.BuildQueryString(getDto)
	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var image string
	if err := json.Unmarshal(resp, &image); err != nil {
		return nil, err
	}

	return &image, nil
}

func GetGuildWelcomeScreen(getDto dto.GetGuildPreviewDto, token string) (*structs.WelcomeScreen, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/welcome-screen"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var screen structs.WelcomeScreen
	if err := json.Unmarshal(resp, &screen); err != nil {
		return nil, err
	}

	return &screen, nil
}

func ModifyGuildWelcomeScreen(patchDto dto.ModifyGuildWelcomeScreenDto, token string) (*structs.WelcomeScreen, error) {
	path := "/guilds/" + patchDto.GuildID.ToString() + "/welcome-screen"
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

	var screen structs.WelcomeScreen
	if err := json.Unmarshal(resp, &screen); err != nil {
		return nil, err
	}

	return &screen, nil
}

func GetGuildOnboarding(getDto dto.GetGuildPreviewDto, token string) (*structs.GuildOnboarding, error) {
	path := "/guilds/" + getDto.GuildID.ToString() + "/onboarding"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var onboarding structs.GuildOnboarding
	if err := json.Unmarshal(resp, &onboarding); err != nil {
		return nil, err
	}

	return &onboarding, nil
}

func ModifyGuildOnboarding(putDto dto.ModifyGuildOnboardingDto, token string) (*structs.GuildOnboarding, error) {
	path := "/guilds/" + putDto.GuildID.ToString() + "/onboarding"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(putDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PUT", path, headers, body)
	if err != nil {
		return nil, err
	}

	var onboarding structs.GuildOnboarding
	if err := json.Unmarshal(resp, &onboarding); err != nil {
		return nil, err
	}

	return &onboarding, nil
}
