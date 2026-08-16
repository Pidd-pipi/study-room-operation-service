# BUG 复现说明（study-room__003）

## Bug 是什么
分页归一化与预约列表 offset 多处联合失效。

## 如何触发
```bash
go test ./internal/dto -run 'TestListQueryNormalizeDefaults' -count=1
```

## 错误信息
```
--- FAIL: TestListQueryNormalizeDefaults
    pagination_combination_test.go:9: Page = 0, want 1
```
