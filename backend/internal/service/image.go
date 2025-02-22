package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type ImageService struct {
	campaignRepo repo.Campaign
	mediaBaseUrl string
	mediaFsPath  string
}

func (s *ImageService) AddCampaignImage(campaign model.Campaign, file *multipart.FileHeader) (model.Campaign, error) {
	if !strings.Contains(file.Filename, ".") {
		return model.Campaign{}, fmt.Errorf("filename without extension: %s", file.Filename)
	}

	ext := file.Filename[strings.LastIndex(file.Filename, "."):]
	filename := uuid.NewString() + ext
	campaign.ImagePath = s.mediaBaseUrl + "/" + filename

	src, err := file.Open()
	if err != nil {
		return model.Campaign{}, fmt.Errorf("open source reader: %w", err)
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	out, err := os.Create(s.mediaFsPath + "/" + filename)
	if err != nil {
		return model.Campaign{}, fmt.Errorf("open destination file: %w", err)
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	_, err = io.Copy(out, src)
	if err != nil {
		return model.Campaign{}, fmt.Errorf("write destination file: %w", err)
	}

	if err := s.campaignRepo.Update(campaign); err != nil {
		return model.Campaign{}, fmt.Errorf("update campaign: %w", err)
	}
	return campaign, nil
}

func (s *ImageService) DeleteCampaignImage(campaign model.Campaign) (model.Campaign, error) {
	campaign.ImagePath = ""
	if err := s.campaignRepo.Update(campaign); err != nil {
		return model.Campaign{}, fmt.Errorf("update campaign: %w", err)
	}
	return campaign, nil
}
