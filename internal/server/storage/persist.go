package storage

import "time"

type PersistSettings struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

func NewPersistSettings(storeInterval time.Duration, storeFile string, restore bool) PersistSettings {
	return PersistSettings{
		StoreInterval: storeInterval,
		StoreFile:     storeFile,
		Restore:       restore,
	}
}

func (ps PersistSettings) ShouldPersist() bool {
	return ps.StoreFile != ""
}

func (ps PersistSettings) ShouldSyncPersist() bool {
	return ps.ShouldPersist() && ps.StoreInterval == 0
}

func (ps PersistSettings) ShouldIntervalPersist() bool {
	return ps.ShouldPersist() && ps.StoreInterval != 0
}

func (ps PersistSettings) ShouldRestore() bool {
	return ps.ShouldPersist() && ps.Restore
}
