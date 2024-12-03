package gametype

import (
	"strconv"
	"strings"
	"time"
)

type DayZ struct {
	BattlEye       bool
	NoThirdPerson  bool
	External       bool
	PrivateHive    bool
	Shard          string
	PlayersQueue   uint8
	TimeDayAccel   float64
	TimeNightAccel float64
	Modded         bool
	GamePort       uint16
	Whitelist      bool
	FlePatching    bool
	DLC            bool
	Time           time.Duration
}

func (d *DayZ) Parse(gametype string) {
	tags := strings.Split(gametype, ",")

	for _, tag := range tags {
		switch {
		case tag == "battleye":
			d.BattlEye = true

		case tag == "no3rd":
			d.NoThirdPerson = true

		case tag == "external":
			d.External = true

		case tag == "privHive":
			d.PrivateHive = true

		case strings.HasPrefix(tag, "shard"):
			d.Shard = strings.TrimPrefix(tag, "shard")

		case strings.HasPrefix(tag, "lqs"):
			if num, err := strconv.ParseUint(tag[3:], 10, 8); err == nil {
				d.PlayersQueue = uint8(num)
			}

		case strings.HasPrefix(tag, "etm"):
			if num, err := strconv.ParseFloat(tag[3:], 64); err == nil {
				d.TimeDayAccel = num
			}

		case strings.HasPrefix(tag, "entm"):
			if num, err := strconv.ParseFloat(tag[4:], 64); err == nil {
				d.TimeNightAccel = num
			}

		case tag == "mod":
			d.Modded = true

		case strings.HasPrefix(tag, "port"):
			if num, err := strconv.ParseUint(tag[4:], 10, 8); err == nil {
				d.GamePort = uint16(num)
			}

		case tag == "whitelisting":
			d.Whitelist = true

		case tag == "allowedFilePatching":
			d.FlePatching = true

		case tag == "isDLC":
			d.DLC = true

		case strings.Contains(tag, ":"):
			if t, err := time.ParseDuration(tag[:2] + "h" + tag[3:] + "m"); err == nil {
				d.Time = t
			}
		}
	}
}
