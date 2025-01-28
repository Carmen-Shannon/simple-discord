package payload

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/util"
)

var SilenceFrame = []byte{0xF8, 0xFF, 0xFE}

type Payload interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
	ToString() string
	Type() string
	Hash() string
}

type SessionPayload struct {
	OpCode    gateway.GatewayOpCode `json:"op"`
	Data      any                   `json:"d"`
	Seq       *int                  `json:"s,omitempty"`
	EventName *string               `json:"t,omitempty"`
}

func (p *SessionPayload) Marshal() ([]byte, error) {
	return json.Marshal(&p)
}

func (p *SessionPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (p *SessionPayload) ToString() string {
	jsonData, _ := json.Marshal(&p)
	return string(jsonData)
}

func (p *SessionPayload) Type() string {
	return "SessionPayload"
}

func (p *SessionPayload) Hash() string {
	h := sha256.New()
	h.Write([]byte(p.ToString()))
	return hex.EncodeToString(h.Sum(nil))
}

type VoicePayload struct {
	OpCode    gateway.VoiceOpCode `json:"op"`
	Data      any                 `json:"d"`
	Seq       *int                `json:"s,omitempty"`
	EventName *string             `json:"t,omitempty"`
}

func (p *VoicePayload) Marshal() ([]byte, error) {
	return json.Marshal(&p)
}

func (p *VoicePayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (p *VoicePayload) ToString() string {
	jsonData, _ := json.Marshal(p)
	return string(jsonData)
}

func (p *VoicePayload) Type() string {
	return "VoicePayload"
}

func (p *VoicePayload) Hash() string {
	h := sha256.New()
	h.Write([]byte(p.ToString()))
	return hex.EncodeToString(h.Sum(nil))
}

type RawMessagePayload struct {
	Data json.RawMessage
}

func (p *RawMessagePayload) Marshal() ([]byte, error) {
	return json.Marshal(&p.Data)
}

func (p *RawMessagePayload) Unmarshal(data []byte) error {
	p.Data = json.RawMessage(data)

	if len(p.Data) == 0 {
		return io.EOF
	}
	return nil
}

func (p *RawMessagePayload) ToString() string {
	jsonData, _ := json.Marshal(&p.Data)
	return string(jsonData)
}

func (p *RawMessagePayload) Hash() string {
	h := sha256.New()
	h.Write(p.Data)
	return hex.EncodeToString(h.Sum(nil))
}

func (p *RawMessagePayload) Type() string {
	return "RawMessagePayload"
}

type DiscoveryPacket struct {
	PacketType uint16
	Length     uint16
	SSRC       uint32
	Address    [64]byte
	Port       uint16
}

func (i *DiscoveryPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, i.PacketType); err != nil {
		return nil, fmt.Errorf("failed to write type: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, i.Length); err != nil {
		return nil, fmt.Errorf("failed to write length: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, i.SSRC); err != nil {
		return nil, fmt.Errorf("failed to write SSRC: %w", err)
	}

	if _, err := buf.Write(i.Address[:]); err != nil {
		return nil, fmt.Errorf("failed to write address: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, i.Port); err != nil {
		return nil, fmt.Errorf("failed to write port: %w", err)
	}

	return buf.Bytes(), nil
}

func (i *DiscoveryPacket) Unmarshal(data []byte) error {
	// if len(data) != 70 {
	// 	return errors.New("invalid data length for Discovery Packet")
	// }
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.BigEndian, &i.PacketType); err != nil {
		return fmt.Errorf("failed to read type: %w", err)
	}

	if err := binary.Read(buf, binary.BigEndian, &i.Length); err != nil {
		return fmt.Errorf("failed to read length: %w", err)
	}

	if err := binary.Read(buf, binary.BigEndian, &i.SSRC); err != nil {
		return fmt.Errorf("failed to read SSRC: %w", err)
	}

	if _, err := buf.Read(i.Address[:]); err != nil {
		return fmt.Errorf("failed to read address: %w", err)
	}

	if err := binary.Read(buf, binary.BigEndian, &i.Port); err != nil {
		return fmt.Errorf("failed to read port: %w", err)
	}

	return nil
}

func (i *DiscoveryPacket) ToString() string {
	// Create a map to hold the JSON representation
	packetMap := map[string]interface{}{
		"Type":    i.PacketType,
		"Length":  i.Length,
		"SSRC":    i.SSRC,
		"Address": string(bytes.Trim(i.Address[:], "\x00")),
		"Port":    i.Port,
	}

	// Marshal the map into a JSON string
	jsonData, err := json.Marshal(packetMap)
	if err != nil {
		return fmt.Sprintf("error marshaling DiscoveryPacket to JSON: %v", err)
	}

	return string(jsonData)
}

func (i *DiscoveryPacket) Type() string {
	return "DiscoveryPacket"
}

func (i *DiscoveryPacket) Hash() string {
	h := sha256.New()
	h.Write([]byte(i.ToString()))
	return hex.EncodeToString(h.Sum(nil))
}

type BinaryVoicePayload struct {
	SequenceNumber *uint16
	OpCode         uint8
	Data           []byte
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (p *BinaryVoicePayload) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write SequenceNumber if present
	if p.SequenceNumber != nil {
		if err := binary.Write(buf, binary.BigEndian, *p.SequenceNumber); err != nil {
			return nil, fmt.Errorf("failed to write sequence number: %w", err)
		}
	} else {
		// Write zero value for SequenceNumber
		if err := binary.Write(buf, binary.BigEndian, uint16(0)); err != nil {
			return nil, fmt.Errorf("failed to write zero sequence number: %w", err)
		}
	}

	// Write OpCode
	if err := binary.Write(buf, binary.BigEndian, p.OpCode); err != nil {
		return nil, fmt.Errorf("failed to write opcode: %w", err)
	}

	// Write payload
	if err := binary.Write(buf, binary.BigEndian, p.Data); err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (p *BinaryVoicePayload) Unmarshal(data []byte) error {
	// Check if the data is a RawMessagePayload
	if isRawMessagePayload(data) {
		return errors.New("data is a RawMessagePayload, cannot unmarshal into BinaryVoicePayload")
	}
	buf := bytes.NewReader(data)

	// Read SequenceNumber
	var seqNum uint16
	if err := binary.Read(buf, binary.BigEndian, &seqNum); err != nil {
		return fmt.Errorf("failed to read sequence number: %w", err)
	}
	if seqNum != 0 {
		p.SequenceNumber = &seqNum
	} else {
		p.SequenceNumber = nil
	}

	// Read OpCode
	if err := binary.Read(buf, binary.BigEndian, &p.OpCode); err != nil {
		return fmt.Errorf("failed to read opcode: %w", err)
	}

	// Read Payload
	payload := make([]byte, buf.Len())
	if err := binary.Read(buf, binary.BigEndian, &payload); err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}
	p.Data = payload

	return nil
}

// Helper function to check if data is a RawMessagePayload
func isRawMessagePayload(data []byte) bool {
	var rawMessage RawMessagePayload
	err := rawMessage.Unmarshal(data)
	return err == nil
}

func (p *BinaryVoicePayload) ToString() string {
	packetMap := map[string]interface{}{
		"SequenceNumber": p.SequenceNumber,
		"OpCode":         p.OpCode,
		"Data":           string(p.Data),
	}

	jsonData, err := json.Marshal(packetMap)
	if err != nil {
		return fmt.Sprintf("error marshaling BinaryVoicePayload to JSON: %v", err)
	}

	return string(jsonData)
}

func (p *BinaryVoicePayload) Type() string {
	return "BinaryVoicePayload"
}

func (p *BinaryVoicePayload) Hash() string {
	h := sha256.New()
	h.Write([]byte(p.ToString()))
	return hex.EncodeToString(h.Sum(nil))
}

type RTPHeader struct {
	Version     uint8
	Padding     bool
	Extension   bool
	CSRCCount   uint8
	Marker      bool
	PayloadType uint8
	Seq         uint16
	Timestamp   uint32
	SSRC        uint32
	CSRC        []uint32
	ExtProfile  uint16
	ExtLength   uint16
	ExtData     []byte
}

func NewRtpHeader(seq uint16, timestamp uint32, ssrc uint32) *RTPHeader {
	return &RTPHeader{
		Version:     2,
		Padding:     false,
		Extension:   false,
		CSRCCount:   0,
		Marker:      false,
		PayloadType: 120,
		Seq:         seq,
		Timestamp:   timestamp,
		SSRC:        ssrc,
	}
}

func (r *RTPHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Combine Version, Padding, Extension, and CSRCCount into the first byte
	firstByte := (r.Version << 6) | (util.BoolToUint8(r.Padding) << 5) | (util.BoolToUint8(r.Extension) << 4) | (r.CSRCCount & 0x0F)
	if err := buf.WriteByte(firstByte); err != nil {
		return nil, fmt.Errorf("failed to write first byte: %w", err)
	}

	// Combine Marker and PayloadType into the second byte
	secondByte := (util.BoolToUint8(r.Marker) << 7) | (r.PayloadType & 0x7F)
	if err := buf.WriteByte(secondByte); err != nil {
		return nil, fmt.Errorf("failed to write second byte: %w", err)
	}

	// Write Seq
	if err := binary.Write(buf, binary.BigEndian, r.Seq); err != nil {
		return nil, fmt.Errorf("failed to write sequence number: %w", err)
	}

	// Write Timestamp
	if err := binary.Write(buf, binary.BigEndian, r.Timestamp); err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Write SSRC
	if err := binary.Write(buf, binary.BigEndian, r.SSRC); err != nil {
		return nil, fmt.Errorf("failed to write SSRC: %w", err)
	}

	// Write CSRC identifiers if present
	for _, csrc := range r.CSRC {
		if err := binary.Write(buf, binary.BigEndian, csrc); err != nil {
			return nil, fmt.Errorf("failed to write CSRC: %w", err)
		}
	}

	// Write extension header if Extension is true
	if r.Extension {
		if err := binary.Write(buf, binary.BigEndian, r.ExtProfile); err != nil {
			return nil, fmt.Errorf("failed to write extension profile: %w", err)
		}
		if err := binary.Write(buf, binary.BigEndian, r.ExtLength); err != nil {
			return nil, fmt.Errorf("failed to write extension length: %w", err)
		}
		if _, err := buf.Write(r.ExtData); err != nil {
			return nil, fmt.Errorf("failed to write extension data: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func (r *RTPHeader) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		fmt.Println(string(data))
		return fmt.Errorf("data too short to contain RTP header")
	}

	buf := bytes.NewReader(data)

	// Read the first byte and extract Version, Padding, Extension, and CSRCCount
	var firstByte uint8
	if err := binary.Read(buf, binary.BigEndian, &firstByte); err != nil {
		return fmt.Errorf("failed to read first byte: %w", err)
	}
	r.Version = firstByte >> 6
	r.Padding = firstByte&0x20 != 0
	r.Extension = firstByte&0x10 != 0
	r.CSRCCount = firstByte & 0x0F

	// Read the second byte and extract Marker and PayloadType
	var secondByte uint8
	if err := binary.Read(buf, binary.BigEndian, &secondByte); err != nil {
		return fmt.Errorf("failed to read second byte: %w", err)
	}
	r.Marker = secondByte&0x80 != 0
	r.PayloadType = secondByte & 0x7F

	// Read Seq
	if err := binary.Read(buf, binary.BigEndian, &r.Seq); err != nil {
		return fmt.Errorf("failed to read sequence number: %w", err)
	}

	// Read Timestamp
	if err := binary.Read(buf, binary.BigEndian, &r.Timestamp); err != nil {
		return fmt.Errorf("failed to read timestamp: %w", err)
	}

	// Read SSRC
	if err := binary.Read(buf, binary.BigEndian, &r.SSRC); err != nil {
		return fmt.Errorf("failed to read SSRC: %w", err)
	}

	// Read CSRC identifiers if present
	r.CSRC = make([]uint32, r.CSRCCount)
	for i := uint8(0); i < r.CSRCCount; i++ {
		if err := binary.Read(buf, binary.BigEndian, &r.CSRC[i]); err != nil {
			return fmt.Errorf("failed to read CSRC: %w", err)
		}
	}

	// Read extension header if Extension is true
	if r.Extension {
		if err := binary.Read(buf, binary.BigEndian, &r.ExtProfile); err != nil {
			return fmt.Errorf("failed to read extension profile: %w", err)
		}
		if err := binary.Read(buf, binary.BigEndian, &r.ExtLength); err != nil {
			return fmt.Errorf("failed to read extension length: %w", err)
		}
		r.ExtData = make([]byte, r.ExtLength*4)
		if _, err := buf.Read(r.ExtData); err != nil {
			return fmt.Errorf("failed to read extension data: %w", err)
		}
	}

	return nil
}

func (r *RTPHeader) ToString() string {
	// Create a map to hold the JSON representation
	headerMap := map[string]interface{}{
		"Version":     r.Version,
		"Padding":     r.Padding,
		"Extension":   r.Extension,
		"CSRCCount":   r.CSRCCount,
		"Marker":      r.Marker,
		"PayloadType": r.PayloadType,
		"Seq":         r.Seq,
		"Timestamp":   r.Timestamp,
		"SSRC":        r.SSRC,
		"CSRC":        r.CSRC,
	}

	// Marshal the map into a JSON string
	jsonData, err := json.Marshal(headerMap)
	if err != nil {
		return fmt.Sprintf("error marshaling RTPHeader to JSON: %v", err)
	}

	return string(jsonData)
}

type VoicePacket struct {
	RTPHeader
	Payload []byte
}

func (v *VoicePacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write RTPHeader
	rtpHeaderBytes, err := v.RTPHeader.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RTP header: %w", err)
	}
	if _, err := buf.Write(rtpHeaderBytes); err != nil {
		return nil, fmt.Errorf("failed to write RTP header: %w", err)
	}

	// Write Payload
	if _, err := buf.Write(v.Payload); err != nil {
		return nil, fmt.Errorf("failed to write payload: %w", err)
	}

	return buf.Bytes(), nil
}

func (v *VoicePacket) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)

	// Read RTPHeader
	if err := v.RTPHeader.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("failed to read RTP header: %w", err)
	}

	// Calculate the length of the RTP header
	rtpHeaderLength := 12 + int(v.RTPHeader.CSRCCount)*4

	// Read the remaining data as Payload
	v.Payload = make([]byte, buf.Len()-rtpHeaderLength)
	if _, err := buf.Read(v.Payload); err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}

	return nil
}

func (v *VoicePacket) ToString() string {
	// Create a map to hold the JSON representation
	packetMap := map[string]interface{}{
		"RTPHeader": v.RTPHeader.ToString(),
		"Payload":   base64.StdEncoding.EncodeToString(v.Payload),
	}

	// Marshal the map into a JSON string
	jsonData, err := json.Marshal(packetMap)
	if err != nil {
		return fmt.Sprintf("error marshaling VoicePacket to JSON: %v", err)
	}

	return string(jsonData)
}

func (v *VoicePacket) Type() string {
	return "VoicePacket"
}

func (v *VoicePacket) Hash() string {
	h := sha256.New()
	h.Write([]byte(v.ToString()))
	return hex.EncodeToString(h.Sum(nil))
}

func NewSenderReportPacket(ntpTimestamp uint64, rtpTimestamp uint32, ssrc uint32, packetCount uint32, octectCount uint32) *SenderReportPacket {
	return &SenderReportPacket{
		Version:              2,
		Padding:              false,
		ReceptionReportCount: 0,
		PacketType:           200,
		Length:               0,
		SSRC:                 ssrc,
		NTPTimestamp:         ntpTimestamp,
		RTPTimestamp:         rtpTimestamp,
		SenderPacketCount:    packetCount,
		SenderOctetCount:     octectCount,
	}
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

func (s *SenderReportPacket) Marshal() ([]byte, error) {
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

func (s *SenderReportPacket) Unmarshal(data []byte) error {
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

func (s *SenderReportPacket) ToString() string {
	// Create a map to hold the JSON representation
	packetMap := map[string]interface{}{
		"Version":              s.Version,
		"Padding":              s.Padding,
		"ReceptionReportCount": s.ReceptionReportCount,
		"PacketType":           s.PacketType,
		"Length":               s.Length,
		"SSRC":                 s.SSRC,
		"NTPTimestamp":         s.NTPTimestamp,
		"RTPTimestamp":         s.RTPTimestamp,
		"SenderPacketCount":    s.SenderPacketCount,
		"SenderOctetCount":     s.SenderOctetCount,
	}

	// Marshal the map into a JSON string
	jsonData, err := json.Marshal(packetMap)
	if err != nil {
		return fmt.Sprintf("error marshaling SenderReportPacket to JSON: %v", err)
	}

	return string(jsonData)
}

func (s *SenderReportPacket) Type() string {
	return "SenderReportPacket"
}

func (s *SenderReportPacket) Hash() string {
	h := sha256.New()
	h.Write([]byte(s.ToString()))
	return hex.EncodeToString(h.Sum(nil))
}
