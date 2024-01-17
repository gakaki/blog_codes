package cspiration

import (
	"fmt"
	"testing"
)

func TestGet400Questions(t *testing.T) {
	dealPage()
	fmt.Println(len(questionMaps))
	fmt.Println(len(question400))
	fmt.Println(len(question250))
	fmt.Println(len(questionDS))
}
