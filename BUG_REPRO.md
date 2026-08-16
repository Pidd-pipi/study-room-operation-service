# BUG 复现说明（study-room__001）

## Bug 是什么
座位分区/状态校验失效：非法分区可录入，非法座位状态可变更。

## 如何触发
```bash
go test ./internal/service -run 'TestSeatServiceCreate|TestSeatServiceChangeStatus' -count=1
```

## 错误信息
```
--- FAIL: TestSeatServiceCreate/invalid_zone
    seat_service_test.go:79: expected error
--- FAIL: TestSeatServiceChangeStatus
    seat_service_test.go:110: expected validation error, got <nil>
```
