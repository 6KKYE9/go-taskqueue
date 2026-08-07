package queue

import (
	"encoding/json"
	"os"
	"time"
)

type Task struct {
	ID      int64      `json:"id"`
	Name    string     `json:"name"`
	Payload string     `json:"payload"`
	Status  string     `json:"status"` // pending, running, done, failed
	Added   time.Time  `json:"added"`
	DoneAt  *time.Time `json:"done_at,omitempty"`
}

type Store struct {
	path string
	list []Task
	next int64
}

func Load(path string) (*Store, error) {
	s := &Store{path: path, next: 1}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return s, nil
	}
	if err := json.Unmarshal(data, &s.list); err != nil {
		return nil, err
	}
	for _, t := range s.list {
		if t.ID >= s.next {
			s.next = t.ID + 1
		}
	}
	return s, nil
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func (s *Store) Push(name, payload string) Task {
	t := Task{
		ID:      s.next,
		Name:    name,
		Payload: payload,
		Status:  "pending",
		Added:   time.Now().Truncate(time.Second),
	}
	s.next++
	s.list = append(s.list, t)
	return t
}

func (s *Store) Pop() (Task, bool) {
	for i, t := range s.list {
		if t.Status == "pending" {
			s.list[i].Status = "running"
			return s.list[i], true
		}
	}
	return Task{}, false
}

func (s *Store) Done(id int64) {
	for i := range s.list {
		if s.list[i].ID == id {
			now := time.Now().Truncate(time.Second)
			s.list[i].Status = "done"
			s.list[i].DoneAt = &now
			return
		}
	}
}

func (s *Store) Fail(id int64) {
	for i := range s.list {
		if s.list[i].ID == id {
			now := time.Now().Truncate(time.Second)
			s.list[i].Status = "failed"
			s.list[i].DoneAt = &now
			return
		}
	}
}

func (s *Store) List(status string) []Task {
	if status == "" {
		out := make([]Task, len(s.list))
		copy(out, s.list)
		return out
	}
	var out []Task
	for _, t := range s.list {
		if t.Status == status {
			out = append(out, t)
		}
	}
	return out
}

func (s *Store) Pending() int {
	n := 0
	for _, t := range s.list {
		if t.Status == "pending" {
			n++
		}
	}
	return n
}
