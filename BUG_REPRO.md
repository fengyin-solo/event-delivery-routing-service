# Bug

失败事件无法重试为成功，详情可能保留旧错误且缺少投递时间。

# Trigger

创建事件并标记失败，随后对同一事件执行成功重试，再读取持久化详情。

# Error

`retry success returned error: status: 失败事件不能转为投递成功`
