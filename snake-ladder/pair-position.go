package main

type PairPosition struct {
	start int
	end   int
}

func NewPairPosition(s, e int) *PairPosition {
	return &PairPosition{
		start: s,
		end:   e,
	}
}
