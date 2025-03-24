# LangGraph SDK for Go

LangGraph SDK là một thư viện Go giúp bạn dễ dàng tương tác với API LangGraph. SDK này cung cấp các client để quản lý các tài nguyên cốt lõi như Assistants, Threads, Runs, Cron Jobs và Store.

## Mục lục

- [Cài đặt](#cài-đặt)
- [Cách sử dụng](#cách-sử-dụng)
  - [Khởi tạo Client](#khởi-tạo-client)
  - [Quản lý Assistants](#quản-lý-assistants)
  - [Quản lý Threads](#quản-lý-threads)
  - [Quản lý Runs](#quản-lý-runs)
  - [Quản lý Cron Jobs](#quản-lý-cron-jobs)
  - [Quản lý Store](#quản-lý-store)
- [Đóng góp](#đóng-góp)
- [Giấy phép](#giấy-phép)

---

## Cài đặt

Để sử dụng LangGraph SDK, bạn cần cài đặt Go phiên bản `1.23.4` trở lên. Thêm SDK vào dự án của bạn bằng cách chạy lệnh sau:

```bash
go get langgraph-sdk