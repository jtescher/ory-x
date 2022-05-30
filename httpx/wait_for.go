package httpx

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	retry "github.com/avast/retry-go/v4"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// WaitForEndpoint waits for the endpoint to be available.
func WaitForEndpoint(ctx context.Context, endpoint string, opts ...retry.Option) error {
	return retry.Do(func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if gjson.GetBytes(body, "status").String() != "ok" {
			return errors.Errorf("status is not yet ok: %s", body)
		}

		return nil
	},
		append([]retry.Option{
			retry.DelayType(retry.BackOffDelay),
			retry.Delay(time.Second),
			retry.MaxDelay(time.Second * 2),
			retry.Attempts(20),
		}, opts...)...)
}
