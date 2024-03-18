package id_generator

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"log"
	"testing"
)

func TestNewSnowflakeIdGenerator(t *testing.T) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatal(err)
	}
	id := node.Generate()
	fmt.Println(id.Base64())
	fmt.Println(id.Base32())
	fmt.Println(id.Int64())
	fmt.Println(id.String())
}
