package client

import (
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestWriteToFile(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	s := &Service{}
	resChan := make(chan *pb.ReceiveResponse)

	sequences := 1000
	wg := sync.WaitGroup{}
	go func() {
		for i := 0; i < sequences; i++ {
			wg.Add(1)
			go func(in int) {
				resChan <- &pb.ReceiveResponse{
					SequenceNumber: int64(in),
					Data:           []byte(strconv.Itoa(in) + "\n"),
				}
				wg.Done()
			}(i + 1)
		}
		wg.Wait()
		resChan <- &pb.ReceiveResponse{
			SequenceNumber: -1,
			Data:           []byte(strconv.Itoa(-1) + "\n"),
		}
	}()

	err := s.writeToFile("/tmp/bm-write-to-file", resChan)
	if err != nil {
		t.Errorf("couldn't write to file %+v", err)
	}

	data, err := os.ReadFile("/tmp/bm-write-to-file")
	if err != nil {
		t.Errorf("couldn't read the file %+v", err)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) != sequences+1 {
		t.Errorf("number of lines should be %d got %d", sequences+1, len(lines))
	}
	for i, line := range lines {
		if line != "" && line != strconv.Itoa(i+1) {
			t.Errorf("line %d is %s, should be %d", i, line, i+1)
		}
	}
}

func BenchmarkWriteToFile(b *testing.B) {
	//b.Run("benchmark writetofile", func(b *testing.B) {
	s := &Service{}
	resChan := make(chan *pb.ReceiveResponse)

	wg := sync.WaitGroup{}
	go func() {
		semaphore := make(chan struct{}, 1000)
		for i := 0; i < b.N; i++ {
			semaphore <- struct{}{}
			wg.Add(1)
			go func(in int) {
				resChan <- &pb.ReceiveResponse{
					SequenceNumber: int64(in),
					Data:           []byte(strconv.Itoa(in) + "\n"),
				}
				<-semaphore
				wg.Done()
			}(i + 1)
		}
		wg.Wait()
		resChan <- &pb.ReceiveResponse{
			SequenceNumber: -1,
			Data:           []byte(strconv.Itoa(-1) + "\n"),
		}
	}()

	err := s.writeToFile("/tmp/bm-write-to-file", resChan)
	if err != nil {
		b.Errorf("couldn't write to file %+v", err)
	}
	//})
}
