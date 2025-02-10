package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/MarSultanius/bookings/pkg/models"
	"github.com/MarSultanius/bookings/pkg/config"
	"github.com/MarSultanius/bookings/pkg/render"
)

// API Key for ipinfo.io (получите свой ключ на ipinfo.io)
const ipInfoAPI = "https://ipinfo.io/%s?token=0f1c3905682fb8"

// Repo  the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// Получение геолокации по IP
	geolocationData, err := getGeolocationByIP(remoteIP)
	if err != nil {
		stringMap["geolocation"] = "Unable to retrieve geolocation"
	} else {

		fmt.Println("Geolocation Data:")
		fmt.Printf("Country: %s, City: %s, Location: %s\n", geolocationData.Country, geolocationData.City, geolocationData.Loc)

		stringMap["geolocation"] = fmt.Sprintf("Country: %s, City: %s, Location: %s", geolocationData.Country, geolocationData.City, geolocationData.Loc)
	}

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Структура для геолокации
type GeolocationData struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Loc     string `json:"loc"` // Coordinates: "latitude,longitude"
}

// Функция для получения геолокации по IP
func getGeolocationByIP(ip string) (*GeolocationData, error) {
	url := fmt.Sprintf(ipInfoAPI, strings.TrimSpace(ip))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data from ipinfo.io, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geolocation GeolocationData
	err = json.Unmarshal(body, &geolocation)
	if err != nil {
		return nil, err
	}

	return &geolocation, nil
}