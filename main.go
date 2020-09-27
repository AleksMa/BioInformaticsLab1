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
	mode string
)

func init() {
	flag.IntVar(&gap, "gap", -10, "gap value")
	flag.StringVar(&outputFile, "out", "", "output file")
	flag.StringVar(&mode, "mode", "simple", "scoring mode")
}

func main() {
	flag.Parse()

	// Селектор режима скоринговой функции по ключу
	var sf ScoringFunc
	switch mode {
	case "simple":
		sf = SimpleFunc
	case "dnafull":
		sf = DNAFull
	case "blosum62":
		sf = Blosum62
	default:
		fmt.Println("Unexpected mode! Expected: {simple | dnafull | blosum62}")
	}

	inputFiles = flag.Args()

	if len(inputFiles) == 0 {
		return
	}
	// Используем парсер fasta из "нулевой" лабораторной работы
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

	// Создание обертки решения и непосредственно решение
	nw := NewNeedlemanWunsch(seqs[0], seqs[1], sf, gap)
	a, b, score := nw.Solve()

	// Вывод в стандартный поток либо в файл, заданный через параметр
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
