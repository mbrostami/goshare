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
	minSequence    int64
	elements       chan *Element
	buffered       []*Element
	bufferIsSorted bool
}

func New() *MemPage {
	return &MemPage{
		minSequence: 1,
		buffered:    make([]*Element, 0),
		elements:    make(chan *Element),
	}
}

func (m *MemPage) Store(e *Element) {
	if m.minSequence == 0 || m.minSequence > e.Sequence {
		m.minSequence = e.Sequence
	}
	m.elements <- e
	return
}

func (m *MemPage) Close() {
	close(m.elements)
}

func (m *MemPage) ReadChannel() chan *Element {
	elementsBySequenceOrder := make(chan *Element)

	go func() {
		for element := range m.elements {
		CheckElement:
			if element.Sequence == m.minSequence {
				elementsBySequenceOrder <- element
				m.minSequence++

				if len(m.buffered) == 0 {
					continue
				}

				if !m.bufferIsSorted {
					sort.Slice(m.buffered, func(i, j int) bool {
						return m.buffered[i].Sequence < m.buffered[j].Sequence
					})
				}

				m.bufferIsSorted = true

				element = m.buffered[0]
				m.buffered = m.buffered[1:]
				goto CheckElement
			}

			m.buffered = append(m.buffered, element)
			m.bufferIsSorted = false
		}

		// check remaining buffered elements
		if !m.bufferIsSorted && len(m.buffered) > 0 {
			sort.Slice(m.buffered, func(i, j int) bool {
				return m.buffered[i].Sequence < m.buffered[j].Sequence
			})
		}

		for _, element := range m.buffered {
			if element.Sequence == m.minSequence {
				elementsBySequenceOrder <- element
				m.minSequence++
				continue
			}
			log.Error().Msgf("missing sequence %d", element.Sequence)
		}

		m.buffered = make([]*Element, 0)
	}()

	return elementsBySequenceOrder
}
