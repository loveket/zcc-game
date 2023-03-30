package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	gc "xiuianserver/config"
	"xiuianserver/elasticsearch/model"
)

var ESClient *elastic.Client

func NewESClient() {
	client, err := elastic.NewClient(elastic.SetURL(gc.GlobalConfig.ElasticConfig.HttpAddr), elastic.SetSniff(gc.GlobalConfig.ElasticConfig.Sniff))
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
		return
	}
	ESClient = client
}
func isExistIndex(client *elastic.Client, ctx context.Context, indexName string) bool {
	exists, _ := client.IndexExists(indexName).Do(ctx)
	if exists {
		return true
	}
	return false
}
func AddIndexData(client *elastic.Client, ctx context.Context, indexName string, id string) {
	if isExistIndex(client, ctx, indexName) {
		return
	}
	put1, err := client.Index().Index(indexName).Id(id).BodyJson(model.WeaponData).Do(ctx)
	if err != nil {
		fmt.Println("add index err", err)
		return
	}
	fmt.Printf("文档Id %s, 索引名 %s\n", put1.Id, put1.Index)
}
func GetIndexData(client *elastic.Client, ctx context.Context, indexName string, id string) {
	get, err := client.Get().Index(indexName).Id(id).Do(ctx)
	if err != nil {
		fmt.Println("add index err", err)
		return
	}
	if get.Found {
		msg := model.WeaponList{}
		data, _ := get.Source.MarshalJSON()
		json.Unmarshal(data, &msg)
		for _, weapon := range msg.Weapons {
			fmt.Println(weapon.WeaponName)
		}
	}
}
