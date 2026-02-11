package api

import (
	"encoding/json"
	"net/http"

	"github.com/ivan-almanza/notes-api/internal/store"
)

type NotesHandler struct {
	store store.NoteStorer
}

func NewNotesHandler(store store.NoteStorer) *NotesHandler {
	return &NotesHandler{store: store}
}

type CreateNoteRequest struct {
	Content string `json:"content"`
}

type CreateNoteResponse struct {
	ID string `json:"id"`
}

func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	// 1. Get UserID from context (middleware injected)
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Parse body
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 3. Create Note object
	note := &store.Note{
		UserID:  userID,
		Content: req.Content,
	}

	// 4. Save to Store
	if err := h.store.CreateNote(r.Context(), note); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 5. Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateNoteResponse{ID: note.ID})
}

func (h *NotesHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notes, err := h.store.ListNotes(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": notes,
		"meta": map[string]interface{}{
			"count": len(notes),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
