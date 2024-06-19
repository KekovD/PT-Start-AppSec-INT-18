package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
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

func FloatSum(numbers []*big.Float, prec uint) big.Float {
	sum := new(big.Float).SetPrec(prec)
	comp := new(big.Float).SetPrec(prec)

	x := new(big.Float).SetPrec(prec)
	y := new(big.Float).SetPrec(prec)

	for _, number := range numbers {
		x.Sub(number, comp)
		y.Add(sum, x)
		comp.Sub(y, sum)
		comp.Sub(comp, x)
		sum.Set(y)
	}

	return *sum
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing filename argument")
	}
	filename := os.Args[1]

	var prec uint = 256
	numbers, err := ReadFile(filename, prec)
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	sum := FloatSum(numbers, prec)

	fmt.Println(&sum)
}
