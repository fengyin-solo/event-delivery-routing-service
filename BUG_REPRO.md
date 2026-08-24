# Bug

事件累计第三次投递失败后仍停在 failed，没有生成包含最后失败说明的死信。

# Trigger

先标记事件失败，再连续执行两次失败重试并查询死信列表。

# Error

`third failure should be dead at attempt 3; status=failed attempts=3`
