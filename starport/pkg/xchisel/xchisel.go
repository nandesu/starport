package xchisel

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	chclient "github.com/jpillora/chisel/client"
	chserver "github.com/jpillora/chisel/server"
)

var DefaultServerPort = "7575"

func ServerAddr() string {
	return os.Getenv("CHISEL_ADDR")
}

func IsEnabled() bool {
	return ServerAddr() != ""
}

func StartServer(ctx context.Context, port string) error {
	s, err := chserver.NewServer(&chserver.Config{
		KeepAlive: time.Second * 3,
	})
	if err != nil {
		return err
	}
	if err := s.StartContext(ctx, "127.0.0.1", port); err != nil {
		log.Fatal(err)
	}
	if err = s.Wait(); err == context.Canceled {
		return nil
	}
	return err
}

func StartClient(ctx context.Context, serverAddr, localPort, remotePort string) error {
	c, err := chclient.NewClient(&chclient.Config{
		MaxRetryInterval: time.Second * 3,
		KeepAlive:        time.Second * 3,
		MaxRetryCount:    -1,
		Server:           serverAddr,
		Remotes:          []string{fmt.Sprintf("127.0.0.1:%s:127.0.0.1:%s", localPort, remotePort)},
	})
	if err != nil {
		return err
	}
	if err := c.Start(ctx); err != nil {
		log.Fatal(err)
	}
	if err = c.Wait(); err == context.Canceled {
		return nil
	}
	return err
}
