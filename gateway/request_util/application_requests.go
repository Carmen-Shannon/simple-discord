package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
)

func GetCurrentApplication(token string) (*structs.Application, error) {
	path := "/oauth2/applications/@me"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var application structs.Application
	err = json.Unmarshal(resp, &application)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

func EditCurrentApplication(updates dto.EditCurrentApplicationDto, token string) (*structs.Application, error) {
	path := "/oauth2/applications/@me"
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

	var application structs.Application
	err = json.Unmarshal(resp, &application)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

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

// THIS IS NOT IMPLEMENTED YET, NO DETAILS IN DOCUMENTATION FOR HOW TO SEND THE PUT REQUEST
func UpdateApplicationRoleConnectionMetadataRecords(token string) ([]structs.ApplicationRoleConnectionMetadata, error) {
	return nil, nil
}