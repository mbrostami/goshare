package client

import (
	"fmt"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"os"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkWriteToFile(b *testing.B) {
	b.Run("benchmark writetofile", func(b *testing.B) {
		s := &Service{}
		resChan := make(chan *pb.ReceiveResponse)
		file, err := os.Create("/tmp/bm-write-to-file")
		if err != nil {
			b.Errorf("couldn't create temp file %+v", err)
		}

		go func() {
			semaphore := make(chan struct{}, 100)
			wg := sync.WaitGroup{}
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
				}(i)
			}
			wg.Wait()
			close(resChan)
		}()

		err = s.writeToFile(file, resChan)
		file.Close()
		//os.Remove(file.Name())
		fmt.Printf("fileName: %s\n", file.Name())
		if err != nil {
			b.Errorf("couldn't write to file %+v", err)
		}
	})
}
