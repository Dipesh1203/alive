package services

import (
	"backend/db"
	"context"
	"fmt"
	"net/http"
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

	for w := 1; w <= 3; w++ {
		go MonitoringWorker(w, taskChannel, resultsChannel, prisma, ctx)
	}

	go func() {
		for res := range resultsChannel {
			fmt.Printf("📊 Result processed for %s: %s (%v)\n", res.WebsiteID, res.Status, res.Latency)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		activeWebsites, err := prisma.Website.FindMany().Exec(ctx)
		if err != nil {
			fmt.Println("Error fetching websites: ", err)
			continue
		}
		for _, site := range activeWebsites {
			regions, err := prisma.Region.FindMany(
				db.Region.Ticks.Some(
					db.WebsiteTicks.WebsiteID.Equals(site.ID),
				),
			).Exec(ctx)
			if err != nil || len(regions) == 0 {
				fmt.Printf("No region found for website %s\n", site.ID)
				continue
			}
			for _, region := range regions {
				taskChannel <- Site{
					ID:       site.ID,
					Website:  site,
					RegionID: region.RegionID,
				}
			}
		}
	}
}

func MonitoringWorker(id int, jobs <-chan Site, results chan<- ResultResponse, prisma *db.PrismaClient, ctx context.Context) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for site := range jobs {
		startTime := time.Now()
		res, err := client.Get(site.Website.URL)
		latency := time.Since(startTime)

		newStatus := db.WebsiteStatusUp
		respMessage := "200 OK"

		if err != nil || res == nil || res.StatusCode >= 400 {
			newStatus = db.WebsiteStatusDown
			if err != nil {
				respMessage = err.Error()
			} else if res != nil {
				respMessage = res.Status
			}
		} else {
			respMessage = res.Status

			res.Body.Close()
		}

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

		results <- ResultResponse{
			WebsiteID: site.ID,
			Status:    newStatus,
			Latency:   latency,
			Response:  respMessage,
		}
	}
}
