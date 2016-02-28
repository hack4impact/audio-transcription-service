package tasks

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskErrorLeadsToErrorStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func() error {
		return errors.New("This is the error text.")
	}
	ex := NewTaskExectuer()
	id := ex.QueueTask(errorTask)
	status := ex.GetTaskStatus(id)
	for status == INPROGRESS {
		status = ex.GetTaskStatus(id)
	}
	assert.Equal(FAILURE, status)
}

func TestTaskPanicLeadsToErrorStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func() error {
		panic("AHHH!!!")
	}
	ex := NewTaskExectuer()
	id := ex.QueueTask(errorTask)
	status := ex.GetTaskStatus(id)
	for status == INPROGRESS {
		status = ex.GetTaskStatus(id)
	}
	assert.Equal(FAILURE, status)
}

func TestTaskOkLeadsToSuccessStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func() error {
		return nil
	}
	ex := NewTaskExectuer()
	id := ex.QueueTask(errorTask)
	status := ex.GetTaskStatus(id)
	for status == INPROGRESS {
		status = ex.GetTaskStatus(id)
	}
	assert.Equal(SUCCESS, status)
}

func TestInProgressStatus(t *testing.T) {
	assert := assert.New(t)
	errorTask := func() error {
		for true {
		}
		return nil
	}
	ex := NewTaskExectuer()
	id := ex.QueueTask(errorTask)
	status := ex.GetTaskStatus(id)
	assert.Equal(INPROGRESS, status)
}
