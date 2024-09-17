package utils

import (
	"fmt"
	"path"
)

func BuildBookFolder(bookID uint64) string {
	return path.Join("book", fmt.Sprintf("%d", bookID))
}

func BuildBookDocumentFolder(bookID uint64) string {
	return path.Join(BuildBookFolder(bookID), "document")
}

func BuildBookDetectTableFolder(bookID uint64) string {
	return path.Join(BuildBookFolder(bookID), "detect_table")
}
