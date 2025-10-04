package main

import (
	"fmt"
	"time"

	"github.com/WhiCu/async"
)

func main() {
	fmt.Println("=== Async Goroutine Examples ===\n")

	// Example 1: Basic SafeGo usage
	fmt.Println("1. Basic SafeGo:")
	g1 := async.SafeGo(func() {
		fmt.Println("   Hello from goroutine!")
		time.Sleep(100 * time.Millisecond)
	})
	err1 := g1.Wait()
	if err1 != nil {
		fmt.Printf("   Error: %v\n", err1)
	} else {
		fmt.Println("   Completed successfully")
	}

	// Example 2: SafeGo with panic handling
	fmt.Println("\n2. SafeGo with panic:")
	g2 := async.SafeGo(func() {
		fmt.Println("   About to panic...")
		panic("intentional panic!")
	})
	err2 := g2.Wait()
	if err2 != nil {
		fmt.Printf("   Caught panic error: %v\n", err2)
		fmt.Printf("   Original panic value: %v\n", g2.Panic())
	}

	// Example 3: SafeGo with timeout
	fmt.Println("\n3. SafeGo with timeout:")
	g3 := async.SafeGoWithTimeout(50*time.Millisecond, func() {
		fmt.Println("   This will timeout...")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("   This won't be reached")
	})
	err3 := g3.Wait()
	if err3 != nil {
		fmt.Printf("   Timeout error: %v\n", err3)
	}

	// Example 4: Go with return value
	fmt.Println("\n4. Go with return value:")
	result := <-async.Go(func() int {
		time.Sleep(50 * time.Millisecond)
		return 42
	})
	fmt.Printf("   Result: %v, Error: %v, Panic: %v\n",
		result.Value, result.Error, result.Panic)

	// Example 5: GoErr with error handling
	fmt.Println("\n5. GoErr with error:")
	result2 := <-async.GoErr(func() (int, error) {
		time.Sleep(50 * time.Millisecond)
		return 0, fmt.Errorf("simulated error")
	})
	fmt.Printf("   Result: %v, Error: %v, Panic: %v\n",
		result2.Value, result2.Error, result2.Panic)

	// Example 6: Parallel execution
	fmt.Println("\n6. Parallel execution:")
	errors := async.Parallel(
		func() {
			fmt.Println("   Task 1 executing...")
			time.Sleep(100 * time.Millisecond)
			fmt.Println("   Task 1 completed")
		},
		func() {
			fmt.Println("   Task 2 executing...")
			time.Sleep(150 * time.Millisecond)
			fmt.Println("   Task 2 completed")
		},
		func() {
			fmt.Println("   Task 3 executing...")
			time.Sleep(80 * time.Millisecond)
			fmt.Println("   Task 3 completed")
		},
	)
	fmt.Printf("   Parallel tasks completed. Errors: %v\n", errors)

	// Example 7: Pool usage
	fmt.Println("\n7. Pool usage:")
	pool := async.NewPool(2)
	for i := 0; i < 5; i++ {
		taskNum := i
		pool.Submit(func() {
			fmt.Printf("   Pool task %d executing...\n", taskNum)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("   Pool task %d completed\n", taskNum)
		})
	}
	pool.Wait()
	pool.Close()
	fmt.Println("   All pool tasks completed")

	// Example 8: Retry with backoff
	fmt.Println("\n8. Retry with backoff:")
	attempt := 0
	value, err := async.RetryWithBackoff(3, 100*time.Millisecond, func() (int, error) {
		attempt++
		fmt.Printf("   Attempt %d...\n", attempt)
		if attempt < 3 {
			return 0, fmt.Errorf("attempt %d failed", attempt)
		}
		return 42, nil
	})
	fmt.Printf("   Final result: %v, Error: %v\n", value, err)

	fmt.Println("\n=== All examples completed ===")
}
