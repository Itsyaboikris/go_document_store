package query

import (
	"errors"

	"github.com/itsyaboikris/go_document_store/models"
)

type Query struct {
	matcher *Matcher
}

func NewQuery() *Query {
	return &Query{
		matcher: NewMatcher(),
	}
}

func (q *Query) Execute(data interface{}, filter map[string]interface{}) (interface{}, error) {
	if filter == nil {
		return data, nil
	}

	if err := q.validateFilter(filter); err != nil {
		return nil, err
	}

	if documents, ok := data.([]*models.Document); ok {

		var results []*models.Document
		for _, doc := range documents {
			if q.matcher.Matches(doc.Data, filter) {
				results = append(results, doc)
			}
		}
		return results, nil
	}

	return nil, errors.New("invalid data type")
}

func (q *Query) validateFilter(filter map[string]interface{}) error {
	for key, value := range filter {
		if key[0] == '$' {
			if !ValidateOperator(Operator(key)) {
				return errors.New("invalid operator: " + key)
			}
		}

		switch v := value.(type) {
		case map[string]interface{}:
			if err := q.validateFilter(v); err != nil {
				return err
			}
		case []interface{}:
			for _, item := range v {
				if subFilter, ok := item.(map[string]interface{}); ok {
					if err := q.validateFilter(subFilter); err != nil {
						return err
					}
				}
			}
		}

	}

	return nil
}
