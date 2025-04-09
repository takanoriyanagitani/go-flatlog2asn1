package json2asn1

import (
	"fmt"
	"time"
)

type TimeParser func(string) (time.Time, error)

func TimeParserFromLayout(layout string) TimeParser {
	return func(val string) (time.Time, error) {
		return time.Parse(layout, val)
	}
}

var TimeParserRfc3339 TimeParser = TimeParserFromLayout(time.RFC3339Nano)

type TimeFromInt func(int64) time.Time

var UnixtimeUsToInt TimeFromInt = time.UnixMicro

type TimeFromAny func(any) (time.Time, error)

type TimeParserConfig struct {
	TimeParser
	TimeFromInt
}

func (t TimeParserConfig) ToTimeFromAny() TimeFromAny {
	return func(a any) (time.Time, error) {
		switch i := a.(type) {
		case float64:
			return t.TimeFromInt(int64(i)), nil
		case string:
			return t.TimeParser(i)
		default:
			return time.UnixMicro(0), fmt.Errorf(
				"%w: %v",
				ErrInvalidArgument,
				i,
			)
		}
	}
}

var TimeParserConfigDefault TimeParserConfig = TimeParserConfig{
	TimeParser:  TimeParserRfc3339,
	TimeFromInt: UnixtimeUsToInt,
}

var TimeFromAnyDefault TimeFromAny = TimeParserConfigDefault.
	ToTimeFromAny()
