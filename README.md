# Wireshark examples

Project to play with packet capture data between docker containers.

## Overview

// TODO: write a short overview

## Steps

1. Start the es01 container:

        $ cd containers
        $ docker-compose up es01

2. Verify all is well, and that we see a single node:

        # In a separate terminal
        $ curl 'localhost:9200/_cat/health?v'
        epoch      timestamp cluster           status node.total node.data shards pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
        1598180192 10:56:32  es-docker-cluster green           1         1      0   0    0    0        0             0                  -                100.0%
        $ curl 'localhost:9200/_cat/nodes?v'
        ip         heap.percent ram.percent cpu load_1m load_5m load_15m node.role master name
        172.24.0.2           46          65   1    0.06    0.32     0.31 dilmrt    *      es01

3. Packet capture the `:9300` (node-to-node) traffic:

        $ docker run --rm -it --net=container:es01 -v $(pwd)/netshoot:/app nicolaka/netshoot tcpdump 'port 9300' -w /app/es02-joining-cluster.pcap

4. Let `es02` join the cluster:

        # In a separate terminal
        $ docker-compose up es02

5. Confirm we see two nodes:
        
        # In a separate terminal
        $ curl 'localhost:9200/_cat/nodes?v'
        ip         heap.percent ram.percent cpu load_1m load_5m load_15m node.role master name
        172.24.0.2           30          96  46    0.96    0.44     0.33 dilmrt    *      es01
        172.24.0.3           25          96  56    0.96    0.44     0.33 dilmrt    -      es02
        $ curl 'localhost:9200/_cat/health?v'
        epoch      timestamp cluster           status node.total node.data shards pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
        1598182416 11:33:36  es-docker-cluster green           2         2      0   0    0    0        0             0                  -                100.0%

6. At this point you can stop the packet capture -- hit `Ctrl-C` in the tcpdump terminal.

7. Open the pcap file in wireshark.


## Attribution

1. Docker compose file is a modified multi-node setup from https://www.elastic.co/guide/en/elasticsearch/reference/7.9/docker.html
2. `nicolaka/netshoot` is a docker network troubleshooting container maintained at https://github.com/nicolaka/netshoot