package main

import (
	"os"
	"testing"
)

// Test the readBIP39FromFile function
func TestReadBIP39FromTempFile(t *testing.T) {
	// Create a temporary file with some BIP39 words for testing
	tmpFile, err := os.CreateTemp("", "bip39_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// sample BIP39 words to the temp file
	sampleWords := `abandon
ability
able
about
above
absent
`
	if _, err := tmpFile.WriteString(sampleWords); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file so it can be read
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Call the function to read the words from the file
	words, err := readBIP39FromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading from file: %v", err)
	}

	expectedWords := []string{"abandon", "ability", "able", "about", "above", "absent"}

	if len(words) != len(expectedWords) {
		t.Fatalf("Expected %d words, got %d", len(expectedWords), len(words))
	}

	for i, word := range expectedWords {
		if words[i] != word {
			t.Errorf("Expected word %d to be %q, but got %q", i, word, words[i])
		}
	}
}

func TestReadBIP39FromActualFile(t *testing.T) {
	words, err := readBIP39FromFile("english.txt")
	if err != nil {
		t.Fatalf("Error reading from file: %v", err)
	}

	expectedLen := 2048
	actualLen := len(words)

	// Compare the result with the expected output
	if actualLen != expectedLen {
		t.Fatalf("Expected %d words, got %d", expectedLen, actualLen)
	}
}
