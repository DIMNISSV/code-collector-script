package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// main - основная функция, которая запускает программу.
func main() {
	// RootDir - директория, которую нужно сканировать. По умолчанию - текущая (.).
	rootDir := "."
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	fmt.Printf("Сканирование и агрегация файлов в директории: %s\n\n", rootDir)

	if err := walkAndAggregate(rootDir); err != nil {
		fmt.Printf("Фатальная ошибка при сканировании: %v\n", err)
		os.Exit(1)
	}
}

// walkAndAggregate рекурсивно обходит директорию и выводит содержимое файлов.
func walkAndAggregate(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Выводим ошибку, но продолжаем обход (например, при ошибках доступа).
			fmt.Printf("Ошибка при доступе к пути %s: %v\n", path, err)
			return nil 
		}

		if isIgnored(path, d) {
			if d.IsDir() {
				// Пропускаем всю игнорируемую директорию.
				return filepath.SkipDir
			}
			// Пропускаем игнорируемый файл.
			return nil 
		}

		// Обрабатываем только обычные файлы.
		if d.Type().IsRegular() {
			return processFile(path)
		}

		return nil
	})
}

// isIgnored проверяет, нужно ли игнорировать файл или директорию.
// Список игнорируемых путей соответствует вашему 'tree' и стандартным исключениям.
func isIgnored(path string, d fs.DirEntry) bool {
	// Игнорируем сам скрипт, чтобы не включать его в вывод.
	if path == "file_aggregator.go" || path == "./file_aggregator.go" {
		return true
	}

	// Игнорируемые имена файлов/директорий (включая служебные)
	name := d.Name()
	if d.IsDir() {
		// Игнорируемые директории.
		if name == "node_modules" || name == ".git" || name == ".idea" || name == "vendor" {
			return true
		}
	} else {
		// Игнорируемые расширения/имена файлов.
		if name == "pnpm-lock.yaml" || strings.HasSuffix(name, ".sum") || strings.HasSuffix(name, ".lock") || strings.HasSuffix(name, ".o") || strings.HasSuffix(name, ".exe") {
			return true
		}
	}

	return false
}

// processFile считывает и выводит содержимое файла в требуемом формате.
func processFile(filePath string) error {
	content, err := os.ReadFile(filePath)

	// Начинаем вывод в требуемом формате: {ПУТЬ}:
	fmt.Printf("%s:\n", filePath)
	fmt.Println("```")

	if err != nil {
		// В случае ошибки чтения, выводим сообщение об ошибке.
		fmt.Printf("!!! ОШИБКА ЧТЕНИЯ ФАЙЛА !!!\n%v", err)
	} else {
		// Выводим содержимое файла.
		fmt.Print(string(content))
	}

	// Закрываем блок содержимого:
	fmt.Println("```")
	return nil
}
