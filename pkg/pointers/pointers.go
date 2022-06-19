package pointers

func FromSlice[S []E, P []*E, E any](s S) P {
	out := P{}
	for i := range s {
		out = append(out, &s[i])
	}
	return out
}

func ToValues[S []E, P []*E, E any](p P) S {
	out := S{}
	for i := range p {
		out = append(out, *p[i])
	}
	return out
}
