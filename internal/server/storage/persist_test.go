package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPersistFlags(t *testing.T) {
	assert.True(t, PersistSettings{storeFile: "foobar"}.ShouldPersist())
	assert.False(t, PersistSettings{storeFile: ""}.ShouldPersist())

	assert.True(t, PersistSettings{storeFile: "foobar", storeInterval: 0}.ShouldSyncPersist())
	assert.False(t, PersistSettings{storeFile: "foobar", storeInterval: time.Duration(1)}.ShouldSyncPersist())
	assert.False(t, PersistSettings{storeFile: "", storeInterval: 0}.ShouldSyncPersist())
	assert.False(t, PersistSettings{storeFile: "", storeInterval: time.Duration(1)}.ShouldSyncPersist())

	assert.True(t, PersistSettings{storeFile: "foobar", storeInterval: time.Duration(1)}.ShouldIntervalPersist())
	assert.False(t, PersistSettings{storeFile: "foobar", storeInterval: 0}.ShouldIntervalPersist())
	assert.False(t, PersistSettings{storeFile: "", storeInterval: time.Duration(1)}.ShouldIntervalPersist())
	assert.False(t, PersistSettings{storeFile: "", storeInterval: 0}.ShouldIntervalPersist())

	assert.True(t, PersistSettings{storeFile: "foobar", restore: true}.ShouldRestore())
	assert.False(t, PersistSettings{storeFile: "foobar", restore: false}.ShouldRestore())
	assert.False(t, PersistSettings{storeFile: "", restore: true}.ShouldRestore())
	assert.False(t, PersistSettings{storeFile: "", restore: false}.ShouldRestore())
}
