package voice

import (
	"encoding/json"
	"fmt"

	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
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

		eventData.Data = binaryPayload
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

		eventData.Data = binaryPayload
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

	switch eventData.OpCode {
	case voice.Ready:
		var event receiveevents.VoiceReadyEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.SessionDescription:
		var event receiveevents.VoiceSessionDescriptionEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.Speaking:
		var event receiveevents.SpeakingEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.HeartbeatAck:
		var event receiveevents.HeartbeatACKEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.Hello:
		var event receiveevents.HelloEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.Resumed:
		var event receiveevents.VoiceResumedEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.ClientsConnect:
		var event receiveevents.VoiceClientsConnectEvent
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case voice.ClientDisconnect:
		// var event receiveevents.VoiceClientDisconnectEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("ClientDisconnect event\n", string(event))

		eventData.Data = event
		return event, nil
	case voice.PrepareTransition:
		// var event receiveevents.VoicePrepareTransitionEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("PrepareTransition event\n", string(event))

		eventData.Data = event
		return event, nil
	case voice.ExecuteTransition:
		// var event receiveevents.VoiceExecuteTransitionEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("ExecuteTransition event\n", string(event))

		eventData.Data = event
		return event, nil
	case voice.PrepareEpoch:
		// var event receiveevents.VoicePrepareEpochEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("PrepareEpoch event\n", string(event))

		eventData.Data = event
		return event, nil
	case voice.MLSAnnounceCommitTransition:
		// var event receiveevents.VoiceMlsAnnounceCommitTransitionEvent
		var event json.RawMessage
		if err = json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		fmt.Println("MLSAnnounceCommitTransition event\n", string(event))

		eventData.Data = event
		return event, nil
	case voice.MLSInvalidCommitWelcome:
		// var event receiveevents.VoiceMlsInvalidCommitWelcomeEvent
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

func NewBinaryReceiveEvent(eventData voice.BinaryVoicePayload) (any, error) {
	_, err := eventData.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// TODO: flesh this function out with the binary receive events

	// switch eventData.OpCode {
	// case uint8(voice.MLSExternalSender):
	// 	var event receiveevents.VoiceMlsExternalSenderEvent
	// 	if err = event.UnmarshalBinary(binaryData); err != nil {
	// 		return nil, err
	// 	}
	// }

	return nil, nil
}
