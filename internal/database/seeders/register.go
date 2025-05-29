package seeders

import (
	"app/internal/database/seeder"
)

var seeders = make(map[string]*seeder.Seeder)

// Register registers a seeder
func Register(name string, s *seeder.Seeder) {
	seeders[name] = s
}

// GetSeeders returns all registered seeders
func GetSeeders() map[string]*seeder.Seeder {
	return seeders
}

// InitSeeders initializes all seeders
func InitSeeders(manager *seeder.SeederManager) {
	for name, s := range seeders {
		manager.Register(name, s)
	}
}
