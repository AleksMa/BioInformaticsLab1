package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	inputFiles []string
	gap        int
	outputFile string
)

func init() {
	flag.IntVar(&gap, "gap", -2, "gap value")
	flag.StringVar(&outputFile, "out", "", "output file")
}

func main() {
	flag.Parse()

	inputFiles = flag.Args()

	if len(inputFiles) == 0 {
		return
	}

	seqs := make([]*Sequence, 0, 2)
	for _, inputFile := range inputFiles {
		f, err := os.Open(inputFile)
		if err != nil {
			log.Fatalf("can not open file: %s", err)
		}
		defer f.Close()

		p := NewFastaParser(f)
		for {
			seq, err := p.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("processing error: %s", err)
			}
			seqs = append(seqs, seq)
		}
	}

	if len(seqs) != 2 {
		log.Fatal("unexpected sequences number")
	}

	fmt.Println(seqs[0].Value)
	fmt.Println(seqs[1].Value)

	nw := NewNeedlemanWunsch(seqs[0], seqs[1], Blosum62, gap)

	a, b, score := nw.Solve()

	if outputFile == "" {
		for i := 0; ; i += 100 {
			if i + 100 > len(a) {
				fmt.Println(a[i:])
				break
			}
			fmt.Println(a[i:i+100])
		}
		for i := 0; i <= len(b); i += 100 {
			if i + 100 > len(b) {
				fmt.Println(b[i:])
				break
			}
			fmt.Println(b[i:i+100])
		}
		fmt.Println("Score:", score)
	} else {
		f, _ := os.Create(outputFile)
		w := bufio.NewWriter(f)

		for i := 0; ; i += 100 {
			if i + 100 > len(a) {
				fmt.Fprintln(w, a[i:])
				break
			}
			fmt.Fprintln(w, a[i:i+100])
		}
		for i := 0; i <= len(b); i += 100 {
			if i + 100 > len(b) {
				fmt.Fprintln(w, b[i:])
				break
			}
			fmt.Fprintln(w, b[i:i+100])
		}
		fmt.Fprintln(w, score)

		w.Flush()
	}

	//nw.Print()
}

/*

Основная задача лабораторной - обдумать, понять и аккуратно реализовать алгоритм NW (нидлмана-вунша), приведенный ниже, а также всю обвязку к нему, которая будет использоваться в остальных лабораторных работах.

Итак, даны две строки, ключ, определяющий алфавит, ключ(и), определяющие скоринговую систему. На выходе мы должны получить оптимальное глобальное пАрное выравнивание и его оценку в данной скоринговой системе.

Что такое глобальное парное выравнивание?
Для строк a и b длиной n и m это матрица из двух  новых строк a' и b' равной длины l (max(m,n)>=l>=n+m), образованных вставками пробелов (gap, '-'). В выравнивании не может быть столбца состоящего из двух пробелов.

a = AATCG, b = AACG
a': AATCG
b': AA-CG

Для нахождения глобального выравнивания мы используем подход динамического программирования.
Задаем матрицу, в каждой ячейке D(i,j) будет храниться скор оптимального выравнивания для префиксов a(1..i) и b(1..j).
Инициализируем первую строку и первый столбец:
D(0,0) = 0
D(i,0) = j*G
D(0,j) = i*G, где G - величина скоринг функции, соответствующая вставки гэпа.
Далее рекурсивно заполняем матрицу:
D(i,j),Ptr(i,j) = max{(D(i-1,j)+G, UP), (D(i,j-1)+G, LEFT), (D(i-1,j-1)+S(i,j), DIAG)},
где Ptr(i,j) - указатель на выбранное направление в матрице для восстановления выравнивания,
S(i,j) - скоринг функция, например BLOSUM62, DNAFull или иная.
Скор оптимального выравнивая для a и b: Score = D(i,j).

Пример:

http://rna.informatik.uni-freiburg.de/Teaching/index.jsp?toolName=Needleman-Wunsch


Что должно быть сделано:
1. Обдумать и осознать как сам алгоритм, так и его следствия, например, как влияет выбор скоринговой системы на результат, как можно модернизировать алгоритм, чтобы не штрафовать за гэпы вначале и/или конце выравнивания, и т.д.
2. Реализовать алгоритм, включая осмысленные тесты к нему, которые надо самостоятельно придумать. Текст должен содержать комментарии, демонстрирующие понимание и мыслительный процесс.
3. На вход принимается один файл, содержащий две последовательности в формате fasta или два файла, содержащих по одной последовательности в формате fasta:
prog_name -i seqs.fasta
prog_name -i seq1.fasta seq2.fasta (расширение, само собой, любое, в т.ч. никакое)
4. Программа должна уметь работать с алфавитами для аминокислот и нуклеотидов (и отказываться работать, если на вход приходит помойка вроде ;%?*_;31%:"). Для аминокислот использовать матрицу BLOSUM62 (во вложении):
https://www.ncbi.nlm.nih.gov/Class/BLAST/BLOSUM62.txt
для нуклеотидов - совпадение символов +5, несовпадение -4 (матрица DNAFull). (Алфавит {A,T,G,C})
Дефолтный случай - совпадение символов +1, несовпадение -1, гэп -2
5. Штраф за гэп задаётся ключом:
prog_name -g -10 или
prog_name --gap=-10
6. Задача простая, потому по возможности должна быть реализована стандартными средствами языка. Весь код должен быть своим. Исключение - можно использовать готовые библиотеки для обработки ключей командной строки,тестирования и измерения времени.
7. Результат должен выдаваться в консоль или опционально (-o out.txt) в файл. Выравнивание должно переноситься на новую строку, если оно не умещается, например, в 100 символов:

seq1: SP-E---TVIHS--GWVIWRELFSH-WPDQCKL-LFGDWFAWIHWTYLVYYSAGPPCQG
seq2: SPSDQFFTVIHSCLYWVIWRDLMSHLFMNGAAIDIHWTWDSIAIGPPLV-YPIEEVFAG

seq1: QSDIVVMMQKKLRTNFCQCYKYWYQ
seq2: PSTIVVMMQKMLRTNFCQCYKPWYQ

Отдельной строкой выводится скор:
Score: 161

Для проверки своих результатов можно воспользоваться:
для белков: https://www.ebi.ac.uk/Tools/psa/emboss_needle/
для ДНК: https://www.ebi.ac.uk/Tools/psa/emboss_needle/nucleotide.html
В настройках  step2-More Options... выбрать BLOSUM62 для белков, DNAFull для ДНК, установить равными GAP OPEN, GAP EXTEND, END GAP OPEN и END GAP EXTEND. END GAP PENALTY=True.

*/
