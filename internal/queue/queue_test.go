package queue

import (
	"testing"
)

func TestPushPop(t *testing.T) {
	s := testStore(t)
	s.Push("task1", `{"url":"https://go.dev"}`)
	s.Push("task2", `{"url":"https://gitee.com"}`)

	t1, ok := s.Pop()
	if !ok || t1.Name != "task1" {
		t.Fatal("pop 应为 task1")
	}
	t2, ok := s.Pop()
	if !ok || t2.Name != "task2" {
		t.Fatal("pop 应为 task2")
	}
}

func TestPopEmpty(t *testing.T) {
	s := testStore(t)
	_, ok := s.Pop()
	if ok {
		t.Fatal("空队列 pop 应返回 false")
	}
}

func TestDone(t *testing.T) {
	s := testStore(t)
	t1 := s.Push("x", "")
	s.Done(t1.ID)
	if s.Pending() != 0 {
		t.Fatal("完成后的队列应无待处理")
	}
}

func TestFail(t *testing.T) {
	s := testStore(t)
	t1 := s.Push("x", "")
	s.Fail(t1.ID)
	list := s.List("failed")
	if len(list) != 1 {
		t.Fatal("应有 1 个失败任务")
	}
}

func TestListByStatus(t *testing.T) {
	s := testStore(t)
	s.Push("a", "")
	s.Push("b", "")
	t1, _ := s.Pop()
	s.Done(t1.ID)

	if s.Pending() != 1 {
		t.Fatal("应有 1 个待处理")
	}
	if len(s.List("done")) != 1 {
		t.Fatal("应有 1 个已完成")
	}
}

func TestSaveReload(t *testing.T) {
	p := t.TempDir() + "/queue.json"
	s, _ := Load(p)
	s.Push("test", "data")
	s.Save()

	s2, _ := Load(p)
	if len(s2.List("pending")) != 1 {
		t.Fatal("重新加载后应有 1 个待处理")
	}
}

func testStore(t *testing.T) *Store {
	t.Helper()
	s, _ := Load(t.TempDir() + "/queue.json")
	return s
}
