# 变更提案: add-ghcr-trigger-for-sync-origin-main-local7-20260318

## 元信息
```yaml
类型: 修复
方案类型: implementation
优先级: P1
状态: 已确认
创建: 2026-04-24
```

---

## 1. 需求

### 背景
当前 `.github/workflows/docker-image-ghcr.yml` 仅在 `main` 分支更新时自动执行。仓库里实际存在 `sync/origin-main-with-local7-20260318` 同步分支，用户希望该分支更新时也能自动触发 GHCR 镜像构建，但又不能覆盖 `main/latest` 主线镜像标签。

### 目标
- 为 `sync/origin-main-with-local7-20260318` 增加自动触发 `docker-image-ghcr.yml` 的能力。
- 保持 `main` 分支的现有 GHCR 标签行为不变。
- 为同步分支发布独立标签，避免覆盖 `main/latest`。

### 约束条件
```yaml
兼容性约束: workflow_dispatch 手动触发也要复用同一套标签解析逻辑
安全约束: sync 分支生成的镜像标签不能覆盖 main/latest
实现约束: 只做最小范围的 CI workflow 与知识库同步，不改动业务代码
```

### 验收标准
- [ ] push 到 `sync/origin-main-with-local7-20260318` 时自动触发 `.github/workflows/docker-image-ghcr.yml`
- [ ] `main` 分支仍发布 `main` / `latest` / `sha-*` 标签
- [ ] `sync/origin-main-with-local7-20260318` 分支发布独立的 `origin-main-with-local7-20260318` / `latest-origin-main-with-local7-20260318` / `sha-*` 标签
- [ ] workflow YAML 结构合法，且标签解析逻辑可读可维护

---

## 2. 方案

### 技术方案
在 `docker-image-ghcr.yml` 的 `push.branches` 中增加 `sync/origin-main-with-local7-20260318`。同时在构建与 manifest 两个 job 中增加统一的标签解析步骤：

- `main` 分支继续映射为 `main` / `latest`
- `sync/origin-main-with-local7-20260318` 特判映射为 `origin-main-with-local7-20260318` / `latest-origin-main-with-local7-20260318`
- 其他手动触发分支使用安全降级分支标签（将 `/` 替换为 `-`）

这样可以在保留主线标签语义的同时，为同步分支输出独立镜像标签。

### 影响范围
```yaml
涉及模块:
  - .github/workflows/docker-image-ghcr.yml: GHCR workflow 触发分支与标签解析
  - .helloagents/modules/ci-ghcr-workflow.md: 记录主线/同步分支标签策略
预计变更文件: 4
```

### 风险评估
| 风险 | 等级 | 应对 |
|------|------|------|
| sync 分支误复用主线标签 | 中 | 在 workflow 中显式分支映射，独立生成 PRIMARY_TAG/LATEST_TAG |
| 手动触发不同分支时标签不可用 | 低 | 增加通用 fallback，将 `/` 规范化为 `-` |
| 文档未同步导致后续维护误判 | 低 | 同步新增知识库模块文档，记录标签隔离规则 |

---

## 3. 技术设计（可选）

### 架构设计
N/A

### API设计
N/A

### 数据模型
| 字段 | 类型 | 说明 |
|------|------|------|
| N/A | N/A | 本次无数据模型变更 |

---

## 4. 核心场景

### 场景: 主线分支持续发布正式 GHCR 标签
**模块**: .github/workflows/docker-image-ghcr.yml
**条件**: `main` 分支 push 或手动从 `main` 触发 workflow
**行为**: 执行 GHCR 单架构构建与 manifest 合并
**结果**: 发布 `main`、`latest` 与 `sha-*` 标签

### 场景: 同步分支自动发布独立 GHCR 标签
**模块**: .github/workflows/docker-image-ghcr.yml
**条件**: `sync/origin-main-with-local7-20260318` 分支 push
**行为**: 执行同一 GHCR workflow，但按分支映射解析镜像标签
**结果**: 发布 `origin-main-with-local7-20260318`、`latest-origin-main-with-local7-20260318` 与 `sha-*` 标签，不覆盖主线标签

---

## 5. 技术决策

### add-ghcr-trigger-for-sync-origin-main-local7-20260318#D001: 同步分支使用独立标签而非复用 main/latest
**日期**: 2026-04-24
**状态**: ✅采纳
**背景**: 用户要求同步分支也自动执行 GHCR workflow，但明确选择了“分支独立标签”，以避免覆盖 `main/latest` 主线镜像。
**选项分析**:
| 选项 | 优点 | 缺点 |
|------|------|------|
| A: 复用 `main/latest` | 实现最简单 | 会覆盖主线镜像，风险高 |
| B: 仅推 `sha-*` | 风险低 | 不利于持续引用同步分支镜像 |
| C: 为同步分支提供独立稳定标签 | 可持续引用，且不影响主线 | 需要增加标签解析逻辑 |
**决策**: 选择方案 C
**理由**: 既满足“自动执行”，又满足“不可覆盖主线镜像”的要求，且便于后续 CI/CD 使用固定分支标签拉取测试镜像。
**影响**: 影响 GHCR workflow 的分支触发条件、标签生成逻辑与知识库文档。

---

## 6. 成果设计

N/A
