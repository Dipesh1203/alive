package services

import (
	"backend/db"
	"context"
	"time"
)

const layout = "2006-01-02"

func ListTicks(ctx context.Context, database *db.PrismaClient, id string, skip *int, take *int, startDate *string, endDate *string) ([]db.WebsiteTicksModel, error) {
	whereFilters := []db.WebsiteTicksWhereParam{
		db.WebsiteTicks.WebsiteID.Equals(id),
	}
	end, err2 := time.Parse(layout, *endDate)
	if startDate != nil {
		start, err := time.Parse(layout, *startDate)
		if err == nil {
			whereFilters = append(whereFilters, db.WebsiteTicks.CreatedAt.Gte(start))
		}
	}
	if endDate != nil {
		if err2 == nil {
			whereFilters = append(whereFilters, db.WebsiteTicks.CreatedAt.Lte(end))
		}
	}
	query := database.WebsiteTicks.FindMany(whereFilters...)

	if skip != nil {
		query = query.Skip(*skip)
	}
	if take != nil {
		query = query.Take(*take)
	}
	query = query.OrderBy(db.WebsiteTicks.CreatedAt.Order(db.SortOrderDesc))

	return query.Exec(ctx)
}

func GetTicks(ctx context.Context, database *db.PrismaClient, tickId int) (*db.WebsiteTicksModel, error) {
	return database.WebsiteTicks.FindUnique(db.WebsiteTicks.ID.Equals(tickId)).Exec(ctx)
}
