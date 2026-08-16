package ttlcache

import "time"

type Clock func() time.Time

type entry struct {
	value     string
	expiresAt time.Time
}

type Cache struct {
	now   Clock
	items map[string]entry
}

func New(now Clock) *Cache {
	if now == nil {
		now = time.Now
	}
	return &Cache{
		now:   now,
		items: make(map[string]entry),
	}
}

func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.items[key] = entry{
		value:     value,
		expiresAt: c.now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	item, ok := c.items[key]
	if !ok {
		return "", false
	}
	if !c.now().Before(item.expiresAt) {
		delete(c.items, key)
		return "", false
	}
	return item.value, true
}
