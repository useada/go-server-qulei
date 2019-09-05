package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"a.com/go-server/service/board/model"
)

func (k *kv) GetSummaries(ctx context.Context, oids []string) (model.Summaries, error) {
	keys := make([]string, 0)
	for _, oid := range oids {
		keys = append(keys, k.genSummaryKey(oid))
	}
	vals, err := k.Pool.MGetBytes(ctx, keys)
	if err != nil {
		return nil, err
	}

	items := make(model.Summaries, 0)
	for _, val := range vals {
		item := model.Summary{}
		if err = json.Unmarshal(val, &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, err
}

func (k *kv) NewSummary(ctx context.Context, pitem *model.Summary) error {
	data, err := json.Marshal(pitem)
	if err != nil {
		return err
	}
	return k.Pool.Set(ctx, k.genSummaryKey(pitem.ID), data, 3600)
}

func (k *kv) DelSummary(ctx context.Context, oid string) error {
	return k.Pool.Delete(ctx, k.genSummaryKey(oid))
}

func (k *kv) genSummaryKey(oid string) string {
	return fmt.Sprintf("BDSUM|%s", oid)
}
