package services

import (
	"errors"
	"os"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
)

type TargetService struct {
	TargetFileService *TargetFileService
	TargetRepo        repositories.TargetRepository
}

func (s *TargetService) Scan() error {
	targets, err := s.TargetRepo.FindAllWithEnabled()
	if err != nil {
		return err
	}

	errs := make([]error, 0)
	for _, target := range targets {
		upgradeErr := s.TargetFileService.UpgradeFileInfoForTarget(target)
		if upgradeErr != nil {
			errs = append(errs, upgradeErr)
		}
		scanErr := s.ScanTarget(target)
		if scanErr != nil {
			errs = append(errs, scanErr)
		}
	}

	if len(errs) > 0 {
		err = errors.Join(errs...)
	}
	return err
}

func (s *TargetService) ScanTarget(target models.Target) error {

	if target.Enabled != true {
		return nil
	}

	stat, err := os.Stat(target.FilePath)
	if err != nil {
		return err
	}

	target.Name = stat.Name()
	target.ModifyTime = stat.ModTime()

	if stat.IsDir() {
		fileInfoTotal, err := s.TargetFileService.ScanFileInfosForDirectory(target)
		if err != nil {
			return err
		}
		target.Size = fileInfoTotal.Size
		target.Total = fileInfoTotal.Total

	} else {
		target.Size = stat.Size()
		target.Total = 1
		err := s.TargetFileService.ScanFileInfoForTarget(target)
		if err != nil {
			return err
		}
	}

	err = s.TargetRepo.Save(target)
	return err
}
