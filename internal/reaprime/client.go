package reaprime

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ClientOptions struct {
	ReconnectMin time.Duration
	ReconnectMax time.Duration
}

type Client struct {
	baseURL string
	store   *Store
	options ClientOptions
}

func NewClient(baseURL string, store *Store, options ClientOptions) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), store: store, options: options}
}

func (c *Client) Run(ctx context.Context) {
	streams := []struct {
		name  string
		path  string
		parse func([]byte) error
	}{
		{name: "machine", path: "/ws/v1/machine/snapshot", parse: func(b []byte) error {
			v, err := ParseMachine(b)
			if err == nil {
				c.store.SetMachine(v)
			}
			return err
		}},
		{name: "scale", path: "/ws/v1/scale/snapshot", parse: func(b []byte) error {
			v, err := ParseScale(b)
			if err == nil {
				c.store.SetScale(v)
			}
			return err
		}},
		{name: "shot_settings", path: "/ws/v1/machine/shotSettings", parse: func(b []byte) error {
			v, err := ParseShotSettings(b)
			if err == nil {
				c.store.SetShot(v)
			}
			return err
		}},
		{name: "water_levels", path: "/ws/v1/machine/waterLevels", parse: func(b []byte) error {
			v, err := ParseWater(b)
			if err == nil {
				c.store.SetWater(v)
			}
			return err
		}},
		{name: "devices", path: "/ws/v1/devices", parse: func(b []byte) error {
			v, err := ParseDevices(b)
			if err == nil {
				c.store.SetDevices(v)
			}
			return err
		}},
		{name: "display", path: "/ws/v1/display", parse: func(b []byte) error {
			v, err := ParseDisplay(b)
			if err == nil {
				c.store.SetDisplay(v)
			}
			return err
		}},
	}

	for _, stream := range streams {
		go c.runStream(ctx, stream.name, stream.path, stream.parse)
	}
}

func (c *Client) runStream(ctx context.Context, name, path string, parse func([]byte) error) {
	backoff := c.options.ReconnectMin
	for {
		if ctx.Err() != nil {
			return
		}
		err := c.readStream(ctx, name, path, parse)
		if err != nil && ctx.Err() == nil {
			c.store.StreamError(name)
			slog.Warn("Reaprime stream disconnected", "stream", name, "error", err)
		}

		jitter := time.Duration(rand.Int64N(int64(backoff / 2)))
		timer := time.NewTimer(backoff + jitter)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
		backoff *= 2
		if backoff > c.options.ReconnectMax {
			backoff = c.options.ReconnectMax
		}
	}
}

func (c *Client) readStream(ctx context.Context, name, path string, parse func([]byte) error) error {
	wsURL, err := c.wsURL(path)
	if err != nil {
		return err
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	c.store.SetStreamConnected(name, true)
	slog.Info("connected Reaprime stream", "stream", name, "url", wsURL)

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			c.store.SetStreamConnected(name, false)
			return err
		}
		if err := parse(payload); err != nil {
			c.store.StreamError(name)
			slog.Debug("failed to parse Reaprime stream message", "stream", name, "error", err)
		}
	}
}

func (c *Client) wsURL(path string) (string, error) {
	parsed, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	switch parsed.Scheme {
	case "https":
		parsed.Scheme = "wss"
	default:
		parsed.Scheme = "ws"
	}
	parsed.Path = path
	parsed.RawQuery = ""
	return parsed.String(), nil
}
