package searchitem

type SearchItemInternalService interface {
	IsNotExist(error) bool
}

type SearchItemService struct {
	SearchItemRepo SearchItemRepository
}

func (s *SearchItemService) IsNotExist(err error) bool {
	return err == ErrSearchItemNotFound
}

func (s *SearchItemService) Select(id uint64) (SearchItem, error) {
	return s.SearchItemRepo.Select(id)
}

func (s *SearchItemService) Search(conditions SearchItemCondition) (int64, []SearchItem, error) {
	return s.SearchItemRepo.Search(conditions)
}

func (s *SearchItemService) Create(fields SearchItemFields) (SearchItem, error) {
	var model SearchItem
	model.Query = fields.Query

	return s.SearchItemRepo.Create(model)
}

func (s *SearchItemService) SelectOrCreate(fields SearchItemFields) (SearchItem, error) {
	var model SearchItem
	model.Query = fields.Query
	item, _, err := s.SearchItemRepo.SelectOrCreate(model)
	return item, err
}

func (s *SearchItemService) Delete(id uint64) error {
	return s.SearchItemRepo.Delete(id)
}
