# go-server

### gateway: 网关服务，基于gin，通过gclient调用service相关服务
### service: 业务服务，基于grpc
### gclient: 所有grpc微服务的接口集合，通过grpc自带的负载均衡，使用consul作为服务发现
### common:  公共组件
### proto:   需要公用的协议 
