package xfmt

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"text/tabwriter"
)

var sink string

type scenario struct {
	name string
	std  func() string // fmt.Sprintf
	fast func() string // xfmt.Sprintf
	gen  func() string // xfmt.SprintfGeneric (nil when arguments are mixed)
}

var scenarios = []scenario{
	{
		name: "int",
		std:  func() string { return fmt.Sprintf("%d", 1234567) },
		fast: func() string { return Sprintf("%d", 1234567) },
		gen:  func() string { return SprintfGeneric[int]("%d", 1234567) },
	},
	{
		name: "ints_x4",
		std:  func() string { return fmt.Sprintf("%d-%d-%d-%d", 1, 22, 333, 4444) },
		fast: func() string { return Sprintf("%d-%d-%d-%d", 1, 22, 333, 4444) },
		gen:  func() string { return SprintfGeneric[int]("%d-%d-%d-%d", 1, 22, 333, 4444) },
	},
	{
		name: "uint_hex",
		std:  func() string { return fmt.Sprintf("%x", uint64(0xdeadbeef)) },
		fast: func() string { return Sprintf("%x", uint64(0xdeadbeef)) },
		gen:  func() string { return SprintfGeneric[uint64]("%x", uint64(0xdeadbeef)) },
	},
	{
		name: "string",
		std:  func() string { return fmt.Sprintf("%s", "hello world") },
		fast: func() string { return Sprintf("%s", "hello world") },
		gen:  func() string { return SprintfGeneric[string]("%s", "hello world") },
	},
	{
		name: "strings_x3",
		std:  func() string { return fmt.Sprintf("%s/%s/%s", "alpha", "beta", "gamma") },
		fast: func() string { return Sprintf("%s/%s/%s", "alpha", "beta", "gamma") },
		gen:  func() string { return SprintfGeneric[string]("%s/%s/%s", "alpha", "beta", "gamma") },
	},
	{
		name: "float",
		std:  func() string { return fmt.Sprintf("%f", 3.14159) },
		fast: func() string { return Sprintf("%f", 3.14159) },
		gen:  func() string { return SprintfGeneric[float64]("%f", 3.14159) },
	},
	{
		name: "bytes_hex",
		std:  func() string { return fmt.Sprintf("%x", []byte("payload")) },
		fast: func() string { return Sprintf("%x", []byte("payload")) },
		gen:  func() string { return SprintfGeneric[[]byte]("%x", []byte("payload")) },
	},
	{
		name: "mixed_hot",
		std:  func() string { return fmt.Sprintf("%s id=%d hex=%x f=%f", "k", 7, 255, 2.5) },
		fast: func() string { return Sprintf("%s id=%d hex=%x f=%f", "k", 7, 255, 2.5) },
	},
	{
		name: "log_line",
		std:  func() string { return fmt.Sprintf("[%s] svc=%s code=%d", "INFO", "auth", 200) },
		fast: func() string { return Sprintf("[%s] svc=%s code=%d", "INFO", "auth", 200) },
	},
	{
		name: "with_fallback",
		std:  func() string { return fmt.Sprintf("%s took %.2fms", "query", 12.5) },
		fast: func() string { return Sprintf("%s took %.2fms", "query", 12.5) },
	},
	{
		name: "all_fallback",
		std:  func() string { return fmt.Sprintf("%v %q %t", []int{1, 2}, "s", true) },
		fast: func() string { return Sprintf("%v %q %t", []int{1, 2}, "s", true) },
	},
}

func BenchmarkCompare(b *testing.B) {
	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			b.Run("fmt", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					sink = sc.std()
				}
			})
			b.Run("xfmt", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					sink = sc.fast()
				}
			})
			if sc.gen != nil {
				b.Run("xfmt_generic", func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						sink = sc.gen()
					}
				})
			}
		})
	}
}

func TestComparisonReport(t *testing.T) {
	if os.Getenv("XFMT_REPORT") == "" {
		t.Skip("set XFMT_REPORT=1 to print the comparison table")
	}

	for _, sc := range scenarios {
		want := sc.std()
		if got := sc.fast(); got != want {
			t.Fatalf("%s: Sprintf=%q want %q", sc.name, got, want)
		}
		if sc.gen != nil {
			if got := sc.gen(); got != want {
				t.Fatalf("%s: SprintfGeneric=%q want %q", sc.name, got, want)
			}
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "scenario\tfmt ns/op\txfmt ns/op\tspeedup\tgeneric ns/op\tspeedup\tallocs fmt/xfmt")
	fmt.Fprintln(w, "--------\t---------\t----------\t-------\t-------------\t-------\t--------------")

	for _, sc := range scenarios {
		rs := benchOne(sc.std)
		rf := benchOne(sc.fast)
		stdNs := float64(rs.NsPerOp())

		genNs, genSpeedup := "-", "-"
		if sc.gen != nil {
			rg := benchOne(sc.gen)
			genNs = fmt.Sprintf("%d", rg.NsPerOp())
			genSpeedup = fmt.Sprintf("%.2fx", stdNs/float64(rg.NsPerOp()))
		}

		fmt.Fprintf(w, "%s\t%d\t%d\t%.2fx\t%s\t%s\t%d/%d\n",
			sc.name,
			rs.NsPerOp(),
			rf.NsPerOp(),
			stdNs/float64(rf.NsPerOp()),
			genNs,
			genSpeedup,
			rs.AllocsPerOp(),
			rf.AllocsPerOp(),
		)
	}
	w.Flush()
}

func benchOne(fn func() string) testing.BenchmarkResult {
	return testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = fn()
		}
	})
}

// ——— overhead benchmarks ———

type overheadScenario struct {
	name string
	raw  func() string
	fmt  func() string
	fast func() string
	gen  func() string
}

var overheadScenarios = []overheadScenario{
	{
		name: "int→str (%d)",
		raw:  func() string { return strconv.Itoa(1234567) },
		fmt:  func() string { return fmt.Sprintf("%d", 1234567) },
		fast: func() string { return Sprintf("%d", 1234567) },
		gen:  func() string { return SprintfGeneric[int]("%d", 1234567) },
	},
	{
		name: "uint→hex (%x)",
		raw:  func() string { return strconv.FormatUint(uint64(0xdeadbeef), 16) },
		fmt:  func() string { return fmt.Sprintf("%x", uint64(0xdeadbeef)) },
		fast: func() string { return Sprintf("%x", uint64(0xdeadbeef)) },
		gen:  func() string { return SprintfGeneric[uint64]("%x", uint64(0xdeadbeef)) },
	},
	{
		name: "int→oct (%o)",
		raw:  func() string { return strconv.FormatInt(int64(420), 8) },
		fmt:  func() string { return fmt.Sprintf("%o", 420) },
		fast: func() string { return Sprintf("%o", 420) },
		gen:  func() string { return SprintfGeneric[int]("%o", 420) },
	},
	{
		name: "int→bin (%b)",
		raw:  func() string { return strconv.FormatInt(int64(5), 2) },
		fmt:  func() string { return fmt.Sprintf("%b", 5) },
		fast: func() string { return Sprintf("%b", 5) },
		gen:  func() string { return SprintfGeneric[int]("%b", 5) },
	},
	{
		name: "float→str (%f)",
		raw:  func() string { return strconv.FormatFloat(3.14159, 'f', 6, 64) },
		fmt:  func() string { return fmt.Sprintf("%f", 3.14159) },
		fast: func() string { return Sprintf("%f", 3.14159) },
		gen:  func() string { return SprintfGeneric[float64]("%f", 3.14159) },
	},
	{
		name: "str (no-op, %s)",
		raw:  func() string { return "hello world" },
		fmt:  func() string { return fmt.Sprintf("%s", "hello world") },
		fast: func() string { return Sprintf("%s", "hello world") },
		gen:  func() string { return SprintfGeneric[string]("%s", "hello world") },
	},
	{
		name: "[]byte→hex (%x)",
		raw:  func() string { return hexEncodeStr("payload") },
		fmt:  func() string { return fmt.Sprintf("%x", []byte("payload")) },
		fast: func() string { return Sprintf("%x", []byte("payload")) },
		gen:  func() string { return SprintfGeneric[[]byte]("%x", []byte("payload")) },
	},
	{
		name: "pre+int (%s=%d)",
		raw:  func() string { return "n=" + strconv.Itoa(1234567) },
		fmt:  func() string { return fmt.Sprintf("n=%d", 1234567) },
		fast: func() string { return Sprintf("n=%d", 1234567) },
		gen:  func() string { return SprintfGeneric[int]("n=%d", 1234567) },
	},
}

var hexLowerTable = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func hexEncodeStr(s string) string {
	buf := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		b := s[i]
		buf = append(buf, hexLowerTable[b>>4], hexLowerTable[b&0xf])
	}
	return string(buf)
}

func BenchmarkOverhead(b *testing.B) {
	for _, sc := range overheadScenarios {
		b.Run(sc.name, func(b *testing.B) {
			b.Run("raw", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					sink = sc.raw()
				}
			})
			b.Run("fmt", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					sink = sc.fmt()
				}
			})
			b.Run("xfmt", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					sink = sc.fast()
				}
			})
			if sc.gen != nil {
				b.Run("xfmt_generic", func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						sink = sc.gen()
					}
				})
			}
		})
	}
}
