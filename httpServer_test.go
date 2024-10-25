package main

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		mode              string
		number            string
		expectedMode      string
		expectedThreshold int
		expectError       bool
	}{
		{"async", "5", "async", 5, false}, // Valid async mode with valid threshold
		{"sync", "10", "sync", 10, false}, // Valid sync mode with valid threshold
		{"async", "invalid", "", 0, true}, // Invalid threshold input
		{"", "5", "", 0, true},            // Empty mode
		{"sync", "", "", 0, true},         // Empty threshold
		{"unknown", "5", "", 0, true},     // Unknown mode
	}

	for _, test := range tests {
		t.Run(test.mode, func(t *testing.T) {
			result, err := initConfig(test.mode, test.number)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error for mode %s and number %s, but did not get one", test.mode, test.number)
				}
				return
			}

			if err != nil {
				t.Errorf("Did not expect an error for mode %s and number %s, but got: %v", test.mode, test.number, err)
				return
			}

			if result.Mode != test.expectedMode {
				t.Errorf("Expected mode %s, got %s", test.expectedMode, result.Mode)
			}
			if result.FailureThreshold != test.expectedThreshold {
				t.Errorf("Expected threshold %d, got %d", test.expectedThreshold, result.FailureThreshold)
			}
		})
	}
}
