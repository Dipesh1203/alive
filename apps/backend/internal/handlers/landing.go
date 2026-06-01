package handlers

import (
	"net/http"
	"strings"

	"github.com/Dipesh1203/alive/apps/backend/internal/utils"

	"github.com/Dipesh1203/alive/apps/backend/db"
)

type LandingStat struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
}

type LandingFeature struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type LandingTestimonial struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Company  string `json:"company"`
	Quote    string `json:"quote"`
	Rating   int    `json:"rating"`
	Location string `json:"location"`
}

type LandingPricingPlan struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Currency    string   `json:"currency"`
	Interval    string   `json:"interval"`
	Popular     bool     `json:"popular"`
	CTA         string   `json:"cta"`
	Features    []string `json:"features"`
}

type LandingFAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LandingOverviewResponse struct {
	Stats        []LandingStat        `json:"stats"`
	Features     []LandingFeature     `json:"features"`
	Testimonials []LandingTestimonial `json:"testimonials"`
	Pricing      []LandingPricingPlan `json:"pricing"`
	FAQs         []LandingFAQ         `json:"faqs"`
}

func landingStats() []LandingStat {
	return []LandingStat{
		{Label: "Checks / day", Value: "4.2M", Hint: "Across all monitored endpoints"},
		{Label: "Median alert time", Value: "23s", Hint: "From failure to notification"},
		{Label: "Global regions", Value: "14", Hint: "Distributed probe network"},
		{Label: "Avg. false positives", Value: "<0.1%", Hint: "Multi-region confirmation"},
	}
}

func landingFeatures() []LandingFeature {
	return []LandingFeature{
		{Title: "Smart incident timelines", Description: "See exactly when service quality dropped, recovered, and how long users were affected."},
		{Title: "Latency from real regions", Description: "Track response behavior by region to uncover geo-specific regressions before customers report them."},
		{Title: "Context-first alerts", Description: "Alert payloads include recent checks and suspected cause so responders can act without context switching."},
		{Title: "Organization-ready access", Description: "Manage teams and organizations with role-based access and clean ownership boundaries."},
	}
}

func landingTestimonials() []LandingTestimonial {
	return []LandingTestimonial{
		{
			Name:     "Rhea Martin",
			Role:     "Staff Engineer",
			Company:  "Northline Commerce",
			Quote:    "Alive helped us catch a regional DNS issue in under a minute. We resolved it before support tickets even started.",
			Rating:   5,
			Location: "Berlin",
		},
		{
			Name:     "Noah Delgado",
			Role:     "Platform Lead",
			Company:  "Pinegrid",
			Quote:    "The incident history and latency trend views gave us concrete proof for a provider migration that cut p95 by 40%.",
			Rating:   5,
			Location: "Toronto",
		},
		{
			Name:     "Mina Solanki",
			Role:     "Head of Reliability",
			Company:  "Sketchlane",
			Quote:    "We replaced two tools with Alive and reduced on-call noise while getting better uptime visibility for leadership.",
			Rating:   5,
			Location: "Bangalore",
		},
	}
}

func landingPricing(interval string) []LandingPricingPlan {
	multiplier := 1
	if interval == "yearly" {
		multiplier = 10
	}

	return []LandingPricingPlan{
		{
			Name:        "Starter",
			Description: "For indie apps and side projects",
			Price:       9 * multiplier,
			Currency:    "USD",
			Interval:    interval,
			Popular:     false,
			CTA:         "Start monitoring",
			Features:    []string{"10 websites", "5 regions", "Email alerts", "7-day retention"},
		},
		{
			Name:        "Growth",
			Description: "For growing teams shipping fast",
			Price:       29 * multiplier,
			Currency:    "USD",
			Interval:    interval,
			Popular:     true,
			CTA:         "Choose Growth",
			Features:    []string{"50 websites", "14 regions", "Slack + webhook alerts", "90-day retention"},
		},
		{
			Name:        "Scale",
			Description: "For mission-critical production systems",
			Price:       99 * multiplier,
			Currency:    "USD",
			Interval:    interval,
			Popular:     false,
			CTA:         "Contact sales",
			Features:    []string{"Unlimited websites", "All regions", "SLO reporting", "Priority support"},
		},
	}
}

func landingFAQs() []LandingFAQ {
	return []LandingFAQ{
		{Question: "How quickly are incidents detected?", Answer: "Alive runs checks at a short interval and confirms failures from multiple regions before opening an incident."},
		{Question: "Can I monitor staging and production separately?", Answer: "Yes. You can group endpoints by organization and environment with independent alert routing."},
		{Question: "Do you support custom alert channels?", Answer: "Yes. Webhook and notification channels are supported so alerts can flow into your existing stack."},
		{Question: "Can we start free and upgrade later?", Answer: "Absolutely. Start on Starter and upgrade as your coverage and team needs grow."},
	}
}

// GetLandingOverview godoc
// @Summary      Landing page overview data
// @Description  Returns dynamic content for landing page sections
// @Tags         landing
// @Produce      json
// @Param        billing  query   string  false  "monthly or yearly"
// @Success      200      {object} LandingOverviewResponse
// @Router       /api/public/landing [get]
func GetLandingOverview(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		billing := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("billing")))
		if billing != "yearly" {
			billing = "monthly"
		}

		utils.WriteJSON(w, http.StatusOK, LandingOverviewResponse{
			Stats:        landingStats(),
			Features:     landingFeatures(),
			Testimonials: landingTestimonials(),
			Pricing:      landingPricing(billing),
			FAQs:         landingFAQs(),
		})
	}
}

// ListLandingTestimonials godoc
// @Summary      Landing testimonials
// @Description  Returns testimonials for the public landing page
// @Tags         landing
// @Produce      json
// @Success      200  {array}  LandingTestimonial
// @Router       /api/public/landing/testimonials [get]
func ListLandingTestimonials(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, landingTestimonials())
	}
}

// ListLandingPricing godoc
// @Summary      Landing pricing
// @Description  Returns pricing plans for the public landing page
// @Tags         landing
// @Produce      json
// @Param        billing  query   string  false  "monthly or yearly"
// @Success      200      {array} LandingPricingPlan
// @Router       /api/public/landing/pricing [get]
func ListLandingPricing(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		billing := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("billing")))
		if billing != "yearly" {
			billing = "monthly"
		}

		utils.WriteJSON(w, http.StatusOK, landingPricing(billing))
	}
}

// ListLandingFAQs godoc
// @Summary      Landing FAQs
// @Description  Returns FAQ items for the public landing page
// @Tags         landing
// @Produce      json
// @Success      200  {array}  LandingFAQ
// @Router       /api/public/landing/faqs [get]
func ListLandingFAQs(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, landingFAQs())
	}
}
