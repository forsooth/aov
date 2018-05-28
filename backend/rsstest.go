package main

import (
    "fmt"
    "github.com/mmcdole/gofeed"
)

func GenerateFeeds() {
        urls := []string{"https://www.parahumans.net/feed/",
                         "https://wanderinginn.com/feed/"}
        fp := gofeed.NewParser()
        for i := 0; i < len(urls); i++ {
                feed, _ := fp.ParseURL(urls[i])
                for _, item := range feed.Items[:3] {
                        fmt.Println(item.Title)
                }
        }
}
