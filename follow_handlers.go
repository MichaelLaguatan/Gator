package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MichaelLaguatan/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("wrong amount of arguments supplied")
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("provided feed URL does not exist in db: %w", err)
	}
	currentTime := time.Now()
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed_follow row: %w", err)
	}
	fmt.Printf("Feed name: %v\nCurrent User: %v", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("wrong amount of arguments supplied")
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("provided feed URL does not exist in db: %w", err)
	}
	err = s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error registering unfollow in db: %w", err)
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followedFeedNames, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error getting feeds followed by current user: %w", err)
	}
	fmt.Printf("Followed Feeds For Current User: %v\n", user.Name)
	for _, feed := range followedFeedNames {
		fmt.Printf("* %v", feed)
	}
	return nil
}
