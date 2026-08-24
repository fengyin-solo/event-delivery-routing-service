# Bug

主题重命名冲突被报告为成功，原主题被覆盖并产生同名列表项。

# Trigger

创建 alerts 和 metrics 两个主题，再把 alerts 更新为 metrics，随后查询原主题与列表。

# Error

`rename collision should return conflict, got <nil>`
