package structs

import (
	"encoding/json"
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

func NewSnowflake(id uint64) Snowflake {
	var sf Snowflake
	sf.ID = id
	return sf.deconstructSnowflake()
}

func (s Snowflake) deconstructSnowflake() Snowflake {
	timestamp := time.UnixMilli(int64((s.ID >> 22) + Epoch))
	workerID := uint8(s.ID & 0x3E0000 >> 17)
	processID := uint8(s.ID & 0x1F000 >> 12)
	increment := uint8(s.ID & 0xFFF)

	s.Timestamp = &timestamp
	s.WorkerID = &workerID
	s.ProcessID = &processID
	s.Increment = &increment

	return s
}

func (s *Snowflake) UnmarshalJSON(data []byte) error {
	type Alias Snowflake
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if s.ID != 0 {
		*s = s.deconstructSnowflake()
	}
	return nil
}
