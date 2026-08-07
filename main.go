package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"go-taskqueue/internal/queue"
)

var storePath = "taskqueue.json"

func run() error {
	push := flag.String("push", "", "添加任务（格式: name=payload）")
	pop := flag.Bool("pop", false, "取出一个待处理任务")
	done := flag.Int64("done", 0, "标记任务完成")
	fail := flag.Int64("fail", 0, "标记任务失败")
	list := flag.Bool("list", false, "列出所有任务")
	status := flag.String("status", "", "按状态筛选: pending/running/done/failed")
	jsonOut := flag.Bool("json", false, "JSON 输出")
	flag.Parse()

	s, err := queue.Load(storePath)
	if err != nil {
		return err
	}

	switch {
	case *push != "":
		name, payload := *push, ""
		for i, c := range *push {
			if c == '=' {
				name = (*push)[:i]
				payload = (*push)[i+1:]
				break
			}
		}
		t := s.Push(name, payload)
		s.Save()
		fmt.Printf("已添加 #%d: %s\n", t.ID, t.Name)
	case *pop:
		t, ok := s.Pop()
		if !ok {
			fmt.Println("无待处理任务")
		} else {
			s.Save()
			data, _ := json.MarshalIndent(t, "", "  ")
			fmt.Println(string(data))
		}
	case *done > 0:
		s.Done(*done)
		s.Save()
		fmt.Printf("#%d 已完成\n", *done)
	case *fail > 0:
		s.Fail(*fail)
		s.Save()
		fmt.Printf("#%d 已失败\n", *fail)
	case *jsonOut:
		data, _ := json.MarshalIndent(s.List(*status), "", "  ")
		fmt.Println(string(data))
	case *list || *status != "":
		tasks := s.List(*status)
		for _, t := range tasks {
			fmt.Printf("#%d [%s] %s\n", t.ID, t.Status, t.Name)
		}
	default:
		fmt.Printf("待处理: %d\n", s.Pending())
		tasks := s.List("")
		for _, t := range tasks {
			fmt.Printf("#%d [%s] %s\n", t.ID, t.Status, t.Name)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
