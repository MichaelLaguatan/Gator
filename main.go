package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/MichaelLaguatan/Gator/internal/config"
	"github.com/MichaelLaguatan/Gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
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

func main() {
	appState, appCommands, err := setup()
	if err != nil {
		fmt.Printf("error occured during setup: %v", err)
		os.Exit(1)
	}
	args := os.Args
	cmd := command{args[1], args[2:]}
	err = appCommands.run(appState, cmd)
	if err != nil {
		fmt.Printf("error running %v command: %v\n", cmd.name, err)
		os.Exit(1)
	}
}

func setup() (*state, *commands, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read config file: %w", err)
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open connection to db: %w", err)
	}
	dbQueries := database.New(db)
	appState := state{dbQueries, &cfg}
	appCommands := commands{map[string]func(*state, command) error{
		"login":     handlerLogin,
		"register":  handlerRegister,
		"reset":     handlerReset,
		"users":     handlerUsers,
		"agg":       handlerAgg,
		"addfeed":   middlewareLoggedIn(handlerAddFeed),
		"feeds":     handlerFeeds,
		"follow":    middlewareLoggedIn(handlerFollow),
		"unfollow":  middlewareLoggedIn(handlerUnfollow),
		"following": middlewareLoggedIn(handlerFollowing),
	}}
	return &appState, &appCommands, nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error: user specified in config is not registered in the db: %w", err)
		}
		err = handler(s, c, user)
		if err != nil {
			return err
		}
		return nil
	}
}
