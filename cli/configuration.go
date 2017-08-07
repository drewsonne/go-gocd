package cli

import (
	"context"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	GetConfigurationCommandName  = "get-configuration"
	GetConfigurationCommandUsage = "Get GoCD server config.xml"
)

// GetConfigurationAction gets a list of agents and return them.
func GetConfigurationAction(c *cli.Context) error {
	pgs, r, err := cliAgent().Configuration.Get(context.Background())
	if err != nil {
		return handleOutput(nil, r, "GetConfiguration", err)
	}

	return handleOutput(pgs, r, "GetConfiguration", err)
}

// GetConfigurationCommand handles the interaction between the cli flags and the action handler for delete-agents
func GetConfigurationCommand() *cli.Command {
	return &cli.Command{
		Name:     GetConfigurationCommandName,
		Usage:    GetConfigurationCommandUsage,
		Action:   GetConfigurationAction,
		Category: "Configuration",
	}
}
