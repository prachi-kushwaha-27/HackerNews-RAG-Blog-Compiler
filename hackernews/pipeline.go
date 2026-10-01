package hackernews

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/prachi-kushwaha-27/tldr/lib/crawler"
	"github.com/prachi-kushwaha-27/tldr/lib/gohn"
	"github.com/prachi-kushwaha-27/tldr/ollama"
	"go.uber.org/fx"
)

var pipelineModule = fx.Provide(newPipeline)

type pipelineParams struct {
	fx.In

	Cfg          config
	OllamaClient ollama.Controller
	HnClient     gohn.Client
}

func newPipeline(p pipelineParams) Pipeline {
	pipeline := &pipeline{
		output: make(chan *message, p.Cfg.StoryQueueSize),
	}
	pipeline.setup(p.Cfg, p.OllamaClient, p.HnClient)
	return pipeline
}

type Pipeline interface {
	Push(id int64)
}

type pipeline struct {
	output chan *message
}

func (p *pipeline) Push(id int64) {
	p.output <- &message{id: id}
}

func (p *pipeline) setup(
	cfg config,
	ollamaClient ollama.Controller,
	hnClient gohn.Client,
) {
	baseFn := baseFunction{
		cfg:          cfg,
		ollamaClient: ollamaClient,
		hnClient:     hnClient,
	}
	functions := []function{}

	filterOutput := make(chan *message, cfg.StoryQueueSize)
	functions = append(functions, &filter{
		baseFunction: baseFn,
		input:        p.output,
		output:       filterOutput,
	})

	fetchStoryOutput := make(chan *message, cfg.StoryQueueSize)
	functions = append(functions, &fetchStory{
		baseFunction: baseFn,
		input:        filterOutput,
		output:       fetchStoryOutput,
	})

	scrapPostOutput := make(chan *message, cfg.StoryQueueSize)
	functions = append(functions, &scrapPost{
		baseFunction: baseFn,
		input:        fetchStoryOutput,
		output:       scrapPostOutput,
	})

	summariseStoryOutput := make(chan *message, cfg.StoryQueueSize)
	functions = append(functions, &summariseStory{
		baseFunction: baseFn,
		input:        scrapPostOutput,
		output:       summariseStoryOutput,
	})

	tagStoryOutput := make(chan *message, cfg.StoryQueueSize)
	functions = append(functions, &tagStory{
		baseFunction: baseFn,
		input:        summariseStoryOutput,
		output:       tagStoryOutput,
	})

	functions = append(functions, &writeOutputToFile{
		baseFunction: baseFn,
		input:        tagStoryOutput,
	})

	for _, fn := range functions {
		go fn.Run()
	}
}

type function interface {
	Run()
}

type message struct {
	id                     int64
	simpleStoryWithComment *gohn.SimpleStoryWithComments
	postContentInMarkdown  string
	summary                string
	tags                   []string
}

type baseFunction struct {
	cfg          config
	ollamaClient ollama.Controller
	hnClient     gohn.Client
}

type filter struct {
	baseFunction

	input  <-chan *message
	output chan<- *message
}

func (f *filter) Run() {
	for m := range f.input {
		alreadyProcessed := false
		dates, err := listAvailablePostsDates(f.cfg)
		if err != nil {
			slog.Error("failed to list available posts dates", slog.Any("error", err))
			continue
		}
		for _, date := range dates {
			postsFileRecords, err := readPosts(date, f.cfg)
			if err != nil {
				slog.Error("failed to read posts", slog.String("date", date), slog.Any("error", err))
				continue
			}
			for _, post := range postsFileRecords.Posts {
				if m.id == post.ID {
					alreadyProcessed = true
				}
			}
		}
		if !alreadyProcessed {
			f.output <- m
		} else {
			slog.Info("post already processed", slog.Int64("id", m.id))
		}
	}
}

type fetchStory struct {
	baseFunction

	input  <-chan *message
	output chan<- *message
}

func (fs *fetchStory) Run() {
	for m := range fs.input {
		start := time.Now()
		story, err := fs.hnClient.GetItemWithKids(m.id)
		if err != nil {
			slog.Error("failed getting story with kids", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
		simpleStory, err := gohn.ToSimpleStoryWithComments(story)
		if err != nil {
			slog.Error("failed converting to simple story", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
		m.simpleStoryWithComment = simpleStory
		slog.Info("fetch story finished", slog.Int64("id", m.id),
			slog.String("time taken", time.Since(start).String()))

		fs.output <- m
	}
}

type scrapPost struct {
	baseFunction

	input  <-chan *message
	output chan<- *message
}

func (sp *scrapPost) Run() {
	for m := range sp.input {
		done := make(chan struct{})
		go func() {
			defer func() { done <- struct{}{} }()
			start := time.Now()

			pdfBytes, err := crawler.UrlToPdf(m.simpleStoryWithComment.URL)
			if err != nil {
				slog.Error("failed to read URL contents",
					slog.String("url", m.simpleStoryWithComment.URL), slog.Int64("id", m.id), slog.Any("error", err))
				return
			}
			markdown, err := crawler.PdfToMarkdown(sp.cfg.TempDir, pdfBytes)
			if err != nil && len(markdown) == 0 {
				slog.Error("failed to convert pdf to markdown", slog.Int64("id", m.id), slog.Any("error", err))
				return
			}
			m.postContentInMarkdown = markdown
			slog.Info("scrap post finished", slog.Int64("id", m.id),
				slog.String("time taken", time.Since(start).String()))

			sp.output <- m
		}()
		select {
		case <-done:
		case <-time.Tick(10 * time.Minute):
			slog.Warn("timeout reading url contents to pdf", slog.Int64("id", m.id), slog.String("url", m.simpleStoryWithComment.URL))
		}
	}
}

type summariseStory struct {
	baseFunction

	input  <-chan *message
	output chan<- *message
}

func (ss *summariseStory) Run() {
	for m := range ss.input {
		start := time.Now()
		simpleStoryBytes, err := json.MarshalIndent(m.simpleStoryWithComment, "", "  ")
		if err != nil {
			slog.Error("failed marshalling simple story", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
		story := "# Post\n " + m.postContentInMarkdown + " \n\n# Conversation:\n " + string(simpleStoryBytes)
		summary, err := ss.ollamaClient.Instruct(context.Background(), &ollama.InstructRequest{
			Content:     story,
			Instruction: ss.cfg.SummarisePrompt,
		})
		if err != nil {
			slog.Error("failed summarising the story", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
		m.summary = summary.Result
		slog.Info("summarise story finished", slog.Int64("id", m.id),
			slog.String("time taken", time.Since(start).String()))

		ss.output <- m
	}
}

type tagStory struct {
	baseFunction

	input  <-chan *message
	output chan<- *message
}

func (ts *tagStory) Run() {
	for m := range ts.input {
		start := time.Now()
		response, err := ts.ollamaClient.Instruct(context.Background(), &ollama.InstructRequest{
			Content:     m.summary,
			Instruction: ts.cfg.TagPrompt + "\nTags: \n" + strings.Join(ts.cfg.Tags, "\n"),
		})
		if err != nil {
			slog.Error("failed to generate story tags", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
		tagsSet := map[string]bool{}
		for _, line := range strings.Split(response.Result, "\n") {
			line = strings.TrimSpace(line)
			for _, tag := range ts.cfg.Tags {
				if strings.Contains(line, tag) {
					tagsSet[tag] = true
				}
			}
		}
		tags := []string{}
		for tag := range tagsSet {
			tags = append(tags, tag)
		}
		m.tags = tags
		slog.Info("tag story finished", slog.Int64("id", m.id),
			slog.String("time taken", time.Since(start).String()))

		ts.output <- m
	}
}

type writeOutputToFile struct {
	baseFunction

	input <-chan *message
}

func (w *writeOutputToFile) Run() {
	for m := range w.input {
		summaryFile := getSummaryFilepath(m.id, w.cfg)
		thinkFile := getSummaryThinkFilepath(m.id, w.cfg)
		thinking, summary := ollama.ParseThinkingAndReply(m.summary)

		err := os.WriteFile(summaryFile, []byte(summary), 0644)
		if err != nil {
			slog.Error("failed writing summary file", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
		err = os.WriteFile(thinkFile, []byte(thinking), 0644)
		if err != nil {
			slog.Error("failed writing thinking file", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
		err = upsertPost(m.id, m.tags, m.simpleStoryWithComment.URL, w.cfg)
		if err != nil {
			slog.Error("failed to update posts file", slog.Int64("id", m.id), slog.Any("error", err))
			continue
		}
	}
}
