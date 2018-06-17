if not exist data mkdir data
start dgraph zero --wal data/zw
start dgraph server --lru_mb 2048 --wal data/w --postings data/p --debugmode --bindall
start dgraph-ratel
