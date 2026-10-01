package gohn

const (
	TypeStory   = "story"
	TypeComment = "comment"
	TypePoll    = "poll"
	TypePollOpt = "pollopt"
	TypeJob     = "job"
)

// Documentation: https://github.com/HackerNews/API/blob/master/README.md
type Item struct {
	ID          int64   `json:"id"`
	Deleted     bool    `json:"deleted,omitempty"`
	Type        string  `json:"type,omitempty"`
	By          string  `json:"by,omitempty"`
	Time        int64   `json:"time,omitempty"`
	Text        string  `json:"text,omitempty"`
	Dead        bool    `json:"dead,omitempty"`
	Parent      int64   `json:"parent,omitempty"`
	Poll        int64   `json:"poll,omitempty"`
	Kids        []int64 `json:"kids,omitempty"`
	URL         string  `json:"url,omitempty"`
	Score       int64   `json:"score,omitempty"`
	Title       string  `json:"title,omitempty"`
	Parts       []int64 `json:"parts,omitempty"`
	Descendants int64   `json:"descendants,omitempty"`
}

type User struct {
	ID        string  `json:"id"`
	Created   int64   `json:"created,omitempty"`
	Karma     int64   `json:"karma,omitempty"`
	About     string  `json:"about,omitempty"`
	Submitted []int64 `json:"submitted,omitempty"`
}

type ItemWithKids struct {
	Item Item            `json:"item,omitempty"`
	Kids []*ItemWithKids `json:"kids,omitempty"`
}

type SimpleComment struct {
	Comment string           `json:"comment,omitempty"`
	Replies []*SimpleComment `json:"replies,omitempty"`
}

type SimpleStoryWithComments struct {
	ID       int64            `json:"id"`
	Score    int64            `json:"score,omitempty"`
	Time     int64            `json:"time,omitempty"`
	URL      string           `json:"url,omitempty"`
	Title    string           `json:"title,omitempty"`
	Text     string           `json:"text,omitempty"`
	Comments []*SimpleComment `json:"comments,omitempty"`
}

type Stories []int64
