package main


type Feed struct {
    Id          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Link        string    `json:"link"`
    Items       []Item    `json:"items"`
}

type Item struct {
    Title       string    `json:"title"`
    Link        string    `json:"link"`
    Updated     string    `json:"updated"`
}

type Feeds []Feed
