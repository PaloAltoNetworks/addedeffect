package utils

import (
	"net/http"
	"time"
)

const retryNumber = 3

// RetryRequest retry to launch the given function each 3 seconds
func RetryRequest(f func() (*http.Response, error)) (resp *http.Response, err error) {
	for index := 0; index < retryNumber; index++ {

		resp, err = f()

		if err == nil {
			return resp, err
		}

		<-time.After(3 * time.Second)
	}

	return resp, err
}
