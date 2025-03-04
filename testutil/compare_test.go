package testutil

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestEqual(t *testing.T) {
	type model struct {
		A string
		B int
		C bool
	}

	type args[T any] struct {
		given   T
		input   T
		isEqual bool
	}
	dataTypes := map[string]any{
		"int": map[string]args[int]{
			"equal": {
				given:   1,
				input:   1,
				isEqual: true,
			},
			"not equal": {
				given:   1,
				input:   2,
				isEqual: false,
			},
		},
		"bool": map[string]args[bool]{
			"equal": {
				given:   true,
				input:   true,
				isEqual: true,
			},
			"not equal": {
				given:   true,
				input:   false,
				isEqual: false,
			},
		},
		"string": map[string]args[string]{
			"equal": {
				given:   "string",
				input:   "string",
				isEqual: true,
			},
			"not equal": {
				given:   "string_1",
				input:   "string_2",
				isEqual: false,
			},
		},
		"slice": map[string]args[[]int]{
			"equal": {
				given:   []int{1, 2, 3},
				input:   []int{1, 2, 3},
				isEqual: true,
			},
			"not equal": {
				given:   []int{1, 2, 3},
				input:   []int{1, 2, 6},
				isEqual: false,
			},
		},
		"map": map[string]args[map[string]int]{
			"equal": {
				given:   map[string]int{"a": 1, "b": 2},
				input:   map[string]int{"a": 1, "b": 2},
				isEqual: true,
			},
			"not equal": {
				given:   map[string]int{"a": 1, "b": 2},
				input:   map[string]int{"a": 5, "b": 6},
				isEqual: false,
			},
		},
		"object": map[string]args[model]{
			"equal": {
				given:   model{A: "a", B: 2},
				input:   model{A: "a", B: 2},
				isEqual: true,
			},
			"not equal": {
				given:   model{A: "a", B: 2},
				input:   model{A: "AAA", B: 100},
				isEqual: false,
			},
		},
	}

	for dataType, tcs := range dataTypes {
		switch dataType {
		case "int":
			tcs := tcs.(map[string]args[int])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(dataType+"_"+scenario, func(t *testing.T) {
					t.Parallel()

					mockTestingT := new(MockTestingT)
					mockTestingT.On("Helper").Return()
					if !tc.isEqual {
						mockTestingT.On("Errorf", mock.AnythingOfType("string"), tc.input, tc.given, mock.AnythingOfType("string")).Return()
						mockTestingT.On("FailNow").Return()
					}

					Equal(mockTestingT, tc.input, tc.given)
				})
			}

		case "string":
			tcs := tcs.(map[string]args[string])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(dataType+"_"+scenario, func(t *testing.T) {
					t.Parallel()

					mockTestingT := new(MockTestingT)
					mockTestingT.On("Helper").Return()
					if !tc.isEqual {
						mockTestingT.On("Errorf", mock.AnythingOfType("string"), tc.input, tc.given, mock.AnythingOfType("string")).Return()
						mockTestingT.On("FailNow").Return()
					}

					Equal(mockTestingT, tc.input, tc.given)
				})
			}

		case "bool":
			tcs := tcs.(map[string]args[bool])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(dataType+"_"+scenario, func(t *testing.T) {
					t.Parallel()

					mockTestingT := new(MockTestingT)
					mockTestingT.On("Helper").Return()
					if !tc.isEqual {
						mockTestingT.On("Errorf", mock.AnythingOfType("string"), tc.input, tc.given, mock.AnythingOfType("string")).Return()
						mockTestingT.On("FailNow").Return()
					}

					Equal(mockTestingT, tc.input, tc.given)
				})
			}

		case "slice":
			tcs := tcs.(map[string]args[[]int])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(dataType+"_"+scenario, func(t *testing.T) {
					t.Parallel()

					mockTestingT := new(MockTestingT)
					mockTestingT.On("Helper").Return()
					if !tc.isEqual {
						mockTestingT.On("Errorf", mock.AnythingOfType("string"), tc.input, tc.given, mock.AnythingOfType("string")).Return()
						mockTestingT.On("FailNow").Return()
					}

					Equal(mockTestingT, tc.input, tc.given)
				})
			}

		case "map":
			tcs := tcs.(map[string]args[map[string]int])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(dataType+"_"+scenario, func(t *testing.T) {
					t.Parallel()

					mockTestingT := new(MockTestingT)
					mockTestingT.On("Helper").Return()
					if !tc.isEqual {
						mockTestingT.On("Errorf", mock.AnythingOfType("string"), tc.input, tc.given, mock.AnythingOfType("string")).Return()
						mockTestingT.On("FailNow").Return()
					}

					Equal(mockTestingT, tc.input, tc.given)
				})
			}

		case "object":
			tcs := tcs.(map[string]args[model])
			for scenario, tc := range tcs {
				tc := tc
				t.Run(dataType+"_"+scenario, func(t *testing.T) {
					t.Parallel()

					mockTestingT := new(MockTestingT)
					mockTestingT.On("Helper").Return()
					if !tc.isEqual {
						mockTestingT.On("Errorf", mock.AnythingOfType("string"), tc.input, tc.given, mock.AnythingOfType("string")).Return()
						mockTestingT.On("FailNow").Return()
					}

					Equal(mockTestingT, tc.input, tc.given)
				})
			}

		default:
		}
	}
}

func TestEqualWithIgnoreMapEntries(t *testing.T) {
	// Given
	given := map[string]any{
		"name":    "Space Marine",
		"faction": "Ultramarines",
		"power":   5,
		"alive":   true,
		"ts":      "2025-03-03T10:00:00Z", // Timestamp field to be ignored
	}
	expected := map[string]any{
		"name":    "Space Marine",
		"faction": "Ultramarines",
		"power":   5,
		"alive":   true,
		"ts":      "2025-03-02T09:00:00Z", // Different timestamp
	}

	// When
	// Then
	Equal(t, expected, given, IgnoreMapEntries(func(k string, v any) bool {
		return k == "ts"
	}))
}

func TestEqualWithIgnoreSliceMapEntries(t *testing.T) {
	// Given
	given := []map[string]any{
		{
			"name":    "Space Marine",
			"faction": "Ultramarines",
			"power":   5,
			"alive":   true,
			"ts":      "2025-03-03T10:00:00Z", // Timestamp field to be ignored
		},
	}
	expected := []map[string]any{
		{
			"name":    "Space Marine",
			"faction": "Ultramarines",
			"power":   5,
			"alive":   true,
			"ts":      "2025-03-02T09:00:00Z", // Different timestamp
		},
	}

	// When
	// Then
	Equal(t, expected, given, IgnoreSliceMapEntries(func(k string, v any) bool {
		return k == "ts"
	}))
}

func TestEqualWithIgnoreSliceElements(t *testing.T) {
	expected := []int{1, 2, 3, 4}
	given := []int{1, 2, 3, 4, 5}

	Equal(t, expected, given, IgnoreSliceElements(func(v int) bool {
		return v == 4 || v == 5 // Ignore values 4 and 5
	}))
}

func TestEqualWithIgnoreUnexported(t *testing.T) {
	type spaceMarine struct {
		Name     string
		Age      int
		Faction  string
		isFallen bool
	}
	expected := spaceMarine{Name: "Gabriel Angelos", Age: 300, Faction: "Blood Ravens", isFallen: false}
	given := spaceMarine{Name: "Gabriel Angelos", Age: 300, Faction: "Blood Ravens", isFallen: true}

	Equal(t, expected, given, IgnoreUnexported[spaceMarine](spaceMarine{}))
}

func TestEqualWithEquateComparable(t *testing.T) {
	type spaceMarine struct {
		Name     string
		Age      int
		Faction  string
		isFallen bool
	}
	expected := spaceMarine{Name: "Gabriel Angelos", Age: 300, Faction: "Blood Ravens", isFallen: true}
	given := spaceMarine{Name: "Gabriel Angelos", Age: 300, Faction: "Blood Ravens", isFallen: true}

	Equal(t, expected, given, EquateComparable[spaceMarine](spaceMarine{}))
}
