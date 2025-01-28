package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
)

func GetApplicationRoleConnectionMetadataRecords(getDto dto.GetApplicationRoleConnectionMetadataRecordsDto, token string) ([]structs.ApplicationRoleConnectionMetadata, error) {
	path := "/applications/" + getDto.ApplicationID.ToString() + "/role-connections/metadata"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var records []structs.ApplicationRoleConnectionMetadata
	err = json.Unmarshal(resp, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func UpdateApplicationRoleConnectionMetadataRecords(updateDto dto.UpdateApplicationRoleConnectionMetadataRecordsDto, token string) ([]structs.ApplicationRoleConnectionMetadata, error) {
	path := "/applications/" + updateDto.ApplicationID.ToString() + "/role-connections/metadata"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	body, err := json.Marshal(updateDto.Records)
	if err != nil {
		return nil, err
	}

	resp, err := HttpRequest("PUT", path, headers, body)
	if err != nil {
		return nil, err
	}

	var records []structs.ApplicationRoleConnectionMetadata
	err = json.Unmarshal(resp, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}
