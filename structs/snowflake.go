package structs

import (
	"encoding/json"
	"strconv"
	"time"
)

const Epoch = 1420070400000

type Snowflake struct {
	ID        uint64     `json:"id"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	WorkerID  *uint8     `json:"worker_id,omitempty"`
	ProcessID *uint8     `json:"process_id,omitempty"`
	Increment *uint8     `json:"increment,omitempty"`
}

func NewSnowflake(id uint64) *Snowflake {
	s := Snowflake{ID: id}
	s.deconstructSnowflake()
	return &s
}

func (s *Snowflake) Equals(other Snowflake) bool {
	if s == nil {
		return false
	}
	return s.ID == other.ID
}

func (s *Snowflake) ToString() string {
	return strconv.Itoa(int(s.ID))
}

func (s *Snowflake) deconstructSnowflake() {
	timestamp := time.UnixMilli(int64((s.ID >> 22) + Epoch))
	workerID := uint8(s.ID & 0x3E0000 >> 17)
	processID := uint8(s.ID & 0x1F000 >> 12)
	increment := uint8(s.ID & 0xFFF)

	s.Timestamp = &timestamp
	s.WorkerID = &workerID
	s.ProcessID = &processID
	s.Increment = &increment
}

func (s *Snowflake) UnmarshalJSON(data []byte) error {
	// try to parse the id field as a uint64, if it fails try to parse it as a string
	var id uint64
	if err := json.Unmarshal(data, &id); err != nil {
		var strId string
		if err := json.Unmarshal(data, &strId); err != nil {
			return err
		}
		id, err = strconv.ParseUint(strId, 10, 64)
		if err != nil {
			return err
		}
	} else if id == 0 {
		s.ID = id
		return nil
	}

	s.ID = id
	s.deconstructSnowflake()
	return nil
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ID)
}
