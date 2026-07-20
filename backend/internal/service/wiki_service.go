package service

import (
	"errors"
	"regexp"
	"strings"
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
)

// WikiService ...
type WikiService interface {
	GetWikis(userID int) ([]models.WikiPage, error)
	GetWiki(slug string, userID int) (models.WikiPage, error)
	CreateWiki(userID int, wiki models.WikiPage) (models.WikiPage, error)
	UpdateWiki(id, userID int, wiki models.WikiPage) (models.WikiPage, error)
	DeleteWiki(id, userID int) error
	GetWikiRevisions(wikiID, userID int) ([]models.WikiPageRevision, error)
	GetWikiRevision(wikiID, revID, userID int) (models.WikiPageRevision, error)
}

type wikiService struct {
	wikiRepo repository.WikiRepository
}

// NewWikiService ...
func NewWikiService(wikiRepo repository.WikiRepository) WikiService {
	return &wikiService{
		wikiRepo: wikiRepo,
	}
}

// GetWikis ...
func (s *wikiService) GetWikis(userID int) ([]models.WikiPage, error) {
	return s.wikiRepo.GetAllByUserID(userID)
}

// GetWiki ...
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

// CreateWiki ...
func (s *wikiService) CreateWiki(userID int, wiki models.WikiPage) (models.WikiPage, error) {
	wiki.UserID = userID
	
	if wiki.Slug == "" {
		wiki.Slug = generateSlug(wiki.Title)
	}

	return s.wikiRepo.Create(wiki)
}

// UpdateWiki ...
func (s *wikiService) UpdateWiki(id, userID int, wiki models.WikiPage) (models.WikiPage, error) {
	wiki.UserID = userID
	if wiki.Slug == "" {
		wiki.Slug = generateSlug(wiki.Title)
	}

	// Check for hierarchy cycle if setting a parent
	if wiki.ParentID != nil {
		if *wiki.ParentID == id {
			return wiki, errors.New("a page cannot be its own parent")
		}
		
		currentParentID := wiki.ParentID
		for currentParentID != nil {
			parentPage, err := s.wikiRepo.GetByIDAndUserID(*currentParentID, userID)
			if err != nil {
				break
			}
			if parentPage.ParentID != nil && *parentPage.ParentID == id {
				return wiki, errors.New("hierarchy cycle detected")
			}
			currentParentID = parentPage.ParentID
		}
	}

	return s.wikiRepo.Update(id, wiki)
}

// DeleteWiki ...
func (s *wikiService) DeleteWiki(id, userID int) error {
	return s.wikiRepo.Delete(id, userID)
}

// GetWikiRevisions ...
func (s *wikiService) GetWikiRevisions(wikiID, userID int) ([]models.WikiPageRevision, error) {
	return s.wikiRepo.GetRevisions(wikiID, userID)
}

// GetWikiRevision ...
func (s *wikiService) GetWikiRevision(wikiID, revID, userID int) (models.WikiPageRevision, error) {
	return s.wikiRepo.GetRevision(wikiID, revID, userID)
}
