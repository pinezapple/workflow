```mermaid
sequenceDiagram
    Definition Programs ->> Workflow Engine: Yêu cầu A
    Workflow Engine -->> Definition Programs: Đã nhận A
    Workflow Engine ->> Event Storage: Đẩy A vào Event Storage
    Event Storage ->> Worker Programs: Đọc A
    Worker Programs ->> Worker Programs: Thực thi A
```