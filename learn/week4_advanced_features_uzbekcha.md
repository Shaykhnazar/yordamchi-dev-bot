# Week 4: Advanced Features & Production Optimization (Ilg'or Xususiyatlar)

## 🎯 O'rganish Maqsadlari

Bu darsda siz professional darajadagi advanced bot features va production optimization usullarini o'rganasiz.

## 📚 Asosiy Tushunchalar

### 1. Advanced Middleware Architecture

Professional botlarda middleware chain'i optimizatsiya qilish muhim:

```go
// Optimal middleware order (dependencies.go)
router.RegisterMiddleware(loggingMiddleware)     // Log first
router.RegisterMiddleware(metricsMiddleware)     // Metrics collection
router.RegisterMiddleware(validationMiddleware)  // Validate input early
router.RegisterMiddleware(cachingMiddleware)     // Cache before expensive operations
router.RegisterMiddleware(authMiddleware)        // Authentication
router.RegisterMiddleware(activityMiddleware)    // Log activity after auth
router.RegisterMiddleware(rateLimitMiddleware)   // Rate limiting last
```

**Middleware Execution Order:**
- **Logging**: Har bir request'ni log qilish
- **Metrics**: Performance ma'lumotlarini to'plash
- **Validation**: Input validation - tez fail qilish
- **Caching**: Expensive operation'lardan oldin cache check
- **Auth**: User authentication va registration
- **Activity**: Database'ga activity yozish
- **Rate Limit**: Abuse protection - oxirida

### 2. Memory Caching System

In-memory caching expensive API call'larni optimize qilish uchun:

```go
// Memory Cache Implementation (internal/cache/memory_cache.go)
type MemoryCache struct {
    items map[string]*CacheItem
    mutex sync.RWMutex
    ttl   time.Duration
}

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}

func NewMemoryCache(defaultTTL time.Duration) *MemoryCache {
    cache := &MemoryCache{
        items: make(map[string]*CacheItem),
        ttl:   defaultTTL,
    }

    // Background cleanup goroutine
    go cache.startCleanupRoutine()
    return cache
}

func (mc *MemoryCache) Set(key string, value interface{}) {
    mc.SetWithTTL(key, value, mc.ttl)
}

func (mc *MemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()

    mc.items[key] = &CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()

    item, exists := mc.items[key]
    if !exists || time.Now().After(item.ExpiresAt) {
        return nil, false
    }

    return item.Value, true
}
```

**Key Concepts:**
- **sync.RWMutex**: Concurrent access uchun reader/writer lock
- **TTL (Time To Live)**: Cache expiration vaqti
- **Background Cleanup**: Memory leak'larni oldini olish
- **Thread Safety**: Concurrent access'da safe

### 3. Caching Middleware

API response'larni cache qilish:

```go
// Caching Middleware (internal/middleware/caching.go)
type CachingMiddleware struct {
    cache             *cache.MemoryCache
    logger            domain.Logger
    cacheTTL          time.Duration
    cacheableCommands map[string]bool
}

func NewCachingMiddleware(logger domain.Logger) *CachingMiddleware {
    cacheableCommands := map[string]bool{
        "/weather": true,
        "/ob-havo": true,
        "/repo":    true,
        "/user":    true,
    }

    return &CachingMiddleware{
        cache:             cache.NewMemoryCache(10 * time.Minute),
        cacheableCommands: cacheableCommands,
        // ...
    }
}

func (m *CachingMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
    return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
        // Check if command should be cached
        if !m.shouldCache(cmd.Text) {
            return next(ctx, cmd)
        }

        // Generate cache key
        cacheKey := m.generateCacheKey(cmd.User.TelegramID, cmd.Text)

        // Try cache first
        if cachedResponse, found := m.cache.Get(cacheKey); found {
            response := cachedResponse.(*domain.Response)
            response.Text = "🔄 " + response.Text // Cache indicator
            return response, nil
        }

        // Execute command
        response, err := next(ctx, cmd)
        if err == nil && response != nil {
            // Cache successful response
            ttl := m.getCacheTTL(cmd.Text)
            m.cache.SetWithTTL(cacheKey, response, ttl)
        }

        return response, err
    }
}
```

**Caching Strategy:**
- **Selective Caching**: Faqat expensive command'lar cache qilinadi
- **User-specific Keys**: Har bir user uchun alohida cache
- **TTL by Command**: Har xil command'lar uchun har xil TTL
- **Cache Indicators**: User'ga cache'dan kelganini ko'rsatish

### 4. Input Validation Middleware

Security va user experience uchun input validation:

```go
// Validation Middleware (internal/middleware/validation.go)
type ValidationMiddleware struct {
    logger     domain.Logger
    maxLength  int
    validators map[string]*CommandValidator
}

type CommandValidator struct {
    Pattern     *regexp.Regexp
    MinArgs     int
    MaxArgs     int
    Description string
    Usage       string
}

func NewValidationMiddleware(logger domain.Logger) *ValidationMiddleware {
    validators := make(map[string]*CommandValidator)

    // Weather command validation
    validators["/weather"] = &CommandValidator{
        Pattern:     regexp.MustCompile(`^/weather\s+[a-zA-Z\s\-']{2,50}$`),
        MinArgs:     2,
        MaxArgs:     5,
        Description: "Weather command requires a city name",
        Usage:       "/weather <city_name>",
    }

    // GitHub repo validation
    validators["/repo"] = &CommandValidator{
        Pattern:     regexp.MustCompile(`^/repo\s+[a-zA-Z0-9\-_.]+/[a-zA-Z0-9\-_.]+$`),
        MinArgs:     2,
        MaxArgs:     2,
        Description: "Repository command requires owner/repo format",
        Usage:       "/repo <owner/repository>",
    }

    return &ValidationMiddleware{
        validators: validators,
        maxLength:  500,
        // ...
    }
}
```

**Validation Features:**
- **Regex Patterns**: Command format validation
- **Argument Count**: Min/max argument checking
- **Length Limits**: DoS attack protection
- **Input Sanitization**: XSS prevention
- **User-friendly Errors**: Clear validation messages

### 5. Performance Metrics Collection

Real-time performance monitoring:

```go
// Metrics Middleware (internal/middleware/metrics.go)
type MetricsMiddleware struct {
    logger             domain.Logger
    totalRequests      int64
    successfulRequests int64
    failedRequests     int64
    commandMetrics     map[string]*CommandMetrics
    mutex             sync.RWMutex
    startTime         time.Time
}

type CommandMetrics struct {
    Count           int64
    TotalDuration   time.Duration
    AverageDuration time.Duration
    LastUsed        time.Time
    ErrorCount      int64
}

func (m *MetricsMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
    return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
        startTime := time.Now()
        
        // Increment total requests
        atomic.AddInt64(&m.totalRequests, 1)

        // Execute command
        response, err := next(ctx, cmd)
        
        // Calculate duration
        duration := time.Since(startTime)
        
        // Update metrics
        m.updateCommandMetrics(cmd.Text, duration, err)
        
        // Update counters
        if err != nil {
            atomic.AddInt64(&m.failedRequests, 1)
        } else {
            atomic.AddInt64(&m.successfulRequests, 1)
        }

        // Warn on slow commands
        if duration > 2*time.Second {
            m.logger.Warn("Slow command execution", 
                "command", cmd.Text,
                "duration", duration)
        }

        return response, err
    }
}
```

**Metrics Collected:**
- **Request Counts**: Total, successful, failed
- **Response Times**: Per command average, total
- **Error Rates**: Per command error percentage
- **Usage Patterns**: Most popular commands
- **Performance Alerts**: Slow command warnings

### 6. Advanced Command: /metrics

Real-time bot performance monitoring command:

```go
// Metrics Command (internal/handlers/commands/metrics.go)
func (h *MetricsCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
    metrics := h.metricsProvider.GetMetrics()
    cacheStats := h.metricsProvider.GetCacheStats()

    message := fmt.Sprintf(
        "📈 <b>Bot Performance Metrics</b>\n\n"+
        "🖥️ <b>System:</b>\n"+
        "   • Uptime: %s\n"+
        "   • Success Rate: %.1f%%\n"+
        "   • Req/min: %.1f\n\n"+
        "📊 <b>Requests:</b>\n"+
        "   • Total: %d\n"+
        "   • Successful: %d\n"+
        "   • Failed: %d\n\n"+
        "💾 <b>Cache:</b>\n"+
        "   • Size: %d items\n"+
        "   • TTL: %d minutes\n\n"+
        "⚡ <b>Performance:</b>\n"+
        "   • Avg Response: %dms\n"+
        "   • Slowest Command: %s (%dms)",
        uptime, successRate, requestsPerMinute,
        totalRequests, successfulRequests, failedRequests,
        cacheSize, cacheTTL,
        avgResponse, slowestCmd, slowestTime,
    )

    return &domain.Response{
        Text:      message,
        ParseMode: "HTML",
    }, nil
}
```

### 7. Testing Framework

Professional testing approach:

```go
// Unit Test Example (internal/cache/memory_cache_test.go)
func TestMemoryCache_SetAndGet(t *testing.T) {
    cache := NewMemoryCache(5 * time.Minute)

    cache.Set("test-key", "test-value")
    
    value, found := cache.Get("test-key")
    if !found {
        t.Error("Expected to find cached value")
    }

    if value != "test-value" {
        t.Errorf("Expected 'test-value', got %v", value)
    }
}

func TestMemoryCache_Expiration(t *testing.T) {
    cache := NewMemoryCache(5 * time.Minute)

    // Set with short TTL
    cache.SetWithTTL("expiring-key", "value", 50*time.Millisecond)

    // Should be available immediately
    _, found := cache.Get("expiring-key")
    if !found {
        t.Error("Value should be available immediately")
    }

    // Wait for expiration
    time.Sleep(100 * time.Millisecond)

    // Should be expired
    _, found = cache.Get("expiring-key")
    if found {
        t.Error("Value should be expired")
    }
}
```

**Testing Best Practices:**
- **Unit Tests**: Individual function testing
- **Middleware Tests**: Request/response flow testing  
- **Mock Objects**: External dependency isolation
- **Table-driven Tests**: Multiple test cases
- **Integration Tests**: End-to-end testing

### 8. Command Handler Testing

```go
// Command Test Example (internal/handlers/commands/start_test.go)
type MockLogger struct{}
func (m *MockLogger) Debug(msg string, args ...interface{}) {}
func (m *MockLogger) Info(msg string, args ...interface{})  {}
// ... other methods

func TestStartCommand_Handle(t *testing.T) {
    logger := &MockLogger{}
    startCmd := NewStartCommand("Welcome!", logger)

    cmd := &domain.Command{
        Text: "/start",
        User: &domain.User{
            TelegramID: 12345,
            FirstName:  "Test",
        },
    }

    response, err := startCmd.Handle(context.Background(), cmd)

    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }

    if response.Text == "" {
        t.Error("Expected non-empty response")
    }
}
```

## 💡 Week 4 Architecture Benefits

### 1. Performance Optimization
- **Memory Caching**: API call'larni 10x tez qilish
- **Metrics Collection**: Bottleneck identification
- **Background Processing**: Non-blocking operations

### 2. Security & Validation
- **Input Validation**: Malicious input protection
- **Rate Limiting**: Abuse prevention  
- **Sanitization**: XSS/injection prevention

### 3. Monitoring & Observability
- **Real-time Metrics**: Live performance monitoring
- **Error Tracking**: Failure analysis
- **Usage Analytics**: User behavior insights

### 4. Code Quality & Testing
- **Unit Testing**: Bug prevention
- **Integration Testing**: Flow validation
- **Mock Objects**: Isolated testing

## 🔧 Production Features

### 1. Graceful Error Handling
```go
if duration > 2*time.Second {
    m.logger.Warn("Slow command execution")
    // Continue operation, don't fail
}
```

### 2. Resource Management
```go
// Background cleanup prevents memory leaks
go cache.startCleanupRoutine()
```

### 3. Concurrent Safety
```go
// Thread-safe operations with mutexes
mc.mutex.Lock()
defer mc.mutex.Unlock()
```

### 4. Performance Alerting
```go
// Automatic slow command detection
if duration > threshold {
    logger.Warn("Performance alert")
}
```

## 🎯 Key Takeaways

1. **Middleware Order**: Performance va security uchun optimal order muhim
2. **Caching Strategy**: Selective caching expensive operations'ni optimize qiladi
3. **Input Validation**: Security va UX uchun zarur
4. **Performance Metrics**: Production monitoring uchun critical
5. **Testing Framework**: Code quality assurance

Week 4 orqali bot'ingiz enterprise-ready bo'lib, high-load production environment'da professional darajada ishlaydi!

## 🚀 Commands Added

- `/metrics` - Real-time performance monitoring dashboard
- Enhanced `/stats` - Detailed analytics with popular commands
- Improved error messages with validation feedback
- Cache indicators in responses (🔄 symbol)

Bu architecture professional Go development'da ishlatiladi va scalable microservice'lar yaratish uchun foundation beradi!