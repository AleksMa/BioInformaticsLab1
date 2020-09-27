package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader *bufio.Reader

func main() {
	f, _ := os.Open("blosum62.txt")
	reader = bufio.NewReader(f)

	letterString, _ := reader.ReadString('\n')
	letters := strings.Split(letterString, " ")

	fmt.Println(letters)

	int_to_amino := make([]string, 0)
	amino_to_int := make(map[string]int)

	for _, letter := range letters {
		if len(letter) < 1 {
			continue
		}
		if len(letter) > 1 {
			letter = letter[:1]
		}
		fmt.Println("<"+letter+">")
		amino_to_int[letter] = len(int_to_amino)
		int_to_amino = append(int_to_amino, letter)
	}

	fmt.Println(int_to_amino)
	fmt.Println(amino_to_int)

	costs1 := make([][]string, 20)
	for i := 0; i < 20; i++ {
		var costsString string
		costsString, _ = reader.ReadString('\n')
		cs := strings.Split(costsString, " ")
		costs1[i] = cs[1:]
	}

	costs := make([][]string, 20)
	for i, costi := range costs1 {
		costs[i] = make([]string, 0)
		for _, costj := range costi {
			if len(costj) >= 1 {
				costs[i] = append(costs[i], costj)
			}
		}
	}

	for i := 0; i < 20; i++ {
		fmt.Println("'" + int_to_amino[i] + "': map[uint8]int{")
		for j := 0; j < 20; j++ {
			fmt.Println("'" + int_to_amino[j] + "': ", costs[i][j], ",")
		}
		fmt.Println("},")
	}
}
