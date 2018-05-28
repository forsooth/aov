package main

import (
    "encoding/json"
    "fmt"
    "time"
    "strconv"
    "math/rand"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/mmcdole/gofeed"
)

var urls = []string{"https://www.parahumans.net/feed/",
                       "https://wanderinginn.com/feed/",
                       "http://themonstersknow.com/feed/",
                       "https://www.sublimetext.com/blog/feed",
                       "https://kimonote.com/@mildbyte/rss/",
                       "https://lwn.net/headlines/rss"}

var startTime = time.Now().UTC()
var started = false

func shuffle(slice []Feed) {
    for i := len(slice) - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        slice[i], slice[j] = slice[j], slice[i]
    }
}

var feeds Feeds

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func AllFeeds(w http.ResponseWriter, r *http.Request) {
    var elapsed = time.Since(startTime)
    if elapsed < 600 * time.Second && started {
    } else {
        started = true
        startTime = time.Now().UTC()
        fp := gofeed.NewParser()
        feeds = nil
        for i := 0; i < len(urls); i++ {
            feed, _ := fp.ParseURL(urls[i])
            var items []Item
            for _, item := range feed.Items[:3] {
                fmt.Println(item.Title)
                newItem := Item{Title: item.Title,
                                Link: item.Link,
                                Updated: item.Published}
                items = append(items, newItem)
            }
            newFeed := Feed{Id: i,
                            Title: feed.Title,
                            Description: feed.Description,
                            Link: feed.Link,
                            Items: items}
            feeds = append(feeds, newFeed)
        }
    }

    shuffle(feeds)
    w.Header().Set("Access-Control-Allow-Origin", "*")

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "    ")
    if err := encoder.Encode(feeds); err != nil {
        panic(err)
    }
}

func FeedByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var feedId int
    var err error
    if feedId, err = strconv.Atoi(vars["feedId"]); err != nil {
        panic(err)
    }

    w.Header().Set("Access-Control-Allow-Origin", "*")

    feed := feeds[feedId]
    if feed.Id >= 0 {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(feed); err != nil {
            panic(err)
        }
        return
    }

    // If we didn't find it, 404
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
        panic(err)
    }

}

