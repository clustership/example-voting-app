package main

import (
	"errors"
  "encoding/json"

  "log"

	"github.com/gomodule/redigo/redis"
)

// Declare a pool variable to hold the pool of Redis connections.
var pool *redis.Pool

var ErrNoVote = errors.New("no vote found")

// Define a custom struct to hold Vote data.
type Vote struct {
	Vote  string  `redis:"vote"`
  VoterId string `redis:"voter_id"`
}

func IncrementVotes(voter string, id string) error {
	conn := pool.Get()
	defer conn.Close()

	// Before we do anything else, check that an vote with the given
	// id exists. The EXISTS command returns 1 if a specific key exists
	// in the database, and 0 if it doesn't.
	// exists, err := redis.Int(conn.Do("EXISTS", "vote:"+id))
	// if err != nil {
	// 	return err
	// } else if exists == 0 {
	// 	return ErrNoVote
	// }

	// Use the MULTI command to inform Redis that we are starting a new
	// transaction. The conn.Send() method writes the command to the
	// connection's output buffer -- it doesn't actually send it to the
	// Redis server... despite it's name!
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	// Increment the number of likes in the vote hash by 1. Because it
	// follows a MULTI command, this HINCRBY command is NOT executed but
	// it is QUEUED as part of the transaction. We still need to check
	// the reply's Err field at this point in case there was a problem
	// queueing the command.
	// err = conn.Send("HINCRBY", "vote:"+id, "likes", 1)
	// if err != nil {
	// 	return err
	// }
	// And we do the same with the increment on our sorted set.
	// err = conn.Send("ZINCRBY", "likes", 1, id)
	// if err != nil {
	// 	return err
	// }

  vote := Vote{id, "zboubiblo"}

  b, err := json.Marshal(vote)
  if err != nil {
    return err
  }

  log.Printf("Pushing '%s' to redis\n", b)

	err = conn.Send("RPUSH", "votes", b)
	if err != nil {
		return err
	}

	// Execute both commands in our transaction together as an atomic
	// group. EXEC returns the replies from both commands but, because
	// we're not interested in either reply in this example, it
	// suffices to simply check for any errors. Note that calling the
	// conn.Do() method flushes the previous commands from the
	// connection output buffer and sends them to the Redis server.
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

func FindVote(id string) (*Vote, error) {
	// Use the connection pool's Get() method to fetch a single Redis
	// connection from the pool.
	conn := pool.Get()

	// Importantly, use defer and the connection's Close() method to
	// ensure that the connection is always returned to the pool before
	// FindVote() exits.
	defer conn.Close()

	// Fetch the details of a specific vote. If no vote is found
	// the given id, the []interface{} slice returned by redis.Values
	// will have a length of zero. So check for this and return an
	// ErrNoVote error as necessary.
	values, err := redis.Values(conn.Do("HGETALL", "vote:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, ErrNoVote
	}

	var vote Vote
	err = redis.ScanStruct(values, &vote)
	if err != nil {
		return nil, err
	}

	return &vote, nil
}
