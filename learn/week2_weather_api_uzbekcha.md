# Week 2: Weather API bilan Ishlash (Weather API Integration)

## 🎯 O'rganish Maqsadlari

Bu darsda siz ob-havo ma'lumotlarini olish uchun tashqi API bilan ishlashni o'rganasiz.

## 📚 Asosiy Tushunchalar

### 1. Weather API nima?

Weather API - bu ob-havo ma'lumotlarini real vaqtda olish uchun ishlatiladigan tashqi service.

**OpenWeatherMap API xususiyatlari:**
- Bepul foydalanish uchun 1000 ta so'rov kuniga
- JSON formatda ma'lumot qaytaradi
- 200,000+ shaharlar ma'lumotlari
- API Key talab qiladi

### 2. Weather Service Yaratish

```go
type WeatherService struct {
    httpClient *HTTPClient    // HTTP so'rovlar uchun
    apiKey     string        // API kaliti
    logger     Logger        // Loglar uchun
}

func NewWeatherService(logger Logger) *WeatherService {
    httpClient := NewHTTPClient(30*time.Second, logger)
    apiKey := os.Getenv("WEATHER_API_KEY")  // Environment'dan oling
    
    return &WeatherService{
        httpClient: httpClient,
        apiKey:     apiKey,
        logger:     logger,
    }
}
```

**Muhim tushunchalar:**
- `os.Getenv()` - Environment variable'lardan ma'lumot olish
- API kalitini kodga yozmang, environment'dan oling
- Demo rejimini qo'shish yaxshi amaliyot

### 3. API Response Struct'lari

```go
// Bizning ichki formatimiz
type WeatherResponse struct {
    Location    string  `json:"name"`
    Temperature float64 `json:"temp"`
    Description string  `json:"description"`
    Humidity    int     `json:"humidity"`
    WindSpeed   float64 `json:"wind_speed"`
}

// OpenWeatherMap API'ning formatiga mos struct
type OpenWeatherResponse struct {
    Name string `json:"name"`
    Main struct {
        Temp     float64 `json:"temp"`
        Humidity int     `json:"humidity"`
    } `json:"main"`
    Weather []struct {
        Description string `json:"description"`
    } `json:"weather"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
}
```

### 4. Nested Struct'lar

Go'da struct ichida struct ishlatish mumkin:

```go
type OpenWeatherResponse struct {
    Main struct {                    // Ichki struct
        Temp     float64 `json:"temp"`
        Humidity int     `json:"humidity"`
    } `json:"main"`
    
    Weather []struct {               // Struct'lar array'i
        Description string `json:"description"`
    } `json:"weather"`
}
```

**Nested struct'lar bilan ishlash:**
```go
// Ma'lumotni olish
temp := apiResp.Main.Temp
humidity := apiResp.Main.Humidity

// Array'dan birinchi element
if len(apiResp.Weather) > 0 {
    description := apiResp.Weather[0].Description
}
```

### 5. Environment Variable'lar

```go
import "os"

// Environment variable'ni olish
apiKey := os.Getenv("WEATHER_API_KEY")

// Tekshirish
if apiKey == "" {
    logger.Printf("⚠️ WEATHER_API_KEY not set, using demo mode")
}
```

**Environment variable'lar nima uchun kerak:**
- API kalitlarini kodda saqlamaslik uchun
- Turli muhitlar uchun turli sozlamalar
- Xavfsizlik uchun muhim

### 6. Demo Rejimi

API kaliti bo'lmaganda demo ma'lumot qaytarish:

```go
func (w *WeatherService) getDemoWeather(city string) *WeatherResponse {
    cityLower := strings.ToLower(city)
    
    switch {
    case strings.Contains(cityLower, "tashkent"):
        return &WeatherResponse{
            Location:    "Toshkent",
            Temperature: 22.5,
            Description: "Clear sky",
            Humidity:    45,
        }
    default:
        return &WeatherResponse{
            Location:    city,
            Temperature: 18.0,
            Description: "Demo weather",
        }
    }
}
```

### 7. String'larni Tarjima Qilish

```go
func (w *WeatherService) translateDescription(desc string) string {
    translations := map[string]string{
        "clear sky":      "ochiq osmon",
        "few clouds":     "kam bulutli",
        "rain":           "yomg'ir yog'moqda",
        "snow":           "qor",
    }
    
    if translation, exists := translations[strings.ToLower(desc)]; exists {
        return translation
    }
    
    return desc // Asl matnni qaytaring
}
```

**Map bilan ishlash:**
- `map[string]string` - string'dan string'ga map
- `translations[key]` - qiymat olish
- `value, exists := map[key]` - mavjudligini tekshirish

### 8. Ma'lumotlarni Formatlash

```go
func (w *WeatherService) FormatWeather(weather *WeatherResponse) string {
    emoji := w.getWeatherEmoji(weather.Icon)
    
    message := fmt.Sprintf(`%s <b>%s ob-havo ma'lumoti</b>

🌡 <b>Harorat:</b> %.1f°C
📝 <b>Holati:</b> %s
💧 <b>Namlik:</b> %d%%
💨 <b>Shamol:</b> %.1f m/s`,
        emoji,
        weather.Location,
        weather.Temperature,
        weather.Description,
        weather.Humidity,
        weather.WindSpeed)
    
    return message
}
```

### 9. Emoji'larni Tanlash

```go
func (w *WeatherService) getWeatherEmoji(icon string) string {
    switch icon {
    case "01d", "01n": // clear sky
        return "☀️"
    case "02d", "02n": // few clouds
        return "🌤"
    case "09d", "09n": // rain
        return "🌧"
    default:
        return "🌤"
    }
}
```

## 🔧 Switch Statement

Go'da `switch` statement bilan qulayroq shart tekshirish:

```go
switch icon {
case "01d", "01n":           // Bir nechta qiymat
    return "☀️"
case "02d":                  // Bitta qiymat
    return "🌤"
default:                     // Default holat
    return "🌤"
}

// yoki boolean bilan
switch {
case temp > 30:
    return "juda issiq"
case temp > 20:
    return "iliq"
default:
    return "sovuq"
}
```

## 💡 Muhim Qoidalar

1. **API kalitlarini yashiring** - Hech qachon kodga qo'ymang
2. **Demo rejimini qo'shing** - API bo'lmaganda ham ishlashi kerak
3. **Ma'lumotlarni tarjima qiling** - Foydalanuvchi tiliga mos
4. **Xatolarni to'g'ri ishlang** - API ishlamasligi mumkin
5. **Timeout qo'ying** - Uzoq kutmaslik uchun

## 🎯 Amaliy Misol: Bot Command

```go
// bot.go da
case strings.HasPrefix(msg.Text, "/weather "):
    parts := strings.Fields(msg.Text)
    if len(parts) < 2 {
        b.sendMessage(chatID, "❌ Shahar nomini kiriting: /weather Tashkent")
        return
    }
    
    city := strings.Join(parts[1:], " ")  // Ko'p so'zli shahar nomlari uchun
    
    weather := services.NewWeatherService(logger)
    data, err := weather.GetWeather(context.Background(), city)
    if err != nil {
        b.sendMessage(chatID, "❌ Ob-havo ma'lumotini olishda xatolik: "+err.Error())
        return
    }
    
    message := weather.FormatWeather(data)
    b.sendMessage(chatID, message)
```

## 📝 Environment Setup

`.env` faylga qo'shing:
```bash
WEATHER_API_KEY=your_openweather_api_key_here
```

API kalitini olish:
1. https://openweathermap.org/ ga boring
2. Ro'yxatdan o'ting
3. API key oling
4. .env faylga qo'shing

## 🚀 Keyingi Qadamlar

- 5 kunlik ob-havo bashorati
- Shahar qidirish funksiyasi
- Ob-havo ogohlantirish
- Grafik va rasmlar

Bu dars orqali siz weather API bilan ishlashni va ma'lumotlarni formatlashni o'rgandingiz!