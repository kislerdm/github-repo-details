// Package to interact with the file system.
package fs

import "os"

// FileRead function to read a file from FS.
func FileRead(path string) []byte {
	o, _ := os.ReadFile(path)
	return o
}

// FileWrite writes an object to the file.
func FileWrite(data []byte, path string) error {
	return os.WriteFile(path, data, 0644)
}
