package searchitem

import (
	"regexp"
	"strings"
)

type SearchItemInternalService interface {
	// IsNotExist returns true if the given error is a ErrSearchItemNotFound error
	IsNotExist(error) bool
}

type SearchItemService struct {
	SearchItemRepo SearchItemRepository
}

// IsNotExist returns true if the given error is a ErrSearchItemNotFound error
func (s *SearchItemService) IsNotExist(err error) bool {
	return err == ErrSearchItemNotFound
}

// Select retrieves the SearchItem associated with the given id.
// It returns the retrieved SearchItem and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
// If the SearchItem was not found, it returns ErrSearchItemNotFound.
func (s *SearchItemService) Select(id uint64) (SearchItem, error) {
	return s.SearchItemRepo.Select(id)
}

// Search retrieves SearchItems based on the given conditions, applying range conditions for pagination.
// It returns the total number of results, a slice of SearchItems, and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.

func (s *SearchItemService) Search(conditions SearchItemCondition) (int64, []SearchItem, error) {
	return s.SearchItemRepo.Search(conditions)
}

// Create creates a new SearchItem with the given query.
// It returns the created SearchItem and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
func (s *SearchItemService) Create(fields SearchItemFields) (SearchItem, error) {
	var model SearchItem
	model.Query = fields.Query

	return s.SearchItemRepo.Create(model)
}

// SelectOrCreate returns the SearchItem associated with the given query.
// If the SearchItem does not exist, it creates a new one.
// It returns the SearchItem and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
func (s *SearchItemService) SelectOrCreate(fields SearchItemFields) (SearchItem, error) {
	var model SearchItem
	model.Query = formatQuery(fields.Query)
	item, _, err := s.SearchItemRepo.SelectOrCreate(model)
	return item, err
}

// Delete removes the SearchItem associated with the given id from the database.
// It returns an error if any issues occur during the deletion process.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
// If the SearchItem was not found, it returns ErrSearchItemNotFound.
func (s *SearchItemService) Delete(id uint64) error {
	return s.SearchItemRepo.Delete(id)
}

func formatQuery(query string) string {
	query = strings.Trim(query, " ")
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(query, " ")
}
