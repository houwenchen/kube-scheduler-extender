# kube-scheduler-extender

## 需求
  
  实现根据 pod 的 annotations 中 nodeName 字段来调度此 pod 到某个 node

## 模块

### 1. controller 模块作用

  监听 pod，根据 pod 的 annotations 字段中的 nodeName 对应的 value，维护一块内存记录 pod: node 的对应信息。

### 2. scheduler 模块作用

  实现 scheduler extender 逻辑，此 demo 只需要实现 filter 功能即可满足需求。

## 使用

  1. 将 /deploy/scheduler-policy-config.yaml 文件放到 /etc/kubernetes/scheduler-policy-config.yaml
  2. 参照 kube-scheduler.yaml 修改 kube-scheduler static pod 的 yaml 文件，有3处修改，具体可以参照 /deploy/kube-scheduler.yaml 中的标注
  3. 运行程序：go run ./cmd/scheduler-extender/scheduler-extender.go，或者打包成镜像。
  4. 参考 /deploy/pod.yaml 部署需要控制的 pod

## 代码分析

  具体分析见详细代码，有问题可以提 issue，会尽快回复