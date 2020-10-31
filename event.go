package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Event is an abstraction for eventing
type Event struct {
	Done      bool
	Execute   func()
	Sound     rl.Sound
	Trigger   func() bool
	Triggered bool
}

// PlaySound plays the event indicator
func (event *Event) PlaySound() {
	rl.PlaySound(event.Sound)
}

// NewEventDuration creates a Duration based event. Pass your function and the duration you want to spawn it in
func NewEventDuration(triggerIn time.Duration, execute func()) *Event {
	endsAt := time.Now().Add(triggerIn)
	trigger := func() bool {
		return time.Now().After(endsAt)
	}

	return &Event{
		Execute: execute,
		Sound:   rl.LoadSound("assets/sounds/ui1.mp3"),
		Trigger: trigger,
	}
}