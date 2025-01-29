package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	receiveevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func handleSendVoiceIdentifyEvent(s VoiceSession, p payload.VoicePayload) error {
	// no DAVE support yet, include DaveProtocolVersion
	voiceIdentifyEvent := sendevents.VoiceIdentifyEvent{
		ServerID:               *s.GetGuildID(),
		UserID:                 s.GetBotData().UserDetails.ID,
		SessionID:              *s.GetSessionID(),
		Token:                  *s.GetToken(),
		MaxDaveProtocolVersion: util.ToPtr(0),
	}
	ackPayload := payload.VoicePayload{
		OpCode: gateway.VoiceOpIdentify,
		Data:   voiceIdentifyEvent,
		Seq:    s.GetSequence(),
	}
	data, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}
	s.Write(data, false)
	return nil
}

func handleSendVoiceSelectProtocolEvent(s VoiceSession, p payload.VoicePayload) error {
	selectProtocolEvent := sendevents.VoiceSelectProtocolEvent{
		Protocol: "udp",
		Data: sendevents.VoiceSelectProtocolData{
			Address: s.GetAudioPlayer().GetSession().GetUdpData().Address,
			Port:    s.GetAudioPlayer().GetSession().GetUdpData().Port,
			Mode:    s.GetAudioPlayer().GetSession().GetUdpData().Mode,
		},
		Codecs:              []structs.Codec{gateway.Opus},
		DaveProtocolVersion: util.ToPtr(0),
	}
	selectProtocolPayload := payload.VoicePayload{
		OpCode: gateway.VoiceOpSelectProtocol,
		Data:   selectProtocolEvent,
		Seq:    s.GetSequence(),
	}

	data, err := json.Marshal(selectProtocolPayload)
	if err != nil {
		return err
	}

	s.Write(data, false)
	return nil
}

func handleVoiceReadyEvent(s VoiceSession, p payload.VoicePayload) error {
	if voiceReadyEvent, ok := p.Data.(receiveevents.VoiceReadyEvent); ok {
		udpConn := &gateway.UdpData{
			Address: voiceReadyEvent.IP,
			Port:    voiceReadyEvent.Port,
			SSRC:    voiceReadyEvent.SSRC,
		}
		if util.SliceContains(voiceReadyEvent.Modes, gateway.AEAD_AES256_GCM) {
			udpConn.Mode = gateway.AEAD_AES256_GCM
		} else if len(voiceReadyEvent.Modes) == 1 {
			udpConn.Mode = voiceReadyEvent.Modes[0]
		} else {
			udpConn.Mode = gateway.AEAD_XCHACHA20_POLY1305
		}

		s.GetAudioPlayer().GetSession().SetUdpData(*udpConn)
		s.GetAudioPlayer().GetSession().CloseConnectReady()
		s.CloseReadyReceived()
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceSendHeartbeatEvent(s VoiceSession, p payload.VoicePayload) error {
	if heartbeatEvent, ok := p.Data.(sendevents.VoiceHeartbeatEvent); ok {
		if heartbeatEvent.SeqAck != nil {
			s.SetSequence(*heartbeatEvent.SeqAck)
		}
		return sendVoiceHeartbeatEvent(s)
	}
	return errors.New("unexpected payload data type")
}

func handleVoiceSessionDescriptionEvent(s VoiceSession, p payload.VoicePayload) error {
	if voiceSessionDescriptionEvent, ok := p.Data.(receiveevents.VoiceSessionDescriptionEvent); ok {
		s.GetAudioPlayer().GetSession().SetSecretKey(voiceSessionDescriptionEvent.SecretKey)
		s.GetAudioPlayer().GetSession().SetEncryption(voiceSessionDescriptionEvent.Mode)
		s.GetAudioPlayer().GetSession().CloseSpeakingReady()
		return nil
	} else {
		return errors.New("unexpected payload data type")
	}
}

func handleVoiceSpeakingEvent(s VoiceSession, p payload.VoicePayload) error {
	if _, ok := p.Data.(receiveevents.SpeakingEvent); ok {
		// do nothing with received speaking events for now
		return nil
	}
	return errors.New("unexpected payload data type")
}

func handleVoiceHeartbeatAckEvent(s VoiceSession, p payload.VoicePayload) error {
	return nil
}

func handleSendVoiceResumeEvent(s VoiceSession, p payload.VoicePayload) error {
	if _, ok := p.Data.(sendevents.VoiceResumeEvent); ok {
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}

		s.Write(data, false)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceHelloEvent(s VoiceSession, p payload.VoicePayload) error {
	if helloEvent, ok := p.Data.(receiveevents.HelloEvent); ok {
		s.SetHeartbeatAck(int(helloEvent.HeartbeatInterval))
	} else {
		return errors.New("unexpected payload data type")
	}
	return startVoiceHeartbeatTimer(s)
}

func handleVoiceResumedEvent(s VoiceSession, p payload.VoicePayload) error {
	if _, ok := p.Data.(receiveevents.VoiceResumedEvent); ok {
		s.CloseResumeReady()
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceClientsConnectEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE CLIENTS CONNECT EVENT")
	if _, ok := p.Data.(receiveevents.VoiceClientsConnectEvent); ok {
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceClientDisconnectEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE CLIENT DISCONNECT EVENT")
	fmt.Println("VOICE CLIENT DISCONNECT NOT IMPLEMENTED")
	return nil
}

func handleVoicePrepareTransitionEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE PREPARE TRANSITION EVENT")
	fmt.Println("VOICE PREPARE TRANSITION NOT IMPLEMENTED")
	return nil
}

func handleVoiceExecuteTransitionEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE EXECUTE TRANSITION EVENT")
	fmt.Println("VOICE EXECUTE TRANSITION NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceTransitionReadyEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE TRANSITION READY EVENT")
	fmt.Println("VOICE TRANSITION READY NOT IMPLEMENTED")
	return nil
}

func handleVoicePrepareEpochEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE PREPARE EPOCH EVENT")
	fmt.Println("VOICE PREPARE EPOCH NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSExternalSenderEvent(s VoiceSession, p payload.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS EXTERNAL SENDER EVENT")
	fmt.Println("VOICE MLS EXTERNAL SENDER NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceMLSKeyPackageEvent(s VoiceSession, p payload.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS KEY PACKAGE EVENT")
	fmt.Println("VOICE MLS KEY PACKAGE NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSProposalsEvent(s VoiceSession, p payload.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS PROPOSALS EVENT")
	fmt.Println("VOICE MLS PROPOSALS NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceMLSCommitWelcomeEvent(s VoiceSession, p payload.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS COMMIT WELCOME EVENT")
	fmt.Println("VOICE MLS COMMIT WELCOME NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSAnnounceCommitTransitionEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE MLS ANNOUNCE COMMIT TRANSITION EVENT")
	fmt.Println("VOICE MLS ANNOUNCE COMMIT TRANSITION NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSWelcomeEvent(s VoiceSession, p payload.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS WELCOME EVENT")
	fmt.Println("VOICE MLS WELCOME NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceMLSInvalidCommitWelcomeEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Println("HANDLING VOICE MLS INVALID COMMIT WELCOME EVENT")
	fmt.Println("VOICE MLS INVALID COMMIT WELCOME NOT IMPLEMENTED")
	return nil
}

func startVoiceHeartbeatTimer(s VoiceSession) error {
	if s.GetHeartbeatAck() == nil {
		return errors.New("no heartbeat interval set")
	}

	ticker := time.NewTicker(time.Duration(*s.GetHeartbeatAck()) * time.Millisecond)
	go voiceHeartbeatLoop(ticker, s)
	return nil
}

func voiceHeartbeatLoop(ticker *time.Ticker, s VoiceSession) {
	defer ticker.Stop()

	for {
		select {
		case <-s.GetCtx().Done():
			return
		case <-ticker.C:
			if err := sendVoiceHeartbeatEvent(s); err != nil {
				return
			}
		}
	}
}

func sendVoiceHeartbeatEvent(s VoiceSession) error {
	heartbeatEvent := sendevents.VoiceHeartbeatEvent{
		Timestamp: time.Now().Unix(),
		SeqAck:    s.GetSequence(),
	}
	ackPayload := payload.VoicePayload{
		OpCode: gateway.VoiceOpHeartbeat,
		Data:   heartbeatEvent,
	}

	heartbeatData, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}

	s.Write(heartbeatData, false)
	return nil
}
