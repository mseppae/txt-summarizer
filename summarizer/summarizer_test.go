package summarizer

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "testfile-*.txt")
	assert.NoError(t, err)
	_, err = tmpfile.WriteString(content)
	assert.NoError(t, err)
	tmpfile.Close()
	return tmpfile.Name()
}

func TestParseAndSumFiles_Simple(t *testing.T) {
	fileContent := "10 apples\n5 oranges\n3 apples\n"
	file := createTempFile(t, fileContent)
	defer os.Remove(file)

	result, err := ParseAndSumFiles([]string{file})
	assert.NoError(t, err)
	assert.Equal(t, 13, result["apples"])
	assert.Equal(t, 5, result["oranges"])
}

func TestParseAndSumFiles_MultipleFiles(t *testing.T) {
	file1 := createTempFile(t, "2 a\n3 b\n")
	file2 := createTempFile(t, "4 a\n1 c\n")
	defer os.Remove(file1)
	defer os.Remove(file2)

	result, err := ParseAndSumFiles([]string{file1, file2})
	assert.NoError(t, err)
	assert.Equal(t, 6, result["a"])
	assert.Equal(t, 3, result["b"])
	assert.Equal(t, 1, result["c"])
}

func TestParseAndSumFiles_BadInput(t *testing.T) {
	file := createTempFile(t, "notanumber apples\n")
	defer os.Remove(file)
	result, err := ParseAndSumFiles([]string{file})
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestParseAndSumFiles_FileNotFound(t *testing.T) {
	_, err := ParseAndSumFiles([]string{"nonexistent.txt"})
	assert.Error(t, err)
}

func TestParseAndSumFiles_IgnoreIrrelevantLines(t *testing.T) {
	fileContent := `== HEADER ==
1234 apples
5678 oranges
== FOOTER ==
=== SUMMARY ===
9999 not-a-fruit
hello world
`
	file := createTempFile(t, fileContent)
	defer os.Remove(file)

	result, err := ParseAndSumFiles([]string{file})
	assert.NoError(t, err)
	assert.Equal(t, 1234, result["apples"])
	assert.Equal(t, 5678, result["oranges"])
	assert.Equal(t, 9999, result["not-a-fruit"])
	assert.NotContains(t, result, "== HEADER ==")
	assert.NotContains(t, result, "== FOOTER ==")
	assert.NotContains(t, result, "=== SUMMARY ===")
	assert.NotContains(t, result, "hello world")
}

func TestSortKeySums(t *testing.T) {
	input := map[string]int{"a": 2, "b": 5, "c": 1}
	result := SortKeySums(input)
	assert.Len(t, result, 3)
	assert.Equal(t, "b", result[0].Key)
	assert.Equal(t, 5, result[0].Sum)
	assert.Equal(t, "c", result[2].Key)
	assert.Equal(t, 1, result[2].Sum)
}

func TestWriteSummary(t *testing.T) {
	pairs := []KeySum{{Key: "apples", Sum: 10}, {Key: "oranges", Sum: 5}}
	var sb strings.Builder
	writer := bufio.NewWriter(&sb)
	assert.NoError(t, WriteSummary(writer, pairs))
	expected := "10 apples\n5 oranges\n"
	assert.Equal(t, expected, sb.String())
}
