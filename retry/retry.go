// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
