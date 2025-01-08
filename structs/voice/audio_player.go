package voice

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Carmen-Shannon/simple-discord/util"
)

type AudioMetadata struct {
	DurationMs float64 `json:"duration_ms"`
	Bitrate    int     `json:"bit_rate"`
	SampleRate int     `json:"sample_rate"`
	Channels   int     `json:"channels"`
}

var (
	Opus = Codec{
		Name:        "opus",
		Type:        "audio",
		Priority:    1000,
		PayloadType: 120,
	}
)

type Codec struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Priority    int    `json:"priority"`
	PayloadType int    `json:"payload_type"`
	Encode      *bool  `json:"encode,omitempty"`
	Decode      *bool  `json:"decode,omitempty"`
}

type SenderReportPacket struct {
	Version              uint8  // 2 bits
	Padding              bool   // 1 bit
	ReceptionReportCount uint8  // 5 bits
	PacketType           uint8  // 8 bits
	Length               uint16 // 16 bits
	SSRC                 uint32 // 32 bits
	NTPTimestamp         uint64 // 64 bits
	RTPTimestamp         uint32 // 32 bits
	SenderPacketCount    uint32 // 32 bits
	SenderOctetCount     uint32 // 32 bits
}

func (s *SenderReportPacket) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// First byte: Version (2 bits), Padding (1 bit), Reception Report Count (5 bits)
	firstByte := (s.Version << 6) | (util.BoolToUint8(s.Padding) << 5) | (s.ReceptionReportCount & 0x1F)
	if err := buf.WriteByte(firstByte); err != nil {
		return nil, fmt.Errorf("failed to write first byte: %w", err)
	}

	// Second byte: Packet Type
	if err := buf.WriteByte(s.PacketType); err != nil {
		return nil, fmt.Errorf("failed to write packet type: %w", err)
	}

	// Placeholder for Length (will be set later)
	if err := binary.Write(buf, binary.BigEndian, uint16(0)); err != nil {
		return nil, fmt.Errorf("failed to write placeholder length: %w", err)
	}

	// Write SSRC
	if err := binary.Write(buf, binary.BigEndian, s.SSRC); err != nil {
		return nil, fmt.Errorf("failed to write SSRC: %w", err)
	}

	// Write NTP Timestamp
	if err := binary.Write(buf, binary.BigEndian, s.NTPTimestamp); err != nil {
		return nil, fmt.Errorf("failed to write NTP timestamp: %w", err)
	}

	// Write RTP Timestamp
	if err := binary.Write(buf, binary.BigEndian, s.RTPTimestamp); err != nil {
		return nil, fmt.Errorf("failed to write RTP timestamp: %w", err)
	}

	// Write Sender Packet Count
	if err := binary.Write(buf, binary.BigEndian, s.SenderPacketCount); err != nil {
		return nil, fmt.Errorf("failed to write sender packet count: %w", err)
	}

	// Write Sender Octet Count
	if err := binary.Write(buf, binary.BigEndian, s.SenderOctetCount); err != nil {
		return nil, fmt.Errorf("failed to write sender octet count: %w", err)
	}

	// Calculate the total length of the packet in bytes
	totalLength := buf.Len()

	// Calculate the length in 32-bit words minus one
	lengthInWords := (totalLength / 4) - 1

	// Set the Length field
	s.Length = uint16(lengthInWords)

	// Update the length field in the buffer
	lengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBytes, s.Length)
	copy(buf.Bytes()[2:4], lengthBytes)

	return buf.Bytes(), nil
}

func (s *SenderReportPacket) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// First byte: Version (2 bits), Padding (1 bit), Reception Report Count (5 bits)
	var firstByte uint8
	if err := binary.Read(buf, binary.BigEndian, &firstByte); err != nil {
		return fmt.Errorf("failed to read first byte: %w", err)
	}
	s.Version = firstByte >> 6
	s.Padding = (firstByte>>5)&0x01 == 1
	s.ReceptionReportCount = firstByte & 0x1F

	// Second byte: Packet Type
	if err := binary.Read(buf, binary.BigEndian, &s.PacketType); err != nil {
		return fmt.Errorf("failed to read packet type: %w", err)
	}

	// Third and fourth bytes: Length
	if err := binary.Read(buf, binary.BigEndian, &s.Length); err != nil {
		return fmt.Errorf("failed to read length: %w", err)
	}

	// Next four bytes: SSRC
	if err := binary.Read(buf, binary.BigEndian, &s.SSRC); err != nil {
		return fmt.Errorf("failed to read SSRC: %w", err)
	}

	// Next eight bytes: NTP Timestamp
	if err := binary.Read(buf, binary.BigEndian, &s.NTPTimestamp); err != nil {
		return fmt.Errorf("failed to read NTP timestamp: %w", err)
	}

	// Next four bytes: RTP Timestamp
	if err := binary.Read(buf, binary.BigEndian, &s.RTPTimestamp); err != nil {
		return fmt.Errorf("failed to read RTP timestamp: %w", err)
	}

	// Next four bytes: Sender Packet Count
	if err := binary.Read(buf, binary.BigEndian, &s.SenderPacketCount); err != nil {
		return fmt.Errorf("failed to read sender packet count: %w", err)
	}

	// Next four bytes: Sender Octet Count
	if err := binary.Read(buf, binary.BigEndian, &s.SenderOctetCount); err != nil {
		return fmt.Errorf("failed to read sender octet count: %w", err)
	}

	return nil
}
