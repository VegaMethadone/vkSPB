package logwriter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func WriteIntoLogs(value ...interface{}) {
	// Получение текущего рабочего каталога
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error occured during Getwd:", err)
	}

	// Формирование пути к каталогу для лог-файлов
	newDirPath := filepath.Join(filepath.Dir(filepath.Dir(currentDir)), "logs")
	// Добавление имени файла к пути к каталогу
	newDirPath += "data.log"

	// Открытие или создание лог-файла для записи
	file, err := os.OpenFile(newDirPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Println("Did not create/find log file")
	}
	// Закрытие файла после завершения работы функции
	defer file.Close()
	// Настройка вывода логов в файл
	log.SetOutput(file)

	log.Println(fmt.Sprint(value...))
}
