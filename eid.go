package eid

import (
	"sync"
)

var (
	generators map[string]*Generator = map[string]*Generator{}
	storage    Storage
	mutex      sync.Mutex
)

type Generator struct {
	key string
	seq uint64
	mut sync.Mutex
}

func (g *Generator) NextID() uint64 {
	g.mut.Lock()
	defer g.mut.Unlock()

	g.seq++
	return g.seq
}

func NextID(key string) uint64 {
	g := GetGenerator(key)
	return g.NextID()
}

func GetGenerator(key string) *Generator {
	mutex.Lock()
	defer mutex.Unlock()

	if g, ok := generators[key]; ok {
		return g
	}

	g := &Generator{key: key, seq: 0}
	generators[key] = g
	return g
}

// 使用方如果需要对序号生成器进行持久化，可以在程序启动时加载，程序退出时保存。
// 特别注意不要在获取序号生成器生成器之后加载，否则使用中的生成器可能被覆盖了。
func Load(s Storage) error {
	if s == nil {
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	storage = s
	loaded, err := storage.Load()
	if err != nil {
		return err
	}
	for k, v := range loaded {
		generators[k] = &Generator{key: k, seq: v}
	}
	return nil
}

func Save() error {
	mutex.Lock()
	defer mutex.Unlock()

	if storage == nil {
		return nil
	}

	saved := map[string]uint64{}
	for k, v := range generators {
		saved[k] = v.seq
	}

	return storage.Save(saved)
}
