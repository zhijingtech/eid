package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Manager(t *testing.T) {
	m, err := NewManager(nil)
	assert.NoError(t, err)
	assert.Empty(t, m.generators)

	g := m.GetGenerator("a")
	assert.Equal(t, uint64(0), g.seq)
	id := g.NextID()
	assert.Equal(t, uint64(1), id)

	g2 := m.GetGenerator("a")
	assert.Equal(t, g, g2)
	assert.Equal(t, uint64(1), g2.seq)

	g3 := m.GetGenerator("b")
	assert.NotEmpty(t, g, g2)
	assert.Equal(t, uint64(0), g3.seq)

	err = m.Close()
	assert.NoError(t, err)
}

func TestGenerator_Storage(t *testing.T) {
	s := &mockStorage{}

	m, err := NewManager(s)
	assert.NoError(t, err)
	assert.Empty(t, m.generators)

	g := m.GetGenerator("a")
	assert.Equal(t, uint64(0), g.seq)
	g.NextID()
	g.NextID()
	assert.Equal(t, uint64(2), g.seq)

	err = m.Close()
	assert.NoError(t, err)

	m2, err := NewManager(s)
	assert.NoError(t, err)
	assert.NotEmpty(t, m2.generators)

	g2 := m2.GetGenerator("a")
	assert.Equal(t, uint64(2), g2.seq)
	g2.NextID()
	g2.NextID()
	assert.Equal(t, uint64(4), g2.seq)

	err = m2.Close()
	assert.NoError(t, err)
}

// 并发测试ID生成，确保没有数据竞争问题
func TestGenerator_Parallel(t *testing.T) {

	m, err := NewManager(nil)
	assert.NoError(t, err)
	defer m.Close()

	var w sync.WaitGroup
	w.Add(2)

	go func() {
		defer w.Done()
		g := m.GetGenerator("a")
		for i := 0; i < 100; i++ {
			g.NextID()
		}
	}()

	go func() {
		defer w.Done()
		g := m.GetGenerator("a")
		for i := 0; i < 100; i++ {
			g.NextID()
		}
	}()

	w.Wait()
	g := m.GetGenerator("a")
	assert.Equal(t, uint64(200), g.seq)
}

type mockStorage struct {
	generators map[string]uint64
}

func (s *mockStorage) Load() (map[string]uint64, error) {
	return s.generators, nil
}

func (s *mockStorage) Save(m map[string]uint64) error {
	s.generators = m
	return nil
}
