服务发现

// Grpc Client
已经把Resolver集成到gclient

// Grpc Server Register
```
func main() {
  .
  .
  .
	if err := consul.NewRegister(Conf.Consul).
		Registe(Conf.Server.Name, LocIP, Conf.Server.Port); err != nil {
		panic(err)
	}
    consul.RegisterGrpcHealth(svr)
 .
 .
 .
}
```
