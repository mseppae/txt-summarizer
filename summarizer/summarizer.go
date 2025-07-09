package summarizer

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type KeySum struct {
	Key string
	Sum int
}

// ParseAndSumFiles reads the given files and returns a map of key to summed values.
func ParseAndSumFiles(filePaths []string) (map[string]int, error) {
	keySums := make(map[string]int)
	re := regexp.MustCompile(`^([0-9]+)\s+(.+)$`)

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if matches := re.FindStringSubmatch(line); len(matches) == 3 {
				num, err := strconv.Atoi(matches[1])
				if err != nil {
					file.Close()
					return nil, fmt.Errorf("error parsing number in line '%s' in file %s: %w", line, filePath, err)
				}
				key := matches[2]
				keySums[key] += num
			}
		}
		if err := scanner.Err(); err != nil {
			file.Close()
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}
		file.Close()
	}
	return keySums, nil
}

// SortKeySums returns a sorted slice of KeySum (descending by sum).
func SortKeySums(keySums map[string]int) []KeySum {
	var keySumPairs []KeySum
	for key, sum := range keySums {
		keySumPairs = append(keySumPairs, KeySum{Key: key, Sum: sum})
	}
	sort.Slice(keySumPairs, func(i, j int) bool {
		return keySumPairs[i].Sum > keySumPairs[j].Sum
	})
	return keySumPairs
}

// WriteSummary writes the sorted key sums to the given writer.
func WriteSummary(writer *bufio.Writer, keySumPairs []KeySum) error {
	for _, pair := range keySumPairs {
		_, err := fmt.Fprintf(writer, "%d %s\n", pair.Sum, pair.Key)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
