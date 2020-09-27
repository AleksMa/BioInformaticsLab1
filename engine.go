package main

import "fmt"

const GapSymbol = "-"

type NeedlemanWunsch struct {
	TopSequence  *Sequence
	LeftSequence *Sequence
	Table        Matrix
	SF           ScoringFunc
	GapValue     int
}

func NewNeedlemanWunsch(first, second *Sequence, sf ScoringFunc, GapValue int) *NeedlemanWunsch {
	nw := &NeedlemanWunsch{
		TopSequence:  second,
		LeftSequence: first,
		Table:        make(Matrix, len(first.Value)+1),
		SF:           sf,
		GapValue:     GapValue,
	}
	// Аллоцируем первую строку
	nw.Table[0] = make(Line, len(second.Value)+1)
	// Обнуляем (0, 0)
	nw.Table[0][0] = &Cell{
		Distance: 0,
		Dir:      NullDirection,
	}
	// Обнуляем первую строку
	for i := range second.Value {
		nw.Table[0][i+1] = &Cell{
			Distance: GapValue * (i + 1),
			Dir:      LeftDirection,
		}
	}
	// Аллоцируем оставшиеся строки, зануляем первый столбец
	for i := range first.Value {
		nw.Table[i+1] = make(Line, len(second.Value)+1)
		nw.Table[i+1][0] = &Cell{
			Distance: GapValue * (i + 1),
			Dir:      TopDirection,
		}
	}
	return nw
}

// Функция вывода таблицы для отладки
func (nw *NeedlemanWunsch) Print() {
	for i := 0; i <= len(nw.LeftSequence.Value); i++ {
		for j := 0; j <= len(nw.TopSequence.Value); j++ {
			fmt.Print(nw.Table[i][j].Distance, ", ", nw.Table[i][j].Dir)
			fmt.Print("   | ")
		}
		fmt.Println()
	}
}

func (nw *NeedlemanWunsch) Solve() (string, string, int) {
	// Рекурсивно определяем значения матрицы, одновременно определяя и направления
	nw.determine(len(nw.LeftSequence.Value), len(nw.TopSequence.Value))

	cell := nw.Table[len(nw.LeftSequence.Value)][len(nw.TopSequence.Value)]

	score := cell.Distance
	firstRes, secondRes := "", ""

	// Двигаемся от правой нижней ячейки матрицы к левой верхней, и строим с конца строки-выравнивания
	fp, sp := len(nw.LeftSequence.Value)-1, len(nw.TopSequence.Value)-1
	for cell.Dir != NullDirection {
		if cell.Dir == DiagonalDirection {
			firstRes = string(rune(nw.LeftSequence.Value[fp])) + firstRes
			secondRes = string(rune(nw.TopSequence.Value[sp])) + secondRes
			sp--
			fp--
		} else if cell.Dir == TopDirection {
			firstRes = string(rune(nw.LeftSequence.Value[fp])) + firstRes
			secondRes = GapSymbol + secondRes
			fp--
		} else if cell.Dir == LeftDirection {
			firstRes = GapSymbol + firstRes
			secondRes = string(rune(nw.TopSequence.Value[sp])) + secondRes
			sp--
		}
		cell = nw.Table[fp+1][sp+1]
	}

	return firstRes, secondRes, score
}

// Рекурсивное заполнение матрицы
func (nw *NeedlemanWunsch) determine(i, j int) {
	if nw.Table[i][j] != nil {
		return
	}
	leftCell, topCell, diagCell := nw.Table[i][j-1], nw.Table[i-1][j], nw.Table[i-1][j-1]
	if leftCell == nil {
		nw.determine(i, j-1)
		leftCell = nw.Table[i][j-1]
	}
	if diagCell == nil {
		nw.determine(i-1, j-1)
		diagCell = nw.Table[i-1][j-1]
	}
	if topCell == nil {
		nw.determine(i-1, j)
		topCell = nw.Table[i-1][j]
	}

	maxVal, maxNum := max3(
		nw.Table[i-1][j-1].Distance+nw.SF[nw.LeftSequence.Value[i-1]][nw.TopSequence.Value[j-1]],
		nw.Table[i][j-1].Distance+nw.GapValue,
		nw.Table[i-1][j].Distance+nw.GapValue,
	)

	nw.Table[i][j] = &Cell{
		Distance: maxVal,
	}
	curCell := nw.Table[i][j]

	switch maxNum {
	case 1:
		curCell.Dir = DiagonalDirection
	case 2:
		curCell.Dir = LeftDirection
	case 3:
		curCell.Dir = TopDirection
	}
}

// Достаточно примитивная функция выбора максимума из трех с указанием номера максимального
func max3(a, b, c int) (int, int) {
	if a >= b {
		if a >= c {
			return a, 1
		}
		return c, 3
	}
	if b >= c {
		return b, 2
	}
	return c, 3
}

func max2(a, b int) (int, bool) {
	if a >= b {
		return a, true
	}
	return b, false
}
