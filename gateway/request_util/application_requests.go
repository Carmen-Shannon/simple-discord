package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
)

func GetCurrentApplication(token string) (*structs.Application, error) {
	path := "/applications/@me"
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
	path := "/applications/@me"
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

func GetApplicationActivityInstance(dto dto.GetApplicationActivityInstanceDto, token string) (*structs.ActivityInstance, error) {
	path := "/applications/" + dto.ApplicationID.ToString() + "/activity-instances/" + dto.InstanceID
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var activityInstance structs.ActivityInstance
	err = json.Unmarshal(resp, &activityInstance)
	if err != nil {
		return nil, err
	}

	return &activityInstance, nil
}
