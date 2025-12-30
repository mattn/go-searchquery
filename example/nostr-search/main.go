package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-searchquery"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func main() {
	var relayURL string
	flag.StringVar(&relayURL, "relay", "wss://relay.damus.io", "Relay URL to connect to")
	flag.Parse()

	var filter nostr.Filter
	searchquery.NewParser(flag.Arg(0), searchquery.WithLookup(func(token string) string {
		if strings.HasPrefix(token, "#") {
			filter.Tags["t"] = append(filter.Tags["t"], token[1:])
			return ""
		}

		if strings.HasPrefix(token, "id:") {
			filter.IDs = append(filter.IDs, token[3:])
			return ""
		}

		if strings.HasPrefix(token, "kind:") {
			if n, err := strconv.ParseInt(token[5:], 10, 64); err == nil {
				filter.Kinds = append(filter.Kinds, int(n))
			}
			return ""
		}

		if strings.HasPrefix(token, "author:") {
			if _, v, err := nip19.Decode(token[7:]); err == nil {
				filter.Authors = append(filter.Authors, v.(string))
			} else {
				filter.Authors = append(filter.Authors, token[7:])
			}
			return ""
		}

		if strings.HasPrefix(token, "until:") {
			if n, err := strconv.ParseInt(token[6:], 10, 64); err == nil && n > 0 {
				var tmp nostr.Timestamp = nostr.Timestamp(n)
				filter.Until = &tmp
			} else if d, err := time.Parse("2006-01-02", token[6:]); err == nil {
				var tmp nostr.Timestamp = nostr.Timestamp(d.Unix())
				filter.Until = &tmp
			}
			return ""
		}

		if strings.HasPrefix(token, "since:") {
			if n, err := strconv.ParseInt(token[6:], 10, 64); err == nil && n > 0 {
				var tmp nostr.Timestamp = nostr.Timestamp(n)
				filter.Since = &tmp
			} else if d, err := time.Parse("2006-01-02", token[6:]); err == nil {
				var tmp nostr.Timestamp = nostr.Timestamp(d.Unix())
				filter.Since = &tmp
			}
			return ""
		}

		if strings.HasPrefix(token, "limit:") {
			if n, err := strconv.ParseInt(token[6:], 10, 64); err == nil && n > 0 {
				filter.Limit = int(n)
			}
			return ""
		}

		if filter.Search != "" {
			filter.Search += " "
		}
		filter.Search += token
		return token
	})).Parse()

	relay, err := nostr.RelayConnect(context.Background(), relayURL)
	if err != nil {
		log.Fatal(err)
	}
	evs, err := relay.QueryEvents(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	for ev := range evs {
		enc.Encode(ev)
	}
}
