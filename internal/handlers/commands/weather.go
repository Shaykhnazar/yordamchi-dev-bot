package commands

import (
	"context"
	"fmt"
	"strings"

	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/services"
)

// WeatherCommand handles /weather command
type WeatherCommand struct {
	weatherService *services.WeatherService
	logger         domain.Logger
}

// NewWeatherCommand creates a new weather command handler
func NewWeatherCommand(weatherService *services.WeatherService, logger domain.Logger) *WeatherCommand {
	return &WeatherCommand{
		weatherService: weatherService,
		logger:         logger,
	}
}

// Handle processes the weather commands
func (h *WeatherCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	// Parse city from command
	parts := strings.Fields(cmd.Text)
	if len(parts) < 2 {
		return &domain.Response{
			Text:      "🌤️ Shahar nomini kiriting!\n\nMisol: `/weather Tashkent`",
			ParseMode: "Markdown",
		}, nil
	}

	city := strings.Join(parts[1:], " ")
	
	// Get weather information
	weather, err := h.weatherService.GetWeather(ctx, city)
	if err != nil {
		h.logger.Error("Failed to get weather", "city", city, "error", err)
		return &domain.Response{
			Text:      fmt.Sprintf("❌ %s shahri uchun ob-havo ma'lumotini olishda xatolik", city),
			ParseMode: "Markdown",
		}, nil
	}

	// Format response based on command language
	var message string
	command := strings.ToLower(parts[0])
	if command == "/weather" {
		// English response
		message = fmt.Sprintf(
			"🌤️ **Weather in %s**\n\n"+
			"🌡️ **Temperature:** %.1f°C\n"+
			"💧 **Humidity:** %d%%\n"+
			"💨 **Wind:** %.1f km/h\n"+
			"📊 **Pressure:** %d hPa\n"+
			"☁️ **Condition:** %s",
			weather.Location,
			weather.Temperature,
			weather.Humidity,
			weather.WindSpeed,
			weather.Pressure,
			weather.Description,
		)
	}

	h.logger.Info("Weather command processed", 
		"user_id", cmd.User.TelegramID,
		"city", city,
		"command", command)

	return &domain.Response{
		Text:      message,
		ParseMode: "Markdown",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *WeatherCommand) CanHandle(command string) bool {
	cmd := strings.ToLower(strings.TrimSpace(command))
	return strings.HasPrefix(cmd, "/weather")
}

// Description returns the command description
func (h *WeatherCommand) Description() string {
	return "Ob-havo ma'lumoti"
}

// Usage returns the command usage
func (h *WeatherCommand) Usage() string {
	return "/weather city - Ob-havo ma'lumoti"
}