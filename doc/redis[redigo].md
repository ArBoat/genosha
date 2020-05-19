### connect to redis server
c, err := redis.Dial("tcp", "127.0.0.1:6379",redis.DialDatabase(9))

### send a command to server [combine Send, Flush, Receive]

c.Do("EXISTS", key)

### pub/sub
listenPubSubChannels
publish
receive


### pipeline vs transactions / pipelined transactions

pipeline: effective(10 times fast) , not acid

transactions: acid， slow

### RDB AOF
config