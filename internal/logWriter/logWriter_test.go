package logwriter

import (
	"testing"
)

func TestWriteIntoLogs(t *testing.T) {
	data := "Check log file"
	WriteIntoLogs(data)
}
