package utils_test

import (
	"testing"

	"github.com/xuhe2/go-IDM/utils"
)

func TestFileDownloader(t *testing.T) {
	// TODO: Implement test for FileDownloader
	fileName := utils.GetFileNameFromUrl("https://example.com/file.txt")
	if fileName != "file.txt" {
		t.Errorf("Expected file name to be 'file.txt', but got '%s'", fileName)
	}
}
