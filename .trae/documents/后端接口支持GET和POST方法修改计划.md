## 修改计划

### 1. 路由设置修改
- 修改 `cmd/main.go` 中的路由定义，为所有接口添加同时支持GET和POST方法的路由
- 对于每个现有路由，添加对应方法的路由映射

### 2. 处理函数修改
- 修改所有处理函数，使其能够同时处理GET和POST请求
- 实现请求数据的自适应解析：
  - 对于GET请求，从URL查询参数获取数据
  - 对于POST请求，从JSON请求体获取数据
- 确保所有响应都是JSON格式

### 3. 具体修改内容

#### 3.1 认证路由（auth.go）
- 修改 `Register` 和 `Login` 函数，支持从GET请求的查询参数获取数据

#### 3.2 租户路由（tenant.go）
- 修改 `CreateTenant` 和 `UpdateTenant` 函数，支持从GET请求的查询参数获取数据
- 为所有租户路由添加GET/POST方法支持

#### 3.3 用户路由（user.go）
- 修改 `UpdateUser` 函数，支持从GET请求的查询参数获取数据
- 为所有用户路由添加GET/POST方法支持

### 4. 实现方法
- 使用Gin的 `Any` 方法或同时注册GET和POST路由
- 在处理函数中根据请求方法选择数据解析方式
- 确保JSON绑定和查询参数绑定都能正常工作

### 5. 验证方法
- 测试所有接口的GET和POST方法是否都能正常工作
- 验证POST请求是否正确接收和返回JSON数据
- 验证GET请求是否正确处理查询参数

### 6. 修改文件列表
- `backend/cmd/main.go` - 修改路由设置
- `backend/internal/handlers/auth.go` - 修改认证处理函数
- `backend/internal/handlers/tenant.go` - 修改租户处理函数
- `backend/internal/handlers/user.go` - 修改用户处理函数