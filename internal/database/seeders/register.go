package seeders

import (
	"app/internal/database/seeder"
)

var globalManager *seeder.SeederManager
var pendingRegistrations = make(map[string]*seeder.Seeder)

// Register registers a seeder with the global manager or stores it temporarily
func Register(name string, s *seeder.Seeder) {
	if globalManager != nil {
		globalManager.Register(name, s)
	} else {
		pendingRegistrations[name] = s
	}
}

// SetGlobalManager sets the global seeder manager and registers all pending seeders
func SetGlobalManager(manager *seeder.SeederManager) {
	globalManager = manager

	// Register all pending seeders
	for name, s := range pendingRegistrations {
		manager.Register(name, s)
	}

	// Clear pending registrations
	pendingRegistrations = make(map[string]*seeder.Seeder)
}
