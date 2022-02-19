package zero

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

var (
	timestampString   = "1356124881000"
	timestampJSON     = []byte(timestampString)
	zeroTimestampJSON = []byte("0")
	zeroTimestamp     = "0"
)

func TestUnmarshalTimestampJSON(t *testing.T) {
	var ti Timestamp
	err := json.Unmarshal(timeObject, &ti)
	if err == nil {
		panic("expected error")
	}

	var blank Timestamp
	err = json.Unmarshal(blankTimeJSON, &blank)
	maybePanic(err)
	assertNullTimestamp(t, blank, "blank time json")

	var zero Timestamp
	err = json.Unmarshal(zeroTimestampJSON, &zero)
	maybePanic(err)
	if !zero.Valid {
		t.Error("zero timestamp is invalid")
	}

	var fromObject Timestamp
	err = json.Unmarshal(timeObject, &fromObject)
	if err == nil {
		panic("expected error")
	}

	var null Timestamp
	err = json.Unmarshal(nullObject, &null)
	if err == nil {
		panic("expected error")
	}

	var invalid Timestamp
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}

	var bad Timestamp
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}

	var wrongString Timestamp
	err = json.Unmarshal(stringJSON, &wrongString)
	if err == nil {
		t.Errorf("expected error: wrong string JSON")
	}
}

func TestMarshalTimestamp(t *testing.T) {
	ti := TimestampFrom(timeValue1)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(timestampJSON), "non-empty json marshal")

	null := TimestampFromPtr(nil)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, string(zeroTimestampJSON), "empty json marshal")
}

func TestUnmarshalTimestampText(t *testing.T) {
	ti := TimestampFrom(timeValue1)
	txt, err := ti.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, timestampString, "marshal text")

	var unmarshal Timestamp
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertTimestamp(t, unmarshal, "unmarshal text")

	var null Timestamp
	err = null.UnmarshalText(nullJSON)
	maybePanic(err)
	assertNullTimestamp(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, zeroTimestamp, "marshal null text")

	var invalid Timestamp
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTimestamp(t, invalid, "bad string")
}

func TestTimestampFrom(t *testing.T) {
	ti := TimestampFrom(timeValue1)
	assertTimestamp(t, ti, "TimestampFrom() time.Timestamp")

	var nt time.Time
	null := TimestampFrom(nt)
	assertNullTimestamp(t, null, "TimestampFrom() empty time.Timestamp")
}

func TestTimestampFromPtr(t *testing.T) {
	ti := TimestampFromPtr(&timeValue1)
	assertTimestamp(t, ti, "TimestampFromPtr() time")

	null := TimestampFromPtr(nil)
	assertNullTimestamp(t, null, "TimestampFromPtr(nil)")
}

func TestTimestampSetValid(t *testing.T) {
	var ti time.Time
	change := TimestampFrom(ti)
	assertNullTimestamp(t, change, "SetValid()")
	change.SetValid(timeValue1)
	assertTimestamp(t, change, "SetValid()")
}

func TestTimestampPointer(t *testing.T) {
	ti := TimestampFrom(timeValue1)
	ptr := ti.Ptr()
	if *ptr != timeValue1 {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, timeValue1)
	}

	var nt time.Time
	null := TimestampFrom(nt)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimestampScan(t *testing.T) {
	var ti Timestamp
	err := ti.Scan(timeValue1)
	maybePanic(err)
	assertTimestamp(t, ti, "scanned time")

	var null Timestamp
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTimestamp(t, null, "scanned null")

	var wrong Timestamp
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
}

func TestTimestampValue(t *testing.T) {
	ti := TimestampFrom(timeValue1)
	_, err := ti.Value()
	maybePanic(err)
	if ti.Time != timeValue1 {
		t.Errorf("bad time.Timestamp value: %v ≠ %v", ti.Time, timeValue1)
	}

	var v driver.Value
	var nt time.Time
	zero := TimestampFrom(nt)
	v, err = zero.Value()
	maybePanic(err)
	if v != nil {
		t.Errorf("bad %s time.Timestamp value: %v ≠ %v", "zero", v, nil)
	}
}

func TestTimestampValueOrZero(t *testing.T) {
	valid := TimestampFrom(timeValue1)
	if valid.ValueOrZero() != valid.Time || valid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := valid
	invalid.Valid = false
	if !invalid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestTimestampIsZero(t *testing.T) {
	str := TimestampFrom(timeValue1)
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	zero := TimestampFrom(time.Time{})
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	null := TimestampFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestTimestampEqual(t *testing.T) {
	t1 := NewTimestamp(timeValue1, false)
	t2 := NewTimestamp(timeValue2, false)
	assertTimestampEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, false)
	t2 = NewTimestamp(timeValue3, false)
	assertTimestampEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue2, true)
	assertTimestampEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, false)
	t2 = NewTimestamp(time.Time{}, true)
	assertTimestampEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue1, true)
	assertTimestampEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue2, false)
	assertTimestampEqualIsFalse(t, t1, t2)

	t1 = NewTimestamp(timeValue1, false)
	t2 = NewTimestamp(timeValue2, true)
	assertTimestampEqualIsFalse(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue3, true)
	assertTimestampEqualIsFalse(t, t1, t2)
}

func TestTimestampExactEqual(t *testing.T) {
	t1 := NewTimestamp(timeValue1, false)
	t2 := NewTimestamp(timeValue1, false)
	assertTimestampExactEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, false)
	t2 = NewTimestamp(timeValue2, false)
	assertTimestampExactEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue1, true)
	assertTimestampExactEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, false)
	t2 = NewTimestamp(time.Time{}, true)
	assertTimestampExactEqualIsTrue(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue1, false)
	assertTimestampExactEqualIsFalse(t, t1, t2)

	t1 = NewTimestamp(timeValue1, false)
	t2 = NewTimestamp(timeValue1, true)
	assertTimestampExactEqualIsFalse(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue2, true)
	assertTimestampExactEqualIsFalse(t, t1, t2)

	t1 = NewTimestamp(timeValue1, true)
	t2 = NewTimestamp(timeValue3, true)
	assertTimestampExactEqualIsFalse(t, t1, t2)
}

func assertTimestamp(t *testing.T, ti Timestamp, from string) {
	if ti.Time != timeValue1 {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Time, timeValue1)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTimestamp(t *testing.T, ti Timestamp, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertTimestampEqualIsTrue(t *testing.T, a, b Timestamp) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Timestamp{%v, Valid:%t} and Timestamp{%v, Valid:%t} should return true", a.Time, a.Valid, b.Time, b.Valid)
	}
}

func assertTimestampEqualIsFalse(t *testing.T, a, b Timestamp) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Timestamp{%v, Valid:%t} and Timestamp{%v, Valid:%t} should return false", a.Time, a.Valid, b.Time, b.Valid)
	}
}

func assertTimestampExactEqualIsTrue(t *testing.T, a, b Timestamp) {
	t.Helper()
	if !a.ExactEqual(b) {
		t.Errorf("ExactEqual() of Timestamp{%v, Valid:%t} and Timestamp{%v, Valid:%t} should return true", a.Time, a.Valid, b.Time, b.Valid)
	}
}

func assertTimestampExactEqualIsFalse(t *testing.T, a, b Timestamp) {
	t.Helper()
	if a.ExactEqual(b) {
		t.Errorf("ExactEqual() of Timestamp{%v, Valid:%t} and Timestamp{%v, Valid:%t} should return false", a.Time, a.Valid, b.Time, b.Valid)
	}
}
