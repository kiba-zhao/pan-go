package services

import (
	"pan/app/errors"
	"pan/app/models"
	"pan/core"
	"strings"
)

type ModuleService struct {
	Registry core.Registry
}

// SetEnabled, Retrieve the module from Registry based on name. If the module implements the Enabling Module interface, set it; otherwise, return an error
func (s *ModuleService) Update(name string, enable bool) (models.Module, error) {
	module := s.Registry.GetModuleByName(name)
	if module == nil {
		return models.Module{}, errors.ErrNotFound
	}
	if enablingModule, ok := module.(EnabledModule); ok {
		err := enablingModule.SetEnable(enable)
		if err != nil {
			return models.Module{}, err
		}
		return generateModule(module), nil
	}
	return models.Module{}, errors.ErrForbidden
}

// Get
func (s *ModuleService) Get(name string) (models.Module, error) {
	module := s.Registry.GetModuleByName(name)
	if module == nil {
		return models.Module{}, errors.ErrNotFound
	}
	return generateModule(module), nil
}

func (s *ModuleService) GetAll() []models.Module {
	modules := s.Registry.GetModules()
	result := make([]models.Module, 0)
	for _, m := range modules {
		result = append(result, generateModule(m))
	}
	return result
}

func (s *ModuleService) Search(conditions models.ModuleSearchCondition) (total int64, items []models.Module, err error) {

	keywords := strings.Split(conditions.Keyword, ",")
	kws := make([]string, 0)
	for _, keyword := range keywords {
		if strings.Trim(keyword, " ") == "" {
			continue
		}
		kws = append(kws, keyword)
	}

	if len(kws) == 0 {
		items = s.GetAll()
		total = int64(len(items))
		return
	}

	modules := s.Registry.GetModules()
	total = 0
	items = make([]models.Module, 0)
	for _, m := range modules {

		matchedStatus := -1
		name := m.Name()
		desc := m.Desc()
		for _, kw := range kws {
			if strings.Contains(name, kw) || strings.Contains(desc, kw) {
				matchedStatus = 1
				break
			}
			matchedStatus = 0
		}

		if matchedStatus == 0 {
			continue
		}

		total++
		items = append(items, generateModule(m))
	}

	return
}

type EnabledModule interface {
	Enabled() bool
	SetEnable(enable bool) error
}

type ReadOnlyModule interface {
	ReadOnly() bool
}

func generateModule(appModule core.AppModule) models.Module {

	module := models.Module{
		Avatar:   appModule.Avatar(),
		Name:     appModule.Name(),
		Desc:     appModule.Desc(),
		Enabled:  true,
		ReadOnly: true,
		HasWeb:   false,
	}

	if web, ok := appModule.(core.WebModule); ok {
		module.HasWeb = web.HasWeb()
	}
	if enabled, ok := appModule.(EnabledModule); ok {
		module.Enabled = enabled.Enabled()
	}
	if readOnly, ok := appModule.(ReadOnlyModule); ok {
		module.ReadOnly = readOnly.ReadOnly()
	}

	return module
}
