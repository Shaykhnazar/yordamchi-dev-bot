# Week 2 Xulosa: External API Integration Patterns

## 🎯 Week 2'da O'rgangan Narsalar

Bu haftada biz tashqi API'lar bilan ishlash, HTTP client yaratish va ma'lumotlarni formatlash bo'yicha muhim ko'nikmalarni o'rgandik.

## 📚 Asosiy Tushunchalar

### 1. Service Layer Architecture

Week 2'da biz Service Layer pattern'ini o'rgandik:

```go
// Service'lar alohida package'da
package services

// Har bir service o'z vazifasini bajaradi
type GitHubService struct { ... }
type WeatherService struct { ... }
type HTTPClient struct { ... }

// Constructor pattern
func NewGitHubService(logger Logger) *GitHubService
func NewWeatherService(logger Logger) *WeatherService
```

**Nega Service Layer kerak:**
- Kodning tashkil etilishi
- Testlashning osonligi
- Qayta foydalanish imkoniyati
- Business logic'ni ajratish

### 2. HTTP Client Pattern

Barcha service'lar bitta HTTP client'ni ishlatadi:

```go
type HTTPClient struct {
    client  *http.Client
    logger  Logger
    baseURL string
}

// Asosiy metodlar
func (h *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
func (h *HTTPClient) GetJSON(ctx context.Context, url string, headers map[string]string, target interface{}) error
```

**HTTP Client'ning afzalliklari:**
- Timeout'lar bilan xavfsizlik
- Logging barcha so'rovlar uchun
- Header'larni markaziy boshqarish
- Context orqali bekor qilish

### 3. JSON Mapping Patterns

External API'lardan ma'lumot olishda ikki turdagi struct kerak:

```go
// API'ning formatiga mos struct
type OpenWeatherResponse struct {
    Name string `json:"name"`
    Main struct {
        Temp float64 `json:"temp"`
    } `json:"main"`
}

// Bizning ichki formatimiz
type WeatherResponse struct {
    Location    string
    Temperature float64
}

// Mapping
func mapToInternal(apiResp OpenWeatherResponse) WeatherResponse {
    return WeatherResponse{
        Location:    apiResp.Name,
        Temperature: apiResp.Main.Temp,
    }
}
```

### 4. Error Handling Patterns

Har bir service'da xatolarni to'g'ri ishlash:

```go
func (g *GitHubService) GetRepository(ctx context.Context, owner, repo string) (*GitHubRepository, error) {
    // 1. Input validation
    if owner == "" || repo == "" {
        return nil, fmt.Errorf("owner va repo name kiritilishi kerak")
    }
    
    // 2. HTTP so'rov
    err := g.httpClient.GetJSON(ctx, url, nil, &repository)
    if err != nil {
        return nil, fmt.Errorf("GitHub repository ma'lumotlarini olishda xatolik: %w", err)
    }
    
    // 3. Success logging
    g.logger.Printf("📦 GitHub repository retrieved: %s/%s", owner, repo)
    return &repository, nil
}
```

### 5. Context Pattern

Barcha tashqi so'rovlarda context ishlatish:

```go
// Bot'da timeout bilan context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Service'da context'ni keyingi qatlamga uzatish
repository, err := github.GetRepository(ctx, owner, repo)
```

**Context'ning foydasi:**
- Timeout'lar
- So'rovlarni bekor qilish
- Request tracing (kelajakda)

### 6. Environment Configuration

Xavfsiz konfiguratsiya uchun environment variable'lar:

```go
// API kalitlarini environment'dan olish
apiKey := os.Getenv("WEATHER_API_KEY")

// Default qiymat
if apiKey == "" {
    logger.Printf("⚠️ WEATHER_API_KEY not set, using demo mode")
}
```

### 7. Demo Mode Pattern

Development uchun demo rejimi:

```go
func (w *WeatherService) GetWeather(ctx context.Context, city string) (*WeatherResponse, error) {
    // API kaliti bo'lmasa demo qaytarish
    if w.apiKey == "" {
        return w.getDemoWeather(city), nil
    }
    
    // Asl API'ga so'rov
    // ...
}
```

### 8. Formatting Pattern

Telegram uchun HTML formatda ma'lumotni tayyorlash:

```go
func (g *GitHubService) FormatRepository(repo *GitHubRepository) string {
    // Default qiymatlarni tekshirish
    description := repo.Description
    if description == "" {
        description = "Tavsif mavjud emas"
    }
    
    // HTML formatda message
    return fmt.Sprintf(`📦 <b>%s</b>

📝 <b>Tavsif:</b> %s
⭐ <b>Yulduzlar:</b> %d`,
        repo.FullName,
        description,
        repo.Stars)
}
```

## 🎯 Asosiy Pattern'lar

### 1. Service Constructor Pattern

```go
func NewServiceName(dependencies...) *ServiceName {
    return &ServiceName{
        dependency1: dep1,
        dependency2: dep2,
    }
}
```

### 2. Error Wrapping Pattern

```go
if err != nil {
    return nil, fmt.Errorf("aniq tavsif: %w", err)
}
```

### 3. Resource Cleanup Pattern

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel() // Resource'ni tozalash

resp, err := http.Get(url)
defer resp.Body.Close() // Resource'ni tozalash
```

### 4. Validation Pattern

```go
// Input validation
if input == "" {
    return nil, fmt.Errorf("input bo'sh bo'lishi mumkin emas")
}

// API response validation
if len(apiResp.Weather) == 0 {
    return nil, fmt.Errorf("ob-havo ma'lumoti topilmadi")
}
```

## 🔧 Code Organization

Week 2'dan keyin kodimiz quyidagicha tashkil etildi:

```
yordamchi-dev-bot/
├── internal/
│   └── services/
│       ├── http_client.go      # HTTP so'rovlar uchun
│       ├── github_service.go   # GitHub API
│       └── weather_service.go  # Weather API
├── learn/
│   ├── week2_http_client_asoslari.md
│   ├── week2_github_api_uzbekcha.md
│   ├── week2_weather_api_uzbekcha.md
│   └── week2_xulosa_va_patterns.md
└── bot.go                      # Command handler'lar
```

## 💡 Best Practices

### 1. Logging
- Har bir service'da meaningful loglar
- Success va error loglarini ajratish
- Strukturali log format

### 2. Error Messages
- Foydalanuvchi uchun tushunarli xabar
- Developer uchun batafsil log
- Original error'ni wrap qilish

### 3. Timeouts
- Har doim timeout belgilash
- API'ga mos timeout qiymati
- Context bilan bekor qilish

### 4. Security
- API kalitlarini kodga yozmaslik
- Environment variable'lar ishlatish
- Input validation

### 5. Testing
- Demo mode test uchun qulay
- Mock'lar uchun interface'lar
- Error scenario'larni test qilish

## 🚀 Keyingi Qadamlar (Week 3-4)

Week 2'da o'rgangan pattern'lar Week 3-4'da quyidagicha kengaytiriladi:

- **Caching**: API so'rovlarni kesh qilish
- **Rate Limiting**: So'rov cheklash
- **Retry Logic**: Xatolikda qayta urinish  
- **Circuit Breaker**: Service himoya qilish
- **Testing**: Comprehensive test yozish

## 📊 Statistika

Week 2'da qo'shilgan:
- ✅ 3 ta yangi service
- ✅ 3 ta yangi command
- ✅ 400+ qator kod
- ✅ 4 ta o'rganish fayli
- ✅ Production-ready patterns

Week 2 orqali biz professional darajadagi external API integration pattern'larini o'rgandik!