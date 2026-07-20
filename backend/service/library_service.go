package service

import (
	"todo-backend/models"
	"todo-backend/repository"
)

// LibraryService ...
type LibraryService interface {
	CreateDocument(doc models.LibraryDocument) (models.LibraryDocument, error)
	GetDocumentsByUserID(userID int) ([]models.LibraryDocument, error)
	GetDocumentByID(id int, userID int) (models.LibraryDocument, error)
	DeleteDocument(id int, userID int) error
}

type libraryService struct {
	repo repository.LibraryRepository
}

// NewLibraryService ...
func NewLibraryService(repo repository.LibraryRepository) LibraryService {
	return &libraryService{repo: repo}
}

// CreateDocument ...
func (s *libraryService) CreateDocument(doc models.LibraryDocument) (models.LibraryDocument, error) {
	return s.repo.CreateDocument(doc)
}

// GetDocumentsByUserID ...
func (s *libraryService) GetDocumentsByUserID(userID int) ([]models.LibraryDocument, error) {
	return s.repo.GetDocumentsByUserID(userID)
}

// GetDocumentByID ...
func (s *libraryService) GetDocumentByID(id int, userID int) (models.LibraryDocument, error) {
	return s.repo.GetDocumentByID(id, userID)
}

// DeleteDocument ...
func (s *libraryService) DeleteDocument(id int, userID int) error {
	return s.repo.DeleteDocument(id, userID)
}
