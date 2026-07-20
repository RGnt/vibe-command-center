package repository

import (
	"database/sql"
	"time"
	"todo-backend/internal/models"
)

// GetAllByProjectID ...
func (r *wikiRepository) GetAllByProjectID(userID int, projectID int) ([]models.WikiPage, error) {
	rows, err := r.db.Query(`SELECT id, project_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE user_id = $1 AND project_id = $2 ORDER BY category, title`, userID, projectID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var wikis []models.WikiPage
	for rows.Next() {
		var wiki models.WikiPage
		var pID sql.NullInt64
		var createdAt, updatedAt time.Time

		err := rows.Scan(&wiki.ID, &pID, &wiki.Category, &wiki.Title, &wiki.Slug, &wiki.Content, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		wiki.UserID = userID
		wiki.CreatedAt = createdAt
		wiki.UpdatedAt = updatedAt
		
		if pID.Valid {
			idVal := int(pID.Int64)
			wiki.ProjectID = &idVal
		}

		wikis = append(wikis, wiki)
	}

	return wikis, nil
}
