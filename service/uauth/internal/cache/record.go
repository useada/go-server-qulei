package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"a.com/go-server/service/uauth/internal/model"
)

func (k *kv) GetRecord(ctx context.Context, uid, deviceType string) (model.Record, error) {
	item := model.Record{}
	data, err := k.Pool.GetBytes(ctx, k.genRecordKey(uid, deviceType))
	if err != nil {
		return item, err
	}
	return item, json.Unmarshal(data, &item)
}

func (k *kv) NewRecord(ctx context.Context, r model.Record) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return k.Pool.Set(ctx, k.genRecordKey(r.Uid, r.DeviceType), data, 0)
}

func (k *kv) genRecordKey(uid, device string) string {
	return fmt.Sprintf("UALR|%s%s", uid, device)
}
