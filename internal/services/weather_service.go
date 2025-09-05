package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// WeatherService provides weather information from OpenWeatherMap API
type WeatherService struct {
	httpClient *HTTPClient
	apiKey     string
	logger     Logger
}

// WeatherResponse represents weather data from OpenWeatherMap
type WeatherResponse struct {
	Location    string  `json:"name"`
	Country     string  `json:"country"`
	Temperature float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	Pressure    int     `json:"pressure"`
	WindSpeed   float64 `json:"wind_speed"`
	WindDeg     int     `json:"wind_deg"`
	Clouds      int     `json:"clouds"`
	Visibility  int     `json:"visibility"`
	Icon        string  `json:"icon"`
}

// OpenWeatherResponse represents the full API response from OpenWeatherMap
type OpenWeatherResponse struct {
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
	} `json:"sys"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Visibility int `json:"visibility"`
}

// NewWeatherService creates a new Weather service
func NewWeatherService(logger Logger) *WeatherService {
	httpClient := NewHTTPClient(30*time.Second, logger)
	apiKey := os.Getenv("WEATHER_API_KEY")
	
	// If no API key is set, use a demo mode with mock data
	if apiKey == "" {
		logger.Printf("⚠️ WEATHER_API_KEY not set, using demo mode")
	}
	
	return &WeatherService{
		httpClient: httpClient,
		apiKey:     apiKey,
		logger:     logger,
	}
}

// GetWeather fetches weather information for a city
func (w *WeatherService) GetWeather(ctx context.Context, city string) (*WeatherResponse, error) {
	// If no API key, return demo data
	if w.apiKey == "" {
		return w.getDemoWeather(city), nil
	}
	
	// Clean city name
	city = strings.TrimSpace(city)
	if city == "" {
		return nil, fmt.Errorf("shahar nomi kiritilmagan")
	}
	
	// Build API URL
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=en", city, w.apiKey)
	
	var apiResp OpenWeatherResponse
	err := w.httpClient.GetJSON(ctx, url, nil, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("ob-havo ma'lumotlarini olishda xatolik: %w", err)
	}
	
	weather := &WeatherResponse{
		Location:    apiResp.Name,
		Country:     apiResp.Sys.Country,
		Temperature: apiResp.Main.Temp,
		FeelsLike:   apiResp.Main.FeelsLike,
		Humidity:    apiResp.Main.Humidity,
		Pressure:    apiResp.Main.Pressure,
		WindSpeed:   apiResp.Wind.Speed,
		WindDeg:     apiResp.Wind.Deg,
		Clouds:      apiResp.Clouds.All,
		Visibility:  apiResp.Visibility,
	}
	
	if len(apiResp.Weather) > 0 {
		weather.Description = apiResp.Weather[0].Description
		weather.Icon = apiResp.Weather[0].Icon
	}
	
	w.logger.Printf("🌤 Weather data retrieved for %s: %.1f°C", city, weather.Temperature)
	return weather, nil
}

// getDemoWeather returns demo weather data when API key is not available
func (w *WeatherService) getDemoWeather(city string) *WeatherResponse {
	// Simulate different weather for different cities
	cityLower := strings.ToLower(city)
	
	switch {
	case strings.Contains(cityLower, "tashkent") || strings.Contains(cityLower, "toshkent"):
		return &WeatherResponse{
			Location:    "Toshkent",
			Country:     "UZ",
			Temperature: 22.5,
			FeelsLike:   25.0,
			Description: "Clear sky",
			Humidity:    45,
			Pressure:    1013,
			WindSpeed:   3.2,
			WindDeg:     180,
			Clouds:      10,
			Visibility:  10000,
			Icon:        "01d",
		}
	case strings.Contains(cityLower, "samarkand") || strings.Contains(cityLower, "samarqand"):
		return &WeatherResponse{
			Location:    "Samarqand",
			Country:     "UZ",
			Temperature: 20.0,
			FeelsLike:   22.0,
			Description: "Few clouds",
			Humidity:    55,
			Pressure:    1010,
			WindSpeed:   2.8,
			WindDeg:     90,
			Clouds:      25,
			Visibility:  10000,
			Icon:        "02d",
		}
	default:
		return &WeatherResponse{
			Location:    city,
			Country:     "DEMO",
			Temperature: 18.0,
			FeelsLike:   20.0,
			Description: "Demo weather",
			Humidity:    60,
			Pressure:    1015,
			WindSpeed:   2.0,
			WindDeg:     45,
			Clouds:      30,
			Visibility:  8000,
			Icon:        "03d",
		}
	}
}

// FormatWeather formats weather info for Telegram message
func (w *WeatherService) FormatWeather(weather *WeatherResponse) string {
	emoji := w.getWeatherEmoji(weather.Icon)
	
	windDirection := w.getWindDirection(weather.WindDeg)
	
	message := fmt.Sprintf(`%s <b>%s ob-havo ma'lumoti</b>

🌡 <b>Harorat:</b> %.1f°C (his qilinishi: %.1f°C)
📝 <b>Holati:</b> %s
💧 <b>Namlik:</b> %d%%
🌪 <b>Bosim:</b> %d hPa
💨 <b>Shamol:</b> %.1f m/s (%s)
☁️ <b>Bulutlar:</b> %d%%`,
		emoji,
		weather.Location,
		weather.Temperature,
		weather.FeelsLike,
		w.translateDescription(weather.Description),
		weather.Humidity,
		weather.Pressure,
		weather.WindSpeed,
		windDirection,
		weather.Clouds)
	
	if weather.Country != "DEMO" && weather.Visibility > 0 {
		message += fmt.Sprintf("\n👁 <b>Ko'rinish:</b> %.1f km", float64(weather.Visibility)/1000)
	}
	
	message += fmt.Sprintf("\n\n🕐 <b>Ma'lumot yangilangan:</b> %s", time.Now().Format("15:04"))
	
	if weather.Country == "DEMO" {
		message += "\n\n💡 <i>Demo rejim: WEATHER_API_KEY sozlamasi kerak</i>"
	}
	
	return message
}

// getWeatherEmoji returns appropriate emoji for weather icon
func (w *WeatherService) getWeatherEmoji(icon string) string {
	switch icon {
	case "01d", "01n": // clear sky
		return "☀️"
	case "02d", "02n": // few clouds
		return "🌤"
	case "03d", "03n": // scattered clouds
		return "⛅"
	case "04d", "04n": // broken clouds
		return "☁️"
	case "09d", "09n": // shower rain
		return "🌧"
	case "10d", "10n": // rain
		return "🌦"
	case "11d", "11n": // thunderstorm
		return "⛈"
	case "13d", "13n": // snow
		return "❄️"
	case "50d", "50n": // mist
		return "🌫"
	default:
		return "🌤"
	}
}

// getWindDirection returns wind direction in Uzbek
func (w *WeatherService) getWindDirection(deg int) string {
	directions := []string{
		"Shimol", "Shimol-sharq", "Sharq", "Janub-sharq",
		"Janub", "Janub-g'arb", "G'arb", "Shimol-g'arb",
	}
	
	index := int((float64(deg)+22.5)/45) % 8
	return directions[index]
}

// translateDescription translates weather description to Uzbek
func (w *WeatherService) translateDescription(desc string) string {
	translations := map[string]string{
		"clear sky":           "ochiq osmon",
		"few clouds":          "kam bulutli",
		"scattered clouds":    "tarqoq bulutlar",
		"broken clouds":       "ko'p bulutlar",
		"shower rain":         "yomg'ir",
		"rain":                "yomg'ir yog'moqda",
		"thunderstorm":        "momaqaldiroq",
		"snow":                "qor",
		"mist":                "tuman",
		"overcast clouds":     "bulutli",
		"light rain":          "engil yomg'ir",
		"moderate rain":       "o'rtacha yomg'ir",
		"heavy rain":          "kuchli yomg'ir",
		"demo weather":        "demo ob-havo",
	}
	
	if translation, exists := translations[strings.ToLower(desc)]; exists {
		return translation
	}
	
	return desc // Return original if no translation found
}