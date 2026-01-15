package main

type Memory[T comparable] struct {
	mem    map[T]int
	Length int
}

func NewMemory[T comparable](n int) *Memory[T] {
	return &Memory[T]{
		make(map[T]int, n),
		n,
	}
}

func (m *Memory[T]) Has(t T) bool {
	_, ok := m.mem[t]
	return ok
}

func (m *Memory[T]) Push(t T) {
	if m.Length <= 0 {
		return
	}

	already := m.Has(t)

	for k, v := range m.mem {
		if already && k == t {
			continue
		}
		if v == m.Length-1 {
			delete(m.mem, k)
			continue
		}
		m.mem[k] = v + 1
	}
	m.mem[t] = 0
}
