package mysql

import (
	"bytes"
	"testing"
)

func TestDialector_QuoteTo(t *testing.T) {
	testdatas := []struct {
		raw    string
		expect string
	}{
		{"datadase.tableUser", "`datadase`.`tableUser`"},
		{"datadase.table`User", "`datadase`.`table``User`"},
	}

	dailor := Open("")
	for _, item := range testdatas {
		buf := &bytes.Buffer{}
		dailor.QuoteTo(buf, item.raw)
		if buf.String() != item.expect {
			t.Errorf("quote %q fail, got %q, expect %q", item.raw, buf.String(), item.expect)
		}
	}
}

// BenchmarkDialector_QuoteTo
// Result:
// goos: darwin
// goarch: amd64
// pkg: gorm.io/driver/mysql
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkDialector_QuoteTo               6298825               179.7 ns/op
// BenchmarkDialector_QuoteTo-2             6729696               173.1 ns/op
// BenchmarkDialector_QuoteTo-4             6804223               173.9 ns/op
// BenchmarkDialector_QuoteTo-8             6885795               173.9 ns/op
// BenchmarkDialector_QuoteTo-12            6879357               174.7 ns/op
func BenchmarkDialector_QuoteTo(b *testing.B) {
	dailor := Open("")
	buf := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dailor.QuoteTo(buf, "datadase.table`User")
		buf.Reset()
	}
}
