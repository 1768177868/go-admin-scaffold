# 缓存系统

本框架提供了灵活的缓存系统，支持多种缓存驱动，包括 Redis、内存缓存和文件缓存。

## 配置

在 `config/config.yaml` 中配置缓存：

```yaml
cache:
  driver: "redis"  # redis, memory, file
  prefix: "cache:"
  ttl: 3600       # 默认过期时间（秒）

  redis:
    host: "localhost"
    port: 6379
    password: ""
    db: 0
    pool_size: 100
    min_idle_conns: 10

  file:
    path: "storage/cache"
    extension: ".cache"
```

## 基本使用

### 1. 获取缓存实例

```go
// internal/services/user_service.go
type UserService struct {
    cache  cache.Cache
    repo   *repositories.UserRepository
}

func NewUserService(cache cache.Cache, repo *repositories.UserRepository) *UserService {
    return &UserService{
        cache: cache,
        repo:  repo,
    }
}
```

### 2. 基本操作

```go
// 设置缓存
err := cache.Set("key", value, time.Hour)

// 获取缓存
value, err := cache.Get("key")

// 删除缓存
err := cache.Delete("key")

// 检查是否存在
exists := cache.Has("key")

// 递增/递减
newValue, err := cache.Increment("counter", 1)
newValue, err := cache.Decrement("counter", 1)

// 永久存储
err := cache.Forever("key", value)

// 获取或设置
value, err := cache.Remember("key", time.Hour, func() interface{} {
    // 如果缓存不存在，这个函数会被调用
    return calculateValue()
})
```

## 缓存标签

支持使用标签组织和管理缓存：

```go
// 使用标签
tags := cache.Tags([]string{"users", "permissions"})

// 设置带标签的缓存
err := tags.Set("user:1", userData, time.Hour)

// 获取带标签的缓存
value, err := tags.Get("user:1")

// 清除特定标签的所有缓存
err := cache.Tags([]string{"users"}).Flush()
```

## 缓存驱动

### 1. Redis 驱动

```go
// pkg/cache/redis_driver.go
type RedisCache struct {
    client *redis.Client
    prefix string
}

func (c *RedisCache) Get(key string) (interface{}, error) {
    value, err := c.client.Get(ctx, c.prefix+key).Result()
    if err == redis.Nil {
        return nil, cache.ErrCacheMiss
    }
    return value, err
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    return c.client.Set(ctx, c.prefix+key, value, ttl).Err()
}
```

### 2. 内存缓存

```go
// pkg/cache/memory_driver.go
type MemoryCache struct {
    store  map[string]*cacheItem
    mutex  sync.RWMutex
    prefix string
}

type cacheItem struct {
    value      interface{}
    expiration time.Time
}

func (c *MemoryCache) Get(key string) (interface{}, error) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    item, exists := c.store[c.prefix+key]
    if !exists || item.isExpired() {
        return nil, cache.ErrCacheMiss
    }
    return item.value, nil
}
```

### 3. 文件缓存

```go
// pkg/cache/file_driver.go
type FileCache struct {
    path      string
    prefix    string
    extension string
}

func (c *FileCache) Get(key string) (interface{}, error) {
    path := c.getPath(key)
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, cache.ErrCacheMiss
    }
    
    item := &cacheItem{}
    if err := json.Unmarshal(data, item); err != nil {
        return nil, err
    }
    
    if item.isExpired() {
        os.Remove(path)
        return nil, cache.ErrCacheMiss
    }
    
    return item.value, nil
}
```

## 缓存策略

### 1. 缓存预热

```go
func (s *UserService) WarmUpCache() error {
    users, err := s.repo.GetAll()
    if err != nil {
        return err
    }

    for _, user := range users {
        key := fmt.Sprintf("user:%d", user.ID)
        if err := s.cache.Set(key, user, time.Hour); err != nil {
            return err
        }
    }
    return nil
}
```

### 2. 缓存穿透保护

```go
func (s *UserService) GetUser(id uint) (*models.User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // 使用空值缓存防止缓存穿透
    user, err := s.cache.Remember(key, time.Hour, func() interface{} {
        user, err := s.repo.Find(id)
        if err != nil {
            if err == gorm.ErrRecordNotFound {
                return nil // 缓存空值
            }
            return err
        }
        return user
    })
    
    if user == nil {
        return nil, gorm.ErrRecordNotFound
    }
    
    return user.(*models.User), nil
}
```

### 3. 缓存击穿保护

```go
func (s *UserService) GetHotData(key string) (interface{}, error) {
    // 使用互斥锁防止缓存击穿
    value, err := s.cache.Get(key)
    if err == nil {
        return value, nil
    }

    s.mutex.Lock()
    defer s.mutex.Unlock()

    // 双重检查
    value, err = s.cache.Get(key)
    if err == nil {
        return value, nil
    }

    // 从数据源获取
    value, err = s.getFromSource(key)
    if err != nil {
        return nil, err
    }

    // 设置缓存
    s.cache.Set(key, value, time.Hour)
    return value, nil
}
```

## 最佳实践

1. 缓存设计：
   - 选择合适的缓存驱动
   - 设置合理的过期时间
   - 使用有意义的键名前缀

2. 性能优化：
   - 批量获取和存储
   - 合理使用标签
   - 实现缓存预热

3. 数据一致性：
   - 及时清理过期数据
   - 实现缓存同步机制
   - 处理并发更新

4. 安全性：
   - 限制缓存大小
   - 防止缓存穿透
   - 保护敏感数据

## 监控和维护

### 1. 缓存统计

```go
type CacheStats struct {
    Hits        uint64
    Misses      uint64
    Size        uint64
    ItemsCount  uint64
}

func (c *Cache) GetStats() *CacheStats {
    return &CacheStats{
        Hits:       atomic.LoadUint64(&c.hits),
        Misses:     atomic.LoadUint64(&c.misses),
        Size:       c.size(),
        ItemsCount: c.count(),
    }
}
```

### 2. 缓存清理

```go
// 清理过期项
func (c *Cache) GC() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    for key, item := range c.store {
        if item.isExpired() {
            delete(c.store, key)
        }
    }
}
``` 