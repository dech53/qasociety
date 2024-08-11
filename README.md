## 问答社区

#### 数据库

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

**回复数表(`question_answer_count`)**:

- `questionID`:问题ID
- `answercount`:问题回复数

**Response返回类型**

```
{
    Code:    code,
    Message: "",
    Result:  data,
}
```

#### **路由分析**

**用户路由组**

`/user/login`用户登录

`/user/register`用户注册

`/user/request_password_reset`重置密码请求

`/user/reset_password`执行重置密码

##### 问题路由组

**POST**`/question/create`创建问题

**GET**`/question/topic10`

**GET**`/question/:id`获取指定id的问题

**PUT**`/question/:id`更新指定id的问题

**DELETE**`/question/:id`删除指定id的问题

**GET**`/question/`正常获取问题列表

##### 回复路由组

**POST**`/question/:id/answer/create`创建回复

**GET**`/question/:id/answer/search`分页搜索回复

**DELETE**`question/:id/answer/:answer_id`删除回复

##### 评论路由组

**POST**`question/:id/answer/:answer_id/comment/create`创建评论

**GET**`question/:id/answer/:answer_id/comment`分页获取评论

**DELETE**`question/:id/answer/:answer_id/:comment_id`删除评论

#### Redis**获取热门问题**返回值

```
{
            "id": 15,
            "user_id": 5,
            "title": "Navicat Data Modeler enables you to build high-quality conceptual, logical and physical             ",
            "content": "If your Internet Service Provider (ISP) does not provide direct access to its server, Secure Tunneling Protocol (SSH) / HTTP is another solution.",
            "created_time": "2010-11-09T08:10:17+08:00",
            "updated_time": "2010-11-09T08:10:17+08:00"
        },
        {
            "id": 20,
            "user_id": 10,
            "title": "To open a query using an external editor, control-click it and select Open with External            ",
            "content": "If you wait, all that happens is you get older. Navicat authorizes you to make connection to remote servers running on different platforms (i.e. Windows, macOS, Linux and UNIX), and                   ",
            "created_time": "2011-10-22T14:28:45+08:00",
            "updated_time": "2011-10-22T14:28:45+08:00"
        },
        {
            "id": 19,
            "user_id": 9,
            "title": "SSH serves to prevent such vulnerabilities and allows you to access a remote server's               ",
            "content": "To clear or reload various internal caches, flush tables, or acquire locks, control-click your connection in the Navigation pane and select Flush and choose the flush option. You must                 ",
            "created_time": "2015-09-26T03:50:51+08:00",
            "updated_time": "2015-09-26T03:50:51+08:00"
        },
        {
            "id": 14,
            "user_id": 4,
            "title": "Difficult circumstances serve as a textbook of life for people.",
            "content": "With its well-designed Graphical User Interface(GUI), Navicat lets you quickly and easily create, organize, access and share information in a secure and easy way.",
            "created_time": "2011-09-29T12:27:27+08:00",
            "updated_time": "2011-09-29T12:27:27+08:00"
        },
        {
            "id": 4,
            "user_id": 4,
            "title": "Instead of wondering when your next vacation is, maybe you should set up a life you                 ",
            "content": "Instead of wondering when your next vacation is, maybe you should set up a life you don’t need to escape from. In a Telnet session, all communications, including username and password,              ",
            "created_time": "2011-02-08T18:20:56+08:00",
            "updated_time": "2011-02-08T18:20:56+08:00"
        },
        {
            "id": 6,
            "user_id": 6,
            "title": "It can also manage cloud databases such as Amazon Redshift, Amazon RDS, Alibaba Cloud.              ",
            "content": "After logged in the Navicat Cloud feature, the Navigation pane will be divided into Navicat Cloud and My Connections sections.",
            "created_time": "2014-10-29T07:24:41+08:00",
            "updated_time": "2014-10-29T07:24:41+08:00"
        },
        {
            "id": 9,
            "user_id": 9,
            "title": "Creativity is intelligence having fun. In a Telnet session, all communications, including           ",
            "content": "It wasn’t raining when Noah built the ark.",
            "created_time": "2008-08-04T06:51:07+08:00",
            "updated_time": "2008-08-04T06:51:07+08:00"
        },
        {
            "id": 16,
            "user_id": 6,
            "title": "If opportunity doesn’t knock, build a door. Optimism is the one quality more associated           ",
            "content": "Anyone who has never made a mistake has never tried anything new. To start working with your server in Navicat, you should first establish a connection or several connections using                    ",
            "created_time": "2011-07-25T15:10:31+08:00",
            "updated_time": "2011-07-25T15:10:31+08:00"
        },
        {
            "id": 10,
            "user_id": 10,
            "title": "Import Wizard allows you to import data to tables/collections from CSV, TXT, XML, DBF and more.",
            "content": "Navicat Cloud could not connect and access your databases. By which it means, it could only store your connection settings, queries, model files, and virtual group; your database passwords            ",
            "created_time": "2022-08-17T05:22:40+08:00",
            "updated_time": "2022-08-17T05:22:40+08:00"
        },
        {
            "id": 12,
            "user_id": 2,
            "title": "Navicat Monitor is a safe, simple and agentless remote server monitoring tool that                  ",
            "content": "It collects process metrics such as CPU load, RAM usage, and a variety of other resources over SSH/SNMP.",
            "created_time": "2014-09-05T14:46:51+08:00",
            "updated_time": "2014-09-05T14:46:51+08:00"
        }
```

### 基本环境

- **语言**: Go
- **Web 框架**: Gin
- **数据库**: MySQL（用于存放持久化数据）
- **缓存**: Redis（用于存放验证码和 Top10 热门问题数据）
- **邮件发送**: Gomail（用于发送邮件验证码，并通过 Redis 验证验证码的正确性）

### 功能概述

1. **用户注册登录**:
   - 使用 MySQL 存储用户信息。
   - 验证用户登录时，验证账号和密码是否正确。
   - 通过Redis查找登录token，通过username+设备信息确认唯一登录
   - 添加正则验证密码格式是否正确
   
2. **验证码功能**:
   - 使用 Gomail 发送邮件验证码。
   - 验证码存储在 Redis 中，并通过 Redis 验证验证码的正确性。

3. **热门问题排行**:
   - 使用 Redis 存储和获取 Top10 热门问题数据。

4. **问题、回答、评论管理**:
   - 使用 MySQL 存储问题、回答、评论的持久化数据。
   - 支持问题、回答、评论的创建、查询和删除功能。
