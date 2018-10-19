# go-server
服务端golang服务 基于gin和grpc
gateway-api: 网关服务，基于gin框架，通过gclient调用下一层微服务
gclient: 所有微服务的接口集合，通过grpc自带负载均衡，使用consul作为服务发现
common: 公共组建
proto: 需要公用的协议 
