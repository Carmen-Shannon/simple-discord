package request_util

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func GetGuildAuditLog(auditLogParams dto.GetGuildAuditLogDto, token string) (*structs.AuditLog, error) {
	path := "/guilds/" + auditLogParams.GuildID.ToString() + "/audit-logs"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	query := util.BuildQueryString(auditLogParams)
	path += query

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var auditLog structs.AuditLog
	err = json.Unmarshal(resp, &auditLog)
	if err != nil {
		return nil, err
	}

	return &auditLog, nil
}
