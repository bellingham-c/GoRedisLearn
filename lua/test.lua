--获取锁中的线程提示 get key
local id =redis.call('get',KEYS[1])
--比较线程标示与锁的标示是否一致
if (id==ARGV[1]) then
    --释放锁 del key
    return redis.call('del',KEYS[1])
end
return 0