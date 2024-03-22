package main

import (
	"context"
	"fmt"

	floodcontrol "task/internal/floodControl"
)

func main() {

	var fc FloodControl = floodcontrol.NewFloodControler()
	ctx := context.Background()
	userID := int64(77)

	for i := 0; i < 15; i++ {
		success, err := fc.Check(ctx, userID)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		if success {
			fmt.Println("Success")
		} else {
			fmt.Println("failed")
		}
	}

}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
