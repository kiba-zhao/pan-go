package searchitem

type SearchItemService struct {
	SearchItemRepo SearchItemRepository
}

func (s *SearchItemService) IsNotExist(err error) bool {
	return err == ErrSearchItemNotFound
}

func (s *SearchItemService) Search(conditions SearchItemCondition) (int64, []SearchItem, error) {
	return s.SearchItemRepo.Search(conditions)
}

func (s *SearchItemService) SelectOrCreate(fields SearchItemFields) (SearchItem, error) {
	var model SearchItem
	model.Query = fields.Query
	return s.SearchItemRepo.SelectOrCreate(model)
}

func (s *SearchItemService) Update(id uint64, fields SearchItemFields) (SearchItem, error) {
	model, err := s.SearchItemRepo.Select(id)
	if err != nil {
		return model, err
	}
	if len(fields.Query) > 0 && model.Query != fields.Query {
		model.Query = fields.Query
	}
	return s.SearchItemRepo.Save(model)
}

func (s *SearchItemService) Delete(id uint64) error {
	return s.SearchItemRepo.Delete(id)
}
