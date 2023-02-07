package mempage

import (
	"sort"
)

const MaxElements = 10

type Element struct {
	Sequence int64
	Data     []byte
}

type MemPage struct {
	maxMemKeys     int32
	maxChunkLength int64
	minSequence    int64
	chanel         chan *Element
	backoff        []*Element
	backoffSorted  bool
}

func New() *MemPage {
	return &MemPage{
		minSequence: 1,
		chanel:      make(chan *Element),
	}
}

// Store
func (m *MemPage) Store(e *Element) {
	if m.minSequence == 0 || m.minSequence > e.Sequence {
		m.minSequence = e.Sequence
	}
	m.chanel <- e
	return
}

func (m *MemPage) Close() {
	close(m.chanel)
}

func (m *MemPage) ReadChannel() chan *Element {
	elementsBySequenceOrder := make(chan *Element)

	go func() {
		for element := range m.chanel {
		CheckElement:
			if element.Sequence == m.minSequence {
				elementsBySequenceOrder <- element
				m.minSequence++

				if len(m.backoff) == 0 {
					continue
				}

				if !m.backoffSorted {
					sort.Slice(m.backoff, func(i, j int) bool {
						return m.backoff[i].Sequence < m.backoff[j].Sequence
					})
				}

				m.backoffSorted = true

				element = m.backoff[0]
				m.backoff = m.backoff[1:]
				goto CheckElement
			}

			m.backoff = append(m.backoff, element)
			m.backoffSorted = false
		}

		for _, element := range m.backoff {
		CheckBackoffElement:
			if element.Sequence == m.minSequence {
				elementsBySequenceOrder <- element
				m.minSequence++

				if len(m.backoff) == 0 {
					continue
				}

				if !m.backoffSorted {
					sort.Slice(m.backoff, func(i, j int) bool {
						return m.backoff[i].Sequence < m.backoff[j].Sequence
					})
				}

				m.backoffSorted = true

				element = m.backoff[0]
				m.backoff = m.backoff[1:]
				goto CheckBackoffElement
			}
		}
	}()

	return elementsBySequenceOrder
}
