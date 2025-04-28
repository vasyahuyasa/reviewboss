package inbound

import (
	"errors"
	"flag"

	"github.com/vasyahuyasa/reviewboss/internal/ports"
)

type CLI struct {
	Store ports.MergeRequestSource
}

func (c *CLI) Run() error {
	id := flag.String("mr", "", "MR ID to process")
	flag.Parse()
	if *id == "" {
		return errors.New("--mr required")
	}
	mr, err := c.Store.Get(*id)
	if err != nil {
		return err
	}
	// push into processing pipeline...
	return nil
}
