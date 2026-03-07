package handlers

import (
	"backend/db"
	"backend/internal/utils"
	"fmt"
	"net/http"
)

type WebsiteInfo struct {
	Status     string         `json:"status"`
	StatusCode int            `json:"status_code"`
	Proto      string         `json:"protocol"`
	Headers    http.Header    `json:"headers"`
	ContentLen int64          `json:"content_length"`
	Cookies    []*http.Cookie `json:"cookies"`
}

func GetTest(database *db.PrismaClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request", r)
		res, err := http.Get("https://google.com/")
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		defer res.Body.Close()
		info := WebsiteInfo{
			Status:     res.Status,
			StatusCode: res.StatusCode,
			Proto:      res.Proto,
			Headers:    res.Header,
			ContentLen: res.ContentLength,
			Cookies:    res.Cookies(),
		}
		fmt.Println("Result", info)
		utils.WriteJSON(w, http.StatusAccepted, info)
	}
}
