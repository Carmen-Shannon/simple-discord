package requestutil

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func CreateInteractionResponse(interactionID, interactionToken, token string, dto dto.CreateInteractionResponseDto, response structs.InteractionResponse) (*structs.InteractionCallbackResponse, error) {
	path := "/interactions/" + interactionID + "/" + interactionToken + "/callback"
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"Content-Type":  "application/json",
	}

	path += util.BuildQueryString(dto)
	body, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	var interactionResponse structs.InteractionCallbackResponse
	res, err := HttpRequest("POST", path, headers, body)
	if err != nil {
		return nil, err
	}

	if dto.WithResponse != nil && *dto.WithResponse {
		err = json.Unmarshal(res, &interactionResponse)
		if err != nil {
			return nil, err
		}
		return &interactionResponse, nil
	}

	return nil, err
}
