package lesson

import (
	"fitbot/internal/fitforfree"
	"fmt"
	"strings"
	"time"
)

// API contains the internal state
type API struct {
	fitforfree.API
}

// Booking contains a booking
type Booking struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Lesson contains a lesson
type Lesson struct {
	ID                   string `json:"id"`
	VenueName            string `json:"venueName"`
	VenueID              string `json:"venueId"`
	StartTimestamp       int64  `json:"startTimestamp"`
	PreCheckinTimeStamp  int64  `json:"preCheckinTimestamp"`
	PostCheckinTimestamp int64  `json:"postCheckinTimestamp"`
	DurationSeconds      int64  `json:"durationSeconds"`
	Instructor           string `json:"instructor"`
	Activity             struct {
		ID          string `json:"id"` // vrijefitness
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
		ImageURL    string `json:"imageUrl"`
	} `json:"activity"`
	Status              string `json:"status"` // AVAILABLE | FULL
	Booked              bool   `json:"booked"`
	ClassType           string `json:"classType"` // free_practise
	SpotsAvailable      int64  `json:"spotsAvailable"`
	Capacity            int64  `json:"capacity"`
	AvailablePercentage int64  `json:"availability_percentage"`
	Roomname            string `json:"roomname"`
}

// Start returns the start time
func (v *Lesson) Start() time.Time {
	return time.Unix(v.StartTimestamp, 0)
}

// End returns the end time
func (v *Lesson) End() time.Time {
	d := time.Duration(v.DurationSeconds)
	return v.Start().Add(d * time.Second)
}

// String returns the string
func (v *Lesson) String() string {
	labels := []string{v.Status}
	if v.Booked {
		labels = append(labels, "BOOKED")
	}

	return fmt.Sprintf("%s-%s %s",
		v.Start().Format("Mon 15:00"),
		v.End().Format("15:00"),
		strings.Join(labels, " "))
}
