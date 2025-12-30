# 组织验证和高级模型权限控制功能说明

## ⚠️ 重要提示

**当前实现策略：监控但不阻止**

- ✅ 自动检测并记录组织验证信息
- ✅ 完全保留现有的重试机制
- ✅ 不会影响你的无感使用体验
- ✅ 通过日志和管理界面查看哪些 Key 通过了组织验证
- ❌ **不会阻止**非组织 Key 使用高级模型

这意味着即使配置了 `premium_models`，非组织 Key 仍然可以尝试使用高级模型，只是会在日志中记录警告。如果失败，会自动重试其他 Key，就像之前一样。

---

## 功能概述

本功能实现了基于组织验证的 API Key 权限控制，允许你：
1. 自动检测 API Key 是否通过 OpenAI 组织验证
2. 配置高级模型列表（如 gpt-5.2）
3. 限制只有组织验证的 Key 才能使用高级模型

## 实现细节

### 1. 数据库字段

在 `api_keys` 表中新增了三个字段：
- `is_organization_key` (boolean): 标识该 Key 是否为组织验证的 Key
- `organization_id` (varchar): 组织 ID
- `organization_name` (varchar): 组织名称

### 2. 自动检测逻辑

在验证 API Key 时（`ValidateKey`），系统会：
1. 发送测试请求到 OpenAI API
2. 检查响应头中的 `openai-organization` 或 `x-openai-organization`
3. 如果存在组织信息，自动更新数据库中的组织字段
4. 每次验证都会更新这些信息

### 3. 配置高级模型

在系统设置中添加 `premium_models` 配置项：

**格式**：用逗号分隔的模型列表

**示例**：
```
gpt-5.2,gpt-5.0,claude-opus-4
```

**配置方式**：
- 通过管理界面：系统设置 -> 密钥配置 -> premium_models
- 通过数据库：在 `system_settings` 表中添加或更新记录

```sql
INSERT INTO system_settings (setting_key, setting_value, description)
VALUES ('premium_models', 'gpt-5.2,gpt-5.0', '需要组织验证的高级模型列表')
ON DUPLICATE KEY UPDATE setting_value = 'gpt-5.2,gpt-5.0';
```

### 4. 权限检查流程（监控但不阻止）

**当前策略**：监控并记录，但不阻止请求，保持重试机制。

当用户请求使用某个模型时：
1. 系统从请求体中提取模型名称
2. 检查该模型是否在 `premium_models` 列表中
3. 如果是高级模型：
   - 检查使用的 API Key 的 `is_organization_key` 字段
   - 如果为 `false`：
     - ⚠️ 记录 **WARN** 日志：使用非组织 Key 访问高级模型
     - ✅ **继续发送请求**（不阻止）
     - ✅ 如果失败，依赖重试机制选择其他 Key
   - 如果为 `true`：
     - ✅ 记录 **DEBUG** 日志：使用组织 Key 访问
     - ✅ 正常发送请求
4. 如果不是高级模型，直接允许通过

**重要**：此功能**不会中断**你现有的重试机制，只是增加了可见性。

## 使用方法

### 步骤 1: 运行数据库迁移

系统启动时会自动运行迁移，添加新的数据库字段。

### 步骤 2: 配置高级模型列表

在系统设置中配置 `premium_models`，例如：
```
gpt-5.2,gpt-o1,claude-opus-4.5
```

### 步骤 3: 验证 API Keys

- 新添加的 Key 在首次验证时会自动检测组织信息
- 已有的 Key 会在下次定时验证时更新组织信息
- 手动验证 Key 也会触发组织信息检测

### 步骤 4: 监控使用情况

配置完成后，系统会自动记录组织验证信息：

**日志示例**：
```
# 使用非组织 Key 尝试访问高级模型
WARN[2024-12-30T10:58:52+08:00] Using non-organization key for premium model - may fail and trigger retry  key_id=123 model=gpt-5.2

# OpenAI 返回 429 错误，自动重试
WARN[2024-12-30T10:58:53+08:00] Using non-organization key for premium model - may fail and trigger retry  key_id=456 model=gpt-5.2

# 选中组织 Key，请求成功
DEBUG[2024-12-30T10:58:54+08:00] Premium model access with organization key  key_id=789 model=gpt-5.2 organization_id=org-abc123
```

**实际效果**：
- 从日志可以清晰看到哪些 Key 是组织验证的
- 不会影响现有的重试机制
- 最终请求会成功（只要有可用的组织 Key）
- 完全无感使用

## 查看 Key 的组织信息

### 方式 1：通过管理界面

在密钥管理页面，每个 Key 的详细信息中会显示：
- **is_organization_key**: `true` 或 `false`（已通过/未通过组织验证）
- **organization_id**: 组织 ID（如 `org-xxxxx`）
- **organization_name**: 组织名称

**标识说明**：
- ✅ 标记为组织 Key：可以使用所有模型（包括高级模型）
- ❌ 未标记为组织 Key：可能无法使用高级模型（取决于 `premium_models` 配置）

### 方式 2：通过 API

```bash
curl -X GET "http://your-server/api/keys?group_id={group_id}" \
  -H "Authorization: Bearer your-admin-key"
```

响应示例：
```json
{
  "items": [
    {
      "id": 1,
      "key_value": "sk-...",
      "status": "active",
      "is_organization_key": true,
      "organization_id": "org-abc123",
      "organization_name": "org-abc123",
      "notes": "",
      "request_count": 150,
      "failure_count": 0
    },
    {
      "id": 2,
      "key_value": "sk-...",
      "status": "active",
      "is_organization_key": false,
      "organization_id": "",
      "organization_name": "",
      "notes": "",
      "request_count": 80,
      "failure_count": 5
    }
  ]
}
```

### 方式 3：查看日志

在请求日志中会显示：
```
WARN[2024-01-15T10:58:52+08:00] Using non-organization key for premium model - may fail and trigger retry  key_id=123 model=gpt-5.2
DEBUG[2024-01-15T10:58:53+08:00] Premium model access with organization key  key_id=456 model=gpt-5.2 organization_id=org-xxxxx
```

## 注意事项

1. **组织信息检测依赖于 OpenAI API 响应头**
   - 只有 OpenAI 官方 API 会返回组织信息
   - 如果使用第三方代理，可能无法检测到组织信息

2. **首次检测**
   - 新添加的 Key 需要经过一次验证才能检测到组织信息
   - 建议在添加 Key 后立即执行手动验证

3. **性能影响**
   - 组织信息检测在验证过程中进行，不会增加额外的 API 调用
   - 权限检查在内存中进行，几乎没有性能开销

4. **空配置行为**
   - 如果 `premium_models` 为空或未配置，所有模型都不受限制
   - 这确保了向后兼容性

## 错误排查

### 问题：所有请求都被拒绝

**解决方案**：
- 检查 `premium_models` 配置是否过于宽泛
- 确认 API Keys 已经过验证并检测到组织信息

### 问题：组织 Key 未被识别

**解决方案**：
- 手动触发 Key 验证
- 检查日志，确认 OpenAI API 是否返回了组织头
- 确认使用的是 OpenAI 官方 API，而非第三方代理

### 问题：验证后仍显示 is_organization_key = false

**解决方案**：
- 确认该 Key 确实属于某个组织
- 检查 OpenAI 响应头中是否包含 `openai-organization` 或 `x-openai-organization`
- 查看验证日志，确认检测逻辑是否正常执行

## 技术细节

### 修改的文件

1. `internal/models/types.go` - 添加数据库字段定义
2. `internal/types/types.go` - 添加系统配置字段
3. `internal/channel/openai_channel.go` - 实现组织信息检测
4. `internal/keypool/provider.go` - 保存组织信息到数据库
5. `internal/config/system_settings.go` - 初始化高级模型配置缓存
6. `internal/proxy/server.go` - 实现模型权限检查
7. `internal/db/migrations/v1_2_0_AddOrganizationFields.go` - 数据库迁移
8. `internal/db/migrations/migration.go` - 注册迁移

### 数据流

```
用户请求
  ↓
提取模型名称
  ↓
检查是否为高级模型？
  ├─ 否 → 允许通过
  └─ 是 → 检查 Key 是否为组织 Key？
           ├─ 是 → 允许通过
           └─ 否 → 拒绝请求
```

## 未来改进

可能的扩展方向：
1. 支持更细粒度的权限控制（基于 Group 配置不同的高级模型列表）
2. 支持其他 AI 提供商的组织验证
3. 添加组织信息的手动设置接口
4. 实现基于组织的使用统计和限额
