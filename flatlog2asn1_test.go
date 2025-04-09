package flatlog2asn1_test

import (
	"testing"

	fa "github.com/takanoriyanagitani/go-flatlog2asn1"
)

func TestFlatLogToAsn1(t *testing.T) {
	t.Parallel()

	t.Run("LogItem", func(t *testing.T) {
		t.Parallel()

		t.Run("ToDerBytes", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				var item fa.LogItem
				der, e := item.ToDerBytes()
				if nil != e {
					t.Fatalf("unexpected error: %v", e)
				}

				if 0 == len(der) {
					t.Fatal("empty bytes got")
				}
			})
		})
	})
}
