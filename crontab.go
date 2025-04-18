package crontask

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// crontab struct representing cron table
type crontab struct {
	ticker *time.Ticker
	jobs   []*job
	sync.RWMutex
}

// job in cron table
type job struct {
	min       map[int]struct{}
	hour      map[int]struct{}
	day       map[int]struct{}
	month     map[int]struct{}
	dayOfWeek map[int]struct{}

	fn   any
	args []any
	sync.RWMutex
}

// tick is individual tick that occures each minute
type tick struct {
	min       int
	hour      int
	day       int
	month     int
	dayOfWeek int
}

// New initializes and returns new cron table
func newCrontab() *crontab {
	return new(time.Minute)
}

// new creates new crontab, arg provided for testing purpose
func new(t time.Duration) *crontab {
	c := &crontab{
		ticker: time.NewTicker(t),
		jobs:   []*job{},
	}

	go func() {
		for t := range c.ticker.C {
			c.runScheduled(t)
		}
	}()

	return c
}

// AddJob to cron table.
//
// Returns error if:
//
// * Cron syntax can't be parsed or out of bounds
//
// * fn is not function
//
// * Provided args don't match the number and/or the type of fn args
func (c *crontab) AddJob(schedule string, fn any, args ...any) error {
	j, err := parseSchedule(schedule)
	c.Lock()
	defer c.Unlock()
	if err != nil {
		return err
	}

	if fn == nil || reflect.ValueOf(fn).Kind() != reflect.Func {
		return newErr("cron job must be func()")
	}

	fnType := reflect.TypeOf(fn)
	if len(args) != fnType.NumIn() {
		return newErr("number of func() params and number of provided params doesn't match")
	}

	for i := range fnType.NumIn() {
		a := args[i]
		t1 := fnType.In(i)
		t2 := reflect.TypeOf(a)

		if t1 != t2 {
			if t1.Kind() != reflect.Interface {
				return newErr("Param with index", i, "shold be", t1, "not", t2)
			}
			if !t2.Implements(t1) {
				return newErr("Param with index", i, "of type", t2, "doesn't implement interface", t1)
			}
		}
	}

	// all checked, add job to cron tab
	j.fn = fn
	j.args = args
	c.jobs = append(c.jobs, j)
	return nil
}

// MustAddJob is like AddJob but panics if there is an problem with job
//
// It simplifies initialization, since we usually add jobs at the beggining so you won't have to check for errors (it will panic when program starts).
// It is a similar aproach as go's std lib package `regexp` and `regexp.Compile()` `regexp.MustCompile()`
// MustAddJob will panic if:
//
// * Cron syntax can't be parsed or out of bounds
//
// * fn is not function
//
// * Provided args don't match the number and/or the type of fn args
func (c *crontab) MustAddJob(schedule string, fn any, args ...any) {
	if err := c.AddJob(schedule, fn, args...); err != nil {
		panic(err)
	}
}

// Shutdown the cron table schedule
//
// Once stopped, it can't be restarted.
// This function is pre-shuttdown helper for your app, there is no Start/Stop functionallity with crontab package.
func (c *crontab) Shutdown() {
	c.ticker.Stop()
}

// Clear all jobs from cron table
func (c *crontab) Clear() {
	c.Lock()
	c.jobs = []*job{}
	c.Unlock()
}

// RunAll jobs in cron table, shcheduled or not
func (c *crontab) RunAll() {
	c.RLock()
	defer c.RUnlock()
	for _, j := range c.jobs {
		go j.run()
	}
}

// RunScheduled jobs
func (c *crontab) runScheduled(t time.Time) {
	tick := getTick(t)
	c.RLock()
	defer c.RUnlock()

	for _, j := range c.jobs {
		if j.tick(tick) {
			go j.run()
		}
	}
}

// run the job using reflection
// Recover from panic although all functions and params are checked by AddJob, but you never know.
func (j *job) run() {
	j.RLock()
	defer func() {
		if r := recover(); r != nil {
			log.Println("crontab error", r)
		}
	}()
	v := reflect.ValueOf(j.fn)
	rargs := make([]reflect.Value, len(j.args))
	for i, a := range j.args {
		rargs[i] = reflect.ValueOf(a)
	}
	j.RUnlock()
	v.Call(rargs)
}

// tick decides should the job be lauhcned at the tick
func (j *job) tick(t tick) bool {
	j.RLock()
	defer j.RUnlock()
	if _, ok := j.min[t.min]; !ok {
		return false
	}

	if _, ok := j.hour[t.hour]; !ok {
		return false
	}

	// cummulative day and dayOfWeek, as it should be
	_, day := j.day[t.day]
	_, dayOfWeek := j.dayOfWeek[t.dayOfWeek]
	if !day && !dayOfWeek {
		return false
	}

	if _, ok := j.month[t.month]; !ok {
		return false
	}

	return true
}

// regexps for parsing schedule string
var (
	matchSpaces = regexp.MustCompile(`\s+`)
	matchN      = regexp.MustCompile(`(.*)/(\d+)`)
	matchRange  = regexp.MustCompile(`^(\d+)-(\d+)$`)
)

// parseSchedule string and creates job struct with filled times to launch, or error if synthax is wrong
func parseSchedule(s string) (*job, error) {
	var err error
	j := &job{}
	j.Lock()
	defer j.Unlock()
	s = matchSpaces.ReplaceAllLiteralString(s, " ")
	parts := strings.Split(s, " ")
	if len(parts) != 5 {
		return j, newErr("Schedule string must have five components like * * * * *")
	}

	j.min, err = parsePart(parts[0], 0, 59)
	if err != nil {
		return j, err
	}

	j.hour, err = parsePart(parts[1], 0, 23)
	if err != nil {
		return j, err
	}

	j.day, err = parsePart(parts[2], 1, 31)
	if err != nil {
		return j, err
	}

	j.month, err = parsePart(parts[3], 1, 12)
	if err != nil {
		return j, err
	}

	j.dayOfWeek, err = parsePart(parts[4], 0, 6)
	if err != nil {
		return j, err
	}

	//  day/dayOfWeek combination
	switch {
	case len(j.day) < 31 && len(j.dayOfWeek) == 7: // day set, but not dayOfWeek, clear dayOfWeek
		j.dayOfWeek = make(map[int]struct{})
	case len(j.dayOfWeek) < 7 && len(j.day) == 31: // dayOfWeek set, but not day, clear day
		j.day = make(map[int]struct{})
	default:
		// both day and dayOfWeek are * or both are set, use combined
		// i.e. don't do anything here
	}

	return j, nil
}

// parsePart parse individual schedule part from schedule string
func parsePart(s string, min, max int) (map[int]struct{}, error) {

	r := make(map[int]struct{})

	// wildcard pattern
	if s == "*" {
		for i := min; i <= max; i++ {
			r[i] = struct{}{}
		}
		return r, nil
	}

	// */2 1-59/5 pattern
	if matches := matchN.FindStringSubmatch(s); matches != nil {
		localMin := min
		localMax := max
		if matches[1] != "" && matches[1] != "*" {
			if rng := matchRange.FindStringSubmatch(matches[1]); rng != nil {
				localMin, _ = strconv.Atoi(rng[1])
				localMax, _ = strconv.Atoi(rng[2])
				if localMin < min || localMax > max {
					return nil, newErr("Out of range for", rng[1], "in", s, rng[1], "must be in range", min, "-", max)
				}
			} else {
				return nil, newErr("Unable to parse", matches[1], "part in", s)
			}
		}
		n, _ := strconv.Atoi(matches[2])
		for i := localMin; i <= localMax; i += n {
			r[i] = struct{}{}
		}
		return r, nil
	}

	// 1,2,4  or 1,2,10-15,20,30-45 pattern
	parts := strings.Split(s, ",")
	for _, x := range parts {
		if rng := matchRange.FindStringSubmatch(x); rng != nil {
			localMin, _ := strconv.Atoi(rng[1])
			localMax, _ := strconv.Atoi(rng[2])
			if localMin < min || localMax > max {
				return nil, newErr("Out of range for", x, "in", s, x, "must be in range", min, "-", max)
			}
			for i := localMin; i <= localMax; i++ {
				r[i] = struct{}{}
			}
		} else if i, err := strconv.Atoi(x); err == nil {
			if i < min || i > max {
				return nil, newErr("Out of range for", i, "in", s, i, "must be in range", min, "-", max)
			}
			r[i] = struct{}{}
		} else {
			return nil, newErr("Unable to parse", x, "part in", s)
		}
	}

	if len(r) == 0 {
		return nil, newErr("Unable to parse", s)
	}

	return r, nil
}

// getTick returns the tick struct from time
func getTick(t time.Time) tick {
	return tick{
		min:       t.Minute(),
		hour:      t.Hour(),
		day:       t.Day(),
		month:     int(t.Month()),
		dayOfWeek: int(t.Weekday()),
	}
}
