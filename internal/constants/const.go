package constants

import "os"

type KeyContext string

const (
	// FileRWPerm - Константа прав доступа к файлу для чтения.
	// Права для файла: Владелец, группа и остальные: чтение и запись.
	FileRWPerm = os.FileMode(0o666)
	// UserIDContextKey - Тип ключа для доступа к данным куки через контекст.
	UserIDContextKey = KeyContext("UserID")
	// ShortIDLength - Длина сокращенной ссылки.
	ShortIDLength = 8
)
