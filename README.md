# cupcake_cache
A lightweight Golang distributed cache.

## Features
* LRU cache eviction algorithm
* Default hashing with consistent hashing
* Offer group caching
* Offer communication between nodes based on HTTP or GRPC protocol
* Single flight mechanism to prevent Cache Breakdown, Cache Avalanche and Cache Penetration.

### To Do
* An indepandent service to monitor and manage nodes.
* Allow dynamically adding new nodes after the cluster has been built.
