package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
)

func GetAutoModerationRule(getAutoModRuleDto dto.GetAutoModerationRuleDto, token string) (*structs.AutoModerationRule, error) {
	path := "/guilds/" + getAutoModRuleDto.GuildID.ToString() + "/auto-moderation/rules/" + getAutoModRuleDto.RuleID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var autoModerationRule structs.AutoModerationRule
	err = json.Unmarshal(resp, &autoModerationRule)
	if err != nil {
		return nil, err
	}

	return &autoModerationRule, nil
}

func CreateAutoModerationRule(createAutoModRuleDto dto.CreateAutoModerationRuleDto, token string) (*structs.AutoModerationRule, error) {
	path := "/guilds/" + createAutoModRuleDto.GuildID.ToString() + "/auto-moderation/rules"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(createAutoModRuleDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	var autoModerationRule structs.AutoModerationRule
	err = json.Unmarshal(resp, &autoModerationRule)
	if err != nil {
		return nil, err
	}

	return &autoModerationRule, nil
}

func ModifyAutoModerationRule(modifyAutoModRuleDto dto.ModifyAutoModerationRuleDto, token string) (*structs.AutoModerationRule, error) {
	path := "/guilds/" + modifyAutoModRuleDto.GuildID.ToString() + "/auto-moderation/rules/" + modifyAutoModRuleDto.RuleID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(modifyAutoModRuleDto)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PATCH", path, headers, body)
	if err != nil {
		return nil, err
	}

	var autoModerationRule structs.AutoModerationRule
	err = json.Unmarshal(resp, &autoModerationRule)
	if err != nil {
		return nil, err
	}

	return &autoModerationRule, nil
}

func DeleteAutoModerationRule(deleteAutoModRuleDto dto.GetAutoModerationRuleDto, token string) error {
	path := "/guilds/" + deleteAutoModRuleDto.GuildID.ToString() + "/auto-moderation/rules/" + deleteAutoModRuleDto.RuleID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}
