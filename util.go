package main

func copyMap(src map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for k, v := range src {
		dst[k] = make([]string, len(v))
		copy(dst[k], v)
	}
	return dst
}
