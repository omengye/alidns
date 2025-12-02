package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	al "github.com/libdns/alidns"
	l "github.com/libdns/libdns"
)

func main() {
	accKeyID := strings.TrimSpace(os.Getenv("ACCESS_KEY_ID"))
	accKeySec := strings.TrimSpace(os.Getenv("ACCESS_KEY_SECRET"))
	if (accKeyID == "") || (accKeySec == "") {
		fmt.Printf("ERROR: %s\n", "ACCESS_KEY_ID or ACCESS_KEY_SECRET missing")
		return
	}
	EchValue := strings.TrimSpace(os.Getenv("ECH_VALUE"))

	zone := ""
	if len(os.Args) > 1 {
		zone = strings.TrimSpace(os.Args[1])
	}
	if zone == "" {
		fmt.Printf("ERROR: %s\n", "First arg zone missing")
		return
	}

	fmt.Printf("Get ACCESS_KEY_ID: %s,ACCESS_KEY_SECRET: %s,ZONE: %s\n", accKeyID, accKeySec, zone)
	provider := al.Provider{
		AccKeyID:     accKeyID,
		AccKeySecret: accKeySec,
	}
	records, err := provider.GetRecords(context.TODO(), zone)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	for _, record := range records {
		fmt.Printf("%s (.%s): %s, %s\n", record.RR().Name, zone, record.RR().Data, record.RR().Type)
	}

	fmt.Println("Press any Key to set the test ech record")
	fmt.Scanln()

	echRecords := make([]l.Record, 0)
	params := l.SvcParams{}
	params["ech"] = []string{EchValue}
	echRecords = append(echRecords, l.ServiceBinding{
		// HTTPS and SVCB RRs: RFC 9460 (https://www.rfc-editor.org/rfc/rfc9460)
		Scheme:   "https",
		Name:     "ech",
		TTL:      10 * time.Minute,
		Priority: 2, // allows a manual override with priority 1
		Target:   ".",
		Params:   params,
	})

	records, err = provider.SetRecords(context.TODO(), zone, echRecords)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	fmt.Println("Press any Key to delete the test entry")
	fmt.Scanln()
	fmt.Println("Deleting the entry")
	_, err = provider.DeleteRecords(context.TODO(), zone, echRecords)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

}
