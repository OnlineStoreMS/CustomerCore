# CustomerCore（客户中心）

OSMS 平台客户底库：客户档案、收货地址、全渠道身份绑定；供门店/订单/商城等服务调用。

| 项 | 值 |
|----|-----|
| Go module | `customercore` |
| API | `:8099` |
| Web | `:5183` |
| Docker 镜像 | `customercore-api`、`customercore-web` |
| UserCore app | `customercore`（`customer:read` / `customer:write`） |
| 端口约定 | [deploy/docs/PORTS.md](../deploy/docs/PORTS.md) |
| 平台编排 | `/home/asialeaf/projects/deploy` |

## 本地开发

```bash
cp configs/config.example.yaml configs/config.yaml
# 编辑 postgres_dsn / jwt_secret

go run ./cmd/api -config configs/config.yaml

cd web && npm i && npm run dev
```

## API

- `GET /health`
- Admin（JWT）：`/api/v1/admin/customers`、地址、绑定、`GET /dashboard/stats`
- Internal（平台内调用）：`POST /api/v1/internal/customers/upsert-by-phone`、`GET .../by-phone`、`GET .../:id`
