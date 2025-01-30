package datepicker

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	config := &gridConfig{
		selectedDate: time.Now(),
	}
	p := New(config)
	fmt.Println(p.View())
}
