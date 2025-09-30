package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MichaelLaguatan/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("wrong amount of arguments supplied")
	}
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user defined in config does not exist in db: %w", err)
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("provided feed URL does not exist in db: %w", err)
	}
	current_time := time.Now()
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed_follow row: %w", err)
	}
	fmt.Printf("Feed name: %v\nCurrent User: %v", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user defined in config does not exist in db: %w", err)
	}
	followedFeedNames, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.Name)
	if err != nil {
		return fmt.Errorf("error getting feeds followed by current user: %w", err)
	}
	fmt.Printf("Followed Feeds For Current User: %v\n", currentUser.Name)
	for _, feed := range followedFeedNames {
		fmt.Printf("* %v", feed)
	}
	return nil
}
