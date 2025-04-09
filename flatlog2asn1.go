package flatlog2asn1

import (
	"encoding/asn1"
)

type Severity asn1.Enumerated

type UnixtimeUs = int64

const (
	SeverityUnspecified Severity = 0
	SeverityTrace       Severity = 1
	SeverityDebug       Severity = 5
	SeverityInfo        Severity = 9
	SeverityWarn        Severity = 13
	SeverityError       Severity = 17
	SeverityFatal       Severity = 21
)

// 64-bit message id(use 0 as undefined id).
type MessageId = int64

type KeyValPair[T any] struct {
	Key string `asn1:"ia5"`
	Val T
}

type IntPair KeyValPair[int64]

type BooPair KeyValPair[bool]

type StrPair struct {
	Key string `asn1:"ia5"`
	Val string `asn1:"utf8"`
}

type IntMap []IntPair

type BooMap []BooPair

type StrMap []StrPair

type GenericItems struct {
	BooMap
	IntMap
	StrMap
}

type Message = string

type LogItem struct {
	UnixtimeUs
	Severity
	Message `asn1:"utf8"`
	MessageId
	Resource   GenericItems
	Attributes GenericItems
	Extra      StrMap
}

func (l LogItem) ToDerBytes() ([]byte, error) { return asn1.Marshal(l) }
