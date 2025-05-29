package commands

import (
	"context"
	"time"

	"app/pkg/console"
)

type HelloWorldCommand struct {
	console.BaseCommand
}

func NewHelloWorldCommand() *HelloWorldCommand {
	return &HelloWorldCommand{}
}

func (c *HelloWorldCommand) Configure(config *console.CommandConfig) {
	config.Name = "hello:world"
	config.Description = "Print Hello World with timestamp"
	config.Usage = "hello:world"
	c.BaseCommand.Configure(config)
}

func (c *HelloWorldCommand) Handle(ctx context.Context) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	c.Line("Hello World! Current time: %s", now)
	return nil
}
