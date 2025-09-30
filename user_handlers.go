package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MichaelLaguatan/Gator/internal/database"
	"github.com/google/uuid"
)

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
	currentTime := time.Now()
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
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
