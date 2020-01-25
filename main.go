package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readlines(path string) []string {
	var fp *os.File
	var err error
	fp, err = os.Open(path)
	check(err)
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	lines := []string{}
	for scanner.Scan() {
		line := string(scanner.Text())
		lines = append(lines, line)
	}
	err = scanner.Err()
	check(err)
	return lines
}

func LinesToData(lines []string) (int, []Item, int) {
	if len(lines) < 4 {
		panic("Not enough lines")
	}
	n, err := strconv.Atoi(lines[0])
	check(err)
	sizes := strings.Split(lines[1], ",")
	values := strings.Split(lines[2], ",")
	if n != len(sizes) || len(sizes) != len(values) {
		panic("Different length")
	}
	items := []Item{}
	for i := 0; i < n; i++ {
		s, err := strconv.Atoi(strings.TrimSpace(sizes[i]))
		check(err)
		v, err := strconv.Atoi(strings.TrimSpace(values[i]))
		check(err)
		item := Item{
			No:    i + 1, // 1-based indexing
			Size:  s,
			Value: v,
		}
		items = append(items, item)
	}
	c, err := strconv.Atoi(lines[3])
	check(err)
	return n, items, c
}

type Item struct {
	No    int
	Size  int
	Value int
}

func (i Item) String() string {
	return fmt.Sprintf("Item {no: %d, size: %d, value: %d}", i.No, i.Size, i.Value)
}

type KnapsackProblemSolver struct {
	Items []Item
}

type KnapsackProblemResult struct {
	*KnapsackProblemSolver
	LastItem []*Item
	Capacity int
}

func NewKnapsackProblemSolver(items []Item) *KnapsackProblemSolver {
	k := new(KnapsackProblemSolver)
	k.Items = items
	return k
}

func (k *KnapsackProblemSolver) Solve(capacity int) *KnapsackProblemResult {
	total := make([]int, capacity+1, capacity+1)
	res := &KnapsackProblemResult{k, make([]*Item, capacity+1, capacity+1), capacity}
	for _, item := range k.Items {
		for j := item.Size; j <= capacity; j++ {
			new_value := total[j-item.Size] + item.Value
			if new_value > total[j] {
				total[j] = new_value
				clone := item
				res.LastItem[j] = &clone
			}
		}
	}
	return res
}

func (k *KnapsackProblemResult) ShowResult(capacity int) ([]int, int, int) {
	if capacity > k.Capacity {
		panic("Over capacity of result")
	}
	num := len(k.Items)
	countItems := make([]int, num, num)
	contents := []*Item{}
	for j := capacity; j > 0 && k.LastItem[j] != nil; {
		item := k.LastItem[j]
		contents = append(contents, item)
		j -= item.Size
		countItems[item.No-1] += 1
	}
	totalSize := 0
	totalValue := 0
	for i, cnt := range countItems {
		item := k.Items[i]
		totalSize += item.Size * cnt
		totalValue += item.Value * cnt
		fmt.Println(item, "x", cnt)
	}
	fmt.Println("Total Size:", totalSize)
	fmt.Println("Total Value:", totalValue)
	return countItems, totalSize, totalValue
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Input filename required")
		os.Exit(1)
	}
	wd, _ := os.Getwd()
	path := filepath.Join(wd, args[1])
	lines := readlines(path)
	num, items, limit := LinesToData(lines)
	fmt.Println("Kind of Items:", num)
	fmt.Println("Capacity:", limit)
	knapsack := NewKnapsackProblemSolver(items)
	solver := knapsack.Solve(limit)
	solver.ShowResult(limit)
}
