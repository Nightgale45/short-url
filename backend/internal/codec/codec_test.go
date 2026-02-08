package codec

import (
	"testing"
)

func TestBase62Encoder(t *testing.T) {
	type input struct {
		id   int64
		salt int64
	}

	tests := []struct {
		testName string
		input    input
		expected string
	}{
		{
			testName: "edge case id=0 salt=10_000_000",
			input: input{
				id:   int64(0),
				salt: int64(10_000_000),
			},
			expected: "FXsk",
		},
		{
			testName: "edge case id=1_000_000_000 salt=99_999_999",
			input: input{
				id:   int64(1_000_000_000),
				salt: int64(99_999_999),
			},
			expected: "9UIymWyzSf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result := Base62Encoder(tt.input.id, tt.input.salt)
			if result != tt.expected {
				t.Errorf("expected: %s, actual: %s", tt.expected, result)
			}
		})
	}

}

func TestBase62Decoder(t *testing.T) {
	type input struct {
		str string
	}

	tests := []struct {
		testName string
		input    input
		expected []int64
	}{
		{
			testName: "edge case id=0 salt=10_000_000",
			input: input{
				str: "FXsk",
			},
			expected: []int64{0, 10_000_000},
		},
		{
			testName: "edge case id=1_000_000_000 salt=99_999_999",
			input: input{
				str: "9UIymWyzSf",
			},
			expected: []int64{1_000_000_000, 99_999_999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			id, salt := Base62Decoder(tt.input.str)
			if id != tt.expected[0] {
				t.Errorf("ID - expected: %d, actual: %d", tt.expected[0], id)
			}

			if salt != tt.expected[1] {
				t.Errorf("SALT - expected: %d, actual: %d", tt.expected[1], salt)
			}
		})
	}
}
