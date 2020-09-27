package main

// Sequence represents description of amino acid or DNA
type Sequence struct {
	ID          string
	Description string
	Value       string
}

type Letter rune
type Direction int

const (
	TopDirection Direction = iota
	LeftDirection
	DiagonalDirection
	NullDirection
)

type Cell struct {
	Distance int
	Dir      Direction
}

type Matrix []Line
type Line []*Cell

type ScoringFunc map[uint8]map[uint8]int
