# go-SimplerMongo

## 前要 
_为更容易操作[Mongodb Atlas](https://www.mongodb.com/atlas)而构建此包，目前不完善。_

_关于Mongodb的操作详情，请前往[MongoDB Go Driver — Go](https://www.mongodb.com/docs/drivers/go/v1.8/)_



## 引入

安装此包

```shell
go get github.com/STARTRACEX/go-SimplerMongo
```

**在项目中引入**

```go
import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
   FUNCNAME "github.com/STARTRACEX/go-SimplerMongo"
)
```

或

拷贝dbfunc.go到您的目录中,安装不存在的包依赖（如有）

```shell
go mod tidy
```


## 初始化

```go

const URI = "Connection uri ....... "

defer func() {
		if r := recover(); r != nil {
			fmt.Println("Received panic:", r)
		}
}()

client := FUNCNAME.InitClientWithOptions(URI)
collection := client.Database("...").Collection("...")
defer FUNCNAME.CloseClient(client)

   /* 
      ...
   */

```

## 部分预览
```go
// 传入连接地址,初始化并返回客户端
func InitClient(URI string) *mongo.Client 
// 传入连接地址,定时初始化并返回客户端
func InitClientWithOptions(URI string) *mongo.Client
// 关闭客户端，返回错误
func CloseClient(client *mongo.Client) error
// 查找满足filter的第一个数据并返回
func Find(collection *mongo.Collection, filter interface{}) primitive.M
// 查找满足filter的数据并返回
func FindM(collection *mongo.Collection, filter interface{}) []primitive.M
// 添加一个数据
func Add(collection *mongo.Collection, data interface{})
// 添加多个数据
func AddM(collection *mongo.Collection, data []interface{}) 
// 对满足filter的数据更新为update，返回修改数
func Update(collection *mongo.Collection, filter, update interface{}) int64
// 删除满足filter的数据,返回删除数
func Delete(collection *mongo.Collection, filter interface{}) int64
```


## 示例
```go
// 查找{ fild: [ { v1:1 } ... ] }
FUNCNAME.Find(collection, bson.D{{Key: "field", Value: bson.D{{Key: "v1", Value: "1"}}}})
FUNCNAME.FindM(collection, bson.D{{Key: "field", Value: bson.D{{Key: "v1", Value: "1"}}}})
FUNCNAME.FindOne(collection, "field",bson.D{{Key: "v1", Value: "1"}}) 


// 添加{ name: "@···" }
FUNCNAME.Add(collection, bson.D{{Key:"name",Value:"@···"}}

// 添加{ name:"@···" } , { name :"#···" }
docs := []interface{}{bson.D{{Key:"name", Value:"@···"}},bson.D{{Key:"name", Value:"#···"}}}
FUNCNAME.AddM(collection, docs)


// 对 [ { v1:1 } ...] ... 更新 {name:"@···"} 为 [ { v1:1 } ...] ... {name:"@···"}
update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: "@···"}}}}
filter := bson.D{{Key: "v1", Value: 1}}
FUNCNAME.FindUpdate(collection, filter, update)

//对[ {name:"@···"} ... ] 替换为 [ { name:"#···" },{ age:1649 } ]
filter := bson.D{{Key: "name", Value: "@···"}}
replacement := bson.D{{Key: "name", Value: "#···"},{Key: "age",Value: 1649}}
FUNCNAME.FindReplace(collection, filter, replacement)
```
