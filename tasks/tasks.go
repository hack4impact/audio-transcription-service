// Package tasks implements a basic task queue.
package tasks

import (
	"math/rand"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

// Status is the status of the task.
type Status int

// TaskExecuter executes a series of task functions.
type TaskExecuter interface {
	QueueTask(task func(string) error, onFailure func(string, string)) string
	GetTaskStatus(id string) Status
	completeTask(id string, task func(string) error, onFailure func(string, string))
}

type taskInfo struct {
	status  Status
	started time.Time
}

type concurrentTaskInfoMap struct {
	sync.RWMutex
	m map[string]taskInfo
}

type defaultExecuter struct {
	cMap       concurrentTaskInfoMap
	expiration time.Duration
}

// These are some enumerated Status constants.
// INPROGRESS: Task is still in progress.
// SUCCESS: Task finished successfully.
// FAILURE: Task finished unsuccessfully.
// NOTFOUND: Task could not be found.
const (
	INPROGRESS Status = iota
	SUCCESS
	FAILURE
	NOTFOUND
)

// DefaultTaskExecuter is an instance of a NewTaskExecuter with a 24-hour
// expiration.
var DefaultTaskExecuter = NewTaskExecuter(time.Hour * 24)

// put(k,v) maps k to v in the map
func (c *concurrentTaskInfoMap) put(k string, v taskInfo) {
	c.Lock()
	c.m[k] = v
	c.Unlock()
}

// get(k) returns the value of k in the map
func (c *concurrentTaskInfoMap) get(k string) (taskInfo, bool) {
	c.RLock()
	v, ok := c.m[k]
	c.RUnlock()
	return v, ok
}

// setStatus sets the status of k if it is already in the map
func (c *concurrentTaskInfoMap) setStatus(k string, s Status) {
	if info, ok := c.get(k); ok {
		info.status = s
		c.put(k, info)
		return
	}
}

func (s Status) String() string {
	var str string

	switch s {
	case INPROGRESS:
		str = "The task is in progress."
	case SUCCESS:
		str = "The task completed successfully."
	case FAILURE:
		str = "The task failed."
	case NOTFOUND:
		str = "Error: task not found."
	}
	return str
}

// NewTaskExecuter returns a TaskExecuter ready to execute. Information for a
// task is deleted after expiration.
func NewTaskExecuter(expiration time.Duration) TaskExecuter {
	ex := &defaultExecuter{
		cMap:       concurrentTaskInfoMap{m: make(map[string]taskInfo)},
		expiration: expiration,
	}
	go ex.deleteExpiredInfo()

	return ex
}

// QueueTask initializes a new task, taking a generic task function. If the
// task panics, the panic will be caught. However, if the task launches another
// goroutine which panics, the panic cannot be caught.
func (ex *defaultExecuter) QueueTask(task func(string) error, onFailure func(string, string)) string {
	id := generateID(20)
	ex.cMap.put(id, taskInfo{
		status:  INPROGRESS,
		started: time.Now(),
	})
	log.WithField("task", id).
		Info("Task started")
	go ex.completeTask(id, task, onFailure)
	return id
}

// GetTaskStatus gets the current status of a task.
func (ex *defaultExecuter) GetTaskStatus(id string) Status {
	if info, ok := ex.cMap.get(id); ok {
		return info.status
	}
	return NOTFOUND
}

func (ex *defaultExecuter) completeTask(id string, task func(string) error, onFailure func(string, string)) {
	defer func() {
		if r := recover(); r != nil {
			log.WithField("task", id).
				Error("Task failed")
			go onFailure(id, "The error message is below. Please check logs for more details."+"\n\n"+"panic occurred")
			ex.cMap.setStatus(id, FAILURE)
		}
	}()

	// Run the task.
	if err := task(id); err != nil {
		log.WithFields(log.Fields{
			"task":  id,
			"error": errors.ErrorStack(err),
		}).Error("Task failed")
		go onFailure(id, "The error message is below. Please check logs for more details."+"\n\n"+errors.ErrorStack(err))
		ex.cMap.setStatus(id, FAILURE)
		return
	}

	log.WithField("task", id).
		Info("Task succeeded")
	ex.cMap.setStatus(id, SUCCESS)
}

func (ex *defaultExecuter) deleteExpiredInfo() {
	for range time.Tick(30 * time.Minute) {
		m := ex.cMap.m
		toDelete := []string{}

		ex.cMap.RLock()
		for k, v := range m {
			if (time.Since(v.started)) > ex.expiration {
				toDelete = append(toDelete, k)
			}
		}
		ex.cMap.RUnlock()

		ex.cMap.Lock()
		for _, k := range toDelete {
			log.WithField("task", k).
				Debug("Expired from info map")
			delete(m, k)
		}
		ex.cMap.Unlock()
	}
}

// Borrowed from https://siongui.github.io/2015/04/13/go-generate-random-string/
func generateID(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
