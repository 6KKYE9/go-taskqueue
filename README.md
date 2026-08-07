# go-taskqueue

FIFO 任务队列。JSON 持久化，支持 push/pop/done/fail 状态流转。

## 用法

```
go-taskqueue -push "download=https://example.com/file.zip"
go-taskqueue -push "convert=output.mp4"
go-taskqueue -pop                        # 取出下一个任务
go-taskqueue -done 1 -fail 2             # 标记完成/失败
go-taskqueue -list                       # 列出全部
go-taskqueue -status pending             # 只看待处理
go-taskqueue -json                       # JSON 输出
```

数据存 `taskqueue.json`。
