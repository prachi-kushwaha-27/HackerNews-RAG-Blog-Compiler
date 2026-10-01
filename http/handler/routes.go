package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/prachi-kushwaha-27/tldr/hackernews"
	"github.com/prachi-kushwaha-27/tldr/http/config"
	"github.com/prachi-kushwaha-27/tldr/http/internal/util"
	"github.com/prachi-kushwaha-27/tldr/ollama"
)

func addRoutes(
	mux *http.ServeMux,
	hnController hackernews.Controller,
	cfg config.Config,
) {
	mux.Handle("/", http.FileServer(http.Dir(cfg.Server.HtmlDir)))
	mux.Handle("/healthz", handleHealthz())
	mux.Handle("/summary", handleGetSummary(hnController))
	mux.Handle("/posts", handleListPosts(hnController))
	mux.Handle("/available-post-dates", handleListAvailablePostDates(hnController))
}

func handleHealthz() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("कार्यं करोति"))
		},
	)
}

func handleSummarise(ollamaController ollama.Controller) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			summariseReq, err := util.DecodeJsonHttpRequest[ollama.InstructRequest](r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			summariseResp, err := ollamaController.Instruct(context.Background(), summariseReq)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			util.EncodeHttpResponseToJson(w, http.StatusOK, summariseResp)
		},
	)
}

func handleGetSummary(hnController hackernews.Controller) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			postIdStr := r.URL.Query().Get("id")
			if len(postIdStr) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("error: id not provided in request params"))
				return
			}
			date := r.URL.Query().Get("date")
			if len(date) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("error: date not provided in request params"))
				return
			}
			postId, err := strconv.ParseInt(postIdStr, 10, 0)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			summary, err := hnController.GetSummary(r.Context(), date, postId)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(err.Error()))
				return
			}
			util.EncodeHttpResponseToJson(w, http.StatusOK, summary)
		},
	)
}

func handleListPosts(hnController hackernews.Controller) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			date := r.URL.Query().Get("date")
			tags := []string{}
			if r.URL.Query().Get("tags") != "" {
				tags = strings.Split(r.URL.Query().Get("tags"), ",")
			}
			if len(date) == 0 {
				dates, err := hnController.ListAvailablePostsDates(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					w.Write([]byte(err.Error()))
					return
				}
				if len(dates.Dates) == 0 {
					util.EncodeHttpResponseToJson(w, http.StatusOK, map[string]any{})
					return
				}
				date = dates.Dates[0]
			}
			posts, err := hnController.ListPosts(r.Context(), date, tags)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(err.Error()))
				return
			}
			util.EncodeHttpResponseToJson(w, http.StatusOK, posts)
		},
	)
}

func handleListAvailablePostDates(hnController hackernews.Controller) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			dates, err := hnController.ListAvailablePostsDates(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(err.Error()))
				return
			}
			util.EncodeHttpResponseToJson(w, http.StatusOK, dates)
		},
	)
}
