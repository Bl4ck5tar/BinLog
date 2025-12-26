## 这是一个基于 Gin 的内容社区后端系统，对外通过类 RESTful API 提供服务

核心模块包括：
- 用户与权限模块（JWT 鉴权）
- 内容管理模块（文章 CRUD、评论体系）
- 搜索模块（Elasticsearch 全文搜索）
- 缓存与性能优化模块（Redis）
- 监控与日志模块（Prometheus + Grafana + Zap）

整体目标是支撑 高并发访问、快速响应和可观测性。
