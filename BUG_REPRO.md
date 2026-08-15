# Bug Reproduction

给子路由启用 `CORSMethodMiddleware` 后，嵌套路径的正常请求会把兄弟路由的方法也写进 `Access-Control-Allow-Methods`。同一路由直接挂在根路由时响应头正常。

复现步骤：

1. 创建 `/test` 子路由。
2. 在子路由上注册 `/hello/{name}` 的 `GET` 和 `OPTIONS` 方法。
3. 给子路由启用 `CORSMethodMiddleware`。
4. 发送 `GET /test/hello/alice` 请求。

期望响应头为 `Access-Control-Allow-Methods: GET,OPTIONS`；当前响应头还包含兄弟 `/hello` 路由的 `POST` 方法。根路由的 CORS 方法响应、路径匹配和其他中间件行为应保持不变。
