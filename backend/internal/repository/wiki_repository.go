package repository

import (
	"database/sql"
	"time"
	"todo-backend/internal/models"
)

// WikiRepository ...
type WikiRepository interface {
	GetAllByUserID(userID int) ([]models.WikiPage, error)
	GetAllByProjectID(userID int, projectID int) ([]models.WikiPage, error)
	GetBySlugAndUserID(slug string, userID int) (models.WikiPage, error)
	GetByIDAndUserID(id int, userID int) (models.WikiPage, error)
	Create(wiki models.WikiPage) (models.WikiPage, error)
	Update(id int, wiki models.WikiPage) (models.WikiPage, error)
	Delete(id int, userID int) error
	GetRevisions(wikiID int, userID int) ([]models.WikiPageRevision, error)
	GetRevision(wikiID int, revID int, userID int) (models.WikiPageRevision, error)
}

type wikiRepository struct {
	db *sql.DB
}

// NewWikiRepository ...
func NewWikiRepository(db *sql.DB) WikiRepository {
	return &wikiRepository{db: db}
}

// GetAllByUserID ...
func (r *wikiRepository) GetAllByUserID(userID int) ([]models.WikiPage, error) {
	rows, err := r.db.Query(`SELECT id, project_id, parent_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE user_id = $1 ORDER BY category, title`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var wikis []models.WikiPage
	for rows.Next() {
		var wiki models.WikiPage
		var projectID sql.NullInt64
		var parentID sql.NullInt64
		var createdAt, updatedAt time.Time

		err := rows.Scan(&wiki.ID, &projectID, &parentID, &wiki.Category, &wiki.Title, &wiki.Slug, &wiki.Content, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		wiki.UserID = userID
		wiki.CreatedAt = createdAt
		wiki.UpdatedAt = updatedAt
		
		if projectID.Valid {
			pid := int(projectID.Int64)
			wiki.ProjectID = &pid
		}
		if parentID.Valid {
			parid := int(parentID.Int64)
			wiki.ParentID = &parid
		}

		wikis = append(wikis, wiki)
	}

	return wikis, nil
}

// GetBySlugAndUserID ...
func (r *wikiRepository) GetBySlugAndUserID(slug string, userID int) (models.WikiPage, error) {
	var wiki models.WikiPage
	var projectID sql.NullInt64
	var parentID sql.NullInt64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`SELECT id, project_id, parent_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE slug = $1 AND user_id = $2`, slug, userID).
		Scan(&wiki.ID, &projectID, &parentID, &wiki.Category, &wiki.Title, &wiki.Slug, &wiki.Content, &createdAt, &updatedAt)
	if err != nil {
		return wiki, err
	}

	wiki.UserID = userID
	wiki.CreatedAt = createdAt
	wiki.UpdatedAt = updatedAt

	if projectID.Valid {
		pid := int(projectID.Int64)
		wiki.ProjectID = &pid
	}
	if parentID.Valid {
		parid := int(parentID.Int64)
		wiki.ParentID = &parid
	}

	return wiki, nil
}

// GetByIDAndUserID ...
func (r *wikiRepository) GetByIDAndUserID(id int, userID int) (models.WikiPage, error) {
	var wiki models.WikiPage
	var projectID sql.NullInt64
	var parentID sql.NullInt64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`SELECT id, project_id, parent_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&wiki.ID, &projectID, &parentID, &wiki.Category, &wiki.Title, &wiki.Slug, &wiki.Content, &createdAt, &updatedAt)
	if err != nil {
		return wiki, err
	}

	wiki.UserID = userID
	wiki.CreatedAt = createdAt
	wiki.UpdatedAt = updatedAt

	if projectID.Valid {
		pid := int(projectID.Int64)
		wiki.ProjectID = &pid
	}
	if parentID.Valid {
		parid := int(parentID.Int64)
		wiki.ParentID = &parid
	}

	return wiki, nil
}

// Create ...
func (r *wikiRepository) Create(wiki models.WikiPage) (models.WikiPage, error) {
	var id int64
	query := `INSERT INTO wiki_pages (user_id, project_id, parent_id, category, title, slug, content) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(query, wiki.UserID, wiki.ProjectID, wiki.ParentID, wiki.Category, wiki.Title, wiki.Slug, wiki.Content).Scan(&id)
	if err != nil {
		return wiki, err
	}
	
	wiki.ID = int(id)
	wiki.CreatedAt = time.Now()
	wiki.UpdatedAt = time.Now()
	
	return wiki, nil
}

// Update ...
func (r *wikiRepository) Update(id int, wiki models.WikiPage) (models.WikiPage, error) {
	// First, fetch the old content
	var oldContent string
	err := r.db.QueryRow(`SELECT content FROM wiki_pages WHERE id = $1 AND user_id = $2`, id, wiki.UserID).Scan(&oldContent)
	if err != nil {
		return wiki, err
	}

	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return wiki, err
	}
	defer tx.Rollback()

	// Insert revision if content changed
	if oldContent != wiki.Content {
		_, err = tx.Exec(`INSERT INTO wiki_page_revisions (wiki_page_id, content) VALUES ($1, $2)`, id, oldContent)
		if err != nil {
			return wiki, err
		}
	}

	// Update page
	query := `UPDATE wiki_pages SET category = $1, title = $2, slug = $3, content = $4, project_id = $5, parent_id = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7 AND user_id = $8`
	_, err = tx.Exec(query, wiki.Category, wiki.Title, wiki.Slug, wiki.Content, wiki.ProjectID, wiki.ParentID, id, wiki.UserID)
	if err != nil {
		return wiki, err
	}
	
	if err := tx.Commit(); err != nil {
		return wiki, err
	}

	wiki.ID = id
	return wiki, nil
}

// Delete ...
func (r *wikiRepository) Delete(id int, userID int) error {
	_, err := r.db.Exec("DELETE FROM wiki_pages WHERE id = $1 AND user_id = $2", id, userID)
	return err
}

// GetRevisions ...
func (r *wikiRepository) GetRevisions(wikiID int, userID int) ([]models.WikiPageRevision, error) {
	// First check if the user has access to this wiki page
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM wiki_pages WHERE id = $1 AND user_id = $2`, wikiID, userID).Scan(&count)
	if err != nil || count == 0 {
		return nil, sql.ErrNoRows
	}

	rows, err := r.db.Query(`SELECT id, wiki_page_id, content, created_at FROM wiki_page_revisions WHERE wiki_page_id = $1 ORDER BY created_at DESC`, wikiID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revisions []models.WikiPageRevision
	for rows.Next() {
		var rev models.WikiPageRevision
		if err := rows.Scan(&rev.ID, &rev.WikiPageID, &rev.Content, &rev.CreatedAt); err != nil {
			continue
		}
		revisions = append(revisions, rev)
	}

	return revisions, nil
}

// GetRevision ...
func (r *wikiRepository) GetRevision(wikiID int, revID int, userID int) (models.WikiPageRevision, error) {
	var count int
	var rev models.WikiPageRevision
	err := r.db.QueryRow(`SELECT COUNT(*) FROM wiki_pages WHERE id = $1 AND user_id = $2`, wikiID, userID).Scan(&count)
	if err != nil || count == 0 {
		return rev, sql.ErrNoRows
	}

	err = r.db.QueryRow(`SELECT id, wiki_page_id, content, created_at FROM wiki_page_revisions WHERE id = $1 AND wiki_page_id = $2`, revID, wikiID).
		Scan(&rev.ID, &rev.WikiPageID, &rev.Content, &rev.CreatedAt)
	return rev, err
}
