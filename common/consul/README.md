服务发现

// Grpc Client
已经把Resolver集成到gclient

// Grpc Server Register
```
consul.RegisterGrpcHealth(svr)

crg := consul.NewConsulRegister(Conf.Consul.Host, Conf.Server.Name, Conf.Server.Port)
if err := crg.Register(locip); err != nil {
    panic(err)                 
}
```
