package dsovs

import (
"context"
"fmt"
"io"
"net/http"
"time"
)

type Client struct {
url  string
http *http.Client
}

func NewClient(url string) *Client {
return &Client{
url:  url,
http: &http.Client{Timeout: 30 * time.Second},
}
}

func (c *Client) Fetch(ctx context.Context) ([]byte, error) {
req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
if err != nil {
return nil, err
}
res, err := c.http.Do(req)
if err != nil {
return nil, err
}
defer res.Body.Close()

if res.StatusCode != http.StatusOK {
return nil, fmt.Errorf("unexpected status: %s", res.Status)
}

payload, err := io.ReadAll(res.Body)
if err != nil {
return nil, err
}
return payload, nil
}
