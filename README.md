# kube-scheduler-extender

# 需求
  
  实现根据 pod 的 annotations 中 nodeName 字段来调度此 pod 到某个 node

# 模块

1. controller 模块作用

  监听 pod 与 node，根据 pod 的 annotations 字段中的 nodeName 对应的 value 值维护到 storage 中

2. scheduler 模块作用

  实现 scheduler extender 逻辑