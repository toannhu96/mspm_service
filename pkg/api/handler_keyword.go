package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/toannhu96/mspm_service/pkg/service"
	"net/http"
)

type KeywordHandler struct {
	*BaseHandler
	keywordService service.KeywordService
}

func NewKeywordHandler(
	baseHandler *BaseHandler,
	keywordService service.KeywordService,
) *KeywordHandler {
	return &KeywordHandler{
		BaseHandler:    baseHandler,
		keywordService: keywordService,
	}
}

func (h *KeywordHandler) Route() chi.Router {
	mux := chi.NewRouter()
	mux.HandleFunc("/", h.getKeywordPageTmpl)
	mux.HandleFunc("/healthcheck", h.healthCheckHandler)
	mux.Post("/keywords", h.searchKeywords)
	return mux
}

func (h *KeywordHandler) getKeywordPageTmpl(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "/resources/index.html")
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *KeywordHandler) searchKeywords(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		req SearchKeywordReq
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("invalid search keyword request, err=%v", err)
		h.BadRequest(w, r, err)
		return
	} else if req.Payload == "" {
		logrus.Errorf("empty payload")
		h.BadRequest(w, r, fmt.Errorf("empty payload"))
		return
	}

	data, err := h.keywordService.MultiTermMatch(ctx, req.Payload)
	if err != nil {
		logrus.Errorf("multi term match failed, err=%v", err)
		h.Internal(w, r, err)
		return
	}

	h.Success(w, r, data)
}

type SearchKeywordReq struct {
	Payload string `json:"payload"`
}
