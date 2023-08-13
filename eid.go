package main

import (
	"sync"

	"github.com/pkg/errors"
)

type Manager struct {
	generators map[string]*Generator
	storage    Storage
	mutex      sync.Mutex
}

func NewManager(s Storage) (*Manager, error) {
	m := &Manager{
		generators: map[string]*Generator{},
		storage:    s,
	}
	if s == nil {
		return m, nil
	}

	generators, err := s.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load stored sequence generators")
	}
	for k, v := range generators {
		m.generators[k] = &Generator{key: k, seq: v}
	}
	return m, nil

}

func (m *Manager) GetGenerator(key string) *Generator {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if g, ok := m.generators[key]; ok {
		return g
	}

	g := &Generator{key: key, seq: 0}
	m.generators[key] = g
	return g
}

func (m *Manager) Close() error {
	if m.storage == nil {
		return nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	generators := map[string]uint64{}
	for k, v := range m.generators {
		generators[k] = v.seq
	}

	err := m.storage.Save(generators)
	return errors.Wrap(err, "failed to save sequence generators")

}

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
