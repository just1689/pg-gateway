package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/just1689/pg-gateway/client"
	"github.com/just1689/pg-gateway/query"
	"sync"
	"time"
)

var svr = "http://localhost:8080"

func main() {
	start := time.Now()
	testInsert()
	//testReadAsync()
	//testRead()
	fmt.Println(time.Since(start))
}

var userEntities = "users"

type user struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func testInsert() {
	totalRows := 5000
	totalClients := 12
	work := make(chan int)
	done := sync.WaitGroup{}
	done.Add(1)

	go func() {
		for i := 1; i <= totalRows; i++ {
			work <- i
		}
		close(work)
	}()

	for g := 1; g <= totalClients; g++ {
		go func() {
			for w := range work {
				u := user{
					ID:        uuid.New().String(),
					FirstName: "Justin",
					LastName:  "Tamblyn",
					Email:     uuid.New().String(),
					Password:  "some_hash",
				}
				if err := client.Insert(svr, userEntities, u); err != nil {
					panic(err)
				}
				if w == totalRows {
					done.Done()
				}
			}
		}()
	}
	done.Wait()

}

func testReadAsync() {
	c, err := client.GetEntityManyAsync(svr, query.Query{
		Entity: userEntities,
		Comparisons: []query.Comparison{
			query.Comparison{
				Field:      "email",
				Comparator: "gt",
				Value:      "5",
			},
		},
		Limit: 5000,
	})
	if err != nil {
		panic(err)
	}
	count := 0
	for _ = range c {
		count++
	}
	fmt.Println("rows", count)

}

func testRead() {
	c, err := client.GetEntityMany(svr, query.Query{
		Entity:      userEntities,
		Comparisons: []query.Comparison{},
		Limit:       5,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(c))

}
