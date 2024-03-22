package floodcontrol

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCacheCleaner(t *testing.T) {

	newFc := NewFloodControler()

	// Подготовка тестовых данных: два тестовых случая
	testCase1 := &UserRequest{
		AmountOfRequests: 1,
		lifetime:         time.Now(),
	}
	testCase2 := &UserRequest{
		AmountOfRequests: 5,
		lifetime:         time.Now().Add(time.Second * -5),
	}

	// Добавление тестовых случаев в кеш нового экземпляра FloodControler
	newFc.cache[1] = *testCase1
	newFc.cache[2] = *testCase2

	// Вызов функции CacheCleaner для нового кеша FloodControler
	newFc.CacheCleaner(newFc.cache)

	// Ожидаемые данные после вызова CacheCleaner и полученные
	expectedData := []int64{1}
	gotData := []int64{}

	// Заполнение gotData значениями ключей из кеша после очистки
	for key := range newFc.cache {
		gotData = append(gotData, key)
	}

	fmt.Println(expectedData)
	fmt.Println(gotData)

	// Сравнение ожидаемых и полученных данных
	if !reflect.DeepEqual(expectedData, gotData) {
		t.Errorf("Expected data: %v, Got data: %v", expectedData, gotData)
	}

}

func TestFloodControler_Check(t *testing.T) {
	// Создание экземпляра FloodControler
	fc := NewFloodControler()

	// Подготовка тестовых данных: выполнение серии запросов от пользователя с идентификатором 123
	for i := 0; i < 5; i++ {
		_, err := fc.Check(context.Background(), 123)
		if err != nil {
			t.Errorf("Error during flood control check: %v", err)
		}
	}

	// Тестирование на превышение лимита запросов в интервале времени от пользователя с идентификатором 666
	result, err := fc.Check(context.Background(), 666)
	if result != true {
		t.Errorf("Error during time exceeding %v", err)
	}

	// Ожидание в течение 4 секунд
	time.Sleep(time.Second * 4)

	// Повторная проверка превышения лимита запросов от пользователя с идентификатором 666 после истечения времени интервала
	result, err = fc.Check(context.Background(), 666)
	if result != true {
		t.Errorf("Error during time exceeding %v", err)
	}

	// Создание тестового случая: пользователя с идентификатором 999, у которого уже есть 50 запросов и который был создан в текущий момент
	testCase := &UserRequest{
		AmountOfRequests: 50,
		lifetime:         time.Now(),
	}
	fc.cache[999] = *testCase

	// Проверка метода Check для пользователя с идентификатором 999
	_, err = fc.Check(context.Background(), 999)

	// Проверка того, что метод Check возвращает ошибку для пользователя с идентификатором 999, у которого превышен лимит запросов
	if err == nil {
		t.Errorf("The handler is not working correctly")
	} else {
		fmt.Println(err)
	}
}
