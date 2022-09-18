package SimplerMongo

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 定时初始化并返回客户端,传入连接地址
func InitClient(URI string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	return client
}

// 定时初始化并返回客户端,传入连接地址
func InitClientWithOptions(URI string) *mongo.Client {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	return client
}

// defer db.CloseClient(client)
// 关闭客户端，返回错误
func CloseClient(client *mongo.Client) error {
	error := client.Disconnect(context.TODO())
	return error
}

// 对primitive.M (是map[string]interface{}的封装) 返回格式化的json
func PrimitiveToJson(primitiveM interface{}) []byte {
	jsonData, err := json.MarshalIndent(primitiveM, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
	return jsonData
}

// 查找*mongo.Collection中key-value，返回搜索到的[]map数据,找不到返回空
func FindOne(collection *mongo.Collection, key string, value interface{}) primitive.M {
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.D{{Key: key, Value: value}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the key %s\n", key)
		return nil
	}
	if err != nil {
		panic(err)
	}
	return result
}

func Find(collection *mongo.Collection, filter interface{}) primitive.M {
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("No document was found with the key")
		return nil
	}
	if err != nil {
		panic(err)
	}
	return result
}

// 查找满足filter的数据
func FindM(collection *mongo.Collection, filter interface{}) []primitive.M {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}
	return results
}

// 在集合中添加一个数据，传入客户端和bson
func Add(collection *mongo.Collection, data interface{}) {
	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}
}

// 在集合中添加多个数据，传入客户端和包含他们的切片
func AddM(collection *mongo.Collection, data []interface{}) {
	result, err := collection.InsertMany(context.TODO(), data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		fmt.Printf("\t%s\n", id)
	}
}

// 更新集合的与filter对应的update数值，返回修改数
func Update(collection *mongo.Collection, filter, update interface{}) int64 {
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result.ModifiedCount
}

// 更新集合的与filter对应的所有update数值，返回修改数
func UpdateM(collection *mongo.Collection, filter, update interface{}) int64 {
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result.ModifiedCount
}

// 删除filter,返回删除数（1）
func Delete(collection *mongo.Collection, filter interface{}) int64 {
	results, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return results.DeletedCount
}

// 删除满足的filter,返回删除数
func DeleteM(collection *mongo.Collection, filter interface{}) int64 {
	results, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return results.DeletedCount
}

// 在filter中以replacement替换,返回修改数
func Replace(collection *mongo.Collection, filter, replacement interface{}) int64 {
	result, err := collection.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		panic(err)
	}
	return result.ModifiedCount
}

// 获取满足filter的文档数并返回
func DocCount(collection *mongo.Collection, filter interface{}) int64 {
	estCount, estCountErr := collection.EstimatedDocumentCount(context.TODO())
	if estCountErr != nil {
		panic(estCountErr)
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d of %d\n", count, estCount)
	return count
}

// 为filter添加字段
func UpdateOneField(collection *mongo.Collection, filter, update interface{}) {
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}

// 查找filter中的不重复字段并返回包含他们的切片
func Distinct(collection *mongo.Collection, key string, filter interface{}) []interface{} {
	results, err := collection.Distinct(context.TODO(), key, filter)
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
	return results
}

// 对数据库执行传入的command命令
func CMD(database *mongo.Database, command interface{}) primitive.M {
	var result bson.M
	err := database.RunCommand(context.TODO(), command).Decode(&result)
	if err != nil {
		panic(err)
	}
	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", output)
	return result
}

// 根基传入的filter删除整个数据
func FindDelete(collection *mongo.Collection, filter interface{}) {
	var deletedDoc bson.D
	err := collection.FindOneAndDelete(context.TODO(), filter).Decode(&deletedDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(deletedDoc)
	fmt.Println(len(deletedDoc))
}

// 对filter更新update，返回更新后的数据
func FindUpdate(collection *mongo.Collection, filter, update interface{}) primitive.D {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedDoc bson.D
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(updatedDoc)
	return updatedDoc
}

// 将满足filter清空并替换为replacement
func FindReplace(collection *mongo.Collection, filter, replacement interface{}) {
	var previousDoc bson.D
	err := collection.FindOneAndReplace(context.TODO(), filter, replacement).Decode(&previousDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(previousDoc)
}
