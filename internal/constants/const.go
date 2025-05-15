// Package constants хранит общие констаны всего проекта.
package constants

import "os"

// KeyContext - Тип ключа для доступа к данным куки через контекст.
type KeyContext string

// Общие константы для internal пакетов.
const (
	// FileRWPerm - Константа прав доступа к файлу для чтения.
	// Права для файла: Владелец, группа и остальные: чтение и запись.
	FileRWPerm = os.FileMode(0o666)
	// UserIDContextKey - Имя ключа для доступа к данным куки через контекст.
	UserIDContextKey = KeyContext("UserID")
	// ShortIDLength - Длина сокращенной ссылки.
	ShortIDLength = 8
)
