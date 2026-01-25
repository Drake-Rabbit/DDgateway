# Gateway Service

一个基于 Gin + GORM 的 Go 微服务网关，具有登录和多租户功能,前端使用 Nuxt.js + vue3.js 实现。
#### 项目结构清晰

## 项目结构

```
gateway-service/
├── backend/              # Go 后端
│   ├── cmd/             # 入口文件
│   ├── internal/        # 内部包
│   │   ├── config/     # 配置
│   │   ├── database/   # 数据库
│   │   ├── define/     # define
│   │   ├── dto/        # dto数据模型
│   │   ├── router/     # router路由 
│   │   ├── middleware/ # 中间件
│   │   └── models/     # 数据模型
│   └── pkg/            # 公共包
│       ├── jwt/        # JWT 工具
│       └── response/   # 响应工具
└── frontend/           # Nuxt.js 前端
    ├── api/           # API 调用
    ├── composables/   # 组合式函数
    ├── components/    # 组件
    ├── layouts/       # 布局
    ├── pages/         # 页面
    └── stores/        # 状态管理
```
## 技术栈

### 后端
- **Gin**: Web 框架
- **GORM**: ORM 框架
- **MySQL**: 数据库
- **JWT**: Token认证

### 前端
- **Nuxt 3**: Vue 3 框架
- **Tailwind CSS**: CSS 框架
- **Pinia**: 状态管理

## API 端点
注册 ✅ 登录 ✅
#### 租户管理 ⚠️
#### 用户管理 ⚠️
#### 服务管理 ✅
    获取服务列表 | ✅ 
    获取服务详情 | ✅ 
    获取服务统计 | ✅ 
    删除服务     | ✅ 
    创建HTTP服务 | ✅ 
    更新HTTP服务 | ✅ 
    创建TCP服务 | ✅ 
    更新TCP服务 | ✅ 
    创建gRPC服务 ❌ 
    更新gRPC服务 ❌
### JWT认证
### 多租户隔离
- 用户表通过 `tenant_id` 字段与租户关联
- 所有用户查询都自动过滤当前租户的数据
- 不同租户之间的用户数据完全隔离

## 前端功能状态

### 📱 页面功能概览

| 登录页 | 用户登录 | ✅ 
| 注册页 | 用户注册 | ✅ 
| Dashboard | 系统概览 | ⚠️ 静态页面 
| 服务管理 | HTTP/TCP服务 | ⚠️ 部分功能 
| 租户管理 | 租户CRUD | ❌ |
| 用户管理 | 用户CRUD | ❌ |


### 📊 网关系统-整体模块
### 🛠️ 模块功能详情
    **认证系统**  
    **租户管理** 
    **用户管理**
    
    **HTTP服务管理** 
    **TCP服务管理** 
    **gRPC服务管理**
    **服务概览** 
    **服务统计** 
    **负载均衡**
    **黑白名单**
    **网络扩展**
    




    


### 🔮 其他规划（AI增强功能）
#### 🤖 AI 智能分析功能
1. **服务健康度分析**
2. **智能负载均衡推荐**
3. **运维智能助手**
4. **洞察报告**




