package main

import (
	"context"
	"time"
	"fmt"
	"math/rand"
	ksql "github.com/rmoff/ksqldb-go"
)

func doSetup(client *ksql.Client) {
	k := `CREATE STREAM riderlocations (profileId VARCHAR, latitude DOUBLE, longitude DOUBLE)
  WITH (kafka_topic='locations', value_format='json', partitions=1);`
  	err := client.Execute(k)
	if err != nil {
		panic(fmt.Sprintf("unable to create riderlocations: %v", err))
	}
}

func doLoop(client *ksql.Client, id int) {
        k := "INSERT INTO riderlocations (profileId, latitude, longitude) VALUES ('%d',%f,%f);"

        var lat float64
        var lon float64
        for {
                lat += rand.Float64()
                lon += rand.Float64()
		err := client.Execute(
                        fmt.Sprintf(k, id, lat,lon),
                )
                if err != nil {
			fmt.Printf("err: %v",err)
                }
                time.Sleep(time.Duration(rand.Int()%100)*time.Millisecond)
                time.Sleep(1*time.Second)
        }
}

func doQuery(client *ksql.Client) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	k := "SELECT * from riderlocations;"
	_, r, err := client.Pull(ctx, k, false)

	if err != nil {
		// handle the error better here, e.g. check for no rows returned
		return fmt.Errorf("Error running Pull request against ksqlDB:\n%v", err)
	}

	for _, row := range r {
		fmt.Printf("got: [%s,%f,%f]\n", row[0],row[1],row[2])
	}
	return nil
}

// Run this inside of the docket network,
// so that we don't have hostname issues with kafka/zk
func main() {
	ksqlDBServer := "http://localhost:8088"
	client := ksql.NewClient(ksqlDBServer, "", "").Debug()
	doSetup(client)

	go doLoop(client, rand.Int())
	go doLoop(client, rand.Int())

	time.Sleep(10 * time.Second)
	err := doQuery(client)
	if err != nil {
		panic(fmt.Sprintf("unable to query ksqlDB: %v", err))
	}
}
