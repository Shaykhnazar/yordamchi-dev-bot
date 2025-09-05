# Week 3: Database Analytics & Activity Tracking (Ma'lumotlar Bazasi Analitikasi)

## 🎯 O'rganish Maqsadlari

Bu darsda siz professional darajadagi database analytics va user activity tracking sistemini o'rganasiz.

## 📚 Asosiy Tushunchalar

### 1. Database Analytics nima?

Database Analytics - bu foydalanuvchi faoliyatini kuzatish va statistik ma'lumotlar to'plash tizimi.

**Asosiy maqsadlar:**
- User behavior tracking (Foydalanuvchi xatti-harakatlarini kuzatish)
- Command popularity analysis (Buyruqlar populyarligi tahlili)
- Daily/weekly analytics (Kunlik/haftalik analitika)
- Performance monitoring (Ish unumdorligini monitoring qilish)

### 2. Enhanced Database Methods

Database analytics uchun maxsus metodlar:

```go
// UserActivity - foydalanuvchi faoliyati entity'si
type UserActivity struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Command   string    `json:"command"`
    CreatedAt time.Time `json:"created_at"`
}

// GetUserActivities - foydalanuvchining oxirgi faoliyatlari
func (db *DB) GetUserActivities(telegramID int64, limit int) ([]UserActivity, error) {
    query := `
    SELECT ua.id, ua.user_id, ua.command, ua.timestamp 
    FROM user_activity ua
    JOIN users u ON ua.user_id = u.id
    WHERE u.telegram_id = ?
    ORDER BY ua.timestamp DESC 
    LIMIT ?`

    rows, err := db.conn.Query(query, telegramID, limit)
    if err != nil {
        return nil, fmt.Errorf("faollikni olishda xatolik: %w", err)
    }
    defer rows.Close()

    var activities []UserActivity
    for rows.Next() {
        var activity UserActivity
        err := rows.Scan(
            &activity.ID,
            &activity.UserID,
            &activity.Command,
            &activity.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("faollik ma'lumotlarini o'qishda xatolik: %w", err)
        }
        activities = append(activities, activity)
    }

    return activities, nil
}
```

**Muhim tushunchalar:**
- `JOIN` - ikki jadval orasidagi bog'lanish
- `ORDER BY ... DESC` - teskari tartibda saralash
- `LIMIT` - natijalar sonini cheklash
- `rows.Next()` - natijalarni ketma-ket o'qish

### 3. Popular Commands Analytics

Eng ko'p ishlatiladigan buyruqlarni aniqlash:

```go
// GetPopularCommands - mashhur buyruqlar ro'yxati
func (db *DB) GetPopularCommands(limit int) (map[string]int, error) {
    query := `
    SELECT command, COUNT(*) as count 
    FROM user_activity 
    GROUP BY command 
    ORDER BY count DESC 
    LIMIT ?`

    rows, err := db.conn.Query(query, limit)
    if err != nil {
        return nil, fmt.Errorf("populyar buyruqlarni olishda xatolik: %w", err)
    }
    defer rows.Close()

    commands := make(map[string]int)
    for rows.Next() {
        var command string
        var count int
        err := rows.Scan(&command, &count)
        if err != nil {
            return nil, fmt.Errorf("buyruq ma'lumotlarini o'qishda xatolik: %w", err)
        }
        commands[command] = count
    }

    return commands, nil
}
```

**SQL Aggregation Functions:**
- `COUNT(*)` - qatorlar sonini hisoblash
- `GROUP BY` - ma'lumotlarni guruhlash
- `ORDER BY count DESC` - son bo'yicha kamayish tartibida
- `map[string]int` - buyruq nomi va soni

### 4. Daily Statistics

Kunlik statistikalar olish:

```go
// GetDailyStats - kunlik statistikalar
func (db *DB) GetDailyStats() (map[string]int, error) {
    stats := make(map[string]int)
    
    // Bugungi yangi foydalanuvchilar
    query := "SELECT COUNT(*) FROM users WHERE DATE(created_at) = DATE('now')"
    var newUsersToday int
    err := db.conn.QueryRow(query).Scan(&newUsersToday)
    if err != nil {
        return nil, fmt.Errorf("bugungi foydalanuvchilar sonini olishda xatolik: %w", err)
    }
    stats["new_users_today"] = newUsersToday
    
    // Bugungi faoliyat
    query = "SELECT COUNT(*) FROM user_activity WHERE DATE(timestamp) = DATE('now')"
    var activitiesToday int
    err = db.conn.QueryRow(query).Scan(&activitiesToday)
    if err != nil {
        return nil, fmt.Errorf("bugungi faollik sonini olishda xatolik: %w", err)
    }
    stats["activities_today"] = activitiesToday
    
    // Bugungi faol foydalanuvchilar
    query = "SELECT COUNT(DISTINCT user_id) FROM user_activity WHERE DATE(timestamp) = DATE('now')"
    var activeUsersToday int
    err = db.conn.QueryRow(query).Scan(&activeUsersToday)
    if err != nil {
        return nil, fmt.Errorf("bugungi faol foydalanuvchilar sonini olishda xatolik: %w", err)
    }
    stats["active_users_today"] = activeUsersToday
    
    return stats, nil
}
```

**SQL Date Functions:**
- `DATE('now')` - hozirgi sana
- `DATE(timestamp)` - vaqt ni sanaga o'zgartirish
- `COUNT(DISTINCT user_id)` - noyob foydalanuvchilar soni

### 5. Activity Tracking Middleware

Foydalanuvchi faoliyatini avtomatik yozib olish:

```go
// ActivityMiddleware - faoliyat kuzatuv middleware
type ActivityMiddleware struct {
    db     *database.DB
    logger domain.Logger
}

func NewActivityMiddleware(db *database.DB, logger domain.Logger) *ActivityMiddleware {
    return &ActivityMiddleware{
        db:     db,
        logger: logger,
    }
}

func (m *ActivityMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
    return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
        // Avval buyruqni bajar
        response, err := next(ctx, cmd)

        // Muvaffaqiyatli bajarilgandan keyin log qil
        if err == nil && cmd.User != nil {
            // Background'da log qil - response'ni bloklama
            go func() {
                logErr := m.db.LogUserActivity(cmd.User.TelegramID, cmd.Text)
                if logErr != nil {
                    m.logger.Warn("Failed to log user activity",
                        "telegram_id", cmd.User.TelegramID,
                        "command", cmd.Text,
                        "error", logErr)
                } else {
                    m.logger.Debug("User activity logged",
                        "telegram_id", cmd.User.TelegramID,
                        "command", cmd.Text)
                }
            }()
        }

        return response, err
    }
}
```

**Middleware Pattern Benefits:**
- Har bir buyruq avtomatik log qilinadi
- Background goroutine - tezlik uchun
- Error handling - xatoliklar boshqariladi
- Clean separation - faoliyat kuzatuv alohida

### 6. Enhanced Statistics Command

Yangilangan `/stats` buyruq:

```go
func (h *StatsCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
    // Asosiy statistikalar
    totalUsers, err := h.db.GetUserStats()
    if err != nil {
        return nil, err
    }

    // Kunlik statistikalar
    dailyStats, err := h.db.GetDailyStats()
    if err != nil {
        dailyStats = make(map[string]int)
    }

    // Mashhur buyruqlar
    popularCommands, err := h.db.GetPopularCommands(5)
    if err != nil {
        popularCommands = make(map[string]int)
    }

    message := fmt.Sprintf(
        "📊 <b>Bot Statistikasi</b>\n\n"+
        "👥 <b>Foydalanuvchilar:</b>\n"+
        "   • Jami: %d\n"+
        "   • Bugun yangi: %d\n"+
        "   • Bugun faol: %d\n\n"+
        "📈 <b>Faollik:</b>\n"+
        "   • Bugun buyruqlar: %d\n\n",
        totalUsers,
        dailyStats["new_users_today"],
        dailyStats["active_users_today"],
        dailyStats["activities_today"],
    )

    // Mashhur buyruqlar qo'shish
    if len(popularCommands) > 0 {
        message += "🔥 <b>Populyar buyruqlar:</b>\n"
        for cmd, count := range popularCommands {
            message += fmt.Sprintf("   • %s: %d\n", cmd, count)
        }
    }

    return &domain.Response{
        Text:      message,
        ParseMode: "HTML",
    }, nil
}
```

## 💡 Week 3 Architecture Benefits

### 1. Data-Driven Decisions
- Qaysi buyruqlar eng ko'p ishlatiladi
- Foydalanuvchilar faoliyati qanday
- Bot qachon eng ko'p ishlatiladi

### 2. Performance Monitoring
- Database query performance
- User engagement metrics
- System health monitoring

### 3. Business Intelligence
- User growth tracking
- Feature adoption rates
- Usage pattern analysis

### 4. Scalability Insights
- Peak usage times
- Resource utilization
- Bottleneck identification

## 🔧 Middleware Architecture Benefits

### 1. Automatic Tracking
- Har bir buyruq avtomatik kuzatiladi
- Manual logging'ga ehtiyoj yo'q
- Consistent data collection

### 2. Performance Optimization
- Background logging
- Non-blocking operations
- Asynchronous processing

### 3. Separation of Concerns
- Analytics logic alohida
- Command handlers tozaroq
- Maintainable code

## 📈 Production Ready Features

### 1. Error Handling
```go
if err != nil {
    m.logger.Warn("Failed to log user activity", "error", err)
    // Continue operation - don't break bot functionality
}
```

### 2. Graceful Degradation
```go
dailyStats, err := h.db.GetDailyStats()
if err != nil {
    dailyStats = make(map[string]int) // Continue with empty stats
}
```

### 3. Background Processing
```go
go func() {
    logErr := m.db.LogUserActivity(cmd.User.TelegramID, cmd.Text)
    // Process in background to avoid blocking
}()
```

## 🎯 Key Takeaways

1. **Analytics Foundation**: Database analytics professional bot development uchun zarur
2. **Middleware Pattern**: Faoliyat kuzatuv avtomatik va samarali
3. **SQL Aggregation**: COUNT, GROUP BY, JOIN - analytics uchun asosiy
4. **Background Processing**: Performance uchun muhim
5. **Error Handling**: Production environment'da zarur

Week 3 orqali siz bot'ing foydalanuvchi faoliyatini professional darajada kuzatishi va tahlil qilishi mumkin bo'ladi!