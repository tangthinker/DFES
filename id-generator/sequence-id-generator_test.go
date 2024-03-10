package id_generator

import (
	"fmt"
	"testing"
)

func TestSequenceIdGenerator_Next(t *testing.T) {
	sg := NewSequenceIdGenerator("shanliao")
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
	fmt.Println(sg.Next())
}
