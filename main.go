package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var BIP39Words []string

// Reads BIP39 words from a file and returns them as a slice of strings
func readBIP39FromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return words, nil
}

// Converts a list of indices to a mnemonic phrase (slice of words)
func indicesToMnemonic(indices []int) []string {
	phrase := make([]string, len(indices))
	for i, idx := range indices {
		phrase[i] = BIP39Words[idx]
	}
	return phrase
}

// generates mnemonic phrases
func mnemonicGenerator(startIndices []int) func() ([]string, bool) {
	// if no starting point is given, start from the first unique combination (0, 1, 2, ..., 11)
	if startIndices == nil {
		startIndices = make([]int, 12)
		for i := range startIndices {
			startIndices[i] = i // Initialize with unique indices: 0, 1, 2, ..., 11
		}
	}

	current := append([]int(nil), startIndices...) // Copy of startIndices
	wordCount := len(BIP39Words)

	return func() ([]string, bool) {
		// Yield the current combination as a mnemonic phrase
		phrase := indicesToMnemonic(current)

		// Increment the current indices (ensuring uniqueness)
		for i := 11; i >= 0; i-- {
			if current[i] < wordCount-1 {
				// Only increment if it's less than the maximum word count
				current[i]++
				// Ensure all previous indices are set to a unique value less than current[i]
				for j := i + 1; j < 12; j++ {
					current[j] = current[j-1] + 1 // Ensure uniqueness by incrementing
				}
				break
			} else {
				current[i] = 0 // Reset current index and carry over to the next higher digit
			}
		}

		// If we've exhausted all combinations (i.e., all indices are unique and at their max)
		if isZeroSlice(current) {
			return phrase, false // False indicates the generator is done
		}

		return phrase, true // True means more combinations to generate
	}
}

// Helper function to check if all elements in the slice are zero
func isZeroSlice(slice []int) bool {
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}

func main() {
	var err error
	BIP39Words, err = readBIP39FromFile("english.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Start generating from a specific point (e.g., all zeros)
	gen := mnemonicGenerator(nil)

	// Simulating a process that stops after generating 100 mnemonics
	for i := 0; i < 100; i++ {
		phrase, more := gen()
		if !more {
			fmt.Println("Reached the end of combinations.")
			break
		}
		fmt.Printf("Mnemonic #%d: %v\n", i+1, phrase)
	}

	// // Example of restarting from a specific point (e.g., indices [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100])
	// fmt.Println("\nResuming from specific point...")
	// specificStart := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100}
	// gen = mnemonicGenerator(specificStart)

	// // Continue generating from the saved state
	// for i := 0; i < 10; i++ { // Generate the next 10 phrases
	// 	phrase, more := gen()
	// 	if !more {
	// 		fmt.Println("Reached the end of combinations.")
	// 		break
	// 	}
	// 	fmt.Printf("Mnemonic (resumed) #%d: %v\n", i+1, phrase)
	// }
}
