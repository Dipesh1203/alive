package services

import (
	"backend/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
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

func StartMonitoring(ctx context.Context, prisma *db.PrismaClient) {
	taskChannel := make(chan Site, 10)
	resultsChannel := make(chan ResultResponse, 10)
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
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

	ticker := time.NewTicker(10 * time.Second)
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
		Timeout: 5 * time.Second,
	}

	for site := range jobs {
		log.Printf("Worker %d: Checking website %s in region %s\n", id, site.Website.WebsiteName, site.RegionID)
		startTime := time.Now()
		res, err := client.Get(site.Website.URL)
		latency := time.Since(startTime)

		newStatus := db.WebsiteStatusUp
		respMessage := "200 OK"
		statusCode := 0

		if err != nil {
			newStatus = db.WebsiteStatusDown
			respMessage = err.Error()
			log.Printf("Worker %d: ERROR checking %s - %v", id, site.Website.URL, err)
		} else {
			respMessage = res.Status
			statusCode = res.StatusCode
			respMessage = res.Status
			defer res.Body.Close()
			if res.StatusCode >= 400 {
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
				prisma.Website.FindUnique(db.Website.ID.Equals(site.ID)).
					Update(db.Website.Status.Set(newStatus)).
					Exec(ctx)
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
