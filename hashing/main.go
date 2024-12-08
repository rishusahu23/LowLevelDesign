package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"sort"
)

type Node struct {
	Id      string
	Address string
	Hash    uint32
}

type ConsistentHashing struct {
	Ring          []uint32
	NodeMap       map[uint32]Node
	VirtualFactor int
	HashFunc      func(data []byte) uint32
}

func NewConsistentHashing(v int, hash func(data []byte) uint32) *ConsistentHashing {
	return &ConsistentHashing{
		Ring:          []uint32{},
		NodeMap:       make(map[uint32]Node),
		VirtualFactor: v,
		HashFunc:      hash,
	}
}

func (ch *ConsistentHashing) AddNode(node Node) {
	for i := 0; i < ch.VirtualFactor; i++ {
		id := fmt.Sprintf("%s-%d", node.Id, i)
		hash := ch.HashFunc([]byte(id))

		ch.Ring = append(ch.Ring, hash)
		ch.NodeMap[hash] = node
	}

	sort.Slice(ch.Ring, func(i, j int) bool {
		return ch.Ring[i] < ch.Ring[j]
	})
}

func (ch *ConsistentHashing) Remove(nodeId string) {
	for i := 0; i < ch.VirtualFactor; i++ {
		id := fmt.Sprintf("%s-%d", nodeId, i)
		hash := ch.HashFunc([]byte(id))

		idx := sort.Search(len(ch.Ring), func(i int) bool {
			return ch.Ring[i] >= hash
		})

		if idx < len(ch.Ring) && ch.Ring[idx] == hash {
			ch.Ring = append(ch.Ring[:idx], ch.Ring[idx+1:]...)
		}

		delete(ch.NodeMap, hash)
	}
}

func (ch *ConsistentHashing) GetNode(key string) Node {
	hash := ch.HashFunc([]byte(key))
	idx := sort.Search(len(ch.Ring), func(i int) bool { return ch.Ring[i] >= hash })
	if idx == len(ch.Ring) {
		idx = 0 // Wrap around the ring
	}
	return ch.NodeMap[ch.Ring[idx]]
}

func SHA1hash(data []byte) uint32 {
	h := sha1.New()
	h.Write(data)
	sum := h.Sum(nil)
	return binary.BigEndian.Uint32(sum[:4])
}
