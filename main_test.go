package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConversion(t *testing.T, dir string) {
	dir = filepath.Join("testdata", dir)

	cmd := exec.Command("lingo")
	cmd.Dir = dir

	err := cmd.Run()
	require.NoError(t, err)

	expectedDir := filepath.Join(dir, "expected")
	expected, err := os.ReadDir(expectedDir)
	require.NoError(t, err)

	expectedFiles := make(map[string][]byte, len(expected))
	for _, entry := range expected {
		if !entry.IsDir() {
			contents, err := os.ReadFile(filepath.Join(expectedDir, entry.Name()))
			require.NoError(t, err)

			expectedFiles[entry.Name()] = contents
		}
	}

	actual, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, entry := range actual {
		if !entry.IsDir() {
			expected, ok := expectedFiles[entry.Name()]
			if !assert.Truef(t, ok, "unexpected file %v", entry.Name()) {
				continue
			}

			contents, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			require.NoError(t, err)

			assert.Equal(t, expected, contents)
			delete(expectedFiles, entry.Name())
		}
	}

	if len(expectedFiles) != 0 {
		missing := make([]string, 0, len(expectedFiles))
		for k := range expectedFiles {
			missing = append(missing, k)
		}
		sort.Strings(missing)

		assert.Failf(t, "missing files: %v", strings.Join(missing, ", "))
	}
}

func TestLingo(t *testing.T) {
	testConversion(t, "lingo")
}

func TestTour(t *testing.T) {
	testConversion(t, "tour")
}

func TestHanoi(t *testing.T) {
	testConversion(t, "hanoi")
}

func TestEmpty(t *testing.T) {
	testConversion(t, "empty")
}
