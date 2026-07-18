package service

import (
	"regexp"
	"strings"
	"todo-backend/models"
	"todo-backend/repository"
)

type WikiService interface {
	GetWikis(userID int) ([]models.WikiPage, error)
	GetWiki(slug string, userID int) (models.WikiPage, error)
	CreateWiki(userID int, wiki models.WikiPage) (models.WikiPage, error)
	UpdateWiki(id, userID int, wiki models.WikiPage) (models.WikiPage, error)
	DeleteWiki(id, userID int) error
}

type wikiService struct {
	wikiRepo repository.WikiRepository
}

func NewWikiService(wikiRepo repository.WikiRepository) WikiService {
	return &wikiService{
		wikiRepo: wikiRepo,
	}
}

func (s *wikiService) GetWikis(userID int) ([]models.WikiPage, error) {
	return s.wikiRepo.GetAllByUserID(userID)
}

func (s *wikiService) GetWiki(slug string, userID int) (models.WikiPage, error) {
	return s.wikiRepo.GetBySlugAndUserID(slug, userID)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	return slug
}

func (s *wikiService) CreateWiki(userID int, wiki models.WikiPage) (models.WikiPage, error) {
	wiki.UserID = userID
	
	if wiki.Slug == "" {
		wiki.Slug = generateSlug(wiki.Title)
	}

	return s.wikiRepo.Create(wiki)
}

func (s *wikiService) UpdateWiki(id, userID int, wiki models.WikiPage) (models.WikiPage, error) {
	wiki.UserID = userID
	if wiki.Slug == "" {
		wiki.Slug = generateSlug(wiki.Title)
	}
	return s.wikiRepo.Update(id, wiki)
}

func (s *wikiService) DeleteWiki(id, userID int) error {
	return s.wikiRepo.Delete(id, userID)
}
