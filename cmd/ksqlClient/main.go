package main

import (
	"context"
	"time"
	"fmt"
	ksql "github.com/rmoff/ksqldb-go"
)

func doQuery() error {
	client := ksql.NewClient("http://localhost:8088", "", "").Debug()
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	k := "SELECT * from riderLocations;"
	_, r, err := client.Pull(ctx, k, false)

	if err != nil {
		// handle the error better here, e.g. check for no rows returned
		return fmt.Errorf("Error running Pull request against ksqlDB:\n%v", err)
	}

	for _, row := range r {
		fmt.Printf("%v\n", row)
	}
	return nil
}

// Run this inside of the docket network,
// so that we don't have hostname issues with kafka/zk
func main() {
	err := doQuery()
	if err != nil {
		panic(fmt.Sprintf("unable to query ksqlDB: %v", err))
	}
}
