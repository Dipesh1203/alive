package services

import (
	"backend/db"
	"context"
	"errors"
	"log"
)

// CreateWebsite creates a website within an organization
func CreateWebsite(ctx context.Context, database *db.PrismaClient, websiteName string, url string, organizationID string) (*db.WebsiteModel, error) {
	log.Printf("[SERVICE] CreateWebsite: Creating website %s (URL: %s) in organization %s", websiteName, url, organizationID)
	if organizationID == "" {
		log.Printf("[SERVICE] CreateWebsite: ERROR - Organization ID is required")
		return nil, errors.New("organization ID is required")
	}

	website, err := database.Website.CreateOne(
		db.Website.WebsiteName.Set(websiteName),
		db.Website.URL.Set(url),
		db.Website.Organization.Link(db.Organization.ID.Equals(organizationID)),
	).Exec(ctx)

	if err != nil {
		log.Printf("[SERVICE] CreateWebsite: ERROR - Failed to create website - %v", err)
		return nil, err
	}

	log.Printf("[SERVICE] CreateWebsite: Website created successfully with ID: %s", website.ID)
	return website, nil
}

// ListWebsites lists all websites for a user (filtered by their organizations)
func ListWebsites(ctx context.Context, database *db.PrismaClient, userID string) ([]db.WebsiteModel, error) {
	log.Printf("[SERVICE] ListWebsites: Fetching websites for user %s", userID)
	// Get user's organizations
	members, err := database.OrganizationMember.FindMany(
		db.OrganizationMember.UserID.Equals(userID),
	).Exec(ctx)

	if err != nil {
		log.Printf("[SERVICE] ListWebsites: ERROR - Failed to fetch organizations - %v", err)
		return []db.WebsiteModel{}, err
	}

	log.Printf("[SERVICE] ListWebsites: Found %d organizations for user %s", len(members), userID)
	// Collect organization IDs
	orgIDs := make([]string, len(members))
	for i, member := range members {
		orgIDs[i] = member.OrganizationID
	}

	if len(orgIDs) == 0 {
		log.Printf("[SERVICE] ListWebsites: User %s has no organizations", userID)
		return []db.WebsiteModel{}, nil
	}

	// Fetch websites from these organizations
	websites := make([]db.WebsiteModel, 0)
	for _, orgID := range orgIDs {
		log.Printf("[SERVICE] ListWebsites: Fetching websites from organization %s", orgID)
		orgsWebsites, err := database.Website.FindMany(
			db.Website.OrganizationID.Equals(orgID),
		).Exec(ctx)

		if err != nil {
			log.Printf("[SERVICE] ListWebsites: WARNING - Failed to fetch websites from org %s - %v", orgID, err)
			continue
		}

		websites = append(websites, orgsWebsites...)
	}

	log.Printf("[SERVICE] ListWebsites: Retrieved %d total websites for user %s", len(websites), userID)
	return websites, nil
}

// GetWebsite gets a website and verifies user has access
func GetWebsite(ctx context.Context, database *db.PrismaClient, id string, userID string) (*db.WebsiteModel, error) {
	log.Printf("[SERVICE] GetWebsite: Fetching website %s for user %s", id, userID)
	website, err := database.Website.FindUnique(db.Website.ID.Equals(id)).Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] GetWebsite: ERROR - Website %s not found - %v", id, err)
		return nil, errors.New("website not found")
	}

	log.Printf("[SERVICE] GetWebsite: Website %s found, verifying access for user %s in org %s", id, userID, website.OrganizationID)
	// Verify user has access to this organization
	if website.OrganizationID != "" {
		role, err := GetMemberRole(ctx, database, website.OrganizationID, userID)
		if err != nil || role == "" {
			log.Printf("[SERVICE] GetWebsite: ERROR - Access denied for user %s to organization %s (role: %s)", userID, website.OrganizationID, role)
			return nil, errors.New("access denied")
		}
	}

	log.Printf("[SERVICE] GetWebsite: Access verified for user %s to website %s", userID, id)
	return website, nil
}

// UpdateWebsite updates a website (with org access check)
func UpdateWebsite(ctx context.Context, database *db.PrismaClient, id string, userID string, updates []db.WebsiteSetParam) (*db.WebsiteModel, error) {
	log.Printf("[SERVICE] UpdateWebsite: Updating website %s for user %s", id, userID)
	// First verify access
	_, err := GetWebsite(ctx, database, id, userID)
	if err != nil {
		log.Printf("[SERVICE] UpdateWebsite: ERROR - Access check failed - %v", err)
		return nil, err
	}

	log.Printf("[SERVICE] UpdateWebsite: Applying %d updates to website %s", len(updates), id)
	updatedWebsite, err := database.Website.FindUnique(db.Website.ID.Equals(id)).Update(updates...).Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] UpdateWebsite: ERROR - Failed to update website %s - %v", id, err)
		return nil, err
	}

	log.Printf("[SERVICE] UpdateWebsite: Website %s updated successfully", id)
	return updatedWebsite, nil
}

// DeleteWebsite deletes a website (with org access check)
func DeleteWebsite(ctx context.Context, database *db.PrismaClient, id string, userID string) (*db.WebsiteModel, error) {
	log.Printf("[SERVICE] DeleteWebsite: Deleting website %s for user %s", id, userID)
	// First verify access
	_, err := GetWebsite(ctx, database, id, userID)
	if err != nil {
		log.Printf("[SERVICE] DeleteWebsite: ERROR - Access check failed - %v", err)
		return nil, err
	}

	log.Printf("[SERVICE] DeleteWebsite: Deleting website %s from database", id)
	deleteWebsite, err := database.Website.FindUnique(db.Website.ID.Equals(id)).Delete().Exec(ctx)
	if err != nil {
		log.Printf("[SERVICE] DeleteWebsite: ERROR - Failed to delete website %s - %v", id, err)
		return nil, err
	}

	log.Printf("[SERVICE] DeleteWebsite: Website %s deleted successfully", id)
	return deleteWebsite, nil
}
