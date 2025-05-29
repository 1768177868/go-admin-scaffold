package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"regexp"
	"strings"
	"time"
)

// StringSliceContains checks if a string slice contains a specific string
func StringSliceContains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, int(math.Ceil(float64(length)/1.33333333333)))
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

// IsValidEmail checks if a string is a valid email address
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// FormatTime formats time.Time to a specified layout
func FormatTime(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}

// ParseTime parses a time string with a specified layout
func ParseTime(timeStr, layout string) (time.Time, error) {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return time.Parse(layout, timeStr)
}

// TruncateString truncates a string to a specified length and adds ellipsis
func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// IsValidIPAddress checks if a string is a valid IP address
func IsValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ToJSON converts an interface to a JSON string
func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON converts a JSON string to an interface
func FromJSON(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// FormatFileSize formats a file size in bytes to a human-readable string
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Slugify converts a string to a URL-friendly slug
func Slugify(str string) string {
	// Convert to lowercase
	str = strings.ToLower(str)

	// Replace spaces with hyphens
	str = strings.ReplaceAll(str, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	str = reg.ReplaceAllString(str, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	str = reg.ReplaceAllString(str, "-")

	// Remove leading and trailing hyphens
	str = strings.Trim(str, "-")

	return str
}

// ExtractNumbers extracts all numbers from a string
func ExtractNumbers(str string) []int {
	reg := regexp.MustCompile(`\d+`)
	matches := reg.FindAllString(str, -1)

	var numbers []int
	for _, match := range matches {
		if num, err := fmt.Sscanf(match, "%d"); err == nil {
			numbers = append(numbers, num)
		}
	}
	return numbers
}

// IsValidPhoneNumber checks if a string is a valid phone number
func IsValidPhoneNumber(phone string) bool {
	pattern := `^\+?[1-9]\d{1,14}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(phone)
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// ReverseString reverses a string
func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
