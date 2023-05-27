package service

import (
	"fmt"
	"os"
	"testing"
)

func TestService(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(wd)
}
