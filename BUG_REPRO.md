# BUG 复现说明（study-room__004）

## Bug 是什么
用户查询/登录/注册 nil 与错误码联合失效：未找到返回 (nil,nil)，登录 nil 用户 panic，重复注册错误码错误。

## 如何触发
```bash
go test ./internal/service -run 'TestGetByIDNilReturnsError|TestLoginNilReturnsUnauthorized|TestRegisterDuplicateConflict' -count=1
```

## 错误信息
```
--- FAIL: TestGetByIDNilReturnsError
    user_combination_test.go:34: expected error for nil user
panic: runtime error: invalid memory address or nil pointer dereference
```
