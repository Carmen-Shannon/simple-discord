package receiveevents

import (
	"encoding/json"
	"fmt"

	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
)

func NewVoiceReceiveEvent(eventData payload.VoicePayload) (any, error) {
	jsonData, err := json.Marshal(eventData.Data)
	if err != nil {
		return nil, err
	}

	switch eventData.OpCode {
	case gateway.VoiceOpReady:
		var event VoiceReadyEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpSessionDescription:
		var event VoiceSessionDescriptionEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpSpeaking:
		var event SpeakingEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpHeartbeatAck:
		var event HeartbeatACKEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpHello:
		var event HelloEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpResumed:
		var event VoiceResumedEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpClientsConnect:
		var event VoiceClientsConnectEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpClientDisconnect:
		var event VoiceClientDisconnectEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpPrepareTransition:
		var event VoicePrepareTransitionEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpExecuteTransition:
		// var event VoiceExecuteTransitionEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("ExecuteTransition event\n", string(event))

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpPrepareEpoch:
		var event VoicePrepareEpochEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpMLSAnnounceCommitTransition:
		// var event VoiceMlsAnnounceCommitTransitionEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("MLSAnnounceCommitTransition event\n", string(event))

		eventData.Data = event
		return event, nil
	case gateway.VoiceOpMLSInvalidCommitWelcome:
		// var event VoiceMlsInvalidCommitWelcomeEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("MLSInvalidCommitWelcome event\n", string(event))

		eventData.Data = event
		return event, nil
	default:
		var event json.RawMessage
		_ = json.Unmarshal(jsonData, &event)
		return nil, fmt.Errorf("unknown opcode for voice receive event: %d", eventData.OpCode)
	}
}
