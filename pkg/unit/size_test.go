package unit

import "testing"

func TestSize_ToString(t *testing.T) {
	tests := []struct {
		name     string
		value    Size
		expected string
	}{
		{
			name:     "TestSizeToStringB",
			value:    B,
			expected: "1B",
		},
		{
			name:     "TestSizeToStringKB",
			value:    KB,
			expected: "1KB",
		},
		{
			name:     "TestSizeToStringMB",
			value:    MB,
			expected: "1MB",
		},
		{
			name:     "TestSizeToStringGB",
			value:    GB,
			expected: "1GB",
		},
		{
			name:     "TestSizeToStringTB",
			value:    TB,
			expected: "1TB",
		},
		{
			name:     "TestSizeToStringGB",
			value:    GB * 50,
			expected: "50GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.ToString(); got != tt.expected {
				t.Errorf("got = %v, expected %v", got, tt.expected)
			}
		})
	}
}
