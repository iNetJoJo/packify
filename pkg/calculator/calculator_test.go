package calculator

import (
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name               string
		itemsOrdered       int
		availablePackSizes []int
		expectedResult     PackResult
		expectError        bool
	}{
		{
			name:               "Order 1 item",
			itemsOrdered:       1,
			availablePackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedResult: PackResult{
				PackCounts:  map[int]int{250: 1},
				TotalPacks:  1,
				TotalItems:  250,
				ExcessItems: 249,
			},
			expectError: false,
		},
		{
			name:               "Order 250 items",
			itemsOrdered:       250,
			availablePackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedResult: PackResult{
				PackCounts:  map[int]int{250: 1},
				TotalPacks:  1,
				TotalItems:  250,
				ExcessItems: 0,
			},
			expectError: false,
		},
		{
			name:               "Order 251 items",
			itemsOrdered:       251,
			availablePackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedResult: PackResult{
				PackCounts:  map[int]int{500: 1},
				TotalPacks:  1,
				TotalItems:  500,
				ExcessItems: 249,
			},
			expectError: false,
		},
		{
			name:               "Order 501 items",
			itemsOrdered:       501,
			availablePackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedResult: PackResult{
				PackCounts:  map[int]int{500: 1, 250: 1},
				TotalPacks:  2,
				TotalItems:  750,
				ExcessItems: 249,
			},
			expectError: false,
		},
		{
			name:               "Order 12001 items",
			itemsOrdered:       12001,
			availablePackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedResult: PackResult{
				PackCounts:  map[int]int{5000: 2, 2000: 1, 250: 1},
				TotalPacks:  4,
				TotalItems:  12250,
				ExcessItems: 249,
			},
			expectError: false,
		},
		{
			name:               "Order 0 items",
			itemsOrdered:       0,
			availablePackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedResult:     PackResult{},
			expectError:        true,
		},
		{
			name:               "Order -1 items",
			itemsOrdered:       -1,
			availablePackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedResult:     PackResult{},
			expectError:        true,
		},
		{
			name:               "No pack sizes available",
			itemsOrdered:       100,
			availablePackSizes: []int{},
			expectedResult:     PackResult{},
			expectError:        true,
		},
		{
			name:               "Custom pack sizes",
			itemsOrdered:       800,
			availablePackSizes: []int{200, 400, 600},
			expectedResult: PackResult{
				PackCounts:  map[int]int{600: 1, 200: 1},
				TotalPacks:  2,
				TotalItems:  800,
				ExcessItems: 0,
			},
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CalculatePacks(tc.itemsOrdered, tc.availablePackSizes)

			// Check error
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Skip further checks if we expected an error
			if tc.expectError {
				return
			}

			// Check total packs
			if result.TotalPacks != tc.expectedResult.TotalPacks {
				t.Errorf("Expected %d total packs, got %d", tc.expectedResult.TotalPacks, result.TotalPacks)
			}

			// Check total items
			if result.TotalItems != tc.expectedResult.TotalItems {
				t.Errorf("Expected %d total items, got %d", tc.expectedResult.TotalItems, result.TotalItems)
			}

			// Check excess items
			if result.ExcessItems != tc.expectedResult.ExcessItems {
				t.Errorf("Expected %d excess items, got %d", tc.expectedResult.ExcessItems, result.ExcessItems)
			}

			// Check pack counts
			if len(result.PackCounts) != len(tc.expectedResult.PackCounts) {
				t.Errorf("Expected %d different pack sizes, got %d", len(tc.expectedResult.PackCounts), len(result.PackCounts))
			}

			for size, count := range tc.expectedResult.PackCounts {
				if result.PackCounts[size] != count {
					t.Errorf("Expected %d packs of size %d, got %d", count, size, result.PackCounts[size])
				}
			}
		})
	}
}

func TestOptimalCalculatePacks(t *testing.T) {
	// Test cases for small orders (should use CalculatePacks)
	smallOrders := []struct {
		name         string
		itemsOrdered int
		packSizes    []int
	}{
		{
			name:         "Small order with standard pack sizes",
			itemsOrdered: 501,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
		},
		{
			name:         "Small order with many pack sizes",
			itemsOrdered: 800,
			packSizes:    []int{100, 200, 300, 400, 500, 600, 700},
		},
	}

	// Test cases for large orders (should use CalculatePacksOptimized)
	largeOrders := []struct {
		name         string
		itemsOrdered int
		packSizes    []int
	}{
		{
			name:         "Large order with standard pack sizes",
			itemsOrdered: 10000,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
		},
		{
			name:         "Large order with many pack sizes",
			itemsOrdered: 5000,
			packSizes:    []int{100, 200, 300, 400, 500, 600, 700},
		},
	}

	// Test small orders - results should match CalculatePacks
	for _, tc := range smallOrders {
		t.Run(tc.name, func(t *testing.T) {
			// Get results from both functions
			optimalResult, err1 := OptimalCalculatePacks(tc.itemsOrdered, tc.packSizes)
			directResult, err2 := CalculatePacks(tc.itemsOrdered, tc.packSizes)

			// Both should have the same error status
			if (err1 != nil) != (err2 != nil) {
				t.Errorf("Error status mismatch: OptimalCalculatePacks error: %v, CalculatePacks error: %v", err1, err2)
			}

			// Skip further checks if we got errors
			if err1 != nil || err2 != nil {
				return
			}

			// Results should be identical
			if optimalResult.TotalItems != directResult.TotalItems {
				t.Errorf("TotalItems mismatch: OptimalCalculatePacks: %d, CalculatePacks: %d", 
					optimalResult.TotalItems, directResult.TotalItems)
			}

			if optimalResult.TotalPacks != directResult.TotalPacks {
				t.Errorf("TotalPacks mismatch: OptimalCalculatePacks: %d, CalculatePacks: %d", 
					optimalResult.TotalPacks, directResult.TotalPacks)
			}

			if optimalResult.ExcessItems != directResult.ExcessItems {
				t.Errorf("ExcessItems mismatch: OptimalCalculatePacks: %d, CalculatePacks: %d", 
					optimalResult.ExcessItems, directResult.ExcessItems)
			}
		})
	}

	// Test large orders - results should match CalculatePacksOptimized
	for _, tc := range largeOrders {
		t.Run(tc.name, func(t *testing.T) {
			// Get results from both functions
			optimalResult, err1 := OptimalCalculatePacks(tc.itemsOrdered, tc.packSizes)
			directResult, err2 := CalculatePacksOptimized(tc.itemsOrdered, tc.packSizes)

			// Both should have the same error status
			if (err1 != nil) != (err2 != nil) {
				t.Errorf("Error status mismatch: OptimalCalculatePacks error: %v, CalculatePacksOptimized error: %v", err1, err2)
			}

			// Skip further checks if we got errors
			if err1 != nil || err2 != nil {
				return
			}

			// Results should be identical
			if optimalResult.TotalItems != directResult.TotalItems {
				t.Errorf("TotalItems mismatch: OptimalCalculatePacks: %d, CalculatePacksOptimized: %d", 
					optimalResult.TotalItems, directResult.TotalItems)
			}

			if optimalResult.TotalPacks != directResult.TotalPacks {
				t.Errorf("TotalPacks mismatch: OptimalCalculatePacks: %d, CalculatePacksOptimized: %d", 
					optimalResult.TotalPacks, directResult.TotalPacks)
			}

			if optimalResult.ExcessItems != directResult.ExcessItems {
				t.Errorf("ExcessItems mismatch: OptimalCalculatePacks: %d, CalculatePacksOptimized: %d", 
					optimalResult.ExcessItems, directResult.ExcessItems)
			}
		})
	}

	// Test error cases
	errorCases := []struct {
		name         string
		itemsOrdered int
		packSizes    []int
	}{
		{
			name:         "Negative order size",
			itemsOrdered: -1,
			packSizes:    []int{250, 500, 1000},
		},
		{
			name:         "Zero order size",
			itemsOrdered: 0,
			packSizes:    []int{250, 500, 1000},
		},
		{
			name:         "Empty pack sizes",
			itemsOrdered: 100,
			packSizes:    []int{},
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := OptimalCalculatePacks(tc.itemsOrdered, tc.packSizes)
			if err == nil {
				t.Errorf("Expected error for %s but got none", tc.name)
			}
		})
	}
}
