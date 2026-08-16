# BUG 复现说明（study-room__002）

## Bug 是什么
预约状态机流转失效：待签到订单不能取消/标记爽约，使用中订单不能爽约。

## 如何触发
```bash
go test ./internal/constants -run 'TestCanBookingTransition' -count=1
```

## 错误信息
```
--- FAIL: TestCanBookingTransition
    booking_test.go:21: CanBookingTransition(pending,cancelled) = false, want true
```
