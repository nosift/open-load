# 如何查看 Key 的组织验证状态

## 方式 1: 通过管理界面（推荐）

### 步骤：
1. 访问管理界面：`http://localhost:3001`
2. 登录后进入"密钥管理"页面
3. 查看 Key 列表

### 显示字段：
每个 Key 会显示以下组织验证相关字段：
- **is_organization_key**: `true`（✅）或 `false`（❌）
- **organization_id**: 组织 ID（如 `org-abc123`），非组织 Key 为空
- **organization_name**: 组织名称，非组织 Key 为空

### 标识说明：
- ✅ `is_organization_key = true`：已通过组织验证，可使用所有模型
- ❌ `is_organization_key = false`：未通过组织验证，使用高级模型时会警告

---

## 方式 2: 通过 API 查询

### 查询某个 Group 的所有 Keys：

```bash
curl -X GET "http://localhost:3001/api/keys?group_id=1" \
  -H "Authorization: Bearer YOUR_AUTH_KEY"
```

### 响应示例：

```json
{
  "items": [
    {
      "id": 1,
      "key_value": "sk-proj-...",
      "status": "active",
      "is_organization_key": true,      ← ✅ 组织验证通过
      "organization_id": "org-abc123",  ← 组织 ID
      "organization_name": "org-abc123",
      "request_count": 150,
      "failure_count": 0
    },
    {
      "id": 2,
      "key_value": "sk-...",
      "status": "active",
      "is_organization_key": false,     ← ❌ 未通过组织验证
      "organization_id": "",            ← 空值
      "organization_name": "",
      "request_count": 80,
      "failure_count": 5
    }
  ]
}
```

---

## 方式 3: 直接查询数据库

### Linux/Mac:

```bash
# 运行查询脚本
./check_org_keys.sh

# 查看详细信息
./check_org_keys.sh -v
```

### Windows:

```cmd
# 运行查询脚本
check_org_keys.bat
```

### 手动 SQL 查询:

```bash
# 进入数据库
sqlite3 data/open-load.db

# 查看所有 Key 的组织状态
SELECT
    id,
    CASE WHEN is_organization_key = 1 THEN '✅ Yes' ELSE '❌ No' END as org_verified,
    organization_id,
    status,
    failure_count,
    request_count
FROM api_keys
ORDER BY is_organization_key DESC, id;
```

### 统计查询:

```sql
-- 统计组织 Key 数量
SELECT
    COUNT(*) as total,
    SUM(CASE WHEN is_organization_key = 1 THEN 1 ELSE 0 END) as org_keys,
    SUM(CASE WHEN is_organization_key = 0 THEN 1 ELSE 0 END) as non_org_keys
FROM api_keys;
```

---

## 方式 4: 查看日志

### Docker 环境:

```bash
# 实时查看日志
docker logs -f open-load

# 查看最近 100 行日志
docker logs --tail 100 open-load
```

### 日志示例:

```
# 使用非组织 Key 访问高级模型
WARN[2025-12-30T11:58:25+08:00] Using non-organization key for premium model - may fail and trigger retry  key_id=123 model=gpt-5.2

# 使用组织 Key 访问高级模型（成功）
DEBUG[2025-12-30T11:58:26+08:00] Premium model access with organization key  key_id=789 model=gpt-5.2 organization_id=org-abc123
```

### 日志过滤:

```bash
# 只看组织验证相关日志
docker logs open-load 2>&1 | grep -i "organization"

# 只看警告日志
docker logs open-load 2>&1 | grep "WARN.*organization"

# 只看成功的组织 Key 使用
docker logs open-load 2>&1 | grep "organization key.*key_id"
```

---

## 快速区分方法总结

| 方式 | 优点 | 缺点 |
|------|------|------|
| **管理界面** | 直观、可视化、易操作 | 需要登录 |
| **API 查询** | 可编程、可集成 | 需要写脚本 |
| **数据库查询** | 完整信息、灵活 | 需要数据库访问权限 |
| **日志查看** | 实时监控、动态查看 | 信息较多，需要过滤 |

---

## 常见问题

### Q1: 为什么某些 Key 显示 `is_organization_key = false`？

**可能原因**：
1. 该 Key 确实不属于任何组织
2. 使用了代理，OpenAI API 没有返回组织头信息
3. Key 还未经过验证（新添加的 Key）

**解决方案**：
- 手动触发 Key 验证
- 等待定时验证任务运行
- 检查是否使用了代理

### Q2: 如何手动触发 Key 验证？

在管理界面中：
1. 进入"密钥管理"页面
2. 选择要验证的 Key
3. 点击"手动验证"按钮

### Q3: 验证后仍显示 `false`，怎么办？

可能该 Key 确实不属于组织。验证方法：
1. 检查日志，看是否检测到组织头
2. 使用该 Key 直接调用 OpenAI API，查看响应头
3. 在 OpenAI 官网确认该 Key 的组织状态

### Q4: 如何批量查看某个 Group 的组织 Key 比例？

```sql
SELECT
    group_id,
    COUNT(*) as total_keys,
    SUM(CASE WHEN is_organization_key = 1 THEN 1 ELSE 0 END) as org_keys,
    ROUND(100.0 * SUM(CASE WHEN is_organization_key = 1 THEN 1 ELSE 0 END) / COUNT(*), 2) as org_percentage
FROM api_keys
GROUP BY group_id;
```

---

## 最佳实践

1. **定期检查**：定期查看 Key 的组织验证状态
2. **监控日志**：关注 WARN 日志，了解哪些 Key 经常失败
3. **分组管理**：将组织 Key 和非组织 Key 分到不同的 Group
4. **标记备注**：在 Key 的 `notes` 字段中标记是否为组织 Key
5. **配置高级模型**：只在 `premium_models` 中配置真正需要组织验证的模型

---

## 示例：完整的检查流程

### 步骤 1: 查看统计

```bash
# Linux/Mac
./check_org_keys.sh

# Windows
check_org_keys.bat
```

### 步骤 2: 查看详细信息

```bash
# API 查询
curl -X GET "http://localhost:3001/api/keys?group_id=1" \
  -H "Authorization: Bearer YOUR_AUTH_KEY" | jq '.items[] | {id, is_organization_key, organization_id, status}'
```

### 步骤 3: 查看实时日志

```bash
# 监控组织 Key 使用情况
docker logs -f open-load 2>&1 | grep --color "organization"
```

### 步骤 4: 更新配置（如果需要）

```sql
-- 配置高级模型
INSERT INTO system_settings (setting_key, setting_value, description)
VALUES ('premium_models', 'gpt-5.2,o1-preview', '需要组织验证的高级模型')
ON DUPLICATE KEY UPDATE setting_value = 'gpt-5.2,o1-preview';
```

---

**提示**: 所有查询脚本和 SQL 示例都已包含在项目中，可直接使用。
