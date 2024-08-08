package utils_test

import (
	"net/http"
	"testing"

	"github.com/xuhe2/go-IDM/utils"
)

func TestFileDownloader(t *testing.T) {
	// TODO: Implement test for FileDownloader
	request, _ := http.NewRequest("Header", "https://download.jetbrains.com/go/goland-2020.2.2.exe", nil)
	response, _ := http.DefaultClient.Do(request)
	fileName := utils.GetFileNameFromUrl(response)
	if fileName != "goland-2020.2.2.exe" {
		t.Errorf("Expected file name to be 'file.txt', but got '%s'", fileName)
	}
}
