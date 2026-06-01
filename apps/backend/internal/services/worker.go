package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	awsservice "github.com/Dipesh1203/alive/apps/backend/internal/services/aws"
	"github.com/Dipesh1203/alive/apps/backend/internal/utils"

	"github.com/Dipesh1203/alive/apps/backend/db"
)

type ResultResponse struct {
	WebsiteID string
	Status    db.WebsiteStatus
	Latency   time.Duration
	Response  string
}

type Site struct {
	ID       string
	Website  db.WebsiteModel
	RegionID string
}

var BASE_URL = utils.GoGetEnv("BASE_URL")
var baseUrl = utils.GetEnv("BASE_URL", "http://localhost:3000/website")

func StartMonitoring(ctx context.Context, prisma *db.PrismaClient) {
	taskChannel := make(chan Site, 10)
	resultsChannel := make(chan ResultResponse, 10)
	var wg sync.WaitGroup
	for w := 1; w <= 1; w++ {
		log.Printf("Starting worker %d\n", w)
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			MonitoringWorker(w, taskChannel, resultsChannel, prisma, ctx)
		}(w)
	}

	go func() {
		for res := range resultsChannel {
			fmt.Printf("📊 Result processed for %s: %s (%v)\n", res.WebsiteID, res.Status, res.Latency)
		}
	}()

	ticker := time.NewTicker(100 * time.Second)
	defer ticker.Stop()
	fetchAndDispatch(ctx, prisma, taskChannel)

	for {
		select {
		case <-ctx.Done():
			log.Println("Monitoring service received shutdown signal")
			close(taskChannel)
			wg.Wait()
			close(taskChannel)
			return
		case <-ticker.C:
			fetchAndDispatch(ctx, prisma, taskChannel)

		}
	}
}

func fetchAndDispatch(ctx context.Context, prisma *db.PrismaClient, taskChannel chan<- Site) {
	activeWebsites, err := prisma.Website.FindMany().Exec(ctx)
	if err != nil {
		log.Println("Error fetching websites:", err)
		return
	}

	for _, site := range activeWebsites {
		markerTicks, err := prisma.WebsiteTicks.FindMany(
			db.WebsiteTicks.WebsiteID.Equals(site.ID),
			db.WebsiteTicks.UpStatus.Equals(db.WebsiteStatusUnknown),
			db.WebsiteTicks.Latency.IsNull(),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error fetching region assignment markers for website %s: %v", site.ID, err)
			continue
		}

		regionIDs := make([]string, 0)
		seen := make(map[string]struct{})
		for _, marker := range markerTicks {
			if _, exists := seen[marker.WebsiteRegionID]; exists {
				continue
			}
			seen[marker.WebsiteRegionID] = struct{}{}
			regionIDs = append(regionIDs, marker.WebsiteRegionID)
		}

		if len(regionIDs) == 0 {
			regions, err := prisma.Region.FindMany().Exec(ctx)
			if err != nil {
				continue
			}
			for _, region := range regions {
				regionIDs = append(regionIDs, region.RegionID)
			}
		}

		for _, regionID := range regionIDs {
			taskChannel <- Site{
				ID:       site.ID,
				Website:  site,
				RegionID: regionID,
			}
		}
	}
}

func MonitoringWorker(id int, jobs <-chan Site, results chan<- ResultResponse, prisma *db.PrismaClient, ctx context.Context) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	for site := range jobs {
		log.Printf("Worker %d: Checking website %s in region %s\n", id, site.Website.WebsiteName, site.RegionID)
		startTime := time.Now()
		res, err := client.Get(site.Website.URL)
		latency := time.Since(startTime)

		newStatus := db.WebsiteStatusUp
		respMessage := "200 OK"
		statusCode := 0
		org, err := prisma.Organization.FindUnique(
			db.Organization.ID.Equals(site.Website.OrganizationID),
		).Exec(ctx)
		if err != nil {
			log.Printf("Worker %d: unable to load organization for website %s: %v", id, site.Website.URL, err)
		}
		user, err := prisma.User.
			FindFirst(
				db.User.ID.Equals(org.AdminID),
			).
			Exec(ctx)
		log.Printf("Start Fetching user")
		if err != nil {
			log.Printf("Worker %d: unable to load user for website %s: %v", id, site.Website.URL, err)
		} else {
			log.Printf("User Email: %s", user.Email)
		}
		userProfile, err := prisma.UserProfile.
			FindFirst(
				db.UserProfile.UserID.Equals(user.ID),
			).
			Exec(ctx)
		if err != nil {
			log.Printf("Worker %d: unable to load user profile for website %s: %v", id, site.Website.URL, err)
		}
		fName, _ := userProfile.FirstName()
		lName, _ := userProfile.LastName()
		bio, _ := userProfile.Bio()
		phone, _ := userProfile.Phone()
		log.Printf("User Name: %s %s", fName, lName)
		log.Printf("User Bio: %s", bio)
		log.Printf("User Phone: %s", phone)

		if err != nil {
			newStatus = db.WebsiteStatusDown
			respMessage = err.Error()
			log.Printf("Worker %d: ERROR checking %s - %v", id, site.Website.URL, err)
		} else {
			respMessage = res.Status
			statusCode = res.StatusCode
			respMessage = res.Status
			defer res.Body.Close()
			log.Printf("Worker %d: Received response from %s - Status: %s, Latency: %v\n", id, site.Website.URL, respMessage, latency)
			if res.StatusCode >= 400 || err != nil {
				newStatus = db.WebsiteStatusDown
			}

		}

		log.Printf("Worker %d: Checked website %s - Status: %s, Latency: %v\n, status code: %d", id, site.Website.URL, newStatus, latency, statusCode)

		if newStatus != site.Website.Status {
			_, err := prisma.WebsiteTicks.CreateOne(
				db.WebsiteTicks.Website.Link(db.Website.ID.Equals(site.ID)),
				db.WebsiteTicks.Region.Link(db.Region.RegionID.Equals(site.RegionID)),
				db.WebsiteTicks.UpStatus.Set(newStatus),
				db.WebsiteTicks.Latency.Set(int(latency.Milliseconds())),
			).Exec(ctx)

			if err == nil {
				_, updateErr := prisma.Website.FindUnique(db.Website.ID.Equals(site.ID)).
					Update(db.Website.Status.Set(newStatus)).
					Exec(ctx)
				if updateErr != nil {
					log.Printf("Worker %d: unable to update website %s status: %v", id, site.Website.URL, updateErr)
				}
			}
			if newStatus == db.WebsiteStatusDown {
				service, err := awsservice.NewNotificationService(ctx)
				if err != nil {
					// utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
					return
				}
				userEmail := user.Email
				log.Print(userProfile)

				log.Printf("Worker %d: Sending notification to %s email, username : %s  for website %s - Status: %s , Subject :%s, file name : %s, URL: %s\n", id, userEmail, fName, site.Website.WebsiteName, newStatus, "CRITICAL ALERT: Website is Down", "anomaly-alert.html", site.Website.URL)
				isSent := service.SendTemplateEmail(ctx, userEmail, "CRITICAL ALERT: Website is Down", "anomaly-alert.html", map[string]any{
					"UserName":       fName,
					"WebsiteName":    site.Website.WebsiteName,
					"WebsiteURL":     site.Website.URL,
					"Status":         newStatus,
					"UserProfileURL": fmt.Sprintf("%s/%s?regionId=%s", baseUrl, site.Website.ID, site.RegionID),
					"DetectedAt":     time.Now().UTC().Format("2006-01-02 15:04:05"),
				})
				if isSent != nil {
					log.Printf("Unable to send email notification for website %s: %v", site.Website.URL, isSent)
				}
				fmt.Printf("Worker %d: Website %s returned status code %d\n", id, site.Website.URL, res.StatusCode)
				fmt.Printf("Send Push notfication")
			} else if newStatus == db.WebsiteStatusUp {
				// Optional: Send a "Your website is back up!" email currently not sending to avoid spamming users with notifications if their website is flapping between up and down status
			}
		}
		log.Printf("Worker %d: Finished checking website %s - Status: %s, Latency: %v\n", id, site.Website.WebsiteName, newStatus, latency)
		results <- ResultResponse{
			WebsiteID: site.ID,
			Status:    newStatus,
			Latency:   latency,
			Response:  respMessage,
		}
	}
}
