package mempage

import (
	"github.com/rs/zerolog/log"
	"sort"
)

type Element struct {
	Sequence int64
	Data     []byte
}

type MemPage struct {
	nextSequence   int64
	elements       chan *Element
	buffered       []*Element
	bufferIsSorted bool
}

func New() *MemPage {
	return &MemPage{
		nextSequence: 1,
		buffered:     make([]*Element, 0),
		elements:     make(chan *Element),
	}
}

func (m *MemPage) Store(e *Element) {
	if m.nextSequence == 0 || m.nextSequence > e.Sequence {
		m.nextSequence = e.Sequence
	}
	m.elements <- e
	return
}

func (m *MemPage) Close() {
	close(m.elements)
}

func (m *MemPage) Write(out chan *Element) {
	for element := range m.elements {
		m.export(element, out)
	}

	m.sortBuffer()

	for _, element := range m.buffered {
		if element.Sequence == m.nextSequence {
			out <- element
			m.nextSequence++
			continue
		}
		log.Error().Msgf("missing sequence %d", element.Sequence)
	}
	close(out) // TODO check if there is any item left
	m.buffered = nil
}

func (m *MemPage) export(element *Element, out chan *Element) {
	if element.Sequence != m.nextSequence {
		m.buffered = append(m.buffered, element)
		m.bufferIsSorted = false
		return
	}

	out <- element
	m.nextSequence++

	if len(m.buffered) == 0 {
		return
	}

	m.sortBuffer()

	element = m.buffered[0]
	m.buffered = m.buffered[1:]
	m.export(element, out)
}

func (m *MemPage) sortBuffer() {
	if !m.bufferIsSorted && len(m.buffered) > 0 {
		sort.Slice(m.buffered, func(i, j int) bool {
			return m.buffered[i].Sequence < m.buffered[j].Sequence
		})
	}
	m.bufferIsSorted = true
}
