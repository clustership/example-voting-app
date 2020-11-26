package main

import (
  "os"
  "fmt"
	"log"
	"net/http"
	"time"
  "encoding/json"

  "github.com/gorilla/mux"
  "github.com/rs/cors"

	"github.com/gomodule/redigo/redis"
)

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {
  redis_host := getEnv("REDIS_HOST", "localhost")
  redis_port := getEnv("REDIS_PORT", "6379")

	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", redis_host, redis_port))
		},
	}

	router := mux.NewRouter()

	router.HandleFunc("/votes", addVote).Methods("POST")

  handler := cors.Default().Handler(router)

  port := fmt.Sprintf(":%s", getEnv("PORT", "8080"))

	log.Printf("Listening on %s...", port)
  log.Fatal(http.ListenAndServe(port, handler))
}

func addVote(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  var vote Vote
  _ = json.NewDecoder(r.Body).Decode(vote)


  // xymox: check within [ 'a', 'b' ] values ?
  //
	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
	// if _, err := strconv.Atoi(id); err != nil {
	//	http.Error(w, http.StatusText(400), 400)
	//	return
	//}

	// Call the IncrementLikes() function passing in the user-provided
	// id. If there's no album found with that id, return a 404 Not
	// Found response. In the event of any other errors, return a 500
	// Internal Server Error response.
	err := IncrementVotes(vote.VoterId, vote.Vote)
	if err == ErrNoVote {
    log.Println(err)
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

  json.NewEncoder(w).Encode(&vote)

	// Redirect the client to the GET /album route, so they can see the
	// impact their like has had.
	// http.Redirect(w, r, "/album?id="+id, 303)
  // response.write
}
