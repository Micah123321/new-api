# 模块: ci-ghcr-workflow

## 职责

- 维护 `.github/workflows/docker-image-ghcr.yml` 的 GHCR 自动构建与推送策略
- 定义 `main` 与 `sync/origin-main-with-local7-20260318` 的 workflow 触发边界
- 约束镜像标签的隔离规则，避免同步分支覆盖主线 `main/latest`

## 行为规范

- `main` 分支 push 时，继续发布 `main`、`latest` 与 `sha-*` 标签
- `sync/origin-main-with-local7-20260318` 分支 push 时，自动触发同一 workflow，但发布独立标签：
  - `origin-main-with-local7-20260318`
  - `latest-origin-main-with-local7-20260318`
  - `sha-*`
- `workflow_dispatch` 也必须复用同一套标签解析逻辑，保证手动触发与自动触发的行为一致
- 当触发分支不是已知分支时，标签解析逻辑应至少将 `/` 规范化为 `-`，避免生成非法 Docker tag

## 依赖关系

- GitHub Actions 上下文变量：`GITHUB_REF_NAME`、`GITHUB_SHA`、`GITHUB_ENV`
- Docker 官方 actions：`docker/setup-buildx-action`、`docker/login-action`、`docker/build-push-action`
- GHCR 仓库路径规范化：`GHCR_REPOSITORY=${GITHUB_REPOSITORY,,}`
