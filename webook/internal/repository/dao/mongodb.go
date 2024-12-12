package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDao struct {
	// 代表的是制作库
	col *mongo.Collection
	// 代表的是线上库
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func (m *MongoDBDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	if err != nil {
		return 0, err
	}
	// 你没有自增主键
	return id, err
}

func (m *MongoDBDao) UpdateById(ctx context.Context, art Article) error {
	// 操作制作库
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.D{bson.E{"$set", bson.M{
		"title":   art.Title,
		"content": art.Content,
		"utime":   time.Now().UnixMilli(),
		"status":  art.Status,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	// 这里校验了author_id
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (m *MongoDBDao) Sync(ctx context.Context, art Article) (int64, error) {
	// 没法引入事务概念
	// 第一步，操作制作库
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	now := time.Now().UnixMilli()
	// 更新
	update := bson.E{"$set", art}
	upsert := bson.E{"$setOnInsert", bson.D{bson.E{"ctime", now}}}
	filter := bson.M{"id": art.Id}
	// 操作线上库, upset语义
	_, err = m.liveCol.UpdateOne(ctx, filter, bson.D{update, upsert}, options.Update().SetUpsert(true))
	return id, err

}

func (m *MongoDBDao) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	panic("implement me")
}

func NewMongoDBArticleDAO(mdb *mongo.Database, node *snowflake.Node) *MongoDBDao {
	return &MongoDBDao{
		col:     mdb.Collection("articles"),
		liveCol: mdb.Collection("published_articles"),
		node:    node,
	}
}
