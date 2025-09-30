package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MichaelLaguatan/Gator/internal/database"
	"github.com/MichaelLaguatan/Gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("wrong amount of arguments supplied")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", cmd.args[0])
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Printf("%v", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed: %w", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:          nextFeed.ID,
		LastFetched: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}
	fetchedFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	fmt.Printf("Fetched feed: %v", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf("Item Title: %v\n", item.Title)
	}
	fmt.Print("\n")
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("wrong amount of arguments supplied")
	}
	currentTime := time.Now()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to db: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error making current user follow added feed: %w", err)
	}
	fmt.Printf("%v\n", feed)
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
	fmt.Print("\n")
	return nil
}
