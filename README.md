# oss
参考[《分布式对象存储：原理 架构及Go语言实现》]（https://github.com/stuarthu/go-implement-your-object-storage）实现的一个简单的对象存储原型系统，包括：分布式数据存储、分布式接口服务、基于redis实现的消息队列广播心跳、基于elasticSearch实现的元数据管理、对上传的对象基于SHA256验证哈希值、基于RS纠删码实现的数据冗余与即时恢复、断点上传与下载、基于gzip算法进行数据压缩、用户鉴权与权限管理及一个简易的客户端。
