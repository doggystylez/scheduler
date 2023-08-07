package scheduler

import (
	"strconv"
	"testing"
	"time"
)

type testEvent struct {
	Name  string
	Index int
}

func (e testEvent) Fire() {
	testChan <- time.Now()
}

func (e testEvent) GetLabel() string {
	return strconv.Itoa(e.Index)
}

var (
	testChan chan time.Time

	event0 = testEvent{
		Name:  "Alice",
		Index: 0,
	}

	event1 = testEvent{
		Name:  "Bob",
		Index: 1,
	}

	event2 = testEvent{
		Name:  "Charlie",
		Index: 2,
	}

	event3 = testEvent{
		Name:  "Dave",
		Index: 3,
	}
)

func init() {
	testChan = make(chan time.Time, 3)
}

func TestQueue(t *testing.T) {
	q := New()
	inFiveSeconds := time.Now().Add(time.Second * 5)
	inSevenSeconds := inFiveSeconds.Add(time.Second * 2)
	inTenSeconds := inSevenSeconds.Add(time.Second * 3)
	q.Add(event0, inFiveSeconds)
	q.Add(event1, inSevenSeconds)
	q.Add(event2, inTenSeconds)

	if !q.Exists(event0) {
		t.Error("event0 failed to schedule")
	}

	if q.Exists(event3) {
		t.Error("event3 not scheduled but appears in queue")
	}

	time.Sleep(time.Second * 11)
	close(testChan)
	if len(testChan) != 3 {
		t.Error(3-len(testChan), "timers failed to fire")
	}
	var testTimes []time.Time
	for t := range testChan {
		testTimes = append(testTimes, t)
	}
	if testTimes[0].Before(inFiveSeconds) {
		t.Error("event0 fired early")
	} else if testTimes[0].After(inFiveSeconds.Add(time.Millisecond * 5)) {
		t.Error("event0 fired more than 5ms late")
	}
	if testTimes[1].Before(inSevenSeconds) {
		t.Error("event2 fired early")
	} else if testTimes[1].After(inSevenSeconds.Add(time.Millisecond * 5)) {
		t.Error("event2 fired more than 5ms late")
	}
	if testTimes[2].Before(inTenSeconds) {
		t.Error("event2 fired early")
	} else if testTimes[2].After(inTenSeconds.Add(time.Millisecond * 5)) {
		t.Error("event2 fired more than 5ms late")
	}
}
