package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/mbrostami/goshare/pkg/tracer"
	"github.com/schollz/progressbar/v3"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
)

const MaxConcurrent = 20

func (c *Client) waitingForReceiver(ctx context.Context) {
	startTime := time.Now()
	waitFor := 120 * time.Second
	bar := progressbar.NewOptions(
		-1,
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetDescription("waiting for receiver..."),
		progressbar.OptionSpinnerType(14),
	)
	defer bar.Clear()

	for {
		if ctx.Err() != nil {
			return
		}

		if startTime.Add(waitFor).Before(time.Now()) {
			return
		}

		bar.Add(1)
		time.Sleep(100 * time.Millisecond)
		continue
	}
}

func (c *Client) ShareInit(ctx context.Context, uid uuid.UUID, fileName string, fileSize int64) error {
	ctx, span := tracer.NewSpan(ctx, "sender")
	defer span.End()

	tctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		c.waitingForReceiver(tctx)
		wg.Done()
	}()

	res, err := c.conn.ShareInit(ctx, &pb.ShareInitRequest{
		Identifier: uid.String(),
		FileName:   fileName,
		FileSize:   fileSize,
	})

	cancel()
	wg.Wait()

	if err != nil {
		return err
	}

	if res.Error != "" {
		return errors.New(res.Error)
	}

	return nil
}

func (c *Client) Share(ctx context.Context, uid uuid.UUID, chunks chan *pb.ShareRequest) error {
	ctx, span := tracer.NewSpan(ctx, "sender")
	defer span.End()

	stream, err := c.conn.Share(ctx)
	if err != nil {
		log.Error().Msgf("failed to stream file: %v", err)
		return err
	}

	for chunk := range chunks {
		log.Debug().Msgf("streaming chunk to server: %d", chunk.SequenceNumber)
		err = stream.Send(chunk)
		if err != nil {
			log.Error().Msgf("failed to send chunk: %v", err)
			return err
		}

		response, err := stream.Recv()

		if err == io.EOF {
			log.Debug().Msg("streaming response finished")
			return nil
		}

		if err != nil {
			return fmt.Errorf("failed to receive response chunk: %v", err)
		}

		log.Debug().Msgf("streaming response received %+v", response)
		if response.Error != "" {
			// TODO retry
			log.Error().Msgf("received response %s", response.Error)
			return fmt.Errorf("received response %s", response.Error)
		}
	}

	return nil
}
