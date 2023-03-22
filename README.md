# x-tiktok
## Description

基于 Gin + Gorm + JWT + MySQL + Redis + RabbitMQ + OSS + Docker 实现的第五届字节跳动后端青训营极简版抖音，完成包含互动方向和社交方向的所有基础功能以及扩展功能，并通过缓存和消息队列等中间件来保证接口性能，提高用户体验。

## Award

<p align="center">
  <img alt="超级码力奖" src="https://user-images.githubusercontent.com/39022409/226916716-ba2c06a0-7e86-4395-8f7f-03bcf5736950.png" width="40%">
</p>

## Video Demo

[第五届字节跳动后端青训营—X-Engineer 小组结业作品演示](https://www.bilibili.com/video/BV1nj411F7L1/?share_source=copy_web&vd_source=fcc347bbf1c653735e29efa79e12dc86)

## Technical Selection

- Web Framework：[Gin](https://gin-gonic.com/)
- ORM Framework：[GORM](https://gorm.io/)
- Security Framework：[JWT](https://jwt.io/)
- DataBase：[MySQL](https://www.mysql.com/cn/)
- Cache：[Redis](https://redis.io/)
- Message Queue：[RabbitMQ](https://www.rabbitmq.com/)
- File Storage Service：[OSS](https://help.aliyun.com/product/31815.html)
- Deployment：[Docker](https://www.docker.com/)
- Project Management：[GitHub Issue/PR/Milestones/Project Board](https://github.com/X-Engineer/x-tiktok)

## System Structure

整个 x-tiktok 项目的整体架构采用经典的 MVC 三层架构，客户端层访问路由层，不同的接口请求路由到不同的业务层，期间经过鉴权层对用户信息进行鉴权认证，通过通过验证的用户请求方可进入到业务层，业务层主要包含了视频流，视频发布，点赞，关注，评论，消息和用户等业务。业务层中不同的业务通过调用服务层中的各项服务来实现接口定义的 Response，服务层主要和数据层以及 Redis，RabbitMQ 等中间件进行交互，获取到数据层的原始数据模型，再做相应的封装，返回给上层服务。数据层主要和数据存储层，其中包括像 MySQL 数据库，OSS 对象存储进行交互。最底层则是基础设施层，包含项目部署所要用到的 Docker，Go Runtime，云服务器等。

![](https://raw.githubusercontent.com/zhicheng-ning/Pic-Go/main/md/whiteboard_exported_image.png)

## Devlopment Document

详细的开发答辩文档：https://bytedancecampus1.feishu.cn/docx/AURFd7pMCourqpxs8YDcriV6n9d

## Project Division

| **团队成员**                                                 | **主要贡献**                                                 |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [@zhicheng-ning](https://github.com/zhicheng-ning)           | 队长，负责 GitHub 项目管理，定期会议同步，代码 review，基础设施搭建，视频模块、消息模块 |
| [@coder-zc](https://github.com/coder-zc)                     | 负责负责用户模块、评论模块、主要包括注册，登录，评论，删评，以及缓存和消息队列的优化 |
| [@MR.XU](https://github.com/Xuuuuuuuuuuuu) [@dongzhengyu](https://github.com/dongzhengyu816) | 负责关注模块、主要包括关注，取关，关注列表，粉丝列表，好友列表，以及缓存和消息队列的优化 |
| [@Jasmine](https://github.com/ruirui-wang-study)             | 负责点赞模块、主要包括点赞，取消点赞，获取点赞列表，以及缓存和消息队列的优化 |

## Contributing

```shell
1. fork x-tiktok to your own namespace
2. git clone https://github.com/<your-username>/x-tiktok.git
3. git remote add upstream https://github.com/X-Engineer/x-tiktok.git
4. git fetch --all
5. git checkout -b <your-local-branch-name>
6. do some changes
7. git add .
8. git commit -m "your commit message"
9. git rebase origin/main (optional)
10. git push origin <your-branch-name>
11. create a pull request
12. wait for code review
13. merge your pull request
14. delete <your-branch-name>
15. git checkout main
16. git pull upstream main
17. Go back to step 5 to develop new features
```

## Deploy

```shell
1. git clone https://github.com/X-Engineer/x-tiktok
2. cd x-tiktok
3. config your database in dao/db.go
4. config your redis in middleware/redis/redis.go
5. config your rabbitmq in middleware/rabbitmq/rabbitmq.go
6. config your oss and jwt-token in config/config.go
7. sh run.sh
```

## Client APK
[客户端 APK 下载地址](https://xlab-open-source.oss-cn-beijing.aliyuncs.com/zhicheng-ning/bytedance-go/x-tiktok/dousheng/app-release_3.apk)

## Thanks

- [第五届字节跳动后端青训营](https://juejin.cn/post/7171281874357059592)