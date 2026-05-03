package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/linksAPI/db"
)

type RecordClickRequest struct {
	LinkID    *int
	UserID    *int
	IPAddress *string
}

func RecordClick(db db.Queryable, req RecordClickRequest) error {
	// validate fields
	if req.LinkID == nil {
		return fmt.Errorf("missing linkID")
	}

	cols := []string{"link_id"}
	vals := []interface{}{*req.LinkID}

	if req.UserID != nil {
		cols = append(cols, "user_id")
		vals = append(vals, req.UserID)
	}

	if req.IPAddress != nil {
		cols = append(cols, "ip_address")
		vals = append(vals, req.IPAddress)
	}

	query, args, err := sq.Insert("link_clicks").
		Columns(cols...).
		Values(vals...).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
