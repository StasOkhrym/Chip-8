package cpu

import "math"

type Integer interface {
	uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64
}

// OverflowAdd performs addition on two integers of any type.
// It returns the result and a remainder if overflow occurs.
func OverflowAdd[Int Integer](a, b Int) (result Int, overflow bool) {
	sum := a + b
	var max uint
	var min int

	switch any(a).(type) {
	case uint8:
		max = math.MaxUint8
		min = 0
	case uint16:
		max = math.MaxUint16
		min = 0
	case uint32:
		max = math.MaxUint32
		min = 0
	case uint64:
		max = math.MaxUint64
		min = 0
	case int8:
		max = math.MaxInt8
		min = math.MinInt8
	case int16:
		max = math.MaxInt16
		min = math.MinInt16
	case int32:
		max = math.MaxInt32
		min = math.MinInt32
	case int64:
		max = math.MaxInt64
		min = math.MinInt64
	}

	if b > 0 {
		if a > Int(max)-b {
			return sum, true
		}
	} else {
		if a < Int(min)-b {
			return sum, true
		}
	}
	return sum, false
}

// OverflowSub performs subtraction on two integers of any type.
// It returns the result and a boolean indicating whether underflow occurs.
func OverflowSub[Int Integer](a, b Int) (result Int, underflow bool) {
	sub := a - b

	var max uint
	var min int

	switch any(a).(type) {
	case uint8:
		max = math.MaxUint8
		min = 0
	case uint16:
		max = math.MaxUint16
		min = 0
	case uint32:
		max = math.MaxUint32
		min = 0
	case uint64:
		max = math.MaxUint64
		min = 0
	case int8:
		max = math.MaxInt8
		min = math.MinInt8
	case int16:
		max = math.MaxInt16
		min = math.MinInt16
	case int32:
		max = math.MaxInt32
		min = math.MinInt32
	case int64:
		max = math.MaxInt64
		min = math.MinInt64
	}

	if b > 0 {
		if a < Int(min)+b {
			return sub, true
		}
	} else {
		if a > Int(max)+b {
			return sub, true
		}
	}
	return sub, false
}
