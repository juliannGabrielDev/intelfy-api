package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SongDuration stores song duration and accepts JSON values in seconds (number)
// or Go duration string format (for example: "3m30s").
type SongDuration struct {
	time.Duration
}

func (d *SongDuration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Duration = 0
		return nil
	}

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		d.Duration = 0
		return nil
	}

	if trimmed[0] == '"' {
		var raw string
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}

		raw = strings.TrimSpace(raw)
		if raw == "" {
			d.Duration = 0
			return nil
		}

		if n, err := strconv.ParseFloat(raw, 64); err == nil {
			d.Duration = time.Duration(n * float64(time.Second))
			return nil
		}

		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("invalid duration format, use seconds as number or a duration string like 3m30s")
		}

		d.Duration = parsed
		return nil
	}

	seconds, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return fmt.Errorf("invalid duration format, use seconds as number or a duration string like 3m30s")
	}

	d.Duration = time.Duration(seconds * float64(time.Second))
	return nil
}

func (d SongDuration) MarshalJSON() ([]byte, error) {
	seconds := d.Duration.Seconds()
	return json.Marshal(seconds)
}

func (d SongDuration) Seconds() float64 {
	return d.Duration.Seconds()
}
