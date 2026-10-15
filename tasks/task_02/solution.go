package main

func rotateRunes(s string, shift int) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return ""
	}
	k := shift % n
	if k < 0 {
		k += n
	}
	if k == 0 {
		return string(runes)
	}
	out := make([]rune, n)
	copy(out, runes[k:])
	copy(out[n-k:], runes[:k])
	return string(out)
}
