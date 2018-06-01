mkdir -p data
start dgraph zero --wal data/zw
start dgraph server --memory_mb 2048 --wal data/w --postings data/p --debugmode --bindall
start dgraph-ratel
