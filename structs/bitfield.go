package structs

import (
	"encoding/json"
	"fmt"
)

type Bitfield[T any] []T

func (b *Bitfield[T]) UnmarshalJSON(data []byte) error {
	var flags int64
	if err := json.Unmarshal(data, &flags); err != nil {
		return err
	}

	parsedFlags, err := convert[T](flags)
	if err != nil {
		return err
	}

	*b = b.ParseFlags(parsedFlags)
	return nil
}

func (b *Bitfield[T]) MarshalJSON() ([]byte, error) {
	flags := b.GetFlags()
	return json.Marshal(flags)
}

func (b *Bitfield[T]) GetFlags() int64 {
	var flag int64
	for _, i := range *b {
		v, err := toInt64(i)
		if err != nil {
			return 0
		}
		flag |= v
	}
	return flag
}

func (b *Bitfield[T]) ParseFlags(flags T) []T {
	var parsedFlags []T
	flagValue := any(flags).(T)
	for _, f := range *b {
		if any(f).(int64)&any(flagValue).(int64) != 0 {
			parsedFlags = append(parsedFlags, f)
		}
	}
	return parsedFlags
}

func (b *Bitfield[T]) SetFlags(flags T) {
	flagsVal := b.ParseFlags(flags)
	*b = any(flagsVal).([]T)
}

func (b *Bitfield[T]) AddFlag(flag T) {
	*b = append(*b, any(flag).(T))
}

func (b *Bitfield[T]) RemoveFlag(flag T) {
	for i, f := range *b {
		if any(f).(int64) == any(flag).(int64) {
			*b = append((*b)[:i], (*b)[i+1:]...)
		}
	}
}

func (b *Bitfield[T]) ToString() string {
	var result string
	for _, f := range *b {
		result += fmt.Sprintf("%d", any(f).(int64))
	}
	return result
}

// convert converts an int64 to type T
func convert[T any](value int64) (T, error) {
	var t T
	switch any(t).(type) {
	case ActivityFlag:
		return any(ActivityFlag(value)).(T), nil
	case ApplicationFlag:
		return any(ApplicationFlag(value)).(T), nil
	case ChannelFlag:
		return any(ChannelFlag(value)).(T), nil
	case GuildMemberFlag:
		return any(GuildMemberFlag(value)).(T), nil
	case SystemChannelFlag:
		return any(SystemChannelFlag(value)).(T), nil
	case AttachmentFlag:
		return any(AttachmentFlag(value)).(T), nil
	case MessageFlag:
		return any(MessageFlag(value)).(T), nil
	case RoleFlag:
		return any(RoleFlag(value)).(T), nil
	case UserFlag:
		return any(UserFlag(value)).(T), nil
	case Permission:
		return any(Permission(value)).(T), nil
	default:
		return t, fmt.Errorf("unsupported type conversion from int64 to %T", t)
	}
}

func toInt64[T any](value T) (int64, error) {
	switch v := any(value).(type) {
	case ActivityFlag:
		return int64(v), nil
	case ApplicationFlag:
		return int64(v), nil
	case ChannelFlag:
		return int64(v), nil
	case GuildMemberFlag:
		return int64(v), nil
	case SystemChannelFlag:
		return int64(v), nil
	case AttachmentFlag:
		return int64(v), nil
	case MessageFlag:
		return int64(v), nil
	case RoleFlag:
		return int64(v), nil
	case UserFlag:
		return int64(v), nil
	case Permission:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type conversion from %T to int64", value)
	}
}
