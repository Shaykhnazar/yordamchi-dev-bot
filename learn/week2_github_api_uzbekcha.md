# Week 2: GitHub API bilan Ishlash (GitHub API Integration)

## 🎯 O'rganish Maqsadlari

Bu darsda siz GitHub API'si bilan ishlashni, repository va foydalanuvchi ma'lumotlarini olishni o'rganasiz.

## 📚 Asosiy Tushunchalar

### 1. GitHub API nima?

GitHub API - bu GitHub'dagi repository'lar, foydalanuvchilar va boshqa ma'lumotlarga dasturiy yo'l bilan kirish imkonini beruvchi vosita.

**GitHub API'ning asosiy xususiyatlari:**
- RESTful API (HTTP metodlarini ishlatadi)
- JSON formatida javob qaytaradi
- Bepul foydalanish uchun soatiga 60 ta so'rov
- Autentifikatsiya bilan ko'proq so'rov mumkin

### 2. GitHub Service Yaratish

```go
type GitHubService struct {
    httpClient *HTTPClient    // HTTP so'rovlar uchun
    logger     Logger        // Loglarni yozish uchun
}

func NewGitHubService(logger Logger) *GitHubService {
    httpClient := NewHTTPClient(30*time.Second, logger)
    
    return &GitHubService{
        httpClient: httpClient,
        logger:     logger,
    }
}
```

**Muhim tushunchalar:**
- `HTTPClient` - oldin yaratgan HTTP client'imizni ishlatamiz
- `30*time.Second` - GitHub API uchun yetarli timeout
- `Logger` - barcha amallarni logga yozish

### 3. Repository Ma'lumotlarini Olish

```go
type GitHubRepository struct {
    Name            string `json:"name"`           // Repository nomi
    FullName        string `json:"full_name"`      // To'liq nom (owner/repo)
    Description     string `json:"description"`    // Tavsif
    Stars           int    `json:"stargazers_count"` // Yulduzlar soni
    Forks           int    `json:"forks_count"`     // Forklar soni
    Language        string `json:"language"`        // Asosiy dasturlash tili
    URL             string `json:"html_url"`        // GitHub sahifa URL'i
}

func (g *GitHubService) GetRepository(ctx context.Context, owner, repo string) (*GitHubRepository, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
    
    var repository GitHubRepository
    err := g.httpClient.GetJSON(ctx, url, nil, &repository)
    if err != nil {
        return nil, fmt.Errorf("GitHub repository ma'lumotlarini olishda xatolik: %w", err)
    }
    
    return &repository, nil
}
```

### 4. Foydalanuvchi Ma'lumotlarini Olish

```go
type GitHubUser struct {
    Login       string `json:"login"`         // Foydalanuvchi nomi
    Name        string `json:"name"`          // To'liq ism
    Company     string `json:"company"`       // Kompaniya
    Location    string `json:"location"`      // Joylashuv
    Bio         string `json:"bio"`           // Bio ma'lumoti
    PublicRepos int    `json:"public_repos"`  // Ochiq repository'lar soni
    Followers   int    `json:"followers"`     // Obunachilar soni
    Following   int    `json:"following"`     // Obunalar soni
}

func (g *GitHubService) GetUser(ctx context.Context, username string) (*GitHubUser, error) {
    url := fmt.Sprintf("https://api.github.com/users/%s", username)
    
    var user GitHubUser
    err := g.httpClient.GetJSON(ctx, url, nil, &user)
    if err != nil {
        return nil, fmt.Errorf("foydalanuvchi ma'lumotlarini olishda xatolik: %w", err)
    }
    
    return &user, nil
}
```

### 5. Ma'lumotlarni Formatlash

```go
func (g *GitHubService) FormatRepository(repo *GitHubRepository) string {
    description := repo.Description
    if description == "" {
        description = "Tavsif mavjud emas"
    }
    
    return fmt.Sprintf(`📦 <b>%s</b>

📝 <b>Tavsif:</b> %s
⭐ <b>Yulduzlar:</b> %d
🍴 <b>Forklar:</b> %d
💻 <b>Til:</b> %s

🔗 <b>Havola:</b> <a href="%s">%s</a>`,
        repo.FullName,
        description,
        repo.Stars,
        repo.Forks,
        repo.Language,
        repo.URL,
        repo.URL)
}
```

## 🔧 JSON Taglari (JSON Tags)

JSON taglari Go struct'larini JSON'ga aylantirish uchun ishlatiladi:

```go
type Repository struct {
    Name  string `json:"name"`        // JSON'da "name" kaliti
    Stars int    `json:"stargazers_count"` // JSON'da "stargazers_count"
}
```

**JSON Tag qoidalari:**
- `json:"field_name"` - JSON kalitini belgilaydi
- `json:"-"` - bu fieldni JSON'da yashiradi
- `json:",omitempty"` - bo'sh bo'lsa JSON'ga qo'shmaydi

## 🌐 API Endpoint'lari

GitHub API'ning asosiy endpoint'lari:

```go
// Repository ma'lumotlari
"https://api.github.com/repos/{owner}/{repo}"

// Foydalanuvchi ma'lumotlari
"https://api.github.com/users/{username}"

// Repository qidirish
"https://api.github.com/search/repositories?q={query}"

// Foydalanuvchi repository'lari
"https://api.github.com/users/{username}/repos"
```

## 💡 Muhim Qoidalar

1. **Rate Limiting** - GitHub API soatiga 60 ta so'rov cheklaydi
2. **User-Agent** - GitHub API User-Agent header'ini talab qiladi
3. **Error Handling** - API 404, 403 kabi xatolar qaytarishi mumkin
4. **Timeout** - Uzoq kutishdan qochish uchun timeout belgilang
5. **Context** - So'rovlarni bekor qilish uchun context ishlating

## 🎯 Amaliy Misol: Bot Commandasi

```go
// bot.go da
case strings.HasPrefix(text, "/repo "):
    parts := strings.Fields(text)
    if len(parts) != 2 {
        b.sendMessage(chatID, "❌ Format: /repo owner/repository")
        return
    }
    
    repoParts := strings.Split(parts[1], "/")
    if len(repoParts) != 2 {
        b.sendMessage(chatID, "❌ Format: /repo owner/repository")
        return
    }
    
    github := services.NewGitHubService(log.New(os.Stdout, "", log.LstdFlags))
    repo, err := github.GetRepository(context.Background(), repoParts[0], repoParts[1])
    if err != nil {
        b.sendMessage(chatID, "❌ Repository topilmadi: "+err.Error())
        return
    }
    
    message := github.FormatRepository(repo)
    b.sendMessage(chatID, message)
```

## 📝 Vazifalar

1. GitHub service'ni yarating
2. Repository va user ma'lumotlarini oling
3. Ma'lumotlarni chiroyli formatda ko'rsating
4. Xatolarni to'g'ri ishlang
5. Rate limiting'ni hisobga oling

## 🚀 Keyingi Qadamlar

- Repository qidirish funksiyasi
- Foydalanuvchining repository'larini ko'rsatish
- Repository'ning commit tarixini olish
- Issue va Pull Request ma'lumotlari

Bu dars orqali siz GitHub API bilan ishlashni va external API'lardan ma'lumot olishni o'rgandingiz!