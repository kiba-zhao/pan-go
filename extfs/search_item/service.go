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

func (s *SearchItemService) Create(fields SearchItemFields) (SearchItem, error) {
	var model SearchItem
	model.Query = fields.Query
	return s.SearchItemRepo.Save(model)
}

func (s *SearchItemService) Update(id uint, fields SearchItemFields) (SearchItem, error) {
	model, err := s.SearchItemRepo.Select(id)
	if err != nil {
		return model, err
	}
	if len(fields.Query) > 0 && model.Query != fields.Query {
		model.Query = fields.Query
	}
	return s.SearchItemRepo.Save(model)
}

func (s *SearchItemService) Delete(id uint) error {
	return s.SearchItemRepo.Delete(id)
}
