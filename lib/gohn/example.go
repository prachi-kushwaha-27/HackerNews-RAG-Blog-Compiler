package gohn

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	_storyId = 8863
	_userId  = "jl"
)

func Test() {
	client := NewClient(10*time.Second, 5)

	printErr := func(err error) {
		fmt.Printf("[ERROR] %s\n", err.Error())
	}

	fmt.Printf("Top Stories...\n")
	topStories, err := client.GetTopStories()
	if err != nil {
		printErr(err)
	} else {
		fmt.Printf("%v\n\n", *topStories)
	}

	fmt.Printf("New Stories...\n")
	newStories, err := client.GetNewStories()
	if err != nil {
		printErr(err)
	} else {
		fmt.Printf("%v\n\n", *newStories)
	}

	fmt.Printf("Best Stories...\n")
	bestStories, err := client.GetBestStories()
	if err != nil {
		printErr(err)
	} else {
		fmt.Printf("%v\n\n", *bestStories)
	}

	fmt.Printf("User...\n")
	user, err := client.GetUser(_userId)
	if err != nil {
		printErr(err)
	} else {
		bytes, _ := json.MarshalIndent(user, "", "  ")
		fmt.Printf("%s\n\n", string(bytes))
	}

	fmt.Printf("Story with comments...\n")
	story, err := client.GetItemWithKids(_storyId)
	if err != nil {
		printErr(err)
	} else {
		simpleStory, _ := ToSimpleStoryWithComments(story)
		bytes, _ := json.MarshalIndent(simpleStory, "", "  ")
		fmt.Printf("%s\n\n", string(bytes))
	}
}
