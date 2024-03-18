package id_generator

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

type SnowflakeIdGenerator struct {
	snowflakeNode *snowflake.Node
}

func NewSnowflakeIdGenerator(id int64) *SnowflakeIdGenerator {
	node, err := snowflake.NewNode(id)
	if err != nil {
		log.Println("init snowflake id generator err:", err)
		return nil
	}
	return &SnowflakeIdGenerator{
		snowflakeNode: node,
	}
}

func (s *SnowflakeIdGenerator) Next() string {
	id := s.snowflakeNode.Generate()
	return id.String()
}
