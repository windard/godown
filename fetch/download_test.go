package fetch

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileLength(t *testing.T) {
	requestUrl := "http://httpbin.org/bytes/%d"

	var length1K int64 = 1024
	lengthResult, err := GetFileLength(fmt.Sprintf(requestUrl, length1K))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, length1K, lengthResult)

	var length10K int64 = 1024 * 10
	lengthResult, err = GetFileLength(fmt.Sprintf(requestUrl, length10K))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, length10K, lengthResult)

}
