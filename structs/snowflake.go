package structs

import (
	"encoding/json"
	"errors"
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
	var id uint64
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	} else if id == 0 {
		return errors.New("ID cannot be 0")
	}

	s.ID = id
	s.deconstructSnowflake()
	return nil
}

func (s *Snowflake) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ID)
}
