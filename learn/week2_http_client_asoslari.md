# Week 2: HTTP Client Asoslari (HTTP Client Fundamentals)

## 🎯 O'rganish Maqsadlari

Bu darsda siz HTTP client yaratish va tashqi API'lar bilan ishlashni o'rganasiz.

## 📚 Asosiy Tushunchalar

### 1. HTTP Client nima?

HTTP Client - bu boshqa serverlar bilan HTTP protokoli orqali aloqa qilish uchun ishlatiladigan vosita.

**Go tilida HTTP Client:**
```go
type HTTPClient struct {
    client  *http.Client    // Go'ning standart HTTP clienti
    logger  Logger         // Xatolarni va loglarni yozish uchun
    baseURL string         // Asosiy URL manzil
}
```

### 2. Asosiy HTTP Metodlari

- **GET**: Ma'lumot olish uchun
- **POST**: Ma'lumot yuborish uchun  
- **PUT**: Ma'lumotni yangilash uchun
- **DELETE**: Ma'lumotni o'chirish uchun

### 3. HTTP Client Yaratish

```go
func NewHTTPClient(timeout time.Duration, logger Logger) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Timeout: timeout,        // Kutish vaqti
        },
        logger: logger,
    }
}
```

**Muhim tushunchalar:**
- `Timeout` - Server javob bermaganida qancha kutish kerakligi
- `Logger` - Barcha amallarni yozib borish uchun

### 4. HTTP So'rov Yuborish

```go
func (h *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
    // So'rov yaratish
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("so'rov yaratishda xatolik: %w", err)
    }

    // Headerlar qo'shish
    for key, value := range headers {
        req.Header.Set(key, value)
    }
    
    // So'rovni yuborish
    resp, err := h.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("so'rov yuborishda xatolik: %w", err)
    }
    defer resp.Body.Close()

    // Javobni o'qish
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("javobni o'qishda xatolik: %w", err)
    }

    return &HTTPResponse{
        StatusCode: resp.StatusCode,
        Body:       body,
        Headers:    resp.Header,
    }, nil
}
```

### 5. JSON bilan ishlash

```go
func (h *HTTPClient) GetJSON(ctx context.Context, url string, headers map[string]string, target interface{}) error {
    resp, err := h.Get(ctx, url, headers)
    if err != nil {
        return err
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP xatolik: %d", resp.StatusCode)
    }

    // JSON ni Go struct'iga aylantirish
    if err := json.Unmarshal(resp.Body, target); err != nil {
        return fmt.Errorf("JSON ni parsing qilishda xatolik: %w", err)
    }

    return nil
}
```

## 🔧 Amaliy Misol: GitHub API

```go
type GitHubRepository struct {
    Name        string `json:"name"`
    FullName    string `json:"full_name"`
    Description string `json:"description"`
    Stars       int    `json:"stargazers_count"`
    Language    string `json:"language"`
}

func (h *HTTPClient) GetRepository(owner, repo string) (*GitHubRepository, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
    
    var repository GitHubRepository
    err := h.GetJSON(context.Background(), url, nil, &repository)
    if err != nil {
        return nil, err
    }
    
    return &repository, nil
}
```

## 💡 Muhim Qoidalar

1. **Har doim timeout belgilang** - Server javob bermasligi mumkin
2. **Context ishlatgung** - So'rovlarni bekor qilish uchun
3. **Xatolarni to'g'ri ishlang** - Barcha xatolarni wraplash kerak
4. **Loglarni yozing** - Debug qilish uchun muhim
5. **Headers qo'shing** - Ba'zi API'lar User-Agent talab qiladi

## 🎯 Keyingi Qadamlar

- GitHub API integration
- Weather API integration  
- Error handling va retry logic
- Rate limiting (so'rov cheklash)

## 📝 Vazifa

1. HTTP Client yarating
2. GitHub repository ma'lumotlarini oling
3. Xatolarni to'g'ri ishlang
4. Loglarni qo'shing

Bu dars orqali siz HTTP client yaratishning asoslarini o'rgandingiz va tashqi API'lar bilan ishlashni boshladingiz.