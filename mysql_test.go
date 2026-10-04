package mysql

import (
	"bytes"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestNew(t *testing.T) {
	dialector := New(Config{DSN: "gorm:gorm@tcp(127.0.0.1:9910)/gorm?charset=utf8&parseTime=True&loc=Local"})
	d, ok := dialector.(*Dialector)
	if !ok {
		t.Fatal("dialector is not *Dialector")
	}
	if d.DSNConfig == nil {
		t.Error("dialector.DSNConfig is nil")
	}

	dialector = New(Config{DSNConfig: mysql.NewConfig()})
	d, ok = dialector.(*Dialector)
	if !ok {
		t.Fatal("dialector is not *Dialector")
	}
	if d.DSN == "" {
		t.Error("dialector.DSN is empty")
	}
}

func TestDialector_QuoteTo(t *testing.T) {
	testdatas := []struct {
		raw    string
		expect string
	}{
		{"datadase.tableUser", "`datadase`.`tableUser`"},
		{"datadase.table`User", "`datadase`.`table``User`"},
		{"`a`.`b`", "`a`.`b`"},
		{"`a`.b`", "`a`.`b```"},
		{"a.`b`", "`a`.`b`"},
		{"`a`.b`c", "`a`.`b``c`"},
		{"`a`.`b`c`", "`a`.`b``c`"},
		{"`a`.b", "`a`.`b`"},
		{"`ab`", "`ab`"},
		{"`a``b`", "`a``b`"},
		{"`a```b`", "`a````b`"},
		{"a`b", "`a``b`"},
		{"ab", "`ab`"},
		{"`a.b`", "`a.b`"},
		{"a.b", "`a`.`b`"},
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
// BenchmarkDialector_QuoteTo               9184232               113.2 ns/op
// BenchmarkDialector_QuoteTo-2             9782818               112.3 ns/op
// BenchmarkDialector_QuoteTo-4            10726722               109.0 ns/op
// BenchmarkDialector_QuoteTo-8             9656778               113.1 ns/op
// BenchmarkDialector_QuoteTo-12           10729615               112.7 ns/op
func BenchmarkDialector_QuoteTo(b *testing.B) {
	dailor := Open("")
	buf := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dailor.QuoteTo(buf, "datadase.table`User")
		buf.Reset()
	}
}

func TestCheckVersion(t *testing.T) {
	versions := map[string]string{
		"5.6.1":  "5.6",
		"5.10.2": "5.6",
		"5.10":   "5.6",
		"10.6.26-MariaDB-1:10.4.26+maria~ubu2004": "10.6",
		"10.6.26-MariaDB-1:10.4.26+maria~ubu2005": "10.6.3",
		"10.4.26-MariaDB-1:10.4.26+maria~ubu2004": "5.6",
	}

	for k, v := range versions {
		if !checkVersion(k, v) || checkVersion(v, k) {
			t.Fatalf("returns %v when comparing %v, %v", checkVersion(k, v), k, v)
		}
	}
}

func TestCheckVersionEqual(t *testing.T) {
	versions := [][2]string{
		{"10.5.2-MariaDB-ubu2404", "10.5.2"},
		{"10.5.2", "10.5.2"},
		{"10.5.29-MariaDB-1:10.5.29+maria~ubu2004", "10.5.2"},
		{"11.4.12-MariaDB-ubu2404", "10.5.2"},
	}

	for _, v := range versions {
		if !checkVersion(v[0], v[1]) {
			t.Fatalf("returns false for %v, expected to be newer or equal to %v", v[0], v[1])
		}
	}
}

func TestCheckVersionOlder(t *testing.T) {
	versions := [][2]string{
		{"10.5.1-MariaDB-ubu2404", "10.5.2"},
		{"10.4.34-MariaDB-1:10.4.34+maria~ubu2004", "10.5.2"},
		{"10.5", "10.5.2"},
		{"5.7.44", "8.0.3"},
	}

	for _, v := range versions {
		if checkVersion(v[0], v[1]) {
			t.Fatalf("returns true for %v, expected to be older than %v", v[0], v[1])
		}
	}
}
