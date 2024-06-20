package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
)

func ReadFile(fp string, prec uint) ([]*big.Float, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("failed to close file: %w; also, failed with error: %v", closeErr, err)
			} else {
				err = fmt.Errorf("failed to close file: %w", closeErr)
			}
		}
	}()

	scanner := bufio.NewScanner(file)
	var numbers []*big.Float

	for scanner.Scan() {
		text := scanner.Text()
		floatVal, _, err := big.ParseFloat(text, 10, prec, big.ToNearestEven)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float from string %q: %w", text, err)
		}
		numbers = append(numbers, floatVal)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("failed to read file: %w", scanErr)
	}

	return numbers, nil
}

func CompensatedSum(numbers []*big.Float, prec uint) big.Float {
	absI := new(big.Float).SetPrec(prec)
	absJ := new(big.Float).SetPrec(prec)

	sort.Slice(numbers, func(i, j int) bool {
		absI = absI.Abs(numbers[i])
		absJ = absJ.Abs(numbers[j])
		return absI.Cmp(absJ) < 0
	})

	sum := new(big.Float).SetPrec(prec)
	comp := new(big.Float).SetPrec(prec)

	temp := new(big.Float).SetPrec(prec)

	for _, number := range numbers {
		temp.Sub(number, comp)
		comp.Add(sum, temp)
		comp.Sub(comp, sum)
		comp.Sub(comp, temp)
		sum.Add(sum, temp)
	}

	return *sum
}

func NaiveSum(numbers []*big.Float, prec uint) big.Float {
	sum := new(big.Float).SetPrec(prec)

	for _, number := range numbers {
		sum.Add(sum, number)
	}

	return *sum
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing filename argument")
	} else if len(os.Args) > 2 {
		log.Fatal("too many arguments")
	}
	filename := os.Args[1]

	var prec uint = 128
	numbers, err := ReadFile(filename, prec)
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	sum := CompensatedSum(numbers, prec)
	naiveSum := NaiveSum(numbers, prec)

	fmt.Printf("Compensated Sum: %s\n", sum.Text('g', -1))
	fmt.Printf("Naive Sum: %s\n", naiveSum.Text('g', -1))
}
