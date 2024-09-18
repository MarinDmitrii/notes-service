package ports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MarinDmitrii/notes-service/internal/common"
	"github.com/MarinDmitrii/notes-service/internal/note/builder"
	"github.com/MarinDmitrii/notes-service/internal/note/domain"
	"github.com/MarinDmitrii/notes-service/internal/note/usecase"
)

type HttpNoteHandler struct {
	app *builder.Application
}

func NewHttpNoteHandler(app *builder.Application) HttpNoteHandler {
	return HttpNoteHandler{app: app}
}

func (h HttpNoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	request := &PostNote{}

	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			h.mapToResponse(w, http.StatusBadRequest, nil, err.Error())
			return
		}
	} else {
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	checkedDescription, err := YaCheckText(request.Description)
	if err != nil {
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	fmt.Println(checkedDescription)

	userID, ok := r.Context().Value(common.UserContextKey).(int)
	if !ok {
		fmt.Println("!!", userID)
		h.mapToResponse(w, http.StatusUnauthorized, nil, http.StatusText(http.StatusUnauthorized))
		return
	}

	request.UserID = userID
	request.Description = checkedDescription

	createdNote, err := h.app.SaveNote.Execute(
		r.Context(),
		usecase.SaveNote{
			UserID:      request.UserID,
			Description: request.Description,
		},
	)
	if err != nil {
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	h.mapToResponse(w, http.StatusOK, NewNote(createdNote), "")
}

func (h HttpNoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userID, ok := r.Context().Value(common.UserContextKey).(int)
	if !ok {
		h.mapToResponse(w, http.StatusUnauthorized, nil, http.StatusText(http.StatusUnauthorized))
		return
	}

	notes, err := h.app.GetNotes.Execute(r.Context(), userID)
	if err != nil {
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	newNotes := make([]Note, len(notes))
	for _, note := range notes {
		newNotes = append(newNotes, NewNote(note))
	}

	h.mapToResponse(w, http.StatusOK, newNotes, "")
}

type PostNote struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userId"`
	Description string `json:"description"`
}

func NewPostNote(note domain.Note) PostNote {
	return PostNote{
		ID:          note.ID,
		Description: note.Description,
	}
}

type Note struct {
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
	UserID      int       `json:"userId"`
	ID          int       `json:"id"`
}

func NewNote(note domain.Note) Note {
	return Note{
		Created:     note.CreateDt,
		Description: note.Description,
		UserID:      note.UserID,
		ID:          note.ID,
	}
}

func (h HttpNoteHandler) mapToResponse(w http.ResponseWriter, statusCode int, data interface{}, errMessage string) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	response := make(map[string]interface{})

	if statusCode >= 200 && statusCode < 300 {
		response["result"] = data
	} else {
		response["error"] = errMessage
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CustomRegisterHandlers(router *http.ServeMux, h HttpNoteHandler, authMiddleware *common.AuthMiddleware) {
	router.HandleFunc("/create_note", common.LogMiddleware(authMiddleware.BasicAuth(h.CreateNote)))
	router.HandleFunc("/notes", common.LogMiddleware(authMiddleware.BasicAuth(h.GetNotes)))
}

type YaCheckResult struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func YaCheckText(text string) (string, error) {
	data := url.Values{}
	data.Set("text", text)
	data.Set("lang", "ru, en")
	data.Set("options", fmt.Sprintf("%d", 14))

	resp, err := http.PostForm("https://speller.yandex.net/services/spellservice.json/checkText", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []YaCheckResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	correctedText := text
	duplicate := ""

	for _, res := range result {
		if duplicate == strings.ToLower(res.Word) {
			correctedText = strings.Replace(correctedText, res.Word, "", 1)
		} else {
			duplicate = strings.ToLower(res.Word)
			correctedText = strings.Replace(correctedText, res.Word, res.S[0], 1)
		}
	}

	for {
		if strings.Contains(correctedText, "  ") {
			correctedText = strings.ReplaceAll(correctedText, "  ", " ")
		} else {
			break
		}
	}
	return correctedText, nil
}
