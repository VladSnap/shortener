package constants

import "os"

const (
	FileRWPerm = os.FileMode(0o666) // Права для файла: Владелец, группа и остальные: чтение и запись
)
