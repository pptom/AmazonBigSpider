package core

import (
	"testing"
)

func TestSetAsinInvalid(t *testing.T) {
	url := "https://www.amazon.com/dp/0786634111"
	err := SetAsinInvalid(url)
	if err != nil {
		t.Error(err.Error())
	}
}
