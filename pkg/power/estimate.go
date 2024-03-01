package power

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetForecast(apiKey string, url string) (Forecasts, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Forecasts{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Forecasts{}, err
	}
	defer resp.Body.Close()

	var forecasts Forecasts
	err = json.NewDecoder(resp.Body).Decode(&forecasts)
	if err != nil {
		return Forecasts{}, err
	}

	return forecasts, nil
}

func GetTotalDayPowerEstimate(forecasts Forecasts, day time.Time) (float64, error) {
	totalPower := 0.0
	for _, forecast := range forecasts.Forecasts {
		periodEnd, err := time.Parse(time.RFC3339, forecast.PeriodEnd)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return totalPower, err
		}
		if periodEnd.Year() == day.Year() && periodEnd.YearDay() == day.YearDay() {
			totalPower += forecast.PVEstimate * 0.5 // Multiply by 0.5 because data is obtained every 30min
		}
	}

	// The calculated totalPower is in Wh
	return totalPower * 1000, nil
}
