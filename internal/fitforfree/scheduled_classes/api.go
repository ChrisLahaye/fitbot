package scheduled_classes

import (
	"github.com/ChrisLahaye/fitbot/internal/fitforfree"
	"fmt"
	"strings"
	"time"
)

// API contains the internal state
type API struct {
	fitforfree.API
}

// ScheduledClass contains a scheduled class
type ScheduledClass struct {
	Activity struct {
		ID string `json:"id"` // vrijtrainen | virtuelespinning
	} `json:"activity"`
	AvailablePercentage int64  `json:"availability_percentage"`
	AvailableText       string `json:"availability_text"`
	Booked              bool   `json:"booked"`
	Capacity            int64  `json:"capacity"`
	ClassType           string `json:"classType"` // free_practise
	DurationSeconds     int64  `json:"durationSeconds"`
	ID                  string `json:"id"`
	Roomname            string `json:"roomname"`
	SpotsAvailable      int64  `json:"spotsAvailable"`
	StartTimestamp      int64  `json:"startTimestamp"`
	Status              string `json:"status"` // AVAILABLE | FULL
}

// Start returns the start time
func (v *ScheduledClass) Start() time.Time {
	return time.Unix(v.StartTimestamp, 0)
}

// End returns the end time
func (v *ScheduledClass) End() time.Time {
	d := time.Duration(v.DurationSeconds)
	return v.Start().Add(d * time.Second)
}

// String returns the string
func (v *ScheduledClass) String() string {
	labels := []string{v.Status}
	if v.Booked {
		labels = append(labels, "BOOKED")
	}

	return fmt.Sprintf("%s-%s %s",
		v.Start().Format("Mon 15:00"),
		v.End().Format("15:00"),
		strings.Join(labels, " "))
}
