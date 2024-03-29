package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/linksAPI/db"
	"github.com/hillview.tv/linksAPI/structs"
)

func LookupRoute(db db.Queryable, route string) (*structs.Route, error) {

	q := sq.Select("links.*").
		From("links").
		Where("links.route = ?", route).
		Where("links.active = ?", true)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	var routeRow structs.Route
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	err = rows.Scan(
		&routeRow.ID,
		&routeRow.Route,
		&routeRow.Destination,
		&routeRow.Active,
		&routeRow.CreatedBy,
		&routeRow.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	return &routeRow, nil
}
