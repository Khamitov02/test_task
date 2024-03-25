package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"task/flood"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	fld := flood.NewFlood(3, 2)

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		wg.Add(1)
		go func(ind int) {
			defer wg.Done()
			userId := int64(ind)
			canMakeCall, err := fld.Check(ctx, userId)
			fmt.Println("user:", userId)
			fmt.Println("answer:", canMakeCall)
			fmt.Println("error:", err)
		}(rand.Intn(2) + 1)

	}

	// wait for all goroutines to finish.
	wg.Wait()

	// wait 15 seconds before checking the status for user 1 again.
	time.Sleep(15 * time.Second)

	// check if user 1 is allowed to perform a task.
	// Will return true, because it's been 15 seconds
	// and he can't find the function calls in the last 15 seconds.
	ch, err := fld.Check(ctx, 1)

	fmt.Println("=========")
	fmt.Println("result for user: ", 1, ch)
	fmt.Println("error for user: ", 1, err)
	fmt.Println("--------")
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
