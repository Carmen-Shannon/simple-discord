package voice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/voice"
	"github.com/Carmen-Shannon/simple-discord/util"
)

var voiceOpCodeNames = map[voice.VoiceOpCode]string{
	voice.Identify:                    "Identify",
	voice.SelectProtocol:              "Select Protocol",
	voice.Ready:                       "Ready",
	voice.Heartbeat:                   "Heartbeat",
	voice.SessionDescription:          "Session Description",
	voice.Speaking:                    "Speaking",
	voice.HeartbeatAck:                "Heartbeat Ack",
	voice.Resume:                      "Resume",
	voice.Hello:                       "Hello",
	voice.Resumed:                     "Resumed",
	voice.ClientsConnect:              "Clients Connect",
	voice.ClientDisconnect:            "Client Disconnect",
	voice.PrepareTransition:           "DAVE Prepare Transition",
	voice.ExecuteTransition:           "DAVE Execute Transition",
	voice.TransitionReady:             "DAVE Transition Ready",
	voice.PrepareEpoch:                "DAVE Prepare Epoch",
	voice.MLSExternalSender:           "DAVE MLS External Sender",
	voice.MLSKeyPackage:               "DAVE MLS Key Package",
	voice.MLSProposals:                "DAVE MLS Proposals",
	voice.MLSCommitWelcome:            "DAVE MLS Commit Welcome",
	voice.MLSAnnounceCommitTransition: "DAVE MLS Announce Commit Transition",
	voice.MLSWelcome:                  "DAVE MLS Welcome",
	voice.MLSInvalidCommitWelcome:     "DAVE MLS Invalid Commit Welcome",
}

type VoiceEventFunc func(VoiceSession, voice.VoicePayload) error
type BinaryVoiceEventFunc func(VoiceSession, voice.BinaryVoicePayload) error
type UdpEventFunc func(UdpSession, voice.VoicePacket) error
type UdpDiscoveryEventFunc func(UdpSession, voice.DiscoveryPacket) error

type VoiceEventHandler struct {
	OpCodeHandlers map[voice.VoiceOpCode]VoiceEventFunc
	BinaryHandlers map[voice.VoiceOpCode]BinaryVoiceEventFunc
}

type UdpEventHandler struct {
	VoicePacketHandlers map[string]UdpEventFunc
	DiscoveryHandlers   map[string]UdpDiscoveryEventFunc
}

func NewVoiceEventHandler() *VoiceEventHandler {
	e := &VoiceEventHandler{
		OpCodeHandlers: map[voice.VoiceOpCode]VoiceEventFunc{
			voice.Identify:                    handleSendVoiceIdentifyEvent,
			voice.SelectProtocol:              handleSendVoiceSelectProtocolEvent,
			voice.Ready:                       handleVoiceReadyEvent,
			voice.Heartbeat:                   handleVoiceSendHeartbeatEvent,
			voice.SessionDescription:          handleVoiceSessionDescriptionEvent,
			voice.Speaking:                    handleVoiceSpeakingEvent,
			voice.HeartbeatAck:                handleVoiceHeartbeatAckEvent,
			voice.Resume:                      handleSendVoiceResumeEvent,
			voice.Hello:                       handleVoiceHelloEvent,
			voice.Resumed:                     handleVoiceResumedEvent,
			voice.ClientsConnect:              handleVoiceClientsConnectEvent,
			voice.ClientDisconnect:            handleVoiceClientDisconnectEvent,
			voice.PrepareTransition:           handleVoicePrepareTransitionEvent,
			voice.ExecuteTransition:           handleVoiceExecuteTransitionEvent,
			voice.TransitionReady:             handleSendVoiceTransitionReadyEvent,
			voice.PrepareEpoch:                handleVoicePrepareEpochEvent,
			voice.MLSAnnounceCommitTransition: handleVoiceMLSAnnounceCommitTransitionEvent,
			voice.MLSInvalidCommitWelcome:     handleSendVoiceMLSInvalidCommitWelcomeEvent,
		},
		BinaryHandlers: map[voice.VoiceOpCode]BinaryVoiceEventFunc{
			voice.MLSExternalSender: handleVoiceMLSExternalSenderEvent,
			voice.MLSKeyPackage:     handleSendVoiceMLSKeyPackageEvent,
			voice.MLSProposals:      handleVoiceMLSProposalsEvent,
			voice.MLSCommitWelcome:  handleSendVoiceMLSCommitWelcomeEvent,
			voice.MLSWelcome:        handleVoiceMLSWelcomeEvent,
		},
	}
	return e
}

func NewUdpEventHandler() *UdpEventHandler {
	e := &UdpEventHandler{
		VoicePacketHandlers: map[string]UdpEventFunc{},
		DiscoveryHandlers: map[string]UdpDiscoveryEventFunc{
			"send-discovery":    handleSendDiscoveryEvent,
			"receive-discovery": handleReceiveDiscoveryEvent,
		},
	}
	return e
}

// HandleEvent handles voice events the same way as HandleEvent for normal events
func (e *VoiceEventHandler) HandleEvent(s VoiceSession, payload voice.VoicePayload) error {
	fmt.Printf("HANDLING VOICE EVENT: %v, %s\n", payload.OpCode, voiceOpCodeNames[payload.OpCode])
	if handler, ok := e.OpCodeHandlers[payload.OpCode]; ok && handler != nil {
		if payload.Seq != nil {
			s.SetSequence(*payload.Seq)
		}
		go func() {
			if err := handler(s, payload); err != nil {
				s.GetSession().errorChan <- fmt.Errorf("error handling voice event: %v\n%s", err, payload.ToString())
			}
		}()
		return nil
	}
	return errors.New("no handler for voice opcode")
}

func (e *VoiceEventHandler) HandleBinaryEvent(s VoiceSession, payload voice.BinaryVoicePayload) error {
	fmt.Println("HANDLING BINARY EVENT")
	if handler, ok := e.BinaryHandlers[voice.VoiceOpCode(payload.OpCode)]; ok && handler != nil {
		if payload.SequenceNumber != nil {
			s.SetSequence(int(*payload.SequenceNumber))
		}

		go func() {
			if err := handler(s, payload); err != nil {
				s.GetSession().errorChan <- err
			}
		}()
		return nil
	}
	return nil
}

func (e *UdpEventHandler) HandleSendEvent(s UdpSession, payload voice.VoicePacket) error {
	fmt.Println("HANDLING UDP SEND EVENT")
	fmt.Println("UDP SEND EVENT NOT IMPLEMENTED")
	return nil
}

func (e *UdpEventHandler) HandleReceiveEvent(s UdpSession, payload voice.VoicePacket) error {
	fmt.Println("HANDLING UDP RECEIVE EVENT")
	fmt.Println("UDP RECEIVE EVENT NOT IMPLEMENTED")
	return nil
}

func (e *UdpEventHandler) HandleSendDiscoveryEvent(s UdpSession, payload voice.DiscoveryPacket) error {
	fmt.Println("HANDLING DISCOVERY SEND EVENT")
	if handler, ok := e.DiscoveryHandlers["send-discovery"]; ok && handler != nil {
		go func() {
			if err := handler(s, payload); err != nil {
				s.GetSession().errorChan <- err
			}
		}()
		return nil
	} else {
		return errors.New("no handler for discovery event")
	}
}

func (e *UdpEventHandler) HandleReceiveDiscoveryEvent(s UdpSession, payload voice.DiscoveryPacket) error {
	fmt.Println("HANDLING DISCOVERY RECEIVE EVENT")
	if handler, ok := e.DiscoveryHandlers["receive-discovery"]; ok && handler != nil {
		go func() {
			if err := handler(s, payload); err != nil {
				s.GetSession().errorChan <- err
			}
		}()
		return nil
	} else {
		return errors.New("no handler for discovery event")
	}
}

func handleSendDiscoveryEvent(s UdpSession, payload voice.DiscoveryPacket) error {
	connData := s.GetConnData()
	if connData == nil {
		return errors.New("UDP connection data is not set")
	}

	var packet voice.DiscoveryPacket
	// Check and populate packet properties with default values if they are zero
	if payload.Type == 0 {
		packet.Type = 0x1
	} else {
		packet.Type = payload.Type
	}

	if payload.Length == 0 {
		packet.Length = 70
	} else {
		packet.Length = payload.Length
	}

	if payload.SSRC == 0 {
		packet.SSRC = uint32(connData.SSRC)
	} else {
		packet.SSRC = payload.SSRC
	}
	fmt.Println("SENDING DISCOVERY SSRC:", packet.SSRC)

	if payload.Address == [64]byte{} {
		packet.Address = [64]byte{}
	} else {
		packet.Address = payload.Address
	}

	packetBytes, err := packet.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal IP discovery packet: %v", err)
	}

	// Send the IP discovery packet
	s.Write(packetBytes)
	return nil
}

func handleReceiveDiscoveryEvent(s UdpSession, payload voice.DiscoveryPacket) error {
	// Extract the external IP and port
	externalIP := string(bytes.Trim(payload.Address[:], "\x00"))
	externalPort := payload.Port
	ssrc := payload.SSRC

	fmt.Println("RECEIVING DISCOVERY SSRC:", ssrc)

	s.GetConnData().Address = externalIP
	s.GetConnData().Port = int(externalPort)
	s.GetConnData().SSRC = int(ssrc)

	// Signal that discovery is complete
	if s.GetSession().discoveryReady != nil {
		close(s.GetSession().discoveryReady)
		s.GetSession().discoveryReady = nil
	}
	return nil
}

func handleSendVoiceIdentifyEvent(s VoiceSession, p voice.VoicePayload) error {
	// no DAVE support yet, include DaveProtocolVersion
	voiceIdentifyEvent := sendevents.VoiceIdentifyEvent{
		ServerID:               *s.GetGuildID(),
		UserID:                 s.GetBotData().UserDetails.ID,
		SessionID:              *s.GetSessionID(),
		Token:                  *s.GetToken(),
		MaxDaveProtocolVersion: util.ToPtr(0),
	}
	ackPayload := voice.VoicePayload{
		OpCode: voice.Identify,
		Data:   voiceIdentifyEvent,
		Seq:    s.GetSequence(),
	}
	data, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}
	s.Write(data)
	return nil
}

func handleSendVoiceSelectProtocolEvent(s VoiceSession, p voice.VoicePayload) error {
	selectProtocolEvent := sendevents.VoiceSelectProtocolEvent{
		Protocol: "udp",
		Data: sendevents.VoiceSelectProtocolData{
			Address: s.GetAudioPlayer().GetUdpSession().GetConnData().Address,
			Port:    s.GetAudioPlayer().GetUdpSession().GetConnData().Port,
			Mode:    s.GetAudioPlayer().GetUdpSession().GetConnData().Mode,
		},
		Codecs:              []voice.Codec{voice.Opus},
		DaveProtocolVersion: util.ToPtr(0),
	}
	selectProtocolPayload := voice.VoicePayload{
		OpCode: voice.SelectProtocol,
		Data:   selectProtocolEvent,
		Seq:    s.GetSequence(),
	}

	data, err := json.Marshal(selectProtocolPayload)
	if err != nil {
		return err
	}

	s.Write(data)
	return nil
}

func handleVoiceReadyEvent(s VoiceSession, p voice.VoicePayload) error {
	if voiceReadyEvent, ok := p.Data.(receiveevents.VoiceReadyEvent); ok {
		udpConn := &voice.UdpData{
			Address: voiceReadyEvent.IP,
			Port:    voiceReadyEvent.Port,
			SSRC:    voiceReadyEvent.SSRC,
		}
		fmt.Println(voiceReadyEvent.Modes)
		if util.SliceContains(voiceReadyEvent.Modes, voice.AEAD_AES256_GCM) {
			udpConn.Mode = voice.AEAD_AES256_GCM
		} else if len(voiceReadyEvent.Modes) == 1 {
			udpConn.Mode = voiceReadyEvent.Modes[0]
		} else {
			udpConn.Mode = voice.AEAD_XCHACHA20_POLY1305
		}

		s.GetAudioPlayer().GetUdpSession().SetConnData(udpConn)
		if s.GetAudioPlayer().GetUdpSession().GetSession().connectionReady != nil {
			close(s.GetAudioPlayer().GetUdpSession().GetSession().connectionReady)
			s.GetAudioPlayer().GetUdpSession().GetSession().connectionReady = nil
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceSendHeartbeatEvent(s VoiceSession, p voice.VoicePayload) error {
	if heartbeatEvent, ok := p.Data.(sendevents.VoiceHeartbeatEvent); ok {
		if heartbeatEvent.SeqAck != nil {
			s.SetSequence(*heartbeatEvent.SeqAck)
		}
		return sendVoiceHeartbeatEvent(s)
	}
	return errors.New("unexpected payload data type")
}

func handleVoiceSessionDescriptionEvent(s VoiceSession, p voice.VoicePayload) error {
	if voiceSessionDescriptionEvent, ok := p.Data.(receiveevents.VoiceSessionDescriptionEvent); ok {
		s.GetAudioPlayer().GetUdpSession().SetSecretKey(voiceSessionDescriptionEvent.SecretKey)
		s.GetAudioPlayer().GetUdpSession().SetEncryptionMode(voiceSessionDescriptionEvent.Mode)

		if s.GetAudioPlayer().GetUdpSession().GetSession().speakingReady != nil {
			close(s.GetAudioPlayer().GetUdpSession().GetSession().speakingReady)
			s.GetAudioPlayer().GetUdpSession().GetSession().speakingReady = nil
		}

		return nil
	} else {
		return errors.New("unexpected payload data type")
	}
}

func handleVoiceSpeakingEvent(s VoiceSession, p voice.VoicePayload) error {
	if p.Data != nil {
		return nil
	}

	var speakingEvent sendevents.SpeakingEvent
	ssrc := s.GetAudioPlayer().GetUdpSession().GetConnData().SSRC
	if s.GetAudioPlayer() != nil && s.GetAudioPlayer().IsPlaying() {
		fmt.Println("SPEAKING STOP")
		speakingEvent.SpeakingEvent = &structs.SpeakingEvent{
			Speaking: structs.Bitfield[structs.SpeakingFlag]{},
			Delay:    0,
			SSRC:     &ssrc,
		}
	} else {
		fmt.Println("SPEAKING START")
		speakingEvent.SpeakingEvent = &structs.SpeakingEvent{
			Speaking: structs.Bitfield[structs.SpeakingFlag]{structs.SpeakingFlagMicrophone},
			Delay:    0,
			SSRC:     &ssrc,
		}
	}
	speakingPayload := voice.VoicePayload{
		OpCode: voice.Speaking,
		Data:   speakingEvent,
		Seq:    s.GetSequence(),
	}

	data, err := json.Marshal(speakingPayload)
	if err != nil {
		return err
	}

	s.Write(data)
	return nil
}

func handleVoiceHeartbeatAckEvent(s VoiceSession, p voice.VoicePayload) error {
	return nil
}

func handleSendVoiceResumeEvent(s VoiceSession, p voice.VoicePayload) error {
	if _, ok := p.Data.(sendevents.VoiceResumeEvent); ok {
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}

		s.Write(data)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceHelloEvent(s VoiceSession, p voice.VoicePayload) error {
	if helloEvent, ok := p.Data.(receiveevents.HelloEvent); ok {
		s.SetHeartbeatACK(int(helloEvent.HeartbeatInterval))
	} else {
		return errors.New("unexpected payload data type")
	}
	return startVoiceHeartbeatTimer(s)
}

func handleVoiceResumedEvent(s VoiceSession, p voice.VoicePayload) error {
	if _, ok := p.Data.(receiveevents.VoiceResumedEvent); ok {
		close(s.GetSession().resumeReady)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceClientsConnectEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE CLIENTS CONNECT EVENT")
	if _, ok := p.Data.(receiveevents.VoiceClientsConnectEvent); ok {
		// fmt.Println(event)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceClientDisconnectEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE CLIENT DISCONNECT EVENT")
	fmt.Println("VOICE CLIENT DISCONNECT NOT IMPLEMENTED")
	return nil
}

func handleVoicePrepareTransitionEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE PREPARE TRANSITION EVENT")
	fmt.Println("VOICE PREPARE TRANSITION NOT IMPLEMENTED")
	return nil
}

func handleVoiceExecuteTransitionEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE EXECUTE TRANSITION EVENT")
	fmt.Println("VOICE EXECUTE TRANSITION NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceTransitionReadyEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE TRANSITION READY EVENT")
	fmt.Println("VOICE TRANSITION READY NOT IMPLEMENTED")
	return nil
}

func handleVoicePrepareEpochEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE PREPARE EPOCH EVENT")
	fmt.Println("VOICE PREPARE EPOCH NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSExternalSenderEvent(s VoiceSession, p voice.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS EXTERNAL SENDER EVENT")
	fmt.Println("VOICE MLS EXTERNAL SENDER NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceMLSKeyPackageEvent(s VoiceSession, p voice.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS KEY PACKAGE EVENT")
	fmt.Println("VOICE MLS KEY PACKAGE NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSProposalsEvent(s VoiceSession, p voice.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS PROPOSALS EVENT")
	fmt.Println("VOICE MLS PROPOSALS NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceMLSCommitWelcomeEvent(s VoiceSession, p voice.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS COMMIT WELCOME EVENT")
	fmt.Println("VOICE MLS COMMIT WELCOME NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSAnnounceCommitTransitionEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE MLS ANNOUNCE COMMIT TRANSITION EVENT")
	fmt.Println("VOICE MLS ANNOUNCE COMMIT TRANSITION NOT IMPLEMENTED")
	return nil
}

func handleVoiceMLSWelcomeEvent(s VoiceSession, p voice.BinaryVoicePayload) error {
	fmt.Println("HANDLING VOICE MLS WELCOME EVENT")
	fmt.Println("VOICE MLS WELCOME NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceMLSInvalidCommitWelcomeEvent(s VoiceSession, p voice.VoicePayload) error {
	fmt.Println("HANDLING VOICE MLS INVALID COMMIT WELCOME EVENT")
	fmt.Println("VOICE MLS INVALID COMMIT WELCOME NOT IMPLEMENTED")
	return nil
}

func startVoiceHeartbeatTimer(s VoiceSession) error {
	if s.GetHeartbeatACK() == nil {
		return errors.New("no heartbeat interval set")
	}

	ticker := time.NewTicker(time.Duration(*s.GetHeartbeatACK()) * time.Millisecond)
	go voiceHeartbeatLoop(ticker, s)
	return nil
}

func voiceHeartbeatLoop(ticker *time.Ticker, s VoiceSession) {
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := sendVoiceHeartbeatEvent(s); err != nil {
				return
			}
		case <-s.GetSession().stopHeartbeat:
			return
		}
	}
}

func sendVoiceHeartbeatEvent(s VoiceSession) error {
	if s.GetConn() == nil {
		return errors.New("connection unavailable")
	}

	heartbeatEvent := sendevents.VoiceHeartbeatEvent{
		Timestamp: time.Now().Unix(),
		SeqAck:    s.GetSequence(),
	}
	ackPayload := voice.VoicePayload{
		OpCode: voice.Heartbeat,
		Data:   heartbeatEvent,
	}

	heartbeatData, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}

	s.Write(heartbeatData)
	return nil
}
