package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/MichaelLaguatan/Gator/internal/config"
	"github.com/MichaelLaguatan/Gator/internal/database"
	"github.com/MichaelLaguatan/Gator/internal/rss"
	"github.com/google/uuid"
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
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("unable to read config file: %v\n", err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Printf("unable to open connection to db: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	appState := state{dbQueries, &cfg}
	appCommands := commands{map[string]func(*state, command) error{
		"login":    handlerLogin,
		"register": handlerRegister,
		"reset":    handlerReset,
		"users":    handlerUsers,
		"agg":      handlerAgg,
		"addfeed":  handlerAddFeed,
		"feeds":    handlerFeeds,
	}}
	args := os.Args
	cmd := command{args[1], args[2:]}
	err = appCommands.run(&appState, cmd)
	if err != nil {
		fmt.Printf("error running %v command: %v\n", cmd.name, err)
		os.Exit(1)
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("unable to read config file: %v\n", err)
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required")
	}
	if _, ok := s.db.GetUser(context.Background(), cmd.args[0]); ok != nil {
		fmt.Printf("user with username %v doesn't exist\n", cmd.args[0])
		os.Exit(1)
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error handling login: %w", err)
	}
	fmt.Printf("User has been set to: %v\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required")
	}
	if _, ok := s.db.GetUser(context.Background(), cmd.args[0]); ok == nil {
		fmt.Printf("user with username %v already exists\n", cmd.args[0])
		os.Exit(1)
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		fmt.Printf("error creating user: %v\n", err)
		os.Exit(1)
	}
	s.config.SetUser(user.Name)
	fmt.Printf("user was created with data: %v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("error reseting db: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting usernames from db: %w", err)
	}
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user defined in config does not exist in db: %w", err)
	}
	fmt.Printf("Current Users:\n")
	for _, user := range users {
		if currentUser.Name == user {
			fmt.Printf("* %v (current)", user)
		} else {
			fmt.Printf("* %v", user)
		}
	}
	fmt.Print("\n")
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	fmt.Printf("%v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("wrong amount of arguments supplied")
	}
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user defined in config does not exist in db: %w", err)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to db: %w", err)
	}
	fmt.Printf("%v", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds info: %w", err)
	}
	fmt.Printf("Current Feeds:\n")
	for _, feed := range feeds {
		fmt.Printf("Name: %v\nCreated By: %v\nURL: %v\n\n", feed.Name_2, feed.Name, feed.Url)
	}
	return nil
}
