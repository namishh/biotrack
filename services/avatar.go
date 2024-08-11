// An avatar service like https://avatar.vercel.sh
package services

import (
	"crypto/sha256"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

func djb2(str string) int {
	hash := 5381
	for _, char := range str {
		hash = ((hash << 5) + hash) + int(char)
	}
	return hash
}

func GenerateGradient(username string) map[string]string {
	colors := []string{
		"#a855f7", // Purple
		"#7c3aed",
		"#f43f5e",
		"#faf5ff", // White
	}

	hash := sha256.Sum256([]byte(username))
	h := float64(int(hash[0]) % len(colors)) // Use the first byte of the hash to select a color
	c1, _ := colorful.Hex(colors[int(h)])    // First color

	// Calculate the next color in the palette
	h2 := math.Mod(h+1.0, float64(len(colors)))
	c2, _ := colorful.Hex(colors[int(h2)]) // Second color

	return map[string]string{
		"fromColor": c1.Hex(),
		"toColor":   c2.Hex(),
	}
}
