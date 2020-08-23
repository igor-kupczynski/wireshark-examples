# Containers

## Overview

Shows how to use `tcpdump` with docker -- to capture network traffic of containerized applications.

### Problem statement

How to `tcpdump` traffic to or from a docker container.

We use dockerized two-node Elasticsearch cluster. Let's say we want to capture traffic on _port 9300_ of the first node. Port 9300 is responsible for the node-to-node communication. 

### Solution

Docker containers can use different network drivers. The most common one is `bridge` (also the default). Each container gets its own IP address and docker creates a bridge network between all containers. In order to have an isolated environment we setup a custom bridge network between the two nodes. See [_Use bridge networks_](https://docs.docker.com/network/bridge/) for details on the bridge networking in docker.

Common alternative is to use `host` network, this shares the network stack with the host, so you can use `tcpdump` as with any other local process. We ignore this case here.

Docker also allows a container B to share a network stack of another container A. This means that both containers A and B get the same IP addresses. They also share network interface, port ranges, etc. See `--net=container:<name-or-id>` command line option.

We will use a _toolbox_ container and make is share the network stack of the Elasticsearch node 1. The _toolbox_ container has `tcpdump` and we'll use it to perform the packet capture.

### Elasticsearch  

We use official [Elasticsearch docker images](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html).

The docker compose file is slightly modified:
- 2 nodes, instead of 3
- `es01` can form a cluster on its own (doesn't need to wait on `es02`)
- data directories mounted locally
- we bind `es01` to localhost

**Note** this is not a production ready setup:
- no transport security
- no access control

(Both are available for free, with the basic licence, we just don't bother for this example).

Refer to the [official docs](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-prod-prerequisites) on how to productionize Elasticsearch in docker setup. 

### Nicolaka/netshoot

Instead of building the _toolbox_ container on our own, we will use [`nicolaka/netshoop`](https://github.com/nicolaka/netshoot) -- a prebuilt container with various network tools.


## Steps to capture network traffic withing docker containers

1. Start the es01 container:

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

7. Open the pcap file in wireshark. (Also attached here under `es02-joining-cluster.pcap`).

Packet capture run on my laptop is available her.


## Attribution

1. Docker compose file is a modified multi-node setup from https://www.elastic.co/guide/en/elasticsearch/reference/7.9/docker.html
2. `nicolaka/netshoot` is a docker network troubleshooting container maintained at https://github.com/nicolaka/netshoot