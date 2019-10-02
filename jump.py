import hashlib
from collections import defaultdict
import mmh3
import socket
import struct


shards = ["shard{0}".format(si) for si in range(0, 1000)]
numserv = [3, 4, 5, 6, 7, 8, 9, 10]

def murmur(key):
    """Return murmur3 hash of the key as 32 bit signed int."""
    return mmh3.hash(key)


def weight(node, key):
    a = 1103515245
    b = 12345
    hash = murmur(key)
    return (a * ((a * node + b) ^ hash) + b) % (2^31)

def rendezvous_hash(key, num_buckets):
    weights = []
    for n in range(0, num_buckets):
        w = weight(murmur("node" + str(n)), key)
        weights.append((w, n))

    _, node = max(weights)
    return node

def consistent_hash(key, num_buckets):
    """
    A Fast, Minimal Memory, Consistent Hash Algorithm (Jump Consistent Hash)

    Hash accepts "a 64-bit key and the number of buckets. It outputs a number
    in the range [0, buckets]." - http://arxiv.org/ftp/arxiv/papers/1406/1406.2294.pdf

    The C++ implementation they provide is as follows:

    int32_t JumpConsistentHash(uint64_t key, int32_t num_buckets) {
        int64_t b = -1, j = 0;
        while (j < num_buckets) {
            b   = j;
            key = key * 2862933555777941757ULL + 1;
            j   = (b + 1) * (double(1LL << 31) / double((key >> 33) + 1));
        }
        return b;
    }

    """
    if not isinstance(key, (int, long)):
        if isinstance(key, unicode):
            key = key.encode('utf-8')
        key = int(hashlib.md5(key).hexdigest(), 16) & 0xffffffffffffffff
    b, j = -1, 0
    if num_buckets < 0:
        num_buckets = 1
    while j < num_buckets:
        b = int(j)
        key = ((key * 2862933555777941757) + 1) & 0xffffffffffffffff
        j = float(b + 1) * (float(1 << 31) / float((key >> 33) + 1))
    return b & 0xffffffff

hash_impl = rendezvous_hash

if False:
    print("shardN\tNloc\t" + "\t".join([str(num) for num in numserv]))
    for ix, shard in enumerate(shards):
        line = [shard, ""]
        sout = set()
        for num in numserv:
            outix = hash_impl(shard, num)
            sout.add(outix)
            line.append(str(outix))
        line[1] = str(len(sout))
        print("\t".join(line))

def print_histogram(hash_fn):
    # Print histogram
    hist = defaultdict(int)
    for ix, shard in enumerate(shards):
        sout = set()
        for num in numserv:
            outix = hash_fn(shard, num)
            sout.add(outix)
        hist[len(sout)] += 1

    lst = list(hist.items())
    lst = sorted(lst, key=lambda kv: kv[1], reverse=True)

    print(lst)

print_histogram(consistent_hash)
print_histogram(rendezvous_hash)