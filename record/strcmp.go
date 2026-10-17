package record

// JaroWinkler scores two strings from 0 (nothing in common) to 1 (equal).
// It forgives transpositions and rewards a shared prefix, which suits
// names: typists slip in the middle and the end more than the start.
func JaroWinkler(a, b string) float64 {
	if a == b {
		return 1
	}
	la, lb := len(a), len(b)
	if la == 0 || lb == 0 {
		return 0
	}
	window := max(la, lb)/2 - 1
	window = max(window, 0)
	var ma, mb [64]bool
	if la > 64 || lb > 64 {
		// Names never get this long; compare a prefix rather than allocate.
		a, b = a[:min(la, 64)], b[:min(lb, 64)]
		la, lb = len(a), len(b)
	}
	m := 0
	for i := range la {
		lo, hi := max(0, i-window), min(lb, i+window+1)
		for j := lo; j < hi; j++ {
			if !mb[j] && a[i] == b[j] {
				ma[i], mb[j] = true, true
				m++
				break
			}
		}
	}
	if m == 0 {
		return 0
	}
	t, j := 0, 0
	for i := range la {
		if !ma[i] {
			continue
		}
		for !mb[j] {
			j++
		}
		if a[i] != b[j] {
			t++
		}
		j++
	}
	fm := float64(m)
	jaro := (fm/float64(la) + fm/float64(lb) + (fm-float64(t)/2)/fm) / 3
	p := 0
	for p < min(4, la, lb) && a[p] == b[p] {
		p++
	}
	return jaro + float64(p)*0.1*(1-jaro)
}

// Soundex codes a cleaned name by how it sounds: Smith and Smyth both
// give S530. Used for blocking, where a cheap key that survives spelling
// variants matters more than precision.
func Soundex(s string) string {
	if s == "" {
		return ""
	}
	code := func(c byte) byte {
		switch c {
		case 'B', 'F', 'P', 'V':
			return '1'
		case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z':
			return '2'
		case 'D', 'T':
			return '3'
		case 'L':
			return '4'
		case 'M', 'N':
			return '5'
		case 'R':
			return '6'
		}
		return 0 // vowels, H, W, Y
	}
	out := []byte{s[0], '0', '0', '0'}
	n, last := 1, code(s[0])
	for i := 1; i < len(s) && n < 4; i++ {
		c := code(s[i])
		if c != 0 && c != last {
			out[n] = c
			n++
		}
		// H and W don't separate equal codes; vowels do.
		if s[i] != 'H' && s[i] != 'W' {
			last = c
		}
	}
	return string(out)
}

// OSA is the optimal string alignment distance: the number of insertions,
// deletions, substitutions and swaps of neighbouring letters to turn a
// into b, each letter edited at most once. One keying mistake is distance 1.
func OSA(a, b string) int {
	la, lb := len(a), len(b)
	prev2 := make([]int, lb+1)
	prev := make([]int, lb+1)
	cur := make([]int, lb+1)
	for j := range lb + 1 {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		cur[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				cur[j] = min(cur[j], prev2[j-2]+1)
			}
		}
		prev2, prev, cur = prev, cur, prev2
	}
	return prev[lb]
}
