package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

// Date accepts date-only strings (YYYY-MM-DD) and RFC3339 timestamps in JSON.
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Time = time.Time{}
		return nil
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		d.Time = time.Time{}
		return nil
	}

	if t, err := time.Parse(dateLayout, raw); err == nil {
		d.Time = t
		return nil
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return fmt.Errorf("invalid date format for release_date, use YYYY-MM-DD or RFC3339")
	}

	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(d.Time.Format(dateLayout))
}
