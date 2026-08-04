package archytas

import "github.com/dgraph-io/badger/v4"

// Stats is a snapshot of database usage. KeyCount is computed by scanning the
// live keys, so callers should avoid collecting it on a hot path for large DBs.
type Stats struct {
	KeyCount          uint64
	LSMSizeBytes      int64
	ValueLogSizeBytes int64
	TotalSizeBytes    int64
	BlockCache        CacheStats
	IndexCache        CacheStats
}

type CacheStats struct {
	Hits        uint64
	Misses      uint64
	KeysAdded   uint64
	KeysUpdated uint64
	KeysEvicted uint64
	CostAdded   uint64
	CostEvicted uint64
}

func (b *Badger) Stats() (Stats, error) {
	var stats Stats
	stats.LSMSizeBytes, stats.ValueLogSizeBytes = b.db.Size()
	stats.TotalSizeBytes = stats.LSMSizeBytes + stats.ValueLogSizeBytes
	if metrics := b.db.BlockCacheMetrics(); metrics != nil {
		stats.BlockCache = cacheStats(metrics)
	}
	if metrics := b.db.IndexCacheMetrics(); metrics != nil {
		stats.IndexCache = cacheStats(metrics)
	}

	err := b.db.View(func(txn *badger.Txn) error {
		iterator := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iterator.Close()
		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			stats.KeyCount++
		}
		return nil
	})
	return stats, err
}

type cacheMetrics interface {
	Hits() uint64
	Misses() uint64
	KeysAdded() uint64
	KeysUpdated() uint64
	KeysEvicted() uint64
	CostAdded() uint64
	CostEvicted() uint64
}

func cacheStats(metrics cacheMetrics) CacheStats {
	return CacheStats{
		Hits: metrics.Hits(), Misses: metrics.Misses(),
		KeysAdded: metrics.KeysAdded(), KeysUpdated: metrics.KeysUpdated(),
		KeysEvicted: metrics.KeysEvicted(), CostAdded: metrics.CostAdded(),
		CostEvicted: metrics.CostEvicted(),
	}
}
