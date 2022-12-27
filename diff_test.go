package diff

import (
	"math/rand"
	"strconv"
	"testing"
)

type edit[T comparable] struct {
	start, end  int
	replacement []T
}

func applyEdits[T comparable](src []T, edits []edit[T]) []T {
	dst := make([]T, 0, len(src))

	i := 0
	for _, edit := range edits {
		dst = append(dst, src[i:edit.start]...)
		dst = append(dst, edit.replacement...)
		i = edit.end
	}

	return append(dst, src[i:]...)
}

func diffAndApplyEdits[T comparable](dst, src []T) []T {
	var edits []edit[T]
	Diff(dst, src, func(start, end int, replacement []T) {
		edits = append(edits, edit[T]{
			start:       start,
			end:         end,
			replacement: replacement})
	})

	return applyEdits(src, edits)
}

func FuzzDiff(f *testing.F) {
	testcases := [][2]string{
		{"ABCABBA", "CBABAC"},
	}

	for _, tc := range testcases {
		f.Add(tc[0], tc[1]) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, dst, src string) {
		var edits []edit[byte]
		Diff([]byte(dst), []byte(src), func(start, end int, replacement []byte) {
			edits = append(edits, edit[byte]{
				start:       start,
				end:         end,
				replacement: replacement})
		})

		if out := string(applyEdits([]byte(src), edits)); out != dst {
			t.Errorf("expected %s, got %s", dst, out)
		}
	})
}

func BenchmarkDiff(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	dst := strconv.FormatUint(r.Uint64(), 2)
	src := strconv.FormatUint(r.Uint64(), 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff([]byte(dst), []byte(src), func(start, end int, replacement []byte) {})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		dst, src string
	}{
		{"ABCABBA", "CBABAC"},
		{"", "0"},
		{"0", "110"},
		{"1111", "0"},
		{"0", "11111"},
		{"0BA", "1BB"},
		{"AAC", "AC"},
		{"0B0", "11BB0"},
		{"000000000", "0001111"},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out := string(diffAndApplyEdits([]byte(test.dst), []byte(test.src)))
			if out != test.dst {
				t.Errorf("expected \"%s\", got \"%s\"", test.dst, out)
			}
		})
	}
}
