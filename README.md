## 问答社区

**问题表 (`questions`)**:

- `id`: 主键，自增
- `title`: 问题标题
- `content`: 问题详细描述
- `created_at`: 创建时间
- `updated_at`: 更新时间，最后一条评论的产生时间

**回答表 (`answers`)**:

- `id`: 主键，自增
- `question_id`: 关联的问题ID
- `content`: 回答内容
- `created_at`: 创建时间

**评论表 (`comments`)**:

- `id`: 主键，自增
- `answer_id`: 关联的回答ID
- `content`: 评论内容
- `created_at`: 创建时间

#### **路由分析**

**用户路由组**

`/user/login`用户登录

`/user/register`用户注册

##### 问题相关路由

**POST**`/question/create`创建问题

**GET**`/question/:id`获取指定id的问题

**PUT**`/question/:id`更新指定id的问题

**DELETE**`/question/:id`删除指定id的问题

#### Redis**获取热门问题**返回值

```
{"id":10,"user_id":7,"title":"八月六日","content":"天涯沈铁","created_time":"2024-08-06T11:49:44.817+08:00","updated_time":"2024-08-06T11:49:44.817+08:00"}
{"id":13,"user_id":1,"title":"八月六日","content":"dech53用户创建了第二个问题","created_time":"2024-08-06T12:30:49.271+08:00","updated_time":"2024-08-06T12:30:49.271+08:00"}
{"id":12,"user_id":1,"title":"八月六日","content":"dech53用户创建了一个问题","created_time":"2024-08-06T12:15:42.123+08:00","updated_time":"2024-08-06T12:15:42.123+08:00"}
{"id":1,"user_id":1,"title":"测试问题1","content":"测试问题1内容","created_time":"2024-08-02T22:33:45.328+08:00","updated_time":"2024-08-06T16:04:13.737+08:00"}
{"id":3,"user_id":1,"title":"测试问题1","content":"测试问题1内容","created_time":"2024-08-02T22:35:40.673+08:00","updated_time":"2024-08-05T21:37:38.879+08:00"}
{"id":11,"user_id":7,"title":"八月六日","content":"天涯沈铁","created_time":"2024-08-06T11:52:27.442+08:00","updated_time":"2024-08-06T11:52:27.442+08:00"}
{"id":9,"user_id":7,"title":"八月五-日","content":"wx修改了问题描述","created_time":"2024-08-05T10:14:47.808+08:00","updated_time":"2024-08-05T21:57:44.113+08:00"}
{"id":6,"user_id":1,"title":"wx问题1","content":"wx问题1内容","created_time":"2024-08-02T22:43:37.511+08:00","updated_time":"2024-08-06T16:04:22.219+08:00"}
```

