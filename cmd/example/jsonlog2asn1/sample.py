import asn1tools
import sys
import json

flatlog = asn1tools.compile_files("./flatlog.asn")
encoded = sys.stdin.buffer.read()
decoded = flatlog.decode(
	"LogItem",
	encoded,
)
json.dump(decoded, fp=sys.stdout)
