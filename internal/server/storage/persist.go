package storage

import "time"

type PersistSettings struct {
	storeInterval time.Duration
	storeFile     string
	restore       bool
}

func NewPersistSettings(storeInterval time.Duration, storeFile string, restore bool) PersistSettings {
	return PersistSettings{
		storeInterval: storeInterval,
		storeFile:     storeFile,
		restore:       restore,
	}
}

func (ps PersistSettings) ShouldPersist() bool {
	return ps.storeFile != ""
}

func (ps PersistSettings) ShouldSyncPersist() bool {
	return ps.ShouldPersist() && ps.storeInterval == 0
}

func (ps PersistSettings) ShouldIntervalPersist() bool {
	return ps.ShouldPersist() && ps.storeInterval != 0
}

func (ps PersistSettings) ShouldRestore() bool {
	return ps.ShouldPersist() && ps.restore
}

func (ps PersistSettings) GetStoreInterval() time.Duration {
	return ps.storeInterval
}

func (ps PersistSettings) GetStoreFile() string {
	return ps.storeFile
}

func (ps PersistSettings) GetRestore() bool {
	return ps.restore
}
