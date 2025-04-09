#!/bin/sh

input=./input.json
output=./output.asn1.der.dat

genjson(){
	echo generating input json...

	jq -c -n '{
		timestamp: "2025-04-08T01:45:13.012345Z",
		severity: "INFO",
		message: "test message",
		id: 299792458,
		tags: "http,read-only",
		"process.pid": 42,
		"service.version": "1.0.0",
		"debug": true,
		"cache.hit": false,
		"db.upserted": 599,
		"http.content.length": 634,
		"http.status.code": 200,
		"user.id": "JD",
		"http.request.method": "GET",
	}' |
		dd \
			if=/dev/stdin \
			of="${input}" \
			status=none
}

test -f "${input}" || genjson

echo converting json to der...
cat "${input}" |
	./jsonlog2asn1 |
	dd \
		if=/dev/stdin \
		of="${output}" \
		status=none

echo converting der to yaml using asn1tools/dasel/bat...
cat "${output}" |
	python3 ./sample.py |
	dasel \
		--read=json \
		--write=yaml |
	bat --language=yaml

ls -l "${input}" "${output}"
