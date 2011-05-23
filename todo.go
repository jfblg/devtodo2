/*
  Copyright 2011 Alec Thomas

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package main

import (
	"time"
	"container/list"
)

type Priority int

const (
	VERYHIGH = Priority(4)
	HIGH = Priority(3)
	MEDIUM = Priority(2)
	LOW = Priority(1)
	VERYLOW = Priority(0)
)

var priorityMapFromString map[string]Priority = map[string]Priority {
	"veryhigh": VERYHIGH,
	"high": HIGH,
	"medium": MEDIUM,
	"low": LOW,
	"verylow": VERYLOW,
}

var priorityToString map[Priority]string = map[Priority]string {
	VERYHIGH: "veryhigh",
	HIGH: "high",
	MEDIUM: "medium",
	VERYLOW: "verylow",
	LOW: "low",
}

func (p Priority) String() string {
	return priorityToString[p]
}

func PriorityFromString(priority string) Priority {
	if p, ok := priorityMapFromString[priority]; ok {
		return p
	}
	return MEDIUM
}

type TaskIterator interface {
	Next() TaskIterator
	Task() Task
}

type Task interface {
	// Return an iterator over child tasks. nil if no children.
	Begin() TaskIterator
	AddTask(text string, priority Priority) Task

	Text() string
	SetText(text string)

	Priority() Priority
	SetPriority(priority Priority)

	SetCreationTime(time *time.Time)
	CreationTime() *time.Time

	SetCompletionTime(time *time.Time)
	CompletionTime() *time.Time
}

type TaskList Task

// Implementation
type taskIteratorImpl struct {
	cursor *list.Element
}

func (i *taskIteratorImpl) Next() TaskIterator {
	i.cursor = i.cursor.Next()
	if i.cursor == nil {
		return nil
	}
	return i
}

func (i *taskIteratorImpl) Task() Task {
	return i.cursor.Value.(*taskImpl)
}

type taskImpl struct {
	text string
	priority Priority
	created, completed *time.Time
	tasks *list.List
}

func newTask(text string, priority Priority) Task {
	return &taskImpl{
		text: text,
		priority: priority,
		created: time.UTC(),
		completed: nil,
		tasks: list.New(),
	}
}

func (t *taskImpl) AddTask(text string, priotity Priority) Task {
	task := newTask(text, priotity)
	t.tasks.PushBack(task)
	return task
}

func (t *taskImpl) Begin() TaskIterator {
	front := t.tasks.Front()
	if front == nil {
		return nil
	}
	return &taskIteratorImpl{cursor: front}
}

func (t *taskImpl) SetCreationTime(time *time.Time) {
	t.created = time
}

func (t *taskImpl) CreationTime() *time.Time {
	return t.created
}

func (t *taskImpl) SetCompletionTime(time *time.Time) {
	t.completed = time
}

func (t *taskImpl) CompletionTime() *time.Time {
	return t.completed
}

func (t *taskImpl) Text() string {
	return t.text
}

func (t *taskImpl) SetText(text string) {
	t.text = text
}

func (t *taskImpl) Priority() Priority {
	return t.priority
}

func (t *taskImpl) SetPriority(priority Priority) {
	t.priority = priority
}

func NewTaskList() TaskList {
	return newTask("", MEDIUM).(TaskList)
}
