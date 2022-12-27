package diff

func mink(m, d int) int {
	if d <= m {
		return -d
	}

	if (m^d)&1 == 1 {
		m--
	}

	return -m
}

func maxk(n, d int) int {
	if d <= n {
		return d
	}

	return n
}

type differ[T comparable] struct {
	a, b   []T
	v0, v1 []int
	edit   func(start, end int, replacement []T)
}

func (df *differ[T]) middlesnake(a, b []T) (int, int, int, int, int) {
	n := len(a)
	m := len(b)
	v0 := df.v0[:n+m+1]
	v1 := df.v1[:len(v0)]
	for i := range v0 {
		v0[i] = 0
		v1[i] = 0
	}

	δ := n - m
	for d := 0; d <= (m+n+1)/2; d++ {
		mink := mink(m, d)
		maxk := maxk(n, d)

		// Forward search
		for k := mink; k <= maxk; k += 2 {
			i := m + k

			var x int
			if k == mink || k != maxk && v0[i-1] < v0[i+1] {
				x = v0[i+1]
			} else {
				x = v0[i-1] + 1
			}

			y := x - k

			// 🐍
			u, v := x, y
			for u < n && v < m && a[u] == b[v] {
				u++
				v++
			}

			if δ&1 == 1 { // ∆ is odd
				// Check that the reverse search has a value for k and that it
				// overlaps with the last forward snake.
				if k := δ - k; k >= mink && k <= maxk && u+v1[m+k] >= n {
					return x, y, u, v, 2*d - 1
				}
			}

			v0[i] = u
		}

		// Reverse search
		for k := maxk; k >= mink; k -= 2 {
			j := m + k

			var x int
			if k == mink || k != maxk && v1[j-1] < v1[j+1] {
				x = v1[j+1]
			} else {
				x = v1[j-1] + 1
			}

			y := x - k

			// 🐍
			u, v := x, y
			for u < n && v < m && a[n-u-1] == b[m-v-1] {
				u++
				v++
			}

			if δ&1 == 0 { // ∆ is even
				// Check that the forward search has a value for k and that it
				// overlaps with the last reverse snake.
				if k := δ - k; k >= mink && k <= maxk && u+v0[m+k] >= n {
					return n - u, m - v, n - x, m - y, 2 * d
				}
			}

			v1[j] = u
		}
	}

	return 0, 0, 0, 0, -1
}

func (df *differ[T]) diff(x, y, u, v int) {
	ma := df.a[x:u]
	if len(ma) == 0 {
		df.edit(y, v, nil)
		return
	}

	mb := df.b[y:v]
	if len(mb) == 0 {
		df.edit(y, y, ma)
		return
	}

	mx, my, mu, mv, d := df.middlesnake(ma, mb)
	if d == 0 {
		return
	}

	if d == 1 {
		// When d == 1, (mx, my) will be the endpoint of the single edit.
		if len(ma) > len(mb) {
			df.edit(y+my, y+my, ma[mx-1:mx])
		} else {
			df.edit(y+my-1, y+my, nil)
		}

		return
	}

	df.diff(x, y, x+mx, y+my)
	df.diff(x+mu, y+mv, u, v)
}

// Diff will call edit in order for each edit of b required to produce a. Start
// and end locations are relative to the begining of b, and the replacement
// slice will either be nil or a subslice of a.
func Diff[T comparable](a, b []T, edit func(start, end int, replacement []T)) {
	max := len(a) + len(b) + 1
	df := &differ[T]{
		a:    a,
		b:    b,
		v0:   make([]int, max),
		v1:   make([]int, max),
		edit: edit,
	}

	df.diff(0, 0, len(a), len(b))
}
