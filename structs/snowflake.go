package structs

import "time"

const Epoch = 1420070400000

type Snowflake struct {
	ID        uint64
	Timestamp *time.Time
	WorkerID  *uint8
	ProcessID *uint8
	Increment *uint8
}

func NewSnowflake(id uint64) Snowflake {
	var sf Snowflake
	sf.ID = id
	return sf.deconstructSnowflake()
}

func (s Snowflake) deconstructSnowflake() Snowflake {
	*s.Timestamp = time.UnixMilli(int64((s.ID >> 22) + Epoch))
	*s.WorkerID = uint8(s.ID & 0x3E0000 >> 17)
	*s.ProcessID = uint8(s.ID & 0x1F000 >> 12)
	*s.Increment = uint8(s.ID & 0xFFF)
	return s
}
