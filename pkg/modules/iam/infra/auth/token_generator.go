package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// TokenGenerator generates secure random tokens for PAT
type TokenGenerator struct{}

// NewTokenGenerator creates a new token generator
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

// GeneratePAT generates a new Personal Access Token
// Returns: plainToken (pat_prefix_random), tokenHash (SHA-256), prefix (pat_prefix), error
func (g *TokenGenerator) GeneratePAT() (string, string, string, error) {
	// Generate 5-character prefix (for user identification)
	prefix, err := g.generateRandomString(5)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate prefix: %w", err)
	}

	// Generate 32-character random token
	randomPart, err := g.generateRandomString(32)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Format: pat_<prefix>_<random>
	plainToken := fmt.Sprintf("pat_%s_%s", prefix, randomPart)
	prefixOnly := "pat_" + prefix

	// Hash the full token with SHA-256
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	return plainToken, tokenHash, prefixOnly, nil
}

// HashToken hashes a plain token with SHA-256
// Used for token verification
func (g *TokenGenerator) HashToken(plainToken string) string {
	hash := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(hash[:])
}

// ValidateTokenFormat checks if a token has the correct format
// Expected format: pat_<5chars>_<32chars>
func (g *TokenGenerator) ValidateTokenFormat(token string) bool {
	parts := strings.Split(token, "_")
	if len(parts) != 3 {
		return false
	}

	if parts[0] != "pat" {
		return false
	}

	if len(parts[1]) != 5 {
		return false
	}

	if len(parts[2]) != 32 {
		return false
	}

	return true
}

// ExtractPrefix extracts the prefix from a full token
// Example: "pat_2Kj9X_abc123..." -> "pat_2Kj9X"
func (g *TokenGenerator) ExtractPrefix(token string) (string, error) {
	parts := strings.Split(token, "_")
	if len(parts) < 2 {
		return "", errors.New("invalid token format")
	}

	return fmt.Sprintf("%s_%s", parts[0], parts[1]), nil
}

// generateRandomString generates a cryptographically secure random string
// Using base64 URL-safe encoding without padding
// Note: Replaces '-' and '_' with alphanumeric chars to avoid delimiter conflicts
func (g *TokenGenerator) generateRandomString(length int) (string, error) {
	// Calculate required bytes (base64 encodes 3 bytes to 4 characters)
	// We need slightly more bytes than the target length
	byteLength := (length*3 + 3) / 4

	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode to base64 URL-safe without padding
	encoded := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Replace '-' and '_' with alphanumeric characters to avoid delimiter conflicts
	// '-' -> 'A', '_' -> 'Z' (arbitrary safe replacements)
	encoded = strings.ReplaceAll(encoded, "-", "A")
	encoded = strings.ReplaceAll(encoded, "_", "Z")

	// Trim to exact length
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}
