package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shirerpeton/subResyncer/internal/syncer"
)

func getOutputPath(input string) string {
	parts := strings.Split(input, ".")
	if len(parts) == 1 {
		return input + "_synced"
	}
	parts[len(parts) - 2] = parts[len(parts) - 2] + "_synced"
	return strings.Join(parts, ".")
}

func main() {
	sub := flag.String("sub", "", "Path to input subtitle file or directory containing them")
	shiftF := flag.Float64("shift", 1.0, "By how much to shift dialog lines in subtitle (in seconds, decimal, could be negative e.g. 1.5, 5.0, -0.25)")
	output := flag.String("out", "", "Path to output subtitle file, defaults to input filename with synced suffix, for diretory processing must be a directory name as well")
	flag.Parse()

	if *sub == "" {
		fmt.Println("provide input subtitle file path")
		os.Exit(1)
	}
	subStat, err := os.Stat(*sub)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var outputPath string
	if *output == "" {
		outputPath = *output
	} else {
		outputPath = getOutputPath(*sub)
	}

	shift := time.Duration(float64(time.Second) * (*shiftF))

	if !subStat.IsDir() {
		result, err := syncer.Sync(*sub, shift)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = os.WriteFile(outputPath, []byte(result), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("done - %s\n", outputPath)
	} else {
		fmt.Println("dir")
	}
}
