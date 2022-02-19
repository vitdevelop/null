package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Timestamp is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Timestamp struct {
	sql.NullTime
}

// Value implements the driver Valuer interface.
func (t Timestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewTimestamp creates a new Timestamp.
func NewTimestamp(t time.Time, valid bool) Timestamp {
	return Timestamp{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

// TimestampFrom creates a new Timestamp that will always be valid.
func TimestampFrom(t time.Time) Timestamp {
	return NewTimestamp(t, true)
}

// TimestampFromPtr creates a new Timestamp that will be null if t is nil.
func TimestampFromPtr(t *time.Time) Timestamp {
	if t == nil {
		return NewTimestamp(time.Time{}, false)
	}
	return NewTimestamp(*t, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Timestamp) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.UnixNano() / int64(time.Millisecond))
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		t.Valid = false
		return nil
	}

	var value Int64
	err := value.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	t.Time = time.UnixMilli(0).UTC().Add(time.Duration(value.Int64 * int64(time.Millisecond)))
	t.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It returns an empty string if invalid, otherwise time.Time's MarshalText.
func (t Timestamp) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	n := Int64From(t.Time.UnixMilli())
	return n.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It has backwards compatibility with v3 in that the string "null" is considered equivalent to an empty string
// and unmarshaling will succeed. This may be removed in a future version.
func (t *Timestamp) UnmarshalText(text []byte) error {
	n := Int64From(0)
	err := n.UnmarshalText(text)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	if !n.Valid {
		return nil
	}

	t.Time = time.UnixMilli(0).UTC().Add(time.Duration(n.Int64 * int64(time.Millisecond)))
	t.Valid = true
	return nil
}

// SetValid changes this Timestamp's value and sets it to be non-null.
func (t *Timestamp) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Timestamp's value, or a nil pointer if this Timestamp is null.
func (t Timestamp) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for invalid Times, hopefully for future omitempty support.
// A non-null Timestamp with a zero value will not be considered zero.
func (t Timestamp) IsZero() bool {
	return !t.Valid
}

// Equal returns true if both Timestamp objects encode the same time or are both null.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
func (t Timestamp) Equal(other Timestamp) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time.Equal(other.Time))
}

// ExactEqual returns true if both Timestamp objects are equal or both null.
// ExactEqual returns false for times that are in different locations or
// have a different monotonic clock reading.
func (t Timestamp) ExactEqual(other Timestamp) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time == other.Time)
}
