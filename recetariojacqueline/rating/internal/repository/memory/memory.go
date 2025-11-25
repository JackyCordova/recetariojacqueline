package memory

import (
	"context"
	"sync"

	"recetariojacqueline.com/rating/pkg/model"
)

type MemoryRepo struct {
	mu     sync.Mutex
	values map[model.RecordID][]float64
}

func New() *MemoryRepo {
	return &MemoryRepo{values: make(map[model.RecordID][]float64)}
}

// Implementa Repository
func (r *MemoryRepo) GetAverage(ctx context.Context, id model.RecordID, t model.RecordType) (float64, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	vals := r.values[id]
	if len(vals) == 0 {
		return 0, 0, nil
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals)), len(vals), nil
}

func (r *MemoryRepo) Put(ctx context.Context, id model.RecordID, t model.RecordType, user string, value float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[id] = append(r.values[id], value)
	return nil
}
