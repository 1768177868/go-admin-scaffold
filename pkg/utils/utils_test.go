package utils

import (
	"testing"
	"time"
)

func TestStringSliceContains(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{"existing item", "banana", true},
		{"non-existing item", "grape", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(slice, tt.str); got != tt.expected {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"zero length", 0},
		{"small length", 10},
		{"medium length", 32},
		{"large length", 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomString(tt.length)
			if err != nil {
				t.Errorf("GenerateRandomString() error = %v", err)
				return
			}
			if len(got) != tt.length {
				t.Errorf("GenerateRandomString() length = %v, want %v", len(got), tt.length)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@sub.example.com", true},
		{"invalid email without @", "testexample.com", false},
		{"invalid email without domain", "test@", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.expected {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		time     time.Time
		layout   string
		expected string
	}{
		{"default layout", testTime, "", "2024-01-01 12:00:00"},
		{"custom layout", testTime, "2006/01/02", "2024/01/01"},
		{"RFC3339", testTime, time.RFC3339, "2024-01-01T12:00:00Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTime(tt.time, tt.layout); got != tt.expected {
				t.Errorf("FormatTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		length   int
		expected string
	}{
		{"no truncation needed", "hello", 10, "hello"},
		{"truncation needed", "hello world", 5, "hello..."},
		{"empty string", "", 5, ""},
		{"zero length", "hello", 0, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateString(tt.str, tt.length); got != tt.expected {
				t.Errorf("TruncateString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsValidIPAddress(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"valid IPv4", "192.168.1.1", true},
		{"valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"invalid IP", "256.256.256.256", false},
		{"empty string", "", false},
		{"invalid format", "192.168.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIPAddress(tt.ip); got != tt.expected {
				t.Errorf("IsValidIPAddress() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"bytes", 500, "500 B"},
		{"kilobytes", 1024 * 1024, "1.0 MB"},
		{"megabytes", 1024 * 1024 * 1024, "1.0 GB"},
		{"zero", 0, "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFileSize(tt.bytes); got != tt.expected {
				t.Errorf("FormatFileSize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{"simple string", "Hello World", "hello-world"},
		{"special characters", "Hello! World@#$%", "hello-world"},
		{"multiple spaces", "Hello   World", "hello-world"},
		{"leading/trailing hyphens", "-hello world-", "hello-world"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.str); got != tt.expected {
				t.Errorf("Slugify() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"valid international", "+1234567890", true},
		{"valid local", "1234567890", true},
		{"invalid characters", "123-456-7890", false},
		{"too short", "123", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPhoneNumber(tt.phone); got != tt.expected {
				t.Errorf("IsValidPhoneNumber() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		expected []string
	}{
		{
			"with duplicates",
			[]string{"apple", "banana", "apple", "orange", "banana"},
			[]string{"apple", "banana", "orange"},
		},
		{
			"no duplicates",
			[]string{"apple", "banana", "orange"},
			[]string{"apple", "banana", "orange"},
		},
		{
			"empty slice",
			[]string{},
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDuplicates(tt.slice)
			if len(got) != len(tt.expected) {
				t.Errorf("RemoveDuplicates() length = %v, want %v", len(got), len(tt.expected))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.expected)
					break
				}
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{"simple string", "hello", "olleh"},
		{"with spaces", "hello world", "dlrow olleh"},
		{"palindrome", "radar", "radar"},
		{"empty string", "", ""},
		{"unicode", "你好世界", "界世好你"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.str); got != tt.expected {
				t.Errorf("ReverseString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
