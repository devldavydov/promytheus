package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPersistFlags(t *testing.T) {
	assert.True(t, PersistSettings{StoreFile: "foobar"}.ShouldPersist())
	assert.False(t, PersistSettings{StoreFile: ""}.ShouldPersist())

	assert.True(t, PersistSettings{StoreFile: "foobar", StoreInterval: 0}.ShouldSyncPersist())
	assert.False(t, PersistSettings{StoreFile: "foobar", StoreInterval: time.Duration(1)}.ShouldSyncPersist())
	assert.False(t, PersistSettings{StoreFile: "", StoreInterval: 0}.ShouldSyncPersist())
	assert.False(t, PersistSettings{StoreFile: "", StoreInterval: time.Duration(1)}.ShouldSyncPersist())

	assert.True(t, PersistSettings{StoreFile: "foobar", StoreInterval: time.Duration(1)}.ShouldIntervalPersist())
	assert.False(t, PersistSettings{StoreFile: "foobar", StoreInterval: 0}.ShouldIntervalPersist())
	assert.False(t, PersistSettings{StoreFile: "", StoreInterval: time.Duration(1)}.ShouldIntervalPersist())
	assert.False(t, PersistSettings{StoreFile: "", StoreInterval: 0}.ShouldIntervalPersist())

	assert.True(t, PersistSettings{StoreFile: "foobar", Restore: true}.ShouldRestore())
	assert.False(t, PersistSettings{StoreFile: "foobar", Restore: false}.ShouldRestore())
	assert.False(t, PersistSettings{StoreFile: "", Restore: true}.ShouldRestore())
	assert.False(t, PersistSettings{StoreFile: "", Restore: false}.ShouldRestore())
}
