package diff

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

type edit[T comparable] struct {
	start, end  Position
	replacement []T
}

func applyEdits[T comparable](src []T, edits []edit[T]) []T {
	dst := make([]T, 0, len(src))

	i := 0
	for _, edit := range edits {
		dst = append(dst, src[i:edit.start.Index]...)
		dst = append(dst, edit.replacement...)
		i = edit.end.Index
	}

	return append(dst, src[i:]...)
}

func diffAndApplyEdits[T comparable](dst, src []T, isNewline func(T) bool) []T {
	var edits []edit[T]
	Diff(dst, src, isNewline, func(start, end Position, replacement []T) {
		edits = append(edits, edit[T]{
			start:       start,
			end:         end,
			replacement: replacement,
		})
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
		Diff([]byte(dst), []byte(src), nil, func(start, end Position, replacement []byte) {
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
	edit := func(start, end Position, replacement []byte) {}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff([]byte(dst), []byte(src), nil, edit)
	}
}

func BenchmarkDiffIsNewline(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	dst := strconv.FormatUint(r.Uint64(), 2)
	src := strconv.FormatUint(r.Uint64(), 2)
	isNewline := func(b byte) bool { return b == '0' }
	edit := func(start, end Position, replacement []byte) {}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff([]byte(dst), []byte(src), isNewline, edit)
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		dst, src string
	}{
		{"ABCABBA", "CBABAC"},
		{"ABCABBA", "ABCABBA"},
		{"", ""},
		{"", "0"},
		{"0", "110"},
		{"1111", "0"},
		{"0", "11111"},
		{"0BA", "1BB"},
		{"AAC", "AC"},
		{"0B0", "11BB0"},
		{"000000000", "0001111"},
		{"01001100", "0000000"},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out := string(diffAndApplyEdits([]rune(test.dst), []rune(test.src), nil))
			if out != test.dst {
				t.Errorf("expected \"%s\", got \"%s\"", test.dst, out)
			}
		})
	}
}

func TestDiffEdits(t *testing.T) {
	var edits []edit[rune]
	Diff([]rune("ABCABBA"), []rune("CBABAC"), nil, func(start, end Position, replacement []rune) {
		edits = append(edits, edit[rune]{
			start:       start,
			end:         end,
			replacement: replacement,
		})
	})

	expectedEdits := []edit[rune]{
		{
			start: Position{Index: 0, Column: 0},
			end:   Position{Index: 1, Column: 1},
		},
		{
			start:       Position{Index: 1, Column: 1},
			end:         Position{Index: 1, Column: 1},
			replacement: []rune{'A'},
		},
		{
			start:       Position{Index: 2, Column: 2},
			end:         Position{Index: 2, Column: 2},
			replacement: []rune{'C'},
		},
		{
			start:       Position{Index: 4, Column: 4},
			end:         Position{Index: 4, Column: 4},
			replacement: []rune{'B'},
		},
		{
			start: Position{Index: 5, Column: 5},
			end:   Position{Index: 6, Column: 6},
		},
	}

	if len(edits) != len(expectedEdits) {
		t.Fatalf("expected %d edits, got %d", len(expectedEdits), len(edits))
	}

	for i, edit := range edits {
		expectedEdit := expectedEdits[i]
		if !reflect.DeepEqual(edit, expectedEdit) {
			t.Errorf("expected edit %v, got %v", expectedEdit, edit)
		}
	}
}
