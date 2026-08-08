# go-taskqueue

想临时起个服务、查个 IP、探个端口，还要装一堆东西？没必要。

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
