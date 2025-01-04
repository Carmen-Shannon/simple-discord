package voice

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type VoiceOpCode int

const (
	Identify VoiceOpCode = iota
	SelectProtocol
	Ready
	Heartbeat
	SessionDescription
	Speaking
	HeartbeatAck
	Resume
	Hello
	Resumed
	ClientsConnect              VoiceOpCode = 11
	ClientDisconnect            VoiceOpCode = 13
	PrepareTransition           VoiceOpCode = 21
	ExecuteTransition           VoiceOpCode = 22
	TransitionReady             VoiceOpCode = 23
	PrepareEpoch                VoiceOpCode = 24
	MLSExternalSender           VoiceOpCode = 25
	MLSKeyPackage               VoiceOpCode = 26
	MLSProposals                VoiceOpCode = 27
	MLSCommitWelcome            VoiceOpCode = 28
	MLSAnnounceCommitTransition VoiceOpCode = 29
	MLSWelcome                  VoiceOpCode = 30
	MLSInvalidCommitWelcome     VoiceOpCode = 31
)

type TransportEncryptionMode string

const (
	AEAD_AES256_GCM         TransportEncryptionMode = "aead_aes256_gcm_rtpsize"
	AEAD_XCHACHA20_POLY1305 TransportEncryptionMode = "aead_xchacha20_poly1305_rtpsize"
)

type UdpData struct {
	SSRC    int                     `json:"ssrc"`
	Address string                  `json:"address"`
	Port    int                     `json:"port"`
	Mode    TransportEncryptionMode `json:"mode"`
}

type DiscoveryPacket struct {
	Type uint16
	Length uint16
	SSRC uint32
	Address [64]byte
	Port uint16
}

func (i *DiscoveryPacket) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, i.Type); err != nil {
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

func (i *DiscoveryPacket) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.BigEndian, &i.Type); err != nil {
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
		"Type":    i.Type,
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

type VoicePacket struct {
	VersionFlags uint8
	PayloadType  uint8
	Seq          uint16
	Timestamp    uint32
	SSRC         uint32
	Payload      []byte
}

func (v *VoicePacket) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write VersionFlags
	if err := binary.Write(buf, binary.BigEndian, v.VersionFlags); err != nil {
		return nil, fmt.Errorf("failed to write version flags: %w", err)
	}

	// Write PayloadType
	if err := binary.Write(buf, binary.BigEndian, v.PayloadType); err != nil {
		return nil, fmt.Errorf("failed to write payload type: %w", err)
	}

	// Write Seq
	if err := binary.Write(buf, binary.BigEndian, v.Seq); err != nil {
		return nil, fmt.Errorf("failed to write sequence number: %w", err)
	}

	// Write Timestamp
	if err := binary.Write(buf, binary.BigEndian, v.Timestamp); err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Write SSRC
	if err := binary.Write(buf, binary.BigEndian, v.SSRC); err != nil {
		return nil, fmt.Errorf("failed to write SSRC: %w", err)
	}

	// Write Payload
	if _, err := buf.Write(v.Payload); err != nil {
		return nil, fmt.Errorf("failed to write payload: %w", err)
	}

	return buf.Bytes(), nil
}

func (v *VoicePacket) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read VersionFlags
	if err := binary.Read(buf, binary.BigEndian, &v.VersionFlags); err != nil {
		return fmt.Errorf("failed to read version flags: %w", err)
	}

	// Read PayloadType
	if err := binary.Read(buf, binary.BigEndian, &v.PayloadType); err != nil {
		return fmt.Errorf("failed to read payload type: %w", err)
	}

	// Read Seq
	if err := binary.Read(buf, binary.BigEndian, &v.Seq); err != nil {
		return fmt.Errorf("failed to read sequence number: %w", err)
	}

	// Read Timestamp
	if err := binary.Read(buf, binary.BigEndian, &v.Timestamp); err != nil {
		return fmt.Errorf("failed to read timestamp: %w", err)
	}

	// Read SSRC
	if err := binary.Read(buf, binary.BigEndian, &v.SSRC); err != nil {
		return fmt.Errorf("failed to read SSRC: %w", err)
	}

	// Read Payload
	v.Payload = make([]byte, buf.Len())
	if _, err := buf.Read(v.Payload); err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}

	return nil
}

func (v *VoicePacket) ToString() string {
	// Create a map to hold the JSON representation
	packetMap := map[string]interface{}{
		"VersionFlags": v.VersionFlags,
		"PayloadType":  v.PayloadType,
		"Seq":          v.Seq,
		"Timestamp":    v.Timestamp,
		"SSRC":         v.SSRC,
		"Payload":      v.Payload,
	}

	// Marshal the map into a JSON string
	jsonData, err := json.Marshal(packetMap)
	if err != nil {
		return fmt.Sprintf("error marshaling VoicePacket to JSON: %v", err)
	}

	return string(jsonData)
}

type VoicePayload struct {
	OpCode    VoiceOpCode `json:"op"`
	Data      any         `json:"d"`
	Seq       *int        `json:"s,omitempty"`
	EventName *string     `json:"t,omitempty"`
}

func (p *VoicePayload) ToString() string {
	jsonData, _ := json.Marshal(p)
	return string(jsonData)
}

type BinaryVoicePayload struct {
	SequenceNumber *uint16
	OpCode         uint8
	Data           any
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (p *BinaryVoicePayload) MarshalBinary() ([]byte, error) {
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

	// Write Payload
	payloadBytes, ok := p.Data.([]byte)
	if !ok {
		return nil, fmt.Errorf("data is not a byte slice")
	}
	if _, err := buf.Write(payloadBytes); err != nil {
		return nil, fmt.Errorf("failed to write payload: %w", err)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (p *BinaryVoicePayload) UnmarshalBinary(data []byte) error {
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
	if _, err := buf.Read(payload); err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}
	p.Data = payload

	return nil
}

// ProtocolVersion represents the protocol version
type ProtocolVersion uint16

const (
	Reserved ProtocolVersion = 0
	MLS10    ProtocolVersion = 1
)

// CipherSuite represents the cipher suite
type CipherSuite uint16

// HPKEPublicKey represents the HPKE public key
type HPKEPublicKey []byte

// LeafNode represents the leaf node
type LeafNode []byte

// Extension represents the extension
type Extension []byte

// KeyPackageTBS represents the KeyPackageTBS struct
type KeyPackageTBS struct {
	Version     ProtocolVersion
	CipherSuite CipherSuite
	InitKey     HPKEPublicKey
	LeafNode    LeafNode
	Extensions  []Extension
}

// KeyPackage represents the KeyPackage struct
type KeyPackage struct {
	KeyPackageTBS
	Signature []byte
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for KeyPackageTBS
func (kp *KeyPackageTBS) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write Version
	if err := binary.Write(buf, binary.BigEndian, kp.Version); err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}

	// Write CipherSuite
	if err := binary.Write(buf, binary.BigEndian, kp.CipherSuite); err != nil {
		return nil, fmt.Errorf("failed to write cipher suite: %w", err)
	}

	// Write InitKey length and InitKey
	if err := binary.Write(buf, binary.BigEndian, uint16(len(kp.InitKey))); err != nil {
		return nil, fmt.Errorf("failed to write init key length: %w", err)
	}
	if _, err := buf.Write(kp.InitKey); err != nil {
		return nil, fmt.Errorf("failed to write init key: %w", err)
	}

	// Write LeafNode length and LeafNode
	if err := binary.Write(buf, binary.BigEndian, uint16(len(kp.LeafNode))); err != nil {
		return nil, fmt.Errorf("failed to write leaf node length: %w", err)
	}
	if _, err := buf.Write(kp.LeafNode); err != nil {
		return nil, fmt.Errorf("failed to write leaf node: %w", err)
	}

	// Write Extensions length and Extensions
	if err := binary.Write(buf, binary.BigEndian, uint16(len(kp.Extensions))); err != nil {
		return nil, fmt.Errorf("failed to write extensions length: %w", err)
	}
	for _, ext := range kp.Extensions {
		if err := binary.Write(buf, binary.BigEndian, uint16(len(ext))); err != nil {
			return nil, fmt.Errorf("failed to write extension length: %w", err)
		}
		if _, err := buf.Write(ext); err != nil {
			return nil, fmt.Errorf("failed to write extension: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for KeyPackageTBS
func (kp *KeyPackageTBS) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read Version
	if err := binary.Read(buf, binary.BigEndian, &kp.Version); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}

	// Read CipherSuite
	if err := binary.Read(buf, binary.BigEndian, &kp.CipherSuite); err != nil {
		return fmt.Errorf("failed to read cipher suite: %w", err)
	}

	// Read InitKey length and InitKey
	var initKeyLen uint16
	if err := binary.Read(buf, binary.BigEndian, &initKeyLen); err != nil {
		return fmt.Errorf("failed to read init key length: %w", err)
	}
	kp.InitKey = make([]byte, initKeyLen)
	if _, err := buf.Read(kp.InitKey); err != nil {
		return fmt.Errorf("failed to read init key: %w", err)
	}

	// Read LeafNode length and LeafNode
	var leafNodeLen uint16
	if err := binary.Read(buf, binary.BigEndian, &leafNodeLen); err != nil {
		return fmt.Errorf("failed to read leaf node length: %w", err)
	}
	kp.LeafNode = make([]byte, leafNodeLen)
	if _, err := buf.Read(kp.LeafNode); err != nil {
		return fmt.Errorf("failed to read leaf node: %w", err)
	}

	// Read Extensions length and Extensions
	var extensionsLen uint16
	if err := binary.Read(buf, binary.BigEndian, &extensionsLen); err != nil {
		return fmt.Errorf("failed to read extensions length: %w", err)
	}
	kp.Extensions = make([]Extension, extensionsLen)
	for i := range kp.Extensions {
		var extLen uint16
		if err := binary.Read(buf, binary.BigEndian, &extLen); err != nil {
			return fmt.Errorf("failed to read extension length: %w", err)
		}
		kp.Extensions[i] = make([]byte, extLen)
		if _, err := buf.Read(kp.Extensions[i]); err != nil {
			return fmt.Errorf("failed to read extension: %w", err)
		}
	}

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for KeyPackage
func (kp *KeyPackage) MarshalBinary() ([]byte, error) {
	tbsData, err := kp.KeyPackageTBS.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.Write(tbsData); err != nil {
		return nil, fmt.Errorf("failed to write KeyPackageTBS: %w", err)
	}

	// Write Signature length and Signature
	if err := binary.Write(buf, binary.BigEndian, uint16(len(kp.Signature))); err != nil {
		return nil, fmt.Errorf("failed to write signature length: %w", err)
	}
	if _, err := buf.Write(kp.Signature); err != nil {
		return nil, fmt.Errorf("failed to write signature: %w", err)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for KeyPackage
func (kp *KeyPackage) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read KeyPackageTBS
	tbsData := make([]byte, buf.Len()-2) // Exclude the signature length
	if _, err := buf.Read(tbsData); err != nil {
		return fmt.Errorf("failed to read KeyPackageTBS: %w", err)
	}
	if err := kp.KeyPackageTBS.UnmarshalBinary(tbsData); err != nil {
		return err
	}

	// Read Signature length and Signature
	var signatureLen uint16
	if err := binary.Read(buf, binary.BigEndian, &signatureLen); err != nil {
		return fmt.Errorf("failed to read signature length: %w", err)
	}
	kp.Signature = make([]byte, signatureLen)
	if _, err := buf.Read(kp.Signature); err != nil {
		return fmt.Errorf("failed to read signature: %w", err)
	}

	return nil
}

// ProposalOrRefType represents the type of ProposalOrRef
type ProposalOrRefType uint8

const (
	ProposalOrRefTypeReserved  ProposalOrRefType = 0
	ProposalOrRefTypeProposal  ProposalOrRefType = 1
	ProposalOrRefTypeReference ProposalOrRefType = 2
)

// Proposal represents a proposal
type Proposal []byte

// ProposalRef represents a proposal reference
type ProposalRef []byte

// ProposalOrRef represents the ProposalOrRef struct
type ProposalOrRef struct {
	Type      ProposalOrRefType
	Proposal  Proposal
	Reference ProposalRef
}

// Commit represents the Commit struct
type Commit struct {
	Proposals []ProposalOrRef
	Path      []byte // optional
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for ProposalOrRef
func (p *ProposalOrRef) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write Type
	if err := binary.Write(buf, binary.BigEndian, p.Type); err != nil {
		return nil, fmt.Errorf("failed to write type: %w", err)
	}

	// Write Proposal or Reference based on Type
	switch p.Type {
	case ProposalOrRefTypeProposal:
		if err := binary.Write(buf, binary.BigEndian, uint16(len(p.Proposal))); err != nil {
			return nil, fmt.Errorf("failed to write proposal length: %w", err)
		}
		if _, err := buf.Write(p.Proposal); err != nil {
			return nil, fmt.Errorf("failed to write proposal: %w", err)
		}
	case ProposalOrRefTypeReference:
		if err := binary.Write(buf, binary.BigEndian, uint16(len(p.Reference))); err != nil {
			return nil, fmt.Errorf("failed to write reference length: %w", err)
		}
		if _, err := buf.Write(p.Reference); err != nil {
			return nil, fmt.Errorf("failed to write reference: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown proposal or reference type: %d", p.Type)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for ProposalOrRef
func (p *ProposalOrRef) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read Type
	if err := binary.Read(buf, binary.BigEndian, &p.Type); err != nil {
		return fmt.Errorf("failed to read type: %w", err)
	}

	// Read Proposal or Reference based on Type
	switch p.Type {
	case ProposalOrRefTypeProposal:
		var proposalLen uint16
		if err := binary.Read(buf, binary.BigEndian, &proposalLen); err != nil {
			return fmt.Errorf("failed to read proposal length: %w", err)
		}
		p.Proposal = make([]byte, proposalLen)
		if _, err := buf.Read(p.Proposal); err != nil {
			return fmt.Errorf("failed to read proposal: %w", err)
		}
	case ProposalOrRefTypeReference:
		var referenceLen uint16
		if err := binary.Read(buf, binary.BigEndian, &referenceLen); err != nil {
			return fmt.Errorf("failed to read reference length: %w", err)
		}
		p.Reference = make([]byte, referenceLen)
		if _, err := buf.Read(p.Reference); err != nil {
			return fmt.Errorf("failed to read reference: %w", err)
		}
	default:
		return fmt.Errorf("unknown proposal or reference type: %d", p.Type)
	}

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for Commit
func (c *Commit) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write Proposals length and Proposals
	if err := binary.Write(buf, binary.BigEndian, uint16(len(c.Proposals))); err != nil {
		return nil, fmt.Errorf("failed to write proposals length: %w", err)
	}
	for _, proposalOrRef := range c.Proposals {
		proposalOrRefData, err := proposalOrRef.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if _, err := buf.Write(proposalOrRefData); err != nil {
			return nil, fmt.Errorf("failed to write proposal or reference: %w", err)
		}
	}

	// Write Path length and Path if present
	if c.Path != nil {
		if err := binary.Write(buf, binary.BigEndian, uint16(len(c.Path))); err != nil {
			return nil, fmt.Errorf("failed to write path length: %w", err)
		}
		if _, err := buf.Write(c.Path); err != nil {
			return nil, fmt.Errorf("failed to write path: %w", err)
		}
	} else {
		if err := binary.Write(buf, binary.BigEndian, uint16(0)); err != nil {
			return nil, fmt.Errorf("failed to write path length: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Commit
func (c *Commit) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read Proposals length and Proposals
	var proposalsLen uint16
	if err := binary.Read(buf, binary.BigEndian, &proposalsLen); err != nil {
		return fmt.Errorf("failed to read proposals length: %w", err)
	}
	c.Proposals = make([]ProposalOrRef, proposalsLen)
	for i := range c.Proposals {
		var proposalOrRef ProposalOrRef
		proposalOrRefData := make([]byte, buf.Len())
		if _, err := buf.Read(proposalOrRefData); err != nil {
			return fmt.Errorf("failed to read proposal or reference: %w", err)
		}
		if err := proposalOrRef.UnmarshalBinary(proposalOrRefData); err != nil {
			return err
		}
		c.Proposals[i] = proposalOrRef
	}

	// Read Path length and Path if present
	var pathLen uint16
	if err := binary.Read(buf, binary.BigEndian, &pathLen); err != nil {
		return fmt.Errorf("failed to read path length: %w", err)
	}
	if pathLen > 0 {
		c.Path = make([]byte, pathLen)
		if _, err := buf.Read(c.Path); err != nil {
			return fmt.Errorf("failed to read path: %w", err)
		}
	} else {
		c.Path = nil
	}

	return nil
}
