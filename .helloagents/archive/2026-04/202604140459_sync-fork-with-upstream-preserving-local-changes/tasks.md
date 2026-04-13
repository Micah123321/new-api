# 任务清单: sync-fork-with-upstream-preserving-local-changes

> **@status:** completed | 2026-04-14 05:13

```yaml
@feature: sync-fork-with-upstream-preserving-local-changes
@created: 2026-04-14
@status: completed
@mode: R2
```

## 进度概览

| 完成 | 失败 | 跳过 | 总数 |
|------|------|------|------|
| 6 | 0 | 0 | 6 |

---

## 任务列表

### 1. Git 同步准备

- [√] 1.1 抓取最新 `origin` / `fork` 引用并确认当前同步分支的分叉状态 | depends_on: []
- [√] 1.2 在当前分支合并 `origin/main`，生成真实冲突现场 | depends_on: [1.1]

### 2. 冲突解析

- [√] 2.1 解析 [`.github/workflows/docker-image-alpha.yml`](/E:/code/go/new-api/.github/workflows/docker-image-alpha.yml) 与 [`.github/workflows/docker-image-arm64.yml`](/E:/code/go/new-api/.github/workflows/docker-image-arm64.yml) 的冲突，保留 fork 的镜像发布定制并吸收上游工作流增强 | depends_on: [1.2]
- [√] 2.2 解析 [`relay/compatible_handler.go`](/E:/code/go/new-api/relay/compatible_handler.go) 的冲突，保留 fork 的兼容修复并吸收上游主线改动 | depends_on: [1.2]

### 3. 验证与交付

- [√] 3.1 运行针对性的 Go 测试或最小验证，确认 relay 兼容路径没有明显回归 | depends_on: [2.1,2.2]
- [√] 3.2 检查 `git status` / `git diff --stat`，确认分支达到可推送到 `fork/sync/origin-main-with-local7-20260318` 的状态 | depends_on: [3.1]

---

## 执行日志

| 时间 | 任务 | 状态 | 备注 |
|------|------|------|------|
| 2026-04-14 05:00 | 方案包创建 | completed | 已记录最新分叉状态、冲突预检结果与实施决策 |
| 2026-04-14 05:02 | 1.1 | completed | 已抓取最新远端引用，确认当前同步分支相对 origin/main 为 behind 171 / ahead 9 |
| 2026-04-14 05:03 | 1.2 | blocked | merge 被未提交改动阻止，涉及 service/convert.go 与同组本地改动 |
| 2026-04-14 05:05 | 1.2 | resumed | 用户确认先做本地保护提交，再继续合并 origin/main |
| 2026-04-14 05:06 | 1.2 | completed | 已创建保护提交 93360dd5，merge origin/main 后进入 3 个文件的真实冲突解析 |
| 2026-04-14 05:09 | 2.1/2.2 | completed | workflow 冲突收敛为 GHCR-only 方案，relay 冲突保留 legacy fallback 与 postConsumeQuota |
| 2026-04-14 05:10 | 3.1 | completed | `go test ./relay ./controller ./service` 相关用例通过，`go test ./relay/channel/vertex` 通过，`git diff --check` 无异常 |
| 2026-04-14 05:12 | 3.2 | completed | merge 提交 af2e84a5 已推送到 fork/sync/origin-main-with-local7-20260318 |

---

## 执行备注

- 最新抓取后，当前同步分支相对 `origin/main` 为 behind 171 / ahead 9。
- 只读冲突预检显示真实内容冲突文件为 3 个：workflow 2 个、relay 1 个。
- 目标是保留 fork 定制，不重写远端历史。
- 额外发现未提交本地改动：[`relay/channel/vertex/adaptor.go`](/E:/code/go/new-api/relay/channel/vertex/adaptor.go)、[`relay/channel/vertex/adaptor_test.go`](/E:/code/go/new-api/relay/channel/vertex/adaptor_test.go)、[`service/convert.go`](/E:/code/go/new-api/service/convert.go)、[`service/convert_test.go`](/E:/code/go/new-api/service/convert_test.go)。需先决定如何纳入同步流程。
