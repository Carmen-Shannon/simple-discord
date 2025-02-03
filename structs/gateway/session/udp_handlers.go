package session

import (
	"bytes"

	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
)

// this is an empty file for now
// the udp session is different than the tcp sessions, we typically want to write directly to the connection
// for this reaso, we will provide interfaces for the udp session writes via the audio player
// TODO: add the proper udp event handlers here, such as receiving voice packets

func handleDiscoveryEvent(s UdpSession, p payload.DiscoveryPacket) error {
	if s.IsDiscovered() {
		return nil
	}

	var packet payload.DiscoveryPacket
	// if we have an empty payload, it means we want to send a discovery packet
	if p.PacketType == 0 {
		packet.PacketType = 0x1
		packet.Length = 70
		packet.SSRC = uint32(s.GetUdpData().SSRC)
		packet.Address = [64]byte{}

		packetBytes, err := packet.Marshal()
		if err != nil {
			return err
		}

		s.Write(packetBytes, true)
		return nil
	}

	// if the payload is not empty, we received it
	s.GetUdpData().Address = string(bytes.Trim(p.Address[:], "\x00"))
	s.GetUdpData().Port = int(p.Port)
	s.GetUdpData().SSRC = int(p.SSRC)
	s.CloseDiscoveryReady()
	return nil
}

func handleVoicePacketEvent(s UdpSession, p payload.VoicePacket) error {
	// receiving voice packets currently just does nothing
	return nil
}
