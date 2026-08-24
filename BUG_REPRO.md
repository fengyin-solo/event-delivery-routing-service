# Bug

待投递事件不能形成 delivered 终态，重新查询仍缺少投递时间。

# Trigger

创建 pending 事件，标记投递成功后读取详情，并尝试再次投递。

# Error

`deliver event: status: 当前事件状态不允许标记投递成功`
