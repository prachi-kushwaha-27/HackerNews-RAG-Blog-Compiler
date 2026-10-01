package hackernews

type Post struct {
	ID       int64    `json:"id,omitempty"`
	Datetime string   `json:"datetime,omitempty"`
	URL      string   `json:"url"`
	Tags     []string `json:"tags,omitempty"`
}

type Posts struct {
	Posts []Post `json:"posts,omitempty"`
}

type Summary struct {
	ID             int64    `json:"id,omitempty"`
	Title          string   `json:"title,omitempty"`
	Summary        string   `json:"summary,omitempty"`
	DateSummarised string   `json:"dateSummarised,omitempty"`
	URL            string   `json:"url"`
	Tags           []string `json:"tags,omitempty"`
}

type PostDates struct {
	Dates []string `json:"dates,omitempty"`
}
