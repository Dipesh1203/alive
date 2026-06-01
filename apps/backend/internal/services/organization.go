package services

import (
	"context"
	"errors"
	"log"

	"github.com/Dipesh1203/alive/apps/backend/db"
)

// CreateOrganization creates a new organization with the provided admin
func CreateOrganization(ctx context.Context, database *db.PrismaClient, name string, adminID string) (*db.OrganizationModel, error) {
	if name == "" {
		return nil, errors.New("organization name is required")
	}

	// Verify admin user exists
	_, err := GetUserByID(ctx, database, adminID)
	if err != nil {
		return nil, errors.New("admin user not found")
	}

	// Create organization
	org, err := database.Organization.CreateOne(
		db.Organization.Name.Set(name),
		db.Organization.Admin.Link(db.User.ID.Equals(adminID)),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to create organization")
	}

	// Add admin as organization member with admin role
	_, err = database.OrganizationMember.CreateOne(
		db.OrganizationMember.User.Link(db.User.ID.Equals(adminID)),
		db.OrganizationMember.Organization.Link(db.Organization.ID.Equals(org.ID)),
		db.OrganizationMember.Role.Set("admin"),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to add admin to organization")
	}

	return org, nil
}

// GetOrganizationByID retrieves an organization by ID
func GetOrganizationByID(ctx context.Context, database *db.PrismaClient, orgID string) (*db.OrganizationModel, error) {
	org, err := database.Organization.FindUnique(db.Organization.ID.Equals(orgID)).Exec(ctx)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	return org, nil
}

func CheckUserOrganizationMembership(ctx context.Context, database *db.PrismaClient, orgID, userID string) (bool, error) {
	member, err := database.OrganizationMember.FindUnique(
		db.OrganizationMember.UserIDOrganizationID(
			db.OrganizationMember.UserID.Equals(userID),
			db.OrganizationMember.OrganizationID.Equals(orgID),
		),
	).Exec(ctx)
	if err != nil {
		return false, errors.New("failed to check user organization membership")
	}
	return member != nil, nil
}

// ListUserOrganizations lists all organizations where user is a member
func ListUserOrganizations(ctx context.Context, database *db.PrismaClient, userID string) ([]db.OrganizationModel, error) {
	members, err := database.OrganizationMember.FindMany(
		db.OrganizationMember.UserID.Equals(userID),
	).With(
		db.OrganizationMember.Organization.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to fetch organizations")
	}

	orgs := make([]db.OrganizationModel, 0)
	for _, member := range members {
		if member.Organization() != nil {
			orgs = append(orgs, *member.Organization())
		}
	}

	return orgs, nil
}

// UpdateOrganization updates only the organization fields that are provided.
func UpdateOrganization(ctx context.Context, database *db.PrismaClient, orgID string, name *string) (*db.OrganizationModel, error) {
	updates := make([]db.OrganizationSetParam, 0)
	if name != nil {
		updates = append(updates, db.Organization.Name.Set(*name))
	}
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	org, err := database.Organization.FindUnique(db.Organization.ID.Equals(orgID)).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to update organization")
	}

	return org, nil
}

// DeleteOrganization removes an organization
func DeleteOrganization(ctx context.Context, database *db.PrismaClient, orgID string) (*db.OrganizationModel, error) {
	org, err := database.Organization.FindUnique(db.Organization.ID.Equals(orgID)).Delete().Exec(ctx)
	if err != nil {
		return nil, errors.New("failed to delete organization")
	}

	return org, nil
}

// AddMemberToOrganization adds a user to an organization
func AddMemberToOrganization(ctx context.Context, database *db.PrismaClient, orgID, userID, role string) (*db.OrganizationMemberModel, error) {
	if role != "admin" && role != "viewer" {
		return nil, errors.New("invalid role")
	}

	// Verify org exists
	_, err := GetOrganizationByID(ctx, database, orgID)
	if err != nil {
		return nil, err
	}

	// Verify user exists
	_, err = GetUserByID(ctx, database, userID)
	if err != nil {
		return nil, err
	}

	// Check if already a member
	existing, _ := database.OrganizationMember.FindUnique(
		db.OrganizationMember.UserIDOrganizationID(
			db.OrganizationMember.UserID.Equals(userID),
			db.OrganizationMember.OrganizationID.Equals(orgID),
		),
	).Exec(ctx)

	if existing != nil {
		return nil, errors.New("user is already a member of this organization")
	}

	member, err := database.OrganizationMember.CreateOne(
		db.OrganizationMember.User.Link(db.User.ID.Equals(userID)),
		db.OrganizationMember.Organization.Link(db.Organization.ID.Equals(orgID)),
		db.OrganizationMember.Role.Set(role),
	).Exec(ctx)

	if err != nil {
		log.Printf("[SERVICE] AddMemberToOrganization: ERROR - Failed to add user %s to organization %s with role %s - %v", userID, orgID, role, err)
		return nil, errors.New("failed to add member to organization")
	}

	return member, nil
}

// RemoveMemberFromOrganization removes a user from an organization
func RemoveMemberFromOrganization(ctx context.Context, database *db.PrismaClient, orgID, memberId string) error {
	_, err := database.OrganizationMember.FindUnique(db.OrganizationMember.ID.Equals(memberId)).Delete().Exec(ctx)

	if err != nil {
		return errors.New("failed to remove member from organization")
	}

	return nil
}

// UpdateMemberRole updates a team member's role when it is provided.
func UpdateMemberRole(ctx context.Context, database *db.PrismaClient, orgID, memberID string, newRole *string) (*db.OrganizationMemberModel, error) {
	if newRole == nil {
		return nil, errors.New("no fields to update")
	}

	if *newRole != "admin" && *newRole != "viewer" {
		return nil, errors.New("invalid role")
	}

	member, err := database.OrganizationMember.
		FindUnique(db.OrganizationMember.ID.Equals(memberID)).
		Update(db.OrganizationMember.Role.Set(*newRole)).
		Exec(ctx)
	log.Printf("Updating member role for memberID %s in orgID %s to new role: %s", memberID, orgID, *newRole)

	if err != nil {
		return nil, errors.New("failed to update member role")
	}

	return member, nil
}

// ListOrganizationMembers lists all members of an organization
func ListOrganizationMembers(ctx context.Context, database *db.PrismaClient, orgID string) ([]db.OrganizationMemberModel, error) {
	members, err := database.OrganizationMember.FindMany(
		db.OrganizationMember.OrganizationID.Equals(orgID),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("failed to fetch organization members")
	}

	return members, nil
}

// GetMemberRole retrieves a member's role in an organization
func GetMemberRole(ctx context.Context, database *db.PrismaClient, orgID, userID string) (string, error) {
	member, err := database.OrganizationMember.FindUnique(
		db.OrganizationMember.UserIDOrganizationID(
			db.OrganizationMember.UserID.Equals(userID),
			db.OrganizationMember.OrganizationID.Equals(orgID),
		),
	).Exec(ctx)

	if err != nil || member == nil {
		return "", errors.New("member not found in organization")
	}

	return member.Role, nil
}
