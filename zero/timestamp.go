package zero

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Timestamp is a nullable time.Time.
// JSON marshals to the zero value for time.Time if null.
// Considered to be null to SQL if zero.
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

// TimestampFrom creates a new Timestamp that will
// be null if t is the zero value.
func TimestampFrom(t time.Time) Timestamp {
	return NewTimestamp(t, !t.IsZero())
}

// TimestampFromPtr creates a new Timestamp that will
// be null if t is nil or *t is the zero value.
func TimestampFromPtr(t *time.Time) Timestamp {
	if t == nil {
		return NewTimestamp(time.Time{}, false)
	}
	return TimestampFrom(*t)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Timestamp) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// MarshalJSON implements json.Marshaler.
// It will encode the zero value of time.Time
// if this time is invalid.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("0"), nil
	}
	return json.Marshal(t.Time.UnixNano() / int64(time.Millisecond))
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "null", `""`:
		t.Valid = false
		return nil
	}

	var value Int64
	err := value.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	t.Time = time.UnixMilli(0).UTC().Add(time.Duration(value.Int64 * int64(time.Millisecond)))
	t.Valid = !t.Time.IsZero()
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode to an empty time.Time if invalid.
func (t Timestamp) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("0"), nil
	}
	n := Int64From(t.Time.UnixMilli())
	return n.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It has compatibility with the null package in that it will accept empty strings as invalid values,
// which will be unmarshaled to an invalid zero value.
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

// SetValid changes this Timestamp's value and
// sets it to be non-null.
func (t *Timestamp) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Timestamp's value,
// or a nil pointer if this Timestamp is zero.
func (t Timestamp) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for null or zero Times, for potential future omitempty support.
func (t Timestamp) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// Equal returns true if both Timestamp objects encode the same time or are both are either null or zero.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
func (t Timestamp) Equal(other Timestamp) bool {
	return t.ValueOrZero().Equal(other.ValueOrZero())
}

// ExactEqual returns true if both Timestamp objects are equal or both are either null or zero.
// ExactEqual returns false for times that are in different locations or
// have a different monotonic clock reading.
func (t Timestamp) ExactEqual(other Timestamp) bool {
	return t.ValueOrZero() == other.ValueOrZero()
}
