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

**路由分析**

