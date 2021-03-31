package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticServerFileSystem(t *testing.T) {
	host := "0.0.0.0"
	port := "8080"
	assert.Equal(t, fmt.Sprintf("%s:%s", host, port), "0.0.0.0:8080")
}
