package constants

import "os"

type KeyContext string

const (
	FileRWPerm       = os.FileMode(0o666) // Права для файла: Владелец, группа и остальные: чтение и запись
	UserIDContextKey = KeyContext("UserID")
)
