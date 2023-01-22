# diff

[![Go Reference](https://pkg.go.dev/badge/github.com/buth/diff.svg)](https://pkg.go.dev/github.com/buth/diff)

This package provides a generic implementation of the [Myers diff
algorithm](http://www.xmailserver.org/diff2.pdf) that produces sequential edits.

## Usage

To produce a list of edits, you can create an `edit` type and build a slice
using a callback provided to the `Diff` method.

```go
import "github.com/buth/diff"

type edit[T comparable] struct {
	start, end  diff.Position
	replacement []T
}

func listOfEdits[T comparable](dst, src []T) []edit[T] {
	var edits []edit[T]
	diff.Diff(dst, src, nil, func(start, end diff.Position, replacement []T) {
		edits = append(edits, edit[T]{
			start:       start,
			end:         end,
			replacement: replacement
		})
	})

	return edits
} 
```

Keep in mind that the provided `replacement` slice will be a subslice of the
destination value if it is not nil.

To apply the edits, you can iterate through the resulting list using an index of
the source slice as the starting point.

```go
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
```
