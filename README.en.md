# ginana

#### Description
- 基于Gin + Gorm + Casbin + Paladin + Logrus + Wire 实现的API开发脚手架，目的是提供一套轻量级的开发框架，致力于结构清晰、方便、快速的完成业务需求的开发。

#### Software Architecture
整体参考bilibili开源的kratos架构，特此感谢其给予的宝贵灵感。
+ Gin：一度在beego、iris等特性突出的框架中徘徊，最终胜出的唯一原因是gin的可扩展性高出许多。
+ Gorm：没什么好描述的，个人常用，没遇到什么理由抛弃它。
+ Casbin：权限控制模块，核心概念（存储RBAC方案中用户与角色之间的映射关系）清晰易用。
+ Paladin：kratos框架里的一个舒适度很高的配置中心模块，支持远程配置中心、热加载。
+ Logrus：github上较活跃的日志框架，体感舒适度也很好，高扩展性，其中Hook机制是亮点。
+ Wire：Google的依赖注入工具。
+ Swag：自动生成Swagger2.0文档的工具，用起来很方便，但复杂一点的功能需求没有写完

#### Instructions

1.  xxxx
2.  xxxx
3.  xxxx

#### Contribution

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request
