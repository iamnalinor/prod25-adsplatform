package service

import (
	"backend/internal/repo"
	"fmt"
)

type SettingsService struct {
	settingsRepo repo.Settings
}

func (s *SettingsService) Date() int {
	return s.settingsRepo.GetCached().CurrentDate
}

func (s *SettingsService) SetDate(date int) error {
	settings := s.settingsRepo.GetCached()
	settings.CurrentDate = date
	if err := s.settingsRepo.Update(settings); err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	return nil
}

func (s *SettingsService) ModerationEnabled() bool {
	return s.settingsRepo.GetCached().ModerationEnabled
}

func (s *SettingsService) SetModerationEnabled(enabled bool) error {
	settings := s.settingsRepo.GetCached()
	settings.ModerationEnabled = enabled
	if err := s.settingsRepo.Update(settings); err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	return nil
}
