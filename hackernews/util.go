package hackernews

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func currentDateTime() string {
	return time.Now().Format(time.RFC3339)
}

func dateFromDatetime(date string) string {
	return date[:10]
}

func getPostsFilepath(date string, cfg config) string {
	return filepath.Join(cfg.PostsDir, date+".json")
}

func getSummaryFilepath(id int64, cfg config) string {
	return filepath.Join(cfg.PostSummaryDir, strconv.FormatInt(id, 10)+".md")
}

func getSummaryThinkFilepath(id int64, cfg config) string {
	return filepath.Join(cfg.PostSummaryDir, strconv.FormatInt(id, 10)+"-think.md")
}

func listAvailablePostsDates(cfg config) ([]string, error) {
	entries, err := os.ReadDir(cfg.PostsDir)
	if err != nil {
		return nil, err
	}
	dates := []string{}
	for _, entry := range entries {
		dates = append(dates, strings.Split(entry.Name(), ".")[0])
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	return dates, nil
}

func readPostSummary(id int64, cfg config) ([]byte, error) {
	postFile := getSummaryFilepath(id, cfg)
	return os.ReadFile(postFile)
}

func readPosts(date string, cfg config) (*Posts, error) {
	postsFilePath := getPostsFilepath(date, cfg)
	bytes, err := os.ReadFile(postsFilePath)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to read posts file: %s", postsFilePath),
			err)
	}

	posts := Posts{}
	err = json.Unmarshal(bytes, &posts)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to unmarshal posts file record into json: %s", postsFilePath),
			err)
	}
	return &posts, nil
}

// TODO: this is currently insert, make it upsert
func upsertPost(id int64, tags []string, url string, cfg config) error {
	currentDateTime := currentDateTime()
	post := Post{
		ID:       id,
		Datetime: currentDateTime,
		URL:      url,
		Tags:     tags,
	}
	date := dateFromDatetime(currentDateTime)
	postsFileRecords, err := readPosts(date, cfg)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		postsFileRecords = &Posts{
			Posts: make([]Post, 0),
		}
	} else if err != nil {
		return err
	}

	postsFileRecords.Posts = append(postsFileRecords.Posts, post)
	jsonData, err := json.MarshalIndent(postsFileRecords, "", "  ")
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed to mashal posts file record to json: %v", postsFileRecords),
			err)
	}

	err = os.WriteFile(getPostsFilepath(date, cfg), jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func extractPostTitle(summary string) string {
	splits := strings.Split(summary, "\n")
	if len(splits) == 0 {
		return ""
	}
	return strings.ReplaceAll(
		strings.ReplaceAll(splits[0], "#", ""),
		"**",
		"",
	)
}

func cleanSummary(summary string) string {
	splits := strings.Split(summary, "\n")
	if len(splits) > 1 && splits[0] == "```markdown" {
		splits = splits[1:]
		if len(splits) > 1 && splits[len(splits)-1] == "```" {
			splits = splits[:len(splits)-1]
		}
	}
	return strings.Join(splits, "\n")
}
