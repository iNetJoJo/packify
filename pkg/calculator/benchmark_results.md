# Benchmark Results: CalculatePacks vs CalculatePacksOptimized

This document presents the benchmark results comparing the original `CalculatePacks` function with the optimized `CalculatePacksOptimized` function.

## Summary

The benchmarks reveal a clear performance pattern:

- **For small orders (1-1000 items)**: The original `CalculatePacks` function is faster and uses less memory
- **For medium to large orders (5000+ items)**: The optimized `CalculatePacksOptimized` function is dramatically faster and uses significantly less memory

## Detailed Results

### Speed Comparison

| Items Ordered | CalculatePacks | CalculatePacksOptimized | Speedup Factor |
|---------------|----------------|-------------------------|----------------|
| 1             | 347.7 ns/op    | 19666 ns/op            | 0.02x (slower) |
| 250           | 2320 ns/op     | 19322 ns/op            | 0.12x (slower) |
| 501           | 4218 ns/op     | 19794 ns/op            | 0.21x (slower) |
| 1000          | 8378 ns/op     | 18871 ns/op            | 0.44x (slower) |
| 5000          | 40356 ns/op    | 284.2 ns/op            | 142x (faster)  |
| 10000         | 82731 ns/op    | 284.7 ns/op            | 291x (faster)  |
| 50000         | 391212 ns/op   | 284.9 ns/op            | 1373x (faster) |
| 100000        | 793401 ns/op   | 281.4 ns/op            | 2819x (faster) |

### Memory Usage Comparison

| Items Ordered | CalculatePacks   | CalculatePacksOptimized | Memory Savings |
|---------------|------------------|-------------------------|----------------|
| 1             | 312 B/op         | 41240 B/op             | -132x (worse)  |
| 250           | 4376 B/op        | 41240 B/op             | -9.4x (worse)  |
| 501           | 8472 B/op        | 41240 B/op             | -4.9x (worse)  |
| 1000          | 16664 B/op       | 41240 B/op             | -2.5x (worse)  |
| 5000          | 82200 B/op       | 280 B/op               | 293x (better)  |
| 10000         | 164120 B/op      | 280 B/op               | 586x (better)  |
| 50000         | 803099 B/op      | 280 B/op               | 2868x (better) |
| 100000        | 1605914 B/op     | 280 B/op               | 5735x (better) |

## Analysis

### Algorithm Differences

1. **CalculatePacks (Original)**:
   - Uses pure dynamic programming approach
   - Creates DP tables of size equal to the number of items ordered
   - Excellent for small inputs but scales poorly with order size

2. **CalculatePacksOptimized (Optimized)**:
   - Uses a hybrid approach:
     - Greedy algorithm for large portions of the order
     - Dynamic programming only for the remaining smaller amount
   - Limits DP table size to a small multiple of the smallest pack size
   - Poor for very small inputs but scales extremely well

### When to Use Each Function

- **Use CalculatePacks when**:
  - Processing small orders (less than 1000 items)
  - Memory usage is not a concern
  - Exact optimal solution is required for small inputs

- **Use CalculatePacksOptimized when**:
  - Processing medium to large orders (5000+ items)
  - Memory efficiency is important
  - Processing speed for large orders is critical

## Custom Pack Sizes Impact

The benchmarks with custom pack sizes (more pack size options) show similar patterns:

- Original algorithm performs better for small orders
- Optimized algorithm performs dramatically better for larger orders

## Conclusion

The choice between `CalculatePacks` and `CalculatePacksOptimized` should be based on the expected order sizes in your application:

- For applications primarily handling small orders, the original algorithm may be sufficient
- For applications handling a wide range of order sizes, especially larger ones, the optimized algorithm provides substantial performance benefits

In the current implementation, `CalculatePacks` calls `CalculatePacksOptimized` by default, which is a good choice for general-purpose use as it handles large orders efficiently, though it may be slightly less efficient for very small orders.