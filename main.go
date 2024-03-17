package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water       int    `json:"water"`
	Wind        int    `json:"wind"`
	WaterStatus string `json:"water_status"`
	WindStatus  string `json:"wind_status"`
}

type Data struct {
	Status Status `json:"status"`
}

func main() {
	go weatherUpdate()

	http.HandleFunc("/data", handleData)

	http.Handle("/", http.FileServer(http.Dir(".")))

	http.ListenAndServe(":8080", nil)
}

func handleData(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("weather.json")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func weatherUpdate() {
	for {
		// acak value water dan wind
		randomWater := rand.Intn(100) + 1
		randomWind := rand.Intn(100) + 1

		var waterStatus string
		if randomWater < 5 {
			waterStatus = "Aman"
		} else if randomWater >= 6 && randomWater <= 8 {
			waterStatus = "Siaga"
		} else {
			waterStatus = "Bahaya"
		}

		var windStatus string
		if randomWind < 6 {
			windStatus = "Aman"
		} else if randomWind >= 7 && randomWind <= 15 {
			windStatus = "Siaga"
		} else {
			windStatus = "Bahaya"
		}

		// data untuk JSON
		data := Data{
			Status: Status{
				Water:       randomWater,
				Wind:        randomWind,
				WaterStatus: waterStatus,
				WindStatus:  windStatus,
			},
		}

		// Konversi data menjadi JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Tulis JSON data ke file weather.json
		err = ioutil.WriteFile("weather.json", jsonData, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Printf("Data updated: water=%d (%s), wind=%d (%s)\n", randomWater, waterStatus, randomWind, windStatus)

		// Tunggu selama 15 detik untuk update kembali
		time.Sleep(15 * time.Second)
	}
}
