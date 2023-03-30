package elasticsearch

import (
	"context"
	"fmt"
	"testing"
	"xiuianserver/utils"
)

func TestElasticsearch(t *testing.T) {
	fmt.Println("***", utils.GetOsPwd())
	NewESClient()
	ctx := context.Background()
	indexName := "weaponlist2"
	indexId := "1000"
	AddIndexData(ESClient, ctx, indexName, indexId)
	GetIndexData(ESClient, ctx, indexName, indexId)
}
