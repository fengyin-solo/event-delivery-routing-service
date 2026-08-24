# Bug

事件投递成功后从 delivered 筛选和概览统计中消失。

# Trigger

创建两个事件并读取概览，将其中一个标记成功后按 delivered 查询并再次读取概览。

# Error

`delivered filter lost event: total=0 items=[]`
