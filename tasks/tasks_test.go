package tasks

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskErrorLeadsToErrorStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func(a string) error {
		return errors.New("This is the error text.")
	}

	ex := NewTaskExecuter()
	id := ex.QueueTask(errorTask, func(a, b string) {})
	status := ex.GetTaskStatus(id)
	for status == INPROGRESS {
		status = ex.GetTaskStatus(id)
	}
	assert.Equal(FAILURE, status)
}

func TestTaskPanicLeadsToErrorStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func(a string) error {
		panic("AHHH!!!")
	}

	ex := NewTaskExecuter()
	id := ex.QueueTask(errorTask, func(a, b string) {})
	status := ex.GetTaskStatus(id)
	for status == INPROGRESS {
		status = ex.GetTaskStatus(id)
	}
	assert.Equal(FAILURE, status)
}

func TestTaskOkLeadsToSuccessStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func(a string) error {
		return nil
	}

	ex := NewTaskExecuter()
	id := ex.QueueTask(errorTask, func(a, b string) {})
	status := ex.GetTaskStatus(id)
	for status == INPROGRESS {
		status = ex.GetTaskStatus(id)
	}
	assert.Equal(SUCCESS, status)
}

func TestInProgressStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func(a string) error {
		for true {
		}
		return nil
	}

	ex := NewTaskExecuter()
	id := ex.QueueTask(errorTask, func(a, b string) {})
	status := ex.GetTaskStatus(id)
	assert.Equal(INPROGRESS, status)
}
