package poser

// Poser defines the interface for an element in a slice
type Poser interface {
	GetPosition() int
	SetPosition(pos int)
	GetIgnore() bool
	SetIgnore(ignore bool)
}

// AppendTo appends the given elements to the slice
func AppendTo[S []E, E Poser](s S, e ...E) S {
	s = append(s, e...)
	return reassignPositions(s)
}

// DeleteAt removes the element from the slice at the given position
func DeleteAt[S []E, E Poser](s S, pos int) (S, E) {
	var out S
	var e E
	for i := range s {
		if s[i].GetPosition() == pos {
			s[i].SetIgnore(true)
			e = s[i]
			continue
		}
		out = append(out, s[i])
	}

	return reassignPositions(out), e
}

// InsertInto inserts elements into a slice at the given position
func InsertInto[S []E, E Poser](s S, pos int, e ...E) S {
	if pos <= 1 {
		aux := s
		s = e
		s = append(s, aux...)
	} else if pos > 1 && pos <= len(s) {
		aux := s
		piv := pos - 1
		s = aux[:piv]
		s = append(s, e...)
		s = append(s, aux[piv:]...)
	} else {
		return AppendTo(s, e...)
	}

	return reassignPositions(s)
}

// MoveTo moves an element from one position to another
func MoveTo[S []E, E Poser](s S, to, from int) S {
	if from == to || from < 1 || len(s) == 0 || from > len(s) {
		return s
	}

	var moved, afterPiv S
	var piv E
	var hasPiv bool
	for i := range s {
		if s[i].GetPosition() == from {
			piv = s[i]
			hasPiv = true
		} else if s[i].GetPosition() < to {
			moved = append(moved, s[i])
		} else if s[i].GetPosition() > to {
			afterPiv = append(afterPiv, s[i])
		} else if s[i].GetPosition() == to {
			if from < to {
				moved = append(moved, s[i])
			} else {
				afterPiv = append(afterPiv, s[i])
			}
		}
	}

	if hasPiv {
		moved = append(moved, piv)
	}
	moved = append(moved, afterPiv...)

	return reassignPositions(moved)
}

// Pop removes the element at position 1 from the slice and returns it
func Pop[S []E, E Poser](s S) (S, E) {
	return DeleteAt(s, 1)
}

// PrependItem inserts elements into the slice at position 1
func PrependItem[S []E, E Poser](s S, e ...E) S {
	return InsertInto(s, 1, e...)
}

func reassignPositions[S []E, E Poser](s S) S {
	pos := 0
	for i := range s {
		if s[i].GetIgnore() {
			continue
		}
		pos++
		s[i].SetPosition(pos)
	}
	return s
}
