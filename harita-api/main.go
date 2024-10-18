package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	apiKey  = ""
	baseURL = "https://maps.googleapis.com/maps/api/distancematrix/json"
)

type DistanceMatrixResponse struct {
	Status string `json:"status"`
	Rows   []struct {
		Elements []struct {
			Status   string `json:"status"`
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
		} `json:"elements"`
	} `json:"rows"`
	ErrorMessage string `json:"error_message"`
}

func getDistance(baslangic, bitis string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("origins", baslangic)
	q.Set("destinations", bitis)
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	fmt.Printf("İstek URL'si: %s\n", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("API Yanıtı: %s\n", string(body))

	var result DistanceMatrixResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Status != "OK" {
		return "", fmt.Errorf("API hatası: %s", result.ErrorMessage)
	}

	if len(result.Rows) > 0 && len(result.Rows[0].Elements) > 0 {
		element := result.Rows[0].Elements[0]
		if element.Status == "OK" {
			return element.Distance.Text, nil
		}
		return "", fmt.Errorf("Element durumu: %s", element.Status)
	}

	return "", fmt.Errorf("mesafe bulunamadı")
}

func main() {
	baslangic := "Bosna Hersek, Mahallesi, Büyük Hizmet Cd. No:37 D:e, 42250 Selçuklu/Konya"
	bitis := "Bosna Hersek, Osmanlı Cd. No:37/A, 42250 Selçuklu/Konya"

	mesafe, err := getDistance(baslangic, bitis)
	if err != nil {
		log.Fatalf("Hata: %v", err)
	}

	fmt.Printf("%s ile %s arasındaki mesafe: %s\n", baslangic, bitis, mesafe)
}
