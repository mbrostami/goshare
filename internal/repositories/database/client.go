package database

import (
	"fmt"
)

func (r *Repository) AddServer(addr string) error {
	return r.insertAndOverride(clientBucketKey, fmt.Sprintf("server-%s", addr), addr)
}
