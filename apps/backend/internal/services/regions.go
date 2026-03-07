package services

import (
	"backend/db"
	"context"
)

func CreateRegion(ctx context.Context, database *db.PrismaClient, regionName string) (*db.RegionModel, error) {
	return database.Region.CreateOne(
		db.Region.RegionName.Set(regionName),
	).Exec(ctx)
}

func ListRegions(ctx context.Context, database *db.PrismaClient) ([]db.RegionModel, error) {
	return database.Region.FindMany().Exec(ctx)
}

func GetRegion(ctx context.Context, database *db.PrismaClient, id string) (*db.RegionModel, error) {
	return database.Region.FindUnique(db.Region.RegionID.Equals(id)).Exec(ctx)
}

func UpdateRegion(ctx context.Context, database *db.PrismaClient, id string, updates []db.RegionSetParam) (*db.RegionModel, error) {
	return database.Region.FindUnique(db.Region.RegionID.Equals(id)).Update(updates...).Exec(ctx)
}

func DeleteRegion(ctx context.Context, database *db.PrismaClient, id string) (*db.RegionModel, error) {
	return database.Region.FindUnique(db.Region.RegionID.Equals(id)).Delete().Exec(ctx)
}
