# gow
🙄 gow!!! the micro web framework.

## Hello Gow

```go
g := gow.Me()
g.Get("/", func(ctx gow.Context) {
    ctx.Text("Hello World")
})
g.Listen(":10777")
```

## Template

```go
g := gow.Me()
g.Get("/", func(ctx gow.Context) {
    fmap := template.FuncMap{
        "Name": "gow",
    }
    ctx.Render("index", fmap)
})
```

## JSON

```go
g := gow.Me()
g.Get("/", func(ctx gow.Context) {
    fmap := template.FuncMap{
        "u1": "jack1",
        "u2": "jack2",
        "u3": "jack3",
    }
    ctx.Json(fmap)
})
```

## 后续计划

- session管理
- 参数优化
- 多模板支持
- 配置优化
