package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/itsyaboikris/go_document_store/config"
	"github.com/itsyaboikris/go_document_store/models"
	"github.com/itsyaboikris/go_document_store/replication"
	"github.com/itsyaboikris/go_document_store/store"
)

type Handler struct {
	store *store.DocumentStore
}

func NewHandler(store *store.DocumentStore) *Handler {
	return &Handler{store: store}
}

func RegisterRoutes(r *mux.Router, store *store.DocumentStore) {
	h := NewHandler(store)

	// Register your routes
	r.HandleFunc("/replicate", h.ReplicationHandler).Methods("POST")
	r.HandleFunc("/{project}/{collection}/document", h.CreateDocument).Methods("POST")
	r.HandleFunc("/{project}/{collection}/document", h.GetAllDocuments).Methods("GET")
	r.HandleFunc("/{project}/{collection}/document/{id}", h.UpdateDocument).Methods("PUT")
	r.HandleFunc("/{project}/{collection}/document/{id}", h.DeleteDocument).Methods("DELETE")

	r.HandleFunc("/{project}/{collection}/query", h.QueryDocuments).Methods("POST")
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]
	collectionID := vars["collection"]

	var document map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	doc, err := h.store.Create(projectID, collectionID, document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	replicationDocument := map[string]interface{}{
		"id":         doc.ID,
		"data":       doc.Data,
		"project":    projectID,
		"collection": collectionID,
		"created_at": doc.CreatedAt,
		"updated_at": doc.UpdatedAt,
	}

	peers := config.GetPeers()
	replication.Replicate(peers, projectID, collectionID, doc.ID, replicationDocument)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h *Handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]
	collectionID := vars["collection"]
	documentID := vars["id"]

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	doc, err := h.store.Update(projectID, collectionID, documentID, updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	replicationDoc := map[string]interface{}{
		"id":         doc.ID,
		"data":       doc.Data,
		"project":    projectID,
		"collection": collectionID,
		"created_at": doc.CreatedAt,
		"updated_at": doc.UpdatedAt,
		"operation":  "update",
	}

	// Replicate to peers
	peers := config.GetPeers()
	replication.Replicate(peers, projectID, collectionID, doc.ID, replicationDoc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h *Handler) GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionID := vars["collection"]
	projectID := vars["project"]

	docs, err := h.store.GetAll(projectID, collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"documents": docs})
}

func (h *Handler) ReplicationHandler(w http.ResponseWriter, r *http.Request) {
	var replicationData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&replicationData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectID, ok := replicationData["project"].(string)
	if !ok {
		http.Error(w, "Missing project ID", http.StatusBadRequest)
		return
	}

	collectionID, ok := replicationData["collection"].(string)
	if !ok {
		http.Error(w, "Missing collection ID", http.StatusBadRequest)
		return
	}

	docID, ok := replicationData["id"].(string)
	if !ok {
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	operation, _ := replicationData["operation"].(string)
	if operation == "delete" {
		err := h.store.Delete(projectID, collectionID, docID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, ok := replicationData["data"].(map[string]interface{})
	if !ok {
		http.Error(w, "Missing document data", http.StatusBadRequest)
		return
	}

	doc := &models.Document{
		ID:   docID,
		Data: data,
	}

	if createdAt, ok := replicationData["created_at"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, createdAt)
		if err == nil {
			doc.CreatedAt = parsedTime
		}
	}

	if updatedAt, ok := replicationData["updated_at"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, updatedAt)
		if err == nil {
			doc.UpdatedAt = parsedTime
		}
	}

	err := h.store.InsertWithID(projectID, collectionID, doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]
	collectionID := vars["collection"]
	documentID := vars["id"]

	_, err := h.store.Get(projectID, collectionID, documentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := h.store.Delete(projectID, collectionID, documentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	replicationDoc := map[string]interface{}{
		"id":         documentID,
		"project":    projectID,
		"collection": collectionID,
		"operation":  "delete",
	}

	peers := config.GetPeers()
	replication.Replicate(peers, projectID, collectionID, documentID, replicationDoc)

	w.WriteHeader(http.StatusNoContent)
}

// query handler
func (h *Handler) QueryDocuments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["project"]
	collectionID := vars["collection"]

	var filter map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	documents, err := h.store.Query(projectID, collectionID, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"documents": documents,
		"count":     len(documents),
	})
}
