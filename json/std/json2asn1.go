package json2asn1

import (
	"errors"
	"fmt"

	fa "github.com/takanoriyanagitani/go-flatlog2asn1"
)

var ErrInvalidArgument error = errors.New("invalid argument")

type JsonLog map[string]any

type LogConfig struct {
	TimestampName string
	SeverityName  string
	MessageName   string
	IdName        string
	ResourceKeys  map[string]struct{}
	AttributeKeys map[string]struct{}

	ParserConfig
}

var LogConfigDefault LogConfig = LogConfig{
	TimestampName: "timestamp",
	SeverityName:  "severity",
	MessageName:   "message",
	IdName:        "id",
	ResourceKeys:  map[string]struct{}{},
	AttributeKeys: map[string]struct{}{},

	ParserConfig: ParserConfigDefault,
}

type Parser func(JsonLog) (fa.LogItem, error)

type ParserConfig struct {
	TimeFromAny
	SeverityFromAny
}

var ParserConfigDefault ParserConfig = ParserConfig{
	TimeFromAny:     TimeFromAnyDefault,
	SeverityFromAny: SeverityFromAnyDefault,
}

func (c LogConfig) ToParser() Parser {
	return func(j JsonLog) (fa.LogItem, error) {
		var item fa.LogItem

		var timestamp any = j[c.TimestampName]
		t, e := c.ParserConfig.TimeFromAny(timestamp)
		if nil != e {
			return item, e
		}
		item.UnixtimeUs = t.UnixMicro()

		var severity any = j[c.SeverityName]
		s := c.ParserConfig.SeverityFromAny(severity)
		item.Severity = s

		var message any = j[c.MessageName]
		msg, ok := message.(string)
		if ok {
			item.Message = msg
		}

		var identifier any = j[c.IdName]
		switch id := identifier.(type) {
		case float64:
			item.MessageId = int64(id)
		default:
		}

		for key, val := range j {
			switch key {
			case c.TimestampName:
				continue
			case c.SeverityName:
				continue
			case c.MessageName:
				continue
			case c.IdName:
				continue
			default:
			}

			_, rkey := c.ResourceKeys[key]
			if rkey {
				switch t := val.(type) {

				case bool:
					item.Resource.BooMap = append(
						item.Resource.BooMap,
						fa.BooPair{
							Key: key,
							Val: t,
						},
					)

				case string:
					item.Resource.StrMap = append(
						item.Resource.StrMap,
						fa.StrPair{
							Key: key,
							Val: t,
						},
					)

				case float64:
					item.Resource.IntMap = append(
						item.Resource.IntMap,
						fa.IntPair{
							Key: key,
							Val: int64(t),
						},
					)

				default:
					item.Extra = append(
						item.Extra,
						fa.StrPair{
							Key: key,
							Val: fmt.Sprintf("%v", val),
						},
					)
				}

				continue
			}

			_, akey := c.AttributeKeys[key]
			if akey {
				switch t := val.(type) {

				case bool:
					item.Attributes.BooMap = append(
						item.Attributes.BooMap,
						fa.BooPair{
							Key: key,
							Val: t,
						},
					)

				case string:
					item.Attributes.StrMap = append(
						item.Attributes.StrMap,
						fa.StrPair{
							Key: key,
							Val: t,
						},
					)

				case float64:
					item.Attributes.IntMap = append(
						item.Attributes.IntMap,
						fa.IntPair{
							Key: key,
							Val: int64(t),
						},
					)

				default:
					item.Extra = append(
						item.Extra,
						fa.StrPair{
							Key: key,
							Val: fmt.Sprintf("%v", val),
						},
					)
				}

				continue
			}

			item.Extra = append(
				item.Extra,
				fa.StrPair{Key: key, Val: fmt.Sprintf("%v", val)},
			)
		}

		return item, nil
	}
}

func (c LogConfig) WithResourceKeys(keys map[string]struct{}) LogConfig {
	c.ResourceKeys = keys
	return c
}

func (c LogConfig) WithAttributeKeys(keys map[string]struct{}) LogConfig {
	c.AttributeKeys = keys
	return c
}
