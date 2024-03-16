package api

import (
	"context"
	"io"
)

type Api interface {
	Push(ctx context.Context, data []byte) (string, error)
	PushStream(ctx context.Context, stream *io.PipeReader) (string, error)
	Get(ctx context.Context, id string) ([]byte, error)
	GetStream(ctx context.Context, id string) (*io.PipeReader, error)
	Delete(ctx context.Context, id string) (bool, error)
}

// 如果要实现Push操作，mate leader节点必须获得所有在线data节点，使用一种分片分配算法将分片存储到各个data节点上
// apply函数实现日志持久化/操作持久化
// Push 操作 将文件分片->加密->存储
// 加密操作：每个data节点使用各自的公私钥，每个分片使用唯一的对称加密密钥，对称加密密钥使用非对称加密密钥保存
// Push 操作需要存储所有分片的地址和ID信息，地址确定机器，ID确定最终分片；一个分片多个副本，多个地址和ID信息
// Delete 操作 获得所有分片的地址和ID，执行删除操作，这里涉及到分布式事务(暂时不考虑事务)
// Get 操作 获得所有分片地址和ID，逐个获取，依次尝试->解密->合并
