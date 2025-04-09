package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	fa "github.com/takanoriyanagitani/go-flatlog2asn1"
	js "github.com/takanoriyanagitani/go-flatlog2asn1/json/std"
	. "github.com/takanoriyanagitani/go-flatlog2asn1/util"
)

var logConfig js.LogConfig = js.
	LogConfigDefault.
	WithResourceKeys(map[string]struct{}{
		"process.pid":     {},
		"service.version": {},
	}).
	WithAttributeKeys(map[string]struct{}{
		"debug":               {},
		"cache.hit":           {},
		"db.upserted":         {},
		"http.content.length": {},
		"http.status.code":    {},
		"cpu.core":            {},
		"user.id":             {},
		"http.request.method": {},
	})

var json2asn1 js.Parser = logConfig.ToParser()

func ReaderToBytesLimited(limit int64) func(io.Reader) IO[[]byte] {
	return Lift(func(rdr io.Reader) ([]byte, error) {
		limited := &io.LimitedReader{
			R: rdr,
			N: limit,
		}
		var buf bytes.Buffer
		_, e := io.Copy(&buf, limited)
		return buf.Bytes(), e
	})
}

const jsonSizeLimit int64 = 1048576

var jsonBytesStdin IO[[]byte] = ReaderToBytesLimited(jsonSizeLimit)(os.Stdin)

func BytesToJsonMap(b []byte) (map[string]any, error) {
	var ret map[string]any
	e := json.Unmarshal(b, &ret)
	return ret, e
}

var jsonMap IO[map[string]any] = Bind(
	jsonBytesStdin,
	Lift(BytesToJsonMap),
)

var jsonLog IO[js.JsonLog] = Bind(
	jsonMap,
	Lift(func(m map[string]any) (js.JsonLog, error) { return m, nil }),
)

var asn1Log IO[fa.LogItem] = Bind(
	jsonLog,
	Lift(json2asn1),
)

func LogItemToDerBytes(i fa.LogItem) ([]byte, error) { return i.ToDerBytes() }

var derBytes IO[[]byte] = Bind(
	asn1Log,
	Lift(LogItemToDerBytes),
)

var BytesToStdout func([]byte) IO[Void] = Lift(
	func(dat []byte) (Void, error) {
		_, e := os.Stdout.Write(dat)
		return Empty, e
	},
)

var stdin2json2asn12der2stdout IO[Void] = Bind(
	derBytes,
	BytesToStdout,
)

func main() {
	_, e := stdin2json2asn12der2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
