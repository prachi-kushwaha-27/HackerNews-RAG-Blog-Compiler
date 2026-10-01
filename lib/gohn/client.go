package gohn

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	_domain         = "https://hacker-news.firebaseio.com"
	_itemUrlTmpl    = _domain + "/v0/item/%d.json"
	_userUrlTmpl    = _domain + "/v0/user/%s.json"
	_topStoriesUrl  = _domain + "/v0/topstories.json"
	_newStoriesUrl  = _domain + "/v0/newstories.json"
	_bestStoriesUrl = _domain + "/v0/beststories.json"
)

type Client interface {
	GetItem(id int64) (*Item, error)
	GetUser(id string) (*User, error)
	GetItemWithKids(id int64) (*ItemWithKids, error)
	GetTopStories() (*Stories, error)
	GetNewStories() (*Stories, error)
	GetBestStories() (*Stories, error)
}

func NewClient(
	timeout time.Duration,
	rps int,
) Client {
	return &client{
		httpClient: newRateControlledHTTPClient(timeout, rps),
	}
}

type client struct {
	httpClient rateControlledHTTPClient
}

func (c *client) GetItem(id int64) (*Item, error) {
	url := fmt.Sprintf(_itemUrlTmpl, id)
	resp, err := c.getSync(url)
	if err != nil {
		return nil, err
	}
	item := &Item{}
	err = json.NewDecoder(resp.Body).Decode(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *client) GetUser(id string) (*User, error) {
	url := fmt.Sprintf(_userUrlTmpl, id)
	resp, err := c.getSync(url)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *client) GetItemWithKids(id int64) (*ItemWithKids, error) {
	item, err := c.GetItem(id)
	if err != nil {
		return nil, err
	}
	result := &ItemWithKids{
		Item: *item,
		Kids: make([]*ItemWithKids, len(item.Kids)),
	}

	wg := sync.WaitGroup{}
	for i, kidId := range item.Kids {
		wg.Add(1)
		idx := i
		id := kidId
		go func() {
			defer wg.Done()
			item, err := c.GetItemWithKids(id)
			if err != nil {
				result.Kids[idx] = nil
			} else {
				result.Kids[idx] = item
			}
		}()
	}
	wg.Wait()

	return result, nil
}

func (c *client) GetTopStories() (*Stories, error) {
	return c.getStories(_topStoriesUrl)
}

func (c *client) GetNewStories() (*Stories, error) {
	return c.getStories(_newStoriesUrl)
}

func (c *client) GetBestStories() (*Stories, error) {
	return c.getStories(_bestStoriesUrl)
}

func (c *client) getSync(url string) (*http.Response, error) {
	resultChan := c.httpClient.Get(url)
	result := <-resultChan
	if result.err != nil {
		return nil, result.err
	}
	if result.resp.StatusCode != http.StatusOK {
		return nil, errors.New(result.resp.Status)
	}
	return result.resp, nil
}

func (c *client) getStories(url string) (*Stories, error) {
	resp, err := c.getSync(url)
	if err != nil {
		return nil, err
	}
	stories := &Stories{}
	err = json.NewDecoder(resp.Body).Decode(stories)
	if err != nil {
		return nil, err
	}
	return stories, nil
}
