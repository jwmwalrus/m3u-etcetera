package poser

type Poser interface {
	GetPosition() int
	SetPosition(pos int)
	GetIgnore() bool
	SetIgnore(ignore bool)
}

func AppendTo[S ~[]E, E Poser](s S, e ...E) S {
	s = append(s, e...)
	return reasignPositions(s)
}

func DeleteAt[S ~[]E, E Poser](s S, pos int) (S, E) {
	var ret S
	var e E
	for _, x := range s {
		if x.GetPosition() == pos {
			x.SetIgnore(true)
			e = x
			continue
		}
		ret = append(ret, x)
	}

	return reasignPositions(ret), e
}

func InsertInto[S ~[]E, E Poser](s S, pos int, e ...E) S {
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

	return reasignPositions(s)
}

func MoveTo[S ~[]E, E Poser](s S, to, from int) S {
	if from == to || from < 1 || len(s) == 0 || from > len(s) {
		return s
	}

	var moved, afterPiv S
	var piv E
	var hasPiv bool
	for i, x := range s {
		if s[i].GetPosition() == from {
			piv = x
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

	return reasignPositions(moved)
}

func Pop[S ~[]E, E Poser](s S) (S, E) {

	var e E
	popped := false
	for i, x := range s {
		if s[i].GetPosition() == 1 {
			s[i].SetIgnore(true)
			e = x
			popped = true
			break
		}
	}
	if !popped {
		return s, e
	}
	return reasignPositions(s), e
}

func PrependItem[S ~[]E, E Poser](s S, e ...E) S {
	return InsertInto(s, 1, e...)
}

func reasignPositions[S ~[]E, E Poser](s S) S {
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
