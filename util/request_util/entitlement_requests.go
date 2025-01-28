package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func ListEntitlements(getDto dto.ListEntitlementsDto, token string) ([]structs.Entitlement, error) {
	path := "/applications/" + getDto.ApplicationID.ToString() + "/entitlements"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	query := path + util.BuildQueryStringDelimitedSlices(getDto)
	resp, err := HttpRequest("GET", query, headers, nil)
	if err != nil {
		return nil, err
	}

	var entitlements []structs.Entitlement
	err = json.Unmarshal(resp, &entitlements)
	if err != nil {
		return nil, err
	}

	return entitlements, nil
}

func GetEntitlement(getDto dto.GetEntitlementDto, token string) (*structs.Entitlement, error) {
	path := "/applications/" + getDto.ApplicationID.ToString() + "/entitlements/" + getDto.EntitlementID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	resp, err := HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var entitlement structs.Entitlement
	err = json.Unmarshal(resp, &entitlement)
	if err != nil {
		return nil, err
	}

	return &entitlement, nil
}

func ConsumeEntitlement(postDto dto.GetEntitlementDto, token string) error {
	path := "/applications/" + postDto.ApplicationID.ToString() + "/entitlements/" + postDto.EntitlementID.ToString() + "/consume"
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("POST", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func CreateTestEntitlement(postDto dto.CreateTestEntitlementDto, token string) (*structs.Entitlement, error) {
	path := "/applications/" + postDto.ApplicationID.ToString() + "/entitlements"
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

	var entitlement structs.Entitlement
	err = json.Unmarshal(resp, &entitlement)
	if err != nil {
		return nil, err
	}

	return &entitlement, nil
}

func DeleteTestEntitlement(deleteDto dto.GetEntitlementDto, token string) error {
	path := "/applications/" + deleteDto.ApplicationID.ToString() + "/entitlements/" + deleteDto.EntitlementID.ToString()
	headers := map[string]string{
		"Authorization": "Bot " + token,
	}

	_, err := HttpRequest("DELETE", path, headers, nil)
	if err != nil {
		return err
	}

	return nil
}
