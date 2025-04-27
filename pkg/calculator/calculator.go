package calculator

import (
	"fmt"
	"sort"
)

// PackResult represents the result of a pack calculation
type PackResult struct {
	PackCounts  map[uint64]uint64 // Map of pack size to count
	TotalPacks  uint64            // Total number of packs
	TotalItems  uint64            // Total number of items
	ExcessItems uint64            // Number of excess items
}

// String returns a string representation of the pack result
func (pr PackResult) String() string {
	return fmt.Sprintf("Packs: %v, Total packs: %d, Total items: %d, Excess items: %d",
		pr.PackCounts, pr.TotalPacks, pr.TotalItems, pr.ExcessItems)
}

// OptimalCalculatePacks chooses between CalculatePacks and CalculatePacksOptimized
// based on order size and pack sizes to ensure optimal memory usage and performance.
// For small orders, it uses CalculatePacks (pure DP approach).
// For large orders, it uses CalculatePacksOptimized (hybrid greedy/DP approach).
func OptimalCalculatePacks(itemsOrdered uint64, availablePackSizes []uint64) (PackResult, error) {
	// Threshold based on benchmark results
	// Below this threshold, the original algorithm is faster and uses less memory
	// Above this threshold, the optimized algorithm is dramatically better
	var threshold uint64 = 2500

	// If we have many pack sizes, lower the threshold as DP becomes more expensive
	if len(availablePackSizes) > 5 {
		threshold = 1000
	}

	// Choose the appropriate algorithm based on order size
	if itemsOrdered <= threshold {
		return CalculatePacks(itemsOrdered, availablePackSizes)
	} else {
		return CalculatePacksOptimized(itemsOrdered, availablePackSizes)
	}
}

// CalculatePacks determines the optimal packing solution
// this implementation uses dynamic programming
// uses slices to store the pack sizes and their counts
// memory usage grows with order size
func CalculatePacks(itemsOrdered uint64, availablePackSizes []uint64) (PackResult, error) {
	if itemsOrdered <= 0 {
		return PackResult{}, fmt.Errorf("items ordered must be positive")
	}

	if len(availablePackSizes) == 0 {
		return PackResult{}, fmt.Errorf("no pack sizes available")
	}

	// Sort pack sizes in descending order
	sort.Slice(availablePackSizes, func(i, j int) bool {
		return availablePackSizes[i] > availablePackSizes[j]
	})

	// Find minimum possible items to ship
	minItems, packsUsed := findMinimumItems(itemsOrdered, availablePackSizes)

	result := PackResult{
		PackCounts:  packsUsed,
		TotalItems:  minItems,
		ExcessItems: minItems - itemsOrdered,
	}

	// Calculate total number of packs
	for _, count := range packsUsed {
		result.TotalPacks += count
	}

	return result, nil
}

// CalculatePacksOptimized calculates pack distribution to fulfill an order using an optimized approach with greedy and DP.
// Uses maps to store the pack sizes and their counts
// This is a more memory efficient solution for large orders
// Uses a hybrid approach with greedy algorithm for large portions and DP for smaller amounts
func CalculatePacksOptimized(itemsOrdered uint64, availablePackSizes []uint64) (PackResult, error) {
	if itemsOrdered <= 0 {
		return PackResult{}, fmt.Errorf("items ordered must be positive")
	}

	if len(availablePackSizes) == 0 {
		return PackResult{}, fmt.Errorf("no pack sizes available")
	}

	// Sort pack sizes in descending order
	sort.Slice(availablePackSizes, func(i, j int) bool {
		return availablePackSizes[i] > availablePackSizes[j]
	})

	smallestPack := availablePackSizes[len(availablePackSizes)-1]

	// First use greedy approach for the bulk of the order
	packCounts := make(map[uint64]uint64)
	remaining := itemsOrdered

	// Limit the DP size to a reasonable value
	dpLimit := smallestPack * 10

	// Use greedy approach for large portions
	if remaining > dpLimit {
		for i := 0; i < len(availablePackSizes)-1; i++ {
			packSize := availablePackSizes[i]
			count := remaining / packSize
			if count > 0 {
				packCounts[packSize] = count
				remaining -= count * packSize
			}
		}
	}

	// Use DP only for the remaining amount (which is smaller than dpLimit)
	if remaining > 0 {
		// DP table: dp[i] = minimum items to fulfill i items
		dp := make([]uint64, dpLimit+1)
		packChoice := make([]uint64, dpLimit+1)
		for i := uint64(1); i <= dpLimit; i++ {
			dp[i] = i + dpLimit // Initialize with a large value
		}
		dp[0] = 0

		// Fill the dp table
		for i := uint64(1); i <= dpLimit; i++ {
			for _, size := range availablePackSizes {
				if size <= i && dp[i-size]+size < dp[i] {
					dp[i] = dp[i-size] + size
					packChoice[i] = size
				} else if i < size && size < dp[i] {
					dp[i] = size
					packChoice[i] = size
				}
			}
		}

		// Find the best target that is >= remaining
		bestTarget := remaining
		for i := remaining; i <= dpLimit; i++ {
			if dp[i] < dp[bestTarget] || (dp[i] == dp[bestTarget] && i < bestTarget) {
				bestTarget = i
			}
		}

		// Reconstruct solution for the remaining amount
		current := bestTarget
		for current > 0 {
			size := packChoice[current]
			packCounts[size]++
			current -= size
		}
	}

	// Calculate total items and packs
	var totalItems uint64 = 0
	var totalPacks uint64 = 0
	for size, count := range packCounts {
		totalItems += size * count
		totalPacks += count
	}

	return PackResult{
		PackCounts:  packCounts,
		TotalPacks:  totalPacks,
		TotalItems:  totalItems,
		ExcessItems: totalItems - itemsOrdered,
	}, nil
}

// findMinimumItems uses logic where finds largest possible pack size
// and then fills the remaining items with smaller packs
// It returns the minimum number of items
func findMinimumItems(target uint64, packSizes []uint64) (uint64, map[uint64]uint64) {
	smallestPack := packSizes[len(packSizes)-1]

	// DP table: dp[i] = minimum number of items to fulfill i items
	// Initialize with a value larger than any possible solution
	// For each number of items, record which pack was used
	// Example default state:
	// dp:       [0, maxPossible, maxPossible, ..., maxPossible] (length = target + 1)

	maxPossible := target + smallestPack + 1
	dp := make([]uint64, target+1)
	for i := range dp {
		dp[i] = maxPossible
	}
	dp[0] = 0

	// For each number of items, record which pack was used
	// packUsed: [0, 0, 0, ..., 0] (length = target + 1)
	packUsed := make([]uint64, target+1)

	// Fill the dp table
	for i := uint64(1); i <= target; i++ { // Iterate through all possible item counts from 1 to the target
		for _, size := range packSizes { // Iterate through each available pack size
			// Check if the current pack size can be used to fulfill the current target (i)
			// and if using this pack results in fewer items than the current best solution
			if size <= i && dp[i-size] != maxPossible && dp[i-size]+size < dp[i] {
				dp[i] = dp[i-size] + size // Update dp[i] with the new minimum number of items
				packUsed[i] = size        // Record the pack size used to achieve this solution
			} else if i < size && size < dp[i] {
				// Special case: If the current target (i) is smaller than the pack size,
				// and using this pack directly results in fewer items than the current best solution
				dp[i] = size       // Use the larger pack directly
				packUsed[i] = size // Record the pack size used
			}
		}
	}

	// Reconstruct the solution
	packCounts := make(map[uint64]uint64)
	current := target

	// If we couldnt fulfill exactly, find the next possible fulfillment
	if dp[target] == maxPossible {
		for i := target + 1; ; i++ {
			if dp[i] != maxPossible {
				current = i
				break
			}
			// Special case: if even the smallest pack is too large
			if i == target+smallestPack {
				packCounts[smallestPack] = 1
				return smallestPack, packCounts
			}
		}
	}

	// Reconstruct which packs were used
	for current > 0 {
		pack := packUsed[current]
		packCounts[pack]++
		current -= pack

		// Handle the case where we need a pack larger than remaining items
		if current > 0 && packUsed[current] == 0 {
			// Find the smallest pack that can fulfill the remaining items
			for i := len(packSizes) - 1; i >= 0; i-- {
				if packSizes[i] >= current {
					packCounts[packSizes[i]]++
					current = 0
					break
				}
			}
		}
	}

	// Calculate total items
	var totalItems uint64 = 0
	for size, count := range packCounts {
		totalItems += size * count
	}

	return totalItems, packCounts
}
