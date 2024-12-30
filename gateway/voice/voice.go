package voice

import (
	"encoding/json"
	"fmt"

	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/structs/voice"
)

func NewSendEvent(eventData voice.VoicePayload) (any, error) {
	jsonData, err := json.Marshal(eventData.Data)
	if err != nil {
		return nil, err
	}

	switch eventData.OpCode {
	case voice.Identify:
		var event sendevents.VoiceIdentifyEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.SelectProtocol:
		var event sendevents.VoiceSelectProtocolEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.Heartbeat:
		var event sendevents.HeartbeatEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.Speaking:
		var event sendevents.SpeakingEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.Resume:
		var event sendevents.VoiceResumeEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.TransitionReady:
		var event sendevents.VoiceDaveReadyForTransitionEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.MLSInvalidCommitWelcome:
		var event sendevents.VoiceDaveReadyForTransitionEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	default:
		return nil, fmt.Errorf("unknown opcode for voice send event: %d", eventData.OpCode)
	}
}

func NewBinarySendEvent(eventData voice.BinaryVoicePayload) (any, error) {
	binaryData, err := eventData.MarshalBinary()
	if err != nil {
		return nil, err
	}

	switch eventData.OpCode {
	case uint8(voice.MLSKeyPackage):
		var event sendevents.VoiceDaveMlsKeyPackageEvent
		if err = event.UnmarshalBinary(binaryData); err != nil {
			return nil, err
		}

		binaryPayload, err := event.MarshalBinary()
		if err != nil {
			return nil, err
		}

		eventData.Payload = binaryPayload
		return event, nil
	case uint8(voice.MLSCommitWelcome):
		var event sendevents.VoiceDaveMlsCommitWelcomeEvent
		if err = event.UnmarshalBinary(binaryData); err != nil {
			return nil, err
		}

		binaryPayload, err := event.MarshalBinary()
		if err != nil {
			return nil, err
		}

		eventData.Payload = binaryPayload
		return event, nil
	default:
		return nil, fmt.Errorf("unknown opcode for voice binary send event: %d", eventData.OpCode)
	}
}

func NewReceiveEvent(eventData voice.VoicePayload) (any, error) {
	jsonData, err := json.Marshal(eventData.Data)
	if err != nil {
		return nil, err
	}
	// Add logic to parse the jsonData into the appropriate receive event type
	// based on the OpCode
	return jsonData, nil
}
