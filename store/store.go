package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/itsyaboikris/go_document_store/models"
	"github.com/itsyaboikris/go_document_store/query"
)

type Collection struct {
	ID        string                      `json:"_id"`
	Documents map[string]*models.Document `json:"documents"`
}

type Project struct {
	ID          string                 `json:"_id"`
	Collections map[string]*Collection `json:"collections"`
}

type DocumentStore struct {
	Projects map[string]*Project `json:"projects"`
	mu       sync.RWMutex
	querier  *query.Query
}

func NewStore() *DocumentStore {
	return &DocumentStore{
		Projects: make(map[string]*Project),
		querier:  query.NewQuery(),
	}
}

func (ds *DocumentStore) Create(projectID, collectionID string, document map[string]interface{}) (*models.Document, error) {
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
			Documents: make(map[string]*models.Document),
		}
		project.Collections[collectionID] = collection
	}

	now := time.Now().UTC()
	doc := &models.Document{
		ID:        uuid.New().String(),
		Data:      document,
		CreatedAt: now,
		UpdatedAt: now,
	}

	collection.Documents[doc.ID] = doc

	return doc, nil
}

func (ds *DocumentStore) Get(projectID, collectionID, documentID string) (*models.Document, error) {
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

func (ds *DocumentStore) GetAll(projectID, collectionID string) ([]*models.Document, error) {
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

	docs := make([]*models.Document, 0, len(collection.Documents))
	for _, doc := range collection.Documents {
		docs = append(docs, doc)
	}

	return docs, nil
}

func (ds *DocumentStore) Update(projectID, collectionID, documentID string, data map[string]interface{}) (*models.Document, error) {
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
func (ds *DocumentStore) InsertWithID(projectID, collectionID string, doc *models.Document) error {
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
			Documents: make(map[string]*models.Document),
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
		Documents: make(map[string]*models.Document),
	}

	project.Collections[collectionID] = collection
	return collection, nil
}

func (ds *DocumentStore) Query(projectID, collectionID string, filter map[string]interface{}) ([]*models.Document, error) {
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

	documents := make([]*models.Document, 0, len(collection.Documents))
	for _, doc := range collection.Documents {
		documents = append(documents, doc)
	}

	results, err := ds.querier.Execute(documents, filter)
	if err != nil {
		return nil, err
	}

	if docs, ok := results.([]*models.Document); ok {
		return docs, nil
	}

	return nil, errors.New("invalid query result type")
}
