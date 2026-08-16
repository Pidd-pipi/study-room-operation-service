# BUG 复现说明（study-room__005）

## Bug 是什么
座位状态/分区/预约状态文案多处错位。

## 如何触发
```bash
go test ./internal/util -run 'TestSeatAndBookingText'
```

## 错误信息
```
--- FAIL: TestSeatAndBookingText
    seat_booking_text_test.go:11: SeatUnavailable.Valid() should be true
```
