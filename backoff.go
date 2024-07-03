package ws

import (
	"context"
	"time"
)

type Backoff struct {
	Initial time.Duration
	Max     time.Duration
	Factor  float64
}

func DefaultBackoff(ctx context.Context, reconnect func() error, backoff *Backoff) error {
	if backoff == nil {
		backoff = &Backoff{
			Initial: 1 * time.Second,
			Max:     30 * time.Second,
			Factor:  2.0,
		}
	}

	delay := backoff.Initial
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := reconnect(); err != nil {
			time.Sleep(delay)
			delay *= time.Duration(backoff.Factor)
			if delay > backoff.Max {
				delay = backoff.Max
			}
		} else {
			return nil
		}
	}
}
