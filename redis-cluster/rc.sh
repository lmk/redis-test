
function cluster_flushall() {
  echo flushall
  redis-cli -c -p 7011 flushall
  redis-cli -c -p 7012 flushall
  redis-cli -c -p 7013 flushall
  redis-cli -c -p 7014 flushall
  redis-cli -c -p 7015 flushall
  redis-cli -c -p 7016 flushall
}

function cluster_reset() {
  echo cluster reset
  redis-cli -c -p 7011 cluster reset
  redis-cli -c -p 7012 cluster reset
  redis-cli -c -p 7013 cluster reset
  redis-cli -c -p 7014 cluster reset
  redis-cli -c -p 7015 cluster reset
  redis-cli -c -p 7016 cluster reset
}

if [ "$1" == "shard" ]; then 

  cluster_flushall
  cluster_reset

  echo cluster addslots
  redis-cli -c -p 7011 cluster addslots {0..5461}
  redis-cli -c -p 7013 cluster addslots {5462..10923}
  redis-cli -c -p 7015 cluster addslots {10924..16383}

  echo cluster meet
  redis-cli -c -p 7011 cluster meet 127.0.0.1 7011
  redis-cli -c -p 7011 cluster meet 127.0.0.1 7012
  redis-cli -c -p 7011 cluster meet 127.0.0.1 7013
  redis-cli -c -p 7011 cluster meet 127.0.0.1 7014
  redis-cli -c -p 7011 cluster meet 127.0.0.1 7015
  redis-cli -c -p 7011 cluster meet 127.0.0.1 7016

  sleep 1

  echo cluster replicate
  redis-cli -c -p 7012 cluster replicate `redis-cli -c -p 7011 cluster nodes | grep ":7011" | awk '{print $1}'` 
  redis-cli -c -p 7014 cluster replicate `redis-cli -c -p 7013 cluster nodes | grep ":7013" | awk '{print $1}'`
  redis-cli -c -p 7016 cluster replicate `redis-cli -c -p 7015 cluster nodes | grep ":7015" | awk '{print $1}'`

  sleep 1

fi 

redis-cli -c -p 7011 cluster nodes
echo ""
redis-cli -c -p 7011 cluster info
