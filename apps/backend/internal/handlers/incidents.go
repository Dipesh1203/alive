package handlers

import (
	"backend/db"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const dateLayout = "2006-01-02"

type AssignWebsiteRegionsRequest struct {
	RegionIDs []string `json:"regionIds"`
}

type AssignWebsiteRegionsResponse struct {
	WebsiteID string   `json:"websiteId"`
	RegionIDs []string `json:"regionIds"`
}

type IncidentResponse struct {
	ID              int              `json:"id"`
	WebsiteID       string           `json:"websiteId"`
	WebsiteName     string           `json:"websiteName,omitempty"`
	WebsiteRegionID string           `json:"websiteRegionId"`
	UpStatus        db.WebsiteStatus `json:"upStatus"`
	Latency         *int             `json:"latency,omitempty"`
	CreatedAt       db.DateTime      `json:"createdAt"`
	UpdatedAt       db.DateTime      `json:"updatedAt"`
}

// AssignWebsiteRegions godoc
// @Summary      Assign monitoring regions for a website
// @Description  Replaces the region assignment used by monitoring workers for the provided website
// @Tags         websites
// @Accept       json
// @Produce      json
// @Param        id    path      string                       true  "Website ID"
// @Param        body  body      AssignWebsiteRegionsRequest  true  "Region IDs"
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/websites/{id}/regions [put]
func AssignWebsiteRegions(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		websiteID := mux.Vars(r)["id"]
		if websiteID == "" {
			http.Error(w, "Website id required", http.StatusBadRequest)
			return
		}

		var req AssignWebsiteRegionsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		regionIDs := normalizeStringSlice(req.RegionIDs)
		if len(regionIDs) == 0 {
			http.Error(w, "regionIds must contain at least one value", http.StatusBadRequest)
			return
		}

		_, err := services.GetWebsite(ctx, database, websiteID)
		if err != nil {
			if db.IsErrNotFound(err) {
				http.Error(w, "Website not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		regions, err := database.Region.FindMany(db.Region.RegionID.In(regionIDs)).Exec(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(regions) != len(regionIDs) {
			http.Error(w, "One or more regionIds do not exist", http.StatusBadRequest)
			return
		}

		_, err = database.WebsiteTicks.FindMany(
			db.WebsiteTicks.WebsiteID.Equals(websiteID),
			db.WebsiteTicks.UpStatus.Equals(db.WebsiteStatusUnknown),
			db.WebsiteTicks.Latency.IsNull(),
		).Delete().Exec(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, regionID := range regionIDs {
			_, err = database.WebsiteTicks.CreateOne(
				db.WebsiteTicks.Website.Link(db.Website.ID.Equals(websiteID)),
				db.WebsiteTicks.Region.Link(db.Region.RegionID.Equals(regionID)),
				db.WebsiteTicks.UpStatus.Set(db.WebsiteStatusUnknown),
			).Exec(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		utils.WriteJSON(w, http.StatusOK, AssignWebsiteRegionsResponse{
			WebsiteID: websiteID,
			RegionIDs: regionIDs,
		})
	}
}

// ListIncidents godoc
// @Summary      List incidents
// @Description  Returns global incident history with optional filters
// @Tags         incidents
// @Produce      json
// @Param        websiteId  query     string  false  "Filter by website ID"
// @Param        regionId   query     string  false  "Filter by region ID"
// @Param        status     query     string  false  "down,unknown,degraded,up"
// @Param        startDate  query     string  false  "YYYY-MM-DD"
// @Param        endDate    query     string  false  "YYYY-MM-DD"
// @Param        skip       query     int     false  "Offset"
// @Param        take       query     int     false  "Limit"
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /api/incidents [get]
func ListIncidents(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		query := r.URL.Query()
		whereFilters := []db.WebsiteTicksWhereParam{}

		if websiteID := strings.TrimSpace(query.Get("websiteId")); websiteID != "" {
			whereFilters = append(whereFilters, db.WebsiteTicks.WebsiteID.Equals(websiteID))
		}

		if regionID := strings.TrimSpace(query.Get("regionId")); regionID != "" {
			whereFilters = append(whereFilters, db.WebsiteTicks.WebsiteRegionID.Equals(regionID))
		}

		if statusQuery := strings.TrimSpace(query.Get("status")); statusQuery != "" {
			statuses, err := parseStatuses(statusQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			whereFilters = append(whereFilters, db.WebsiteTicks.UpStatus.In(statuses))
		} else {
			whereFilters = append(whereFilters, db.WebsiteTicks.UpStatus.NotIn([]db.WebsiteStatus{db.WebsiteStatusUp}))
		}

		if startDateStr := strings.TrimSpace(query.Get("startDate")); startDateStr != "" {
			startDate, err := time.Parse(dateLayout, startDateStr)
			if err != nil {
				http.Error(w, "Invalid startDate, expected YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			whereFilters = append(whereFilters, db.WebsiteTicks.CreatedAt.Gte(startDate))
		}

		if endDateStr := strings.TrimSpace(query.Get("endDate")); endDateStr != "" {
			endDate, err := time.Parse(dateLayout, endDateStr)
			if err != nil {
				http.Error(w, "Invalid endDate, expected YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			whereFilters = append(whereFilters, db.WebsiteTicks.CreatedAt.Lte(endDate.Add(24*time.Hour-time.Nanosecond)))
		}

		ticksQuery := database.WebsiteTicks.FindMany(whereFilters...).OrderBy(db.WebsiteTicks.CreatedAt.Order(db.SortOrderDesc))

		if skipStr := strings.TrimSpace(query.Get("skip")); skipStr != "" {
			skip, err := strconv.Atoi(skipStr)
			if err != nil || skip < 0 {
				http.Error(w, "Invalid skip", http.StatusBadRequest)
				return
			}
			ticksQuery = ticksQuery.Skip(skip)
		}

		if takeStr := strings.TrimSpace(query.Get("take")); takeStr != "" {
			take, err := strconv.Atoi(takeStr)
			if err != nil || take <= 0 {
				http.Error(w, "Invalid take", http.StatusBadRequest)
				return
			}
			ticksQuery = ticksQuery.Take(take)
		}

		ticks, err := ticksQuery.Exec(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		websiteNames, err := buildWebsiteNameMap(ctx, database, ticks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := make([]IncidentResponse, 0, len(ticks))
		for _, tick := range ticks {
			response = append(response, IncidentResponse{
				ID:              tick.ID,
				WebsiteID:       tick.WebsiteID,
				WebsiteName:     websiteNames[tick.WebsiteID],
				WebsiteRegionID: tick.WebsiteRegionID,
				UpStatus:        tick.UpStatus,
				Latency:         tick.InnerWebsiteTicks.Latency,
				CreatedAt:       tick.CreatedAt,
				UpdatedAt:       tick.UpdatedAt,
			})
		}

		utils.WriteJSON(w, http.StatusOK, response)
	}
}

func parseStatuses(raw string) ([]db.WebsiteStatus, error) {
	parts := strings.Split(raw, ",")
	statuses := make([]db.WebsiteStatus, 0, len(parts))
	seen := map[db.WebsiteStatus]struct{}{}

	for _, part := range parts {
		normalized := strings.ToLower(strings.TrimSpace(part))
		if normalized == "" {
			continue
		}

		var status db.WebsiteStatus
		switch normalized {
		case "up":
			status = db.WebsiteStatusUp
		case "down":
			status = db.WebsiteStatusDown
		case "unknown", "degraded":
			status = db.WebsiteStatusUnknown
		default:
			return nil, fmt.Errorf("invalid status '%s'", normalized)
		}

		if _, exists := seen[status]; exists {
			continue
		}
		seen[status] = struct{}{}
		statuses = append(statuses, status)
	}

	if len(statuses) == 0 {
		return nil, fmt.Errorf("status filter is empty")
	}

	return statuses, nil
}

func normalizeStringSlice(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}

	sort.Strings(result)
	return result
}

func buildWebsiteNameMap(ctx context.Context, database *db.PrismaClient, ticks []db.WebsiteTicksModel) (map[string]string, error) {
	websiteIDs := make([]string, 0)
	seen := map[string]struct{}{}

	for _, tick := range ticks {
		if _, exists := seen[tick.WebsiteID]; exists {
			continue
		}
		seen[tick.WebsiteID] = struct{}{}
		websiteIDs = append(websiteIDs, tick.WebsiteID)
	}

	if len(websiteIDs) == 0 {
		return map[string]string{}, nil
	}

	websites, err := database.Website.FindMany(db.Website.ID.In(websiteIDs)).Exec(ctx)
	if err != nil {
		return nil, err
	}

	websiteNames := make(map[string]string, len(websites))
	for _, website := range websites {
		websiteNames[website.ID] = website.WebsiteName
	}

	return websiteNames, nil
}
