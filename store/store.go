package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        string                 `json:"_id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type Collection struct {
	ID        string               `json:"_id"`
	Documents map[string]*Document `json:"documents"`
}

type Project struct {
	ID          string                 `json:"_id"`
	Collections map[string]*Collection `json:"collections"`
}

type DocumentStore struct {
	Projects map[string]*Project `json:"projects"`
	mu       sync.RWMutex
}

func NewStore() *DocumentStore {
	return &DocumentStore{
		Projects: make(map[string]*Project),
	}
}

func (ds *DocumentStore) Create(projectID, collectionID string, document map[string]interface{}) (*Document, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	project, exist := ds.Projects[projectID]
	if !exist {
		project = &Project{
			ID:          projectID,
			Collections: make(map[string]*Collection),
		}
		ds.Projects[projectID] = project
	}

	collection, exists := project.Collections[collectionID]
	if !exists {
		collection = &Collection{
			ID:        collectionID,
			Documents: make(map[string]*Document),
		}
		project.Collections[collectionID] = collection
	}

	now := time.Now().UTC()
	doc := &Document{
		ID:        uuid.New().String(),
		Data:      document,
		CreatedAt: now,
		UpdatedAt: now,
	}

	collection.Documents[doc.ID] = doc

	return doc, nil
}

func (ds *DocumentStore) Get(projectID, collectionID, documentID string) (*Document, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	project, exists := ds.Projects[projectID]
	if !exists {
		return nil, errors.New("project not found")
	}

	collection, exists := project.Collections[collectionID]
	if !exists {
		return nil, errors.New("collection not found")
	}

	doc, exists := collection.Documents[documentID]
	if !exists {
		return nil, errors.New("document not found")
	}

	return doc, nil
}

func (ds *DocumentStore) GetAll(projectID, collectionID string) ([]*Document, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	project, exists := ds.Projects[projectID]
	if !exists {
		return nil, errors.New("project not found")
	}

	collection, exists := project.Collections[collectionID]
	if !exists {
		return nil, errors.New("collection not found")
	}

	docs := make([]*Document, 0, len(collection.Documents))
	for _, doc := range collection.Documents {
		docs = append(docs, doc)
	}

	return docs, nil
}

func (ds *DocumentStore) Update(projectID, collectionID, documentID string, data map[string]interface{}) (*Document, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	project, exists := ds.Projects[projectID]
	if !exists {
		return nil, errors.New("project not found")
	}

	collection, exists := project.Collections[collectionID]
	if !exists {
		return nil, errors.New("collection not found")
	}

	doc, exists := collection.Documents[documentID]
	if !exists {
		return nil, errors.New("document not found")
	}

	doc.Data = data
	doc.UpdatedAt = time.Now().UTC()

	return doc, nil
}

func (ds *DocumentStore) Delete(projectID, collectionID, documentID string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	project, exists := ds.Projects[projectID]
	if !exists {
		return errors.New("project not found")
	}

	collection, exists := project.Collections[collectionID]
	if !exists {
		return errors.New("collection not found")
	}

	if _, exists := collection.Documents[documentID]; !exists {
		return errors.New("document not found")
	}

	delete(collection.Documents, documentID)
	return nil
}

// replication
func (ds *DocumentStore) InsertWithID(projectID, collectionID string, doc *Document) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Ensure project exists
	project, exists := ds.Projects[projectID]
	if !exists {
		project = &Project{
			ID:          projectID,
			Collections: make(map[string]*Collection),
		}
		ds.Projects[projectID] = project
	}

	// Ensure collection exists
	collection, exists := project.Collections[collectionID]
	if !exists {
		collection = &Collection{
			ID:        collectionID,
			Documents: make(map[string]*Document),
		}
		project.Collections[collectionID] = collection
	}

	existingDoc, exists := collection.Documents[doc.ID]
	if exists {
		existingDoc.Data = doc.Data
		existingDoc.UpdatedAt = time.Now().UTC()
	} else {
		if doc.CreatedAt.IsZero() {
			doc.CreatedAt = time.Now().UTC()
		}
		if doc.UpdatedAt.IsZero() {
			doc.UpdatedAt = doc.CreatedAt
		}
		collection.Documents[doc.ID] = doc
	}

	return nil
}

// helpers

// Helper functions for managing projects and collections
func (ds *DocumentStore) CreateProject(projectID string) (*Project, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, exists := ds.Projects[projectID]; exists {
		return nil, errors.New("project already exists")
	}

	project := &Project{
		ID:          projectID,
		Collections: make(map[string]*Collection),
	}

	ds.Projects[projectID] = project
	return project, nil
}

func (ds *DocumentStore) CreateCollection(projectID, collectionID string) (*Collection, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	project, exists := ds.Projects[projectID]
	if !exists {
		return nil, errors.New("project not found")
	}

	if _, exists := project.Collections[collectionID]; exists {
		return nil, errors.New("collection already exists")
	}

	collection := &Collection{
		ID:        collectionID,
		Documents: make(map[string]*Document),
	}

	project.Collections[collectionID] = collection
	return collection, nil
}
