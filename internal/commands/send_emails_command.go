package commands

import (
	"context"

	"app/pkg/console"
)

type SendEmailsCommand struct {
	console.BaseCommand
}

func NewSendEmailsCommand() *SendEmailsCommand {
	return &SendEmailsCommand{}
}

func (c *SendEmailsCommand) Configure(config *console.CommandConfig) {
	config.Name = "send:emails"
	config.Description = "Send queued emails"
	config.Usage = "send:emails [options]"

	config.Options = []console.Option{
		{
			Name:        "queue",
			Shortcut:    "q",
			Description: "Queue name",
			Required:    false,
		},
		{
			Name:        "limit",
			Shortcut:    "l",
			Description: "Number of emails to send",
			Required:    false,
		},
	}

	c.BaseCommand.Configure(config)
}

func (c *SendEmailsCommand) Handle(ctx context.Context) error {
	queue := c.GetOption("queue")
	if queue == "" {
		queue = "default"
	}

	limit := c.GetOption("limit")
	if limit == "" {
		limit = "50"
	}

	c.Info("Starting to send emails from queue: %s", queue)
	c.Info("Processing limit: %s", limit)

	// Here you would implement your email sending logic
	// For example:
	// 1. Connect to your queue system
	// 2. Fetch emails to be sent
	// 3. Send them through your email service
	// 4. Update their status in the queue

	c.Success("Emails sent successfully")
	return nil
}
