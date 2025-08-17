package main

import (
	"context"
	"net"
	"time"
)

func aLittleBit() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	return ctx, cancel
}

func removePort(ip string) string {
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return ip
	}
	return host
}
