// Package helpers реализует вспомогательные функции которые не были определены в конкретный пакет.
package helpers

import (
	crypto "crypto/rand"
	"fmt"
	"math/big"
	"os"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStringRunes - Генерирует случайную строку из заданного набора символов заданной длины.
func RandStringRunes(n int) (string, error) {
	b := make([]rune, n)
	maxLetters := big.NewInt(int64(len(letters)))
	for i := range b {
		rndIndex, err := crypto.Int(crypto.Reader, maxLetters)
		if err != nil {
			return "", fmt.Errorf("failed create new letter: %w", err)
		}
		b[i] = letters[rndIndex.Int64()]
	}
	return string(b), nil
}

// DirectoryExists - Проверяет существует ли директория на диске.
func DirectoryExists(path string) (bool, error) {
	// Пытаемся получить информацию о пути.
	_, err := os.Stat(path)
	if err == nil {
		// Если ошибки нет, проверяем, что это именно директория.
		fileInfo, _ := os.Stat(path)
		return fileInfo.IsDir(), nil
	} else if os.IsNotExist(err) {
		// Если ошибка указывает на отсутствие пути.
		return false, nil
	}
	// Возвращаем любую другую ошибку.
	return false, fmt.Errorf("failed check directory exists: %w", err)
}
