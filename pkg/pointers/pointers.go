package pointers

// FromSlice returns a slice of pointers from the given slice s
func FromSlice[S []E, P []*E, E any](s S) P {
	out := P{}
	for i := range s {
		out = append(out, &s[i])
	}
	return out
}

// ToValues returns a slice of values from the given slice of pointers, p
func ToValues[S []E, P []*E, E any](p P) S {
	out := S{}
	for i := range p {
		out = append(out, *p[i])
	}
	return out
}
