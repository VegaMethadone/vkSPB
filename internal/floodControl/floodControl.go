package floodcontrol

import (
	"context"
	"fmt"
	"time"

	"task/config"
	logwriter "task/internal/logWriter"
)

func NewFloodControler() *FloodControler {
	return &FloodControler{
		cache: make(map[int64]UserRequest),
	}
}

func (fc *FloodControler) Check(ctx context.Context, userID int64) (bool, error) {
	// Блокировка мьютекса перед выполнением операции
	fc.mu.Lock()
	defer fc.mu.Unlock()

	// Проверка, существует ли пользователь в кеше
	user, exists := fc.cache[userID]

	// Если пользователя нет в кеше, создаем новую запись
	if !exists {
		newUser := &UserRequest{
			AmountOfRequests: 1,
			lifetime:         time.Now(),
		}
		// Запись информации о пользователе в лог, если включена опция записи в лог
		if config.WriteLogs {
			logwriter.WriteIntoLogs(fmt.Sprintf("UsedID %d, time: %v, Requests: %d", userID, newUser.lifetime, newUser.AmountOfRequests))
		}
		// Добавление нового пользователя в кеш
		fc.cache[userID] = *newUser
		// Вызов метода CacheCleaner для очистки устаревших записей в кеше
		fc.CacheCleaner(fc.cache)
		return true, nil
	}

	// Проверка, не превышает ли количество запросов пользователя максимальное количество в заданном временном интервале
	if user.AmountOfRequests >= config.MaxRequestsPerTimeInterval && time.Since(user.lifetime) <= config.TimeInterval {
		// Если превышено, пользователь временно заблокирован
		if config.WriteLogs {
			logwriter.WriteIntoLogs(fmt.Sprintf("UsedID %d temporary banned at %v", userID, time.Now()))
		}
		// Удаление пользователя из кеша
		delete(fc.cache, userID)
		// Вызов метода CacheCleaner для очистки устаревших записей в кеше
		fc.CacheCleaner(fc.cache)
		// Возвращаем false и ошибку, указывающую на превышение количества запросов
		err := RequestExceededError(userID)
		return false, err
	}

	// Обновление данных пользователя
	newData := &UserRequest{
		AmountOfRequests: user.AmountOfRequests + 1,
		lifetime:         user.lifetime,
	}
	fc.cache[userID] = *newData
	// Запись информации о пользователе в лог, если включена опция записи в лог
	if config.WriteLogs {
		logwriter.WriteIntoLogs(fmt.Sprintf("UsedID %d, time: %v, Requests: %d", userID, newData.lifetime, newData.AmountOfRequests))
	}
	// Вызов метода CacheCleaner для очистки устаревших записей в кеше
	fc.CacheCleaner(fc.cache)
	return true, nil
}

func (fc *FloodControler) CacheCleaner(cache map[int64]UserRequest) {
	// Создание списка ключей для удаления
	keysToDelete := []int64{}

	// Проверка каждого элемента кеша на превышение времени жизни
	for key, value := range cache {
		if time.Since(value.lifetime) >= config.CacheCleanerDuration {
			keysToDelete = append(keysToDelete, key)
		}
	}

	// Если есть ключи для удаления, запись об этом в лог
	if len(keysToDelete) != 0 {
		logwriter.WriteIntoLogs(fmt.Sprint("Delited keys:", keysToDelete))
	}

	// Удаление устаревших записей из кеша
	for _, value := range keysToDelete {
		delete(cache, value)
		fmt.Println("Delited key is:", value)
	}
}
