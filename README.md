# redis 分布式锁测试实验

这个仓库展示了针对 redis 分布式锁的测试实验细节

## 宕机重启实验

操作 docker 容器的行为模拟节点宕机的危害，相关代码在[redis.go](https://github.com/lawrshen/redlock-test/blob/master/redis/redis.go)中

## 时钟偏移实验

借助混沌工程工具，生成模拟时钟偏移的命令行工具：

```BASH
git clone https://github.com/chaos-mesh/chaos-mesh
cd chaos-mesh; make watchmaker
```

使用 `watchmaker` 对 redis 实例模拟时钟偏移

## 测试脚本

在 [python脚本](https://gist.github.com/JJGO/0d73540ef7cc2f066cb535156b7cbdab) 基础上增加了每次实验的 setup 和 teardown 工作，使大量重复测试自动化的在一个全新的 redis 集中工作。测试前需要保证 server 服务器在运行状态，脚本指定 client\_test 中的测试代码。

```bash
go run server/server.go &
cd client && python dstest.py TestNormal -n 5
```
