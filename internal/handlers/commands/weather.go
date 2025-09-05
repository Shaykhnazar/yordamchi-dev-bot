package commands

import (
	"context"
	"fmt"
	"strings"

	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/services"
)

// WeatherCommand handles /weather and /ob-havo commands
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
			Text:      "🌤️ Shahar nomini kiriting!\n\nMisol: <code>/weather Tashkent</code>",
			ParseMode: "HTML",
		}, nil
	}

	city := strings.Join(parts[1:], " ")
	
	// Get weather information
	weather, err := h.weatherService.GetWeather(ctx, city)
	if err != nil {
		h.logger.Error("Failed to get weather", "city", city, "error", err)
		return &domain.Response{
			Text:      fmt.Sprintf("❌ %s shahri uchun ob-havo ma'lumotini olishda xatolik", city),
			ParseMode: "HTML",
		}, nil
	}

	// Format response based on command language
	var message string
	command := strings.ToLower(parts[0])
	if command == "/ob-havo" {
		// Uzbek response
		message = fmt.Sprintf(
			"🌤️ <b>%s shahrida ob-havo</b>\n\n"+
			"🌡️ <b>Harorat:</b> %.1f°C\n"+
			"💧 <b>Namlik:</b> %d%%\n"+
			"💨 <b>Shamol:</b> %.1f km/soat\n"+
			"📊 <b>Bosim:</b> %d hPa\n"+
			"☁️ <b>Holat:</b> %s",
			weather.Location,
			weather.Temperature,
			weather.Humidity,
			weather.WindSpeed,
			weather.Pressure,
			weather.Description,
		)
	} else {
		// English response
		message = fmt.Sprintf(
			"🌤️ <b>Weather in %s</b>\n\n"+
			"🌡️ <b>Temperature:</b> %.1f°C\n"+
			"💧 <b>Humidity:</b> %d%%\n"+
			"💨 <b>Wind:</b> %.1f km/h\n"+
			"📊 <b>Pressure:</b> %d hPa\n"+
			"☁️ <b>Condition:</b> %s",
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
		ParseMode: "HTML",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *WeatherCommand) CanHandle(command string) bool {
	cmd := strings.ToLower(strings.TrimSpace(command))
	return strings.HasPrefix(cmd, "/weather") || strings.HasPrefix(cmd, "/ob-havo")
}

// Description returns the command description
func (h *WeatherCommand) Description() string {
	return "Ob-havo ma'lumoti"
}

// Usage returns the command usage
func (h *WeatherCommand) Usage() string {
	return "/weather <shahar> yoki /ob-havo <shahar> - Ob-havo ma'lumoti"
}