package converter

import (
	"github.com/K-Phoen/sdk"
)

func stringInSlice(search string, haystack []string) bool {
	for _, item := range haystack {
		if item == search {
			return true
		}
	}

	return false
}

func boolPtr(input bool) *bool {
	return &input
}

func intPtr(input int) *int {
	return &input
}

func float64Ptr(input float64) *float64 {
	return &input
}

func strPtr(input string) *string {
	return &input
}

func sdkRepeatDirectionToYAML(repeatDirection sdk.RepeatDirection) string {
	if repeatDirection == sdk.RepeatDirectionHorizontal {
		return "vertical"
	}

	return "horizontal"
}
