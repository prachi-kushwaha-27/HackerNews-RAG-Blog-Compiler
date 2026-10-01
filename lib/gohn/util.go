package gohn

import (
	"errors"
	"log/slog"
)

func ToSimpleStoryWithComments(item *ItemWithKids) (*SimpleStoryWithComments, error) {
	if item == nil {
		return nil, errors.New("nil item provided to convert to simple story")
	}
	story := &SimpleStoryWithComments{}
	if item.Item.Type != TypeStory {
		return nil, errors.New("provided is not story")
	}
	if item.Item.Dead || item.Item.Deleted {
		return nil, errors.New("story is dead or deleted")
	}
	story.ID = item.Item.ID
	story.Score = item.Item.Score
	story.Time = item.Item.Time
	story.URL = item.Item.URL
	story.Title = item.Item.Title
	story.Text = item.Item.Text

	comments := make([]*SimpleComment, 0)
	for _, kid := range item.Kids {
		comment, err := ToSimpleComment(kid)
		if err != nil {
			slog.Debug("failed converting comment to simple comment", slog.Any("error", err))
		} else {
			comments = append(comments, comment)
		}
	}
	story.Comments = comments

	return story, nil
}

func ToSimpleComment(item *ItemWithKids) (*SimpleComment, error) {
	if item == nil {
		return nil, errors.New("nil item provided to convert to simple comment")
	}
	comment := &SimpleComment{}
	if item.Item.Type != TypeComment {
		return nil, errors.New("provided is not comment")
	}
	if item.Item.Dead || item.Item.Deleted {
		return nil, errors.New("comment is dead or deleted")
	}

	comment.Comment = item.Item.Text
	replies := make([]*SimpleComment, 0)
	for _, kid := range item.Kids {
		reply, err := ToSimpleComment(kid)
		if err == nil {
			replies = append(replies, reply)
		}
	}
	comment.Replies = replies

	return comment, nil
}
