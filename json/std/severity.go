package json2asn1

import (
	fa "github.com/takanoriyanagitani/go-flatlog2asn1"
)

type SeverityFromAny func(any) fa.Severity

type SeverityMap map[int64]fa.Severity

var SeverityMapDefault SeverityMap = map[int64]fa.Severity{
	0:  fa.SeverityUnspecified,
	1:  fa.SeverityTrace,
	5:  fa.SeverityDebug,
	9:  fa.SeverityInfo,
	13: fa.SeverityWarn,
	17: fa.SeverityError,
	21: fa.SeverityFatal,
}

type SeverityMapStr map[string]fa.Severity

var SeverityMapStrDefault SeverityMapStr = map[string]fa.Severity{
	"TRACE": fa.SeverityTrace,
	"DEBUG": fa.SeverityDebug,
	"INFO":  fa.SeverityInfo,
	"WARN":  fa.SeverityWarn,
	"ERROR": fa.SeverityError,
	"FATAL": fa.SeverityFatal,
}

type SeverityFromInt func(int64) fa.Severity

type SeverityFromStr func(string) fa.Severity

func (m SeverityMap) ToSeverityFromInt() SeverityFromInt {
	return func(s int64) fa.Severity {
		val, found := m[s]
		switch found {
		case true:
			return val
		default:
			return fa.SeverityUnspecified
		}
	}
}

func (m SeverityMapStr) ToSeverityFromStr() SeverityFromStr {
	return func(s string) fa.Severity {
		val, found := m[s]
		switch found {
		case true:
			return val
		default:
			return fa.SeverityUnspecified
		}
	}
}

var SeverityFromIntDefault SeverityFromInt = SeverityMapDefault.
	ToSeverityFromInt()

var SeverityFromStrDefault SeverityFromStr = SeverityMapStrDefault.
	ToSeverityFromStr()

type SeverityConfig struct {
	SeverityFromInt
	SeverityFromStr
}

func (s SeverityConfig) ToSeverityFromAny() SeverityFromAny {
	return func(a any) fa.Severity {
		switch t := a.(type) {
		case float64:
			return s.SeverityFromInt(int64(t))
		case string:
			return s.SeverityFromStr(t)
		default:
			return fa.SeverityUnspecified
		}
	}
}

var SeverityConfigDefault SeverityConfig = SeverityConfig{
	SeverityFromInt: SeverityFromIntDefault,
	SeverityFromStr: SeverityFromStrDefault,
}

var SeverityFromAnyDefault SeverityFromAny = SeverityConfigDefault.
	ToSeverityFromAny()
