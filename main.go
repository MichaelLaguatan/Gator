package main

import (
	"fmt"
	"os"

	"github.com/MichaelLaguatan/Gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if function, ok := c.cmds[cmd.name]; ok {
		err := function(s, cmd)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("command not found")
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.cmds[name] = f
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("unable to read config file: %v\n", err)
		os.Exit(1)
	}
	appState := state{&cfg}
	appCommands := commands{map[string]func(*state, command) error{}}
	appCommands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Print("error: wrong amount of arguments provided\n")
		os.Exit(1)
	}
	cmd := command{args[1], args[2:]}
	err = appCommands.run(&appState, cmd)
	if err != nil {
		fmt.Printf("error running login command: %v\n", err)
		os.Exit(1)
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("unable to read config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n%v\n", cfg.DbURL, cfg.CurrentUserName)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required")
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error handling login: %w", err)
	}
	fmt.Printf("User has been set to: %v\n", cmd.args[0])
	return nil
}
