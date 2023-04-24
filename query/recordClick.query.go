package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/linksAPI/db"
)

type RecordClickRequest struct {
	LinkID *int
}

func RecordClick(db db.Queryable, req RecordClickRequest) error {
	// validate fields
	if req.LinkID == nil {
		return fmt.Errorf("missing linkID")
	}

	// update the link
	query, args, err := sq.Insert(
		"link_clicks",
	).Columns(
		"link_id",
	).Values(
		*req.LinkID,
	).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
