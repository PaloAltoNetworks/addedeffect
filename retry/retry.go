package retry

import (
	"context"
	"time"
)

// Retry retries the given jobFunc until the context cancels.
// If a retry is needed and retryFunc is not nil, it will be called.
// retryFunc will get the error that caused the retry. If retryFunc
// retries an error, then the retry procedure stops, and the error is
// returned.
func Retry(
	ctx context.Context,
	jobFunc func() (interface{}, error),
	retryFunc func(error) error,
) (out interface{}, err error) {

	for {
		if out, err = jobFunc(); err == nil {
			return out, nil
		}

		if retryFunc != nil {
			if rerr := retryFunc(err); rerr != nil {
				return nil, rerr
			}
		}

		select {
		case <-time.After(3 * time.Second):
		case <-ctx.Done():
			return nil, err
		}
	}
}
