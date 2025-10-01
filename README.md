# Gator

**Gator** is a command-line RSS feed aggregator built in Go. It supports user accounts, feed management, aggregation, and browsing of posts from subscribed feeds.

## Table of Contents

- [Features](#features)  
- [Requirements](#requirements)  
- [Installation](#installation)
- [Setup](#setup)
- [Commands](#commands)  
---

## Features

- Register and manage multiple users  
- Add, list, follow, and unfollow RSS feeds  
- Periodic aggregation (fetching new posts)  
- Browse posts from followed feeds  
- Persistent configuration & session management  

## Requirements

- Go (version 1.XX or newer)  
- PostgreSQL
- (Optional) A migration tool like `goose` or similar  

## Installation

### 1. **Install via `go install`**

   ```bash
   go install github.com/MichaelLaguatan/Gator@latest
   ```

### 2. Build from source

   You will need Go installed locally to build the executable
   ```bash
   git clone https://github.com/MichaelLaguatan/Gator.git
   cd Gator
   ```

## Setup

### 1. Register yourself as a user

   ```bash
    ./gator register <username>
   ```

### 2. Set yourself as the current user

   ```bash
   ./gator login <username>
   ```

You are now able to use the rest of the commands to fetch and browse posts from RSS feeds.

## Commands

### Register
   
   Used for registering a new user 

   ```bash
   ./gator register <username>
   ```

### Login

   Used for switching between registered users

   ```bash
   ./gator login <username> 
   ```

### Addfeed

   Used for storing an RSS feed in the database. The current user automatically follows the newly added feed

   ```bash
   ./gator addfeed <feed_name> <feed_url>
   ```

### Follow

   Used for following an RSS feed that has already been added by another user.

   ```bash
   ./gator follow <feed_url>
   ```

### Agg

   Used for collecting posts from the RSS feeds that the current user is following. The time between requests dictates the amount of time between fetching posts from each feed.

   ```bash
   ./gator agg <time_between_requests>
   ```

### Browse

   Used for browsing posts from the RSS feeds that the current user is following, after they have been fetched using `agg`. An optional `limit` parameter can be used to limit the amount of posts outputted.

   ```bash
   ./gator browse [limit]
   ```