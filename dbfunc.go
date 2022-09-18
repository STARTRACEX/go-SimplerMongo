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

// 传入连接地址,初始化并返回客户端
func InitClient(URI string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	return client
}

// 传入连接地址,定时初始化并返回客户端
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

// 关闭客户端，返回错误
func CloseClient(client *mongo.Client) error {
	error := client.Disconnect(context.TODO())
	return error
}

// 对primitive.M (map[string]interface{}) 格式化返回json格式化的byte切片
func PrimitiveToJson(primitive interface{}) []byte {
	jsonData, err := json.MarshalIndent(primitive, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
	return jsonData
}

// 查找Key: “key” : Value: value的数据并返回
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

// 查找满足filter的第一个数据并返回
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

// 查找满足filter的数据并返回
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

// 添加一个数据
func Add(collection *mongo.Collection, data interface{}) {
	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}
}

// 添加多个数据
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

// 对满足filter的数据更新为update，返回修改数
func Update(collection *mongo.Collection, filter, update interface{}) int64 {
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result.ModifiedCount
}

// 对满足filter的数据更新为update，返回修改数
func UpdateM(collection *mongo.Collection, filter, update interface{}) int64 {
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return result.ModifiedCount
}

// 删除满足filter的数据,返回删除数
func Delete(collection *mongo.Collection, filter interface{}) int64 {
	results, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return results.DeletedCount
}

// 删除满足filter的数据,返回删除数
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

// 为满足filter的数据添加字段
func UpdateOneField(collection *mongo.Collection, filter, update interface{}) {
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}

// 查找满足filter的数据中的不重复字段，返回包含他们的切片
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

// 对数据库执行command命令，返回结果
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

// 删除满足filter的整个数据
func FindDelete(collection *mongo.Collection, filter interface{}) {
	var deletedDoc bson.D
	err := collection.FindOneAndDelete(context.TODO(), filter).Decode(&deletedDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(deletedDoc)
	fmt.Println(len(deletedDoc))
}

// 将满足filter的数据更新为update，返回更新后的数据
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

// 将满足filter的数据替换为replacement
func FindReplace(collection *mongo.Collection, filter, replacement interface{}) {
	var previousDoc bson.D
	err := collection.FindOneAndReplace(context.TODO(), filter, replacement).Decode(&previousDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(previousDoc)
}
