if (redis.call('exists',KEYS[1])==0) then // 锁是否存在
    redis.call('hincrby',KEYS[1],ARGV[2],1); // 不存在 记录锁标识 次数加一
    redis.call('pexpire',KEYS[1],ARGV[1]); // 设置有效期
    return nil;
end;
if (redis.call('hexists',KEYS[1],ARGV[2])==1) then // 存在 判断锁是否是自己的
    redis.call('hincrby',KEYS[1],ARGV[2],1); // 是自己的 也要次数加一
    redis.call('pexpire',KEYS[1],ARGV[1]); // 然后设置有效期
    return nil;
end;
return redis.call('pttl',KEYS[1]); // 获取锁失败 返回锁的剩余有效期