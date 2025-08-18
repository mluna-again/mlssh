package main

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mluna-again/luna/luna"
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

func getLunaPet(name string) luna.LunaPet {
	switch name {
	case "Cat":
		return luna.CAT

	case "Turtle":
		return luna.TURTLE

	case "Bunny":
		return luna.BUNNY

	default:
		log.Warnf("unkown pet: %s", name)
		return luna.CAT
	}
}

func getLunaVariant(name string) luna.LunaVariant {
	return luna.LunaVariant(strings.ToLower(name))
}
