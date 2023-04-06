package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_GetSlots(t *testing.T) {
	store := &Store{
		Slots: []map[string]bool{
			{"10:00-10:30": true, "10:30-11:00": false},
			{},
			{},
			{},
			{},
			{},
			{},
		},
		Holidays: nil,
	}
	slots := store.GetSlots("2023-01-22")
	assert.Equal(t, []string{"10:00-10:30"}, slots)
}

func TestStore_GetSlots_Holidays(t *testing.T) {
	store := &Store{
		Slots: []map[string]bool{
			{"10:00-10:30": true, "10:30-11:00": false},
			{},
			{},
			{},
			{},
			{},
			{},
		},
		Holidays: []string{"2023-01-22"},
	}
	slots := store.GetSlots("2023-01-22")
	assert.Empty(t, slots)
}

func TestStore_GetSlots_Special(t *testing.T) {
	store := &Store{
		Slots: []map[string]bool{
			{},
			{},
			{},
			{},
			{},
			{},
			{},
		},
		SpecialDays: map[string]map[string]bool{
			"2023-01-22": {"10:00-10:30": true, "10:30-11:00": false},
		},
	}
	slots := store.GetSlots("2023-01-22")
	assert.Equal(t, 1, len(slots))
	assert.Equal(t, "10:00-10:30", slots[0])
}

func TestStore_GetSlots_Monday(t *testing.T) {
	store := &Store{
		Slots: []map[string]bool{
			{"10:00-10:30": true, "10:30-11:00": false},
			{},
			{},
			{},
			{},
			{},
			{},
		},
		Holidays: []string{"2023-01-22"},
	}
	slots := store.GetSlots("2023-01-23")
	assert.Empty(t, slots)
}
