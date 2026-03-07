package services

import (
	"backend/db"
	"context"
)

func CreateWebsite(ctx context.Context, database *db.PrismaClient, websiteName string, url string) (*db.WebsiteModel, error) {
	return database.Website.CreateOne(
		db.Website.WebsiteName.Set(websiteName),
		db.Website.URL.Set(url),
	).Exec(ctx)
}

func ListWebsites(ctx context.Context, database *db.PrismaClient) ([]db.WebsiteModel, error) {
	return database.Website.FindMany().Exec(ctx)
}

func GetWebsite(ctx context.Context, database *db.PrismaClient, id string) (*db.WebsiteModel, error) {
	return database.Website.FindUnique(db.Website.ID.Equals(id)).Exec(ctx)
}

func UpdateWebsite(ctx context.Context, database *db.PrismaClient, id string, updates []db.WebsiteSetParam) (*db.WebsiteModel, error) {
	return database.Website.FindUnique(db.Website.ID.Equals(id)).Update(updates...).Exec(ctx)
}

func DeleteWebsite(ctx context.Context, database *db.PrismaClient, id string) (
	*db.WebsiteModel, error) {
	return database.Website.FindUnique(db.Website.ID.Equals(id)).Delete().Exec(ctx)
}
