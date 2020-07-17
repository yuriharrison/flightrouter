package flightsdb

import "strings"

func keyForOriginDestination(origin, destination string) string {
	return strings.ToUpper(origin + destination)
}

// Cache simple map cache
type Cache struct {
	hits, misses  int
	cheapestRoute map[string][]*Flight
}

// Clean cache memory
func (c *Cache) Clean() {
	c.cheapestRoute = nil
}

func (c *Cache) resetCache() {
	c.cheapestRoute = make(map[string][]*Flight)
}

// GetCheapestRoute return cached route if exist
func (c *Cache) GetCheapestRoute(origin, destination string) []*Flight {
	key := keyForOriginDestination(origin, destination)
	if c.cheapestRoute == nil {
		c.resetCache()
	}
	if v, ok := c.cheapestRoute[key]; ok {
		c.hits++
		return v
	}
	c.misses++
	return nil
}

// SetCheapestRoute return cached route if exist
func (c *Cache) SetCheapestRoute(origin, destination string, value []*Flight) {
	key := keyForOriginDestination(origin, destination)
	if c.cheapestRoute == nil {
		c.resetCache()
	}
	c.cheapestRoute[key] = value
}
