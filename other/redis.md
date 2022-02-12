Redis 是一个开源（BSD许可）的，内存中的数据结构存储系统，它可以用作数据库、缓存和消息中间件。 它支持多种类型的数据结构，如 [字符串（strings）](http://redis.cn/topics/data-types-intro.html#strings)， [散列（hashes）](http://redis.cn/topics/data-types-intro.html#hashes)， [列表（lists）](http://redis.cn/topics/data-types-intro.html#lists)， [集合（sets）](http://redis.cn/topics/data-types-intro.html#sets)， [有序集合（sorted sets）](http://redis.cn/topics/data-types-intro.html#sorted-sets) 与范围查询， [bitmaps](http://redis.cn/topics/data-types-intro.html#bitmaps)， [hyperloglogs](http://redis.cn/topics/data-types-intro.html#hyperloglogs) 和 [地理空间（geospatial）](http://redis.cn/commands/geoadd.html) 索引半径查询。 Redis 内置了 [复制（replication）](http://redis.cn/topics/replication.html)，[LUA脚本（Lua scripting）](http://redis.cn/commands/eval.html)， [LRU驱动事件（LRU eviction）](http://redis.cn/topics/lru-cache.html)，[事务（transactions）](http://redis.cn/topics/transactions.html) 和不同级别的 [磁盘持久化（persistence）](http://redis.cn/topics/persistence.html)， 并通过 [Redis哨兵（Sentinel）](http://redis.cn/topics/sentinel.html)和自动 [分区（Cluster）](http://redis.cn/topics/cluster-tutorial.html)提供高可用性（high availability）。

# 0. 安装

默认会在 `6379` 端口启动数据库。

```
$ docker run --name some-redis -d -p 6379:6379 redis
```

另外还可以启用 [持久存储](https://redis.io/topics/persistence)。

```
$ docker run --name some-redis -d -p 6379:6379 redis redis-server --appendonly yes
```

默认数据存储位置在 `VOLUME/data`。可以使用 `--volumes-from some-volume-container` 或 `-v /docker/host/dir:/data` 将数据存放到本地。

使用其他应用连接到容器，可以用

```
$ docker run --name some-app --link some-redis:redis -d application-that-uses-redis
```

或者通过 `redis-cli`

```
$ docker run -it --rm \
    --link some-redis:redis \
    redis \
    sh -c 'exec redis-cli -h "$REDIS_PORT_6379_TCP_ADDR" -p "$REDIS_PORT_6379_TCP_PORT"'
```



# 1. redis 命令

## 1.1 key

| 命令                                  | 描述                                      | 备注                                                         |
| ------------------------------------- | ----------------------------------------- | ------------------------------------------------------------ |
| redis-cli                             | 启动命令，连接本地的redis服务器           | redis-cli 后面加上 --raw可以避免中文乱码                     |
| redis-cli -h host -p port -a password | 链接远程服务器                            |                                                              |
| select 1                              | 选择数据库，默认开启16个，编号从0-15      |                                                              |
| flushdb 1 或者 flushall               | 清空数据库                                |                                                              |
| keys *                                | 查询符合正则表达的key                     |                                                              |
| del key                               | key存在则删除key                          |                                                              |
| dump key                              | 序列化给定key,并返回被序列化的值          | 带有64位的校验用与检测错误                                   |
| restore key                           | 反序列化                                  |                                                              |
| exists key                            | 检查给定key是否存在                       |                                                              |
| expire key n                          | 为给定key 设则过期时间，以秒计            | 设置成功返回1,设置失败返回0，pexpire 单位毫秒，在后加at为unix时间戳 |
| persist key                           | 移除key的过期时间                         |                                                              |
| ttl key                               | 以秒为单位，返回给定key的剩余时间         | pttl key ，返回值为毫秒                                      |
| move key db                           | 将当前数据库的key移动到给定的数据库db当中 |                                                              |
| randomkey                             | 从当前数据库中随机返回一个key             |                                                              |
| rename key newkey                     | 修改当前key的名称                         | 当newkey已经存在时，覆盖                                     |
| renamenx key newkey                   | 仅当newkey不存在时，将key改名为newkey     |                                                              |
| type key                              | 返回key所存储的值的类型                   |                                                              |
|                                       |                                           |                                                              |
|                                       |                                           |                                                              |

## string

| 命令                       | 描述                                                         | 备注                                   |
| -------------------------- | ------------------------------------------------------------ | -------------------------------------- |
| set key value              | 设置指定key的值                                              | setnx只有在key不存在时设置key的值      |
| mset key value key2 value  | 同时设置一个或多个key-value对                                | msetnx只有在所有key不存在时设置key的值 |
| get key                    | 获取指定 key的值                                             |                                        |
| getrange key start end     | 返回start-end的字符串                                        |                                        |
| setrange key  offset value | 将指定的key的value从offset开始替换                           |                                        |
| getset key value           | 将给定key的值设为value，并返回key的旧值                      |                                        |
| setex key n value          | 将值 value 关联到 key ，并将 key 的过期时间设为 n (以秒为单位) | 即设置初始时间                         |
| strlen key                 | 返回key所存储的字符串值的长度                                |                                        |
| incr key                   | 将key所存储的数字值增一                                      |                                        |
| incrby key n               | 将key所存储的值加n                                           |                                        |
| decr key                   | 将key所存储的数字值减一                                      |                                        |
| decrby key n               | 将key所存储的值减n                                           |                                        |
| append key value           | 如果key已经存在并且是一个字符串，append命令将指定的value追加到该key原来value的末尾，如果key不存在则新建 |                                        |
|                            |                                                              |                                        |
|                            |                                                              |                                        |
|                            |                                                              |                                        |
|                            |                                                              |                                        |
|                            |                                                              |                                        |
|                            |                                                              |                                        |

## list

Redis列表是简单的字符串列表，按照插入顺序排序。你可以添加一个元素到列表的头部（左边）或者尾部（右边）

一个列表最多可以包含  $ 2^{32}-1 $个元素 (4294967295, 每个列表超过40亿个元素)

| 命令                                  | 描述                                                         | 备注 |
| ------------------------------------- | ------------------------------------------------------------ | ---- |
| blpop key1 [key2] timeout             | 移出并获取列表的第一个元素，如果列表没有元素阻塞列表直到等待超时或发现可弹出元素为止 |      |
| brpop key1 [key2] timeout             | 移出并获取列表的最后的一个元素，如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止 |      |
| brpoplpush source destination timeout | 从列表中弹出一个值，将弹出的元素插入到另一个列表中返回它;如果元素没有元素会阻塞列表直到等待超时或发现可弹出元素为止 |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |
|                                       |                                                              |      |



