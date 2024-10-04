package mongodb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/win30221/core/http/catch"
	"github.com/win30221/core/syserrno"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrNoDocuments = catch.New(syserrno.Mongo, "no documents in result", "no documents in result")
)

// ToDoc 會將 struct 轉為 bson 格式
//
// example:
//
//	type Sample struct {
//			Id       string    `bson:"-"`
//			Title    string    `bson:"title,omitempty"`
//			ExpireAt time.Time `bson:"expireAt,omitempty"`
//	}
//
//	s := Sample{
//		ExpireAt: time.Now()
//	}
//
// m, _ := ToDoc(s)
//
// // m = &map[expireAt:1641516571365]
// fmt.Println(m)
func ToDoc(v any) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		err = catch.NewWitStack(syserrno.Mongo, "marshal bson error", fmt.Sprintf("marshal bson error. err: %s", err.Error()), 3)
		return
	}

	err = bson.Unmarshal(data, &doc)
	if err != nil {
		err = catch.NewWitStack(syserrno.Mongo, "unmarshal bson error", fmt.Sprintf("marshal bson error. err: %s", err.Error()), 3)
		return
	}

	return
}

// StructToBSONMap 會將 struct 轉為 bson 格式
//
// example:
//
// type Shipping struct {
// 		Id            *primitive.ObjectId `bson:"_id"`
// 		Type          *uint8              `bson:"type"`
// 		ItemIds       *[]int              `bson:"itemIds"`
// 		Address       *string             `bson:"address,omitempty"`
// 		StoreNumber   *string             `bson:"storeNumber,omitempty"`
// 		StoreName     *string             `bson:"storeName,omitempty"`
// 		RecipientName *string             `bson:"recipientName"`
// 		Phone         *string             `bson:"phone"`
//		Status        *uint8              `bson:"status"`
// 		CreatedAt     *time.Time          `bson:"createdAt"`
// 		UpdatedAt     *time.Time          `bson:"updatedAt"`
// }
//
// s := Shipping{
// 		Type:      1,
// 		Address:   "台北市中正區重慶南路一段122號",
// 		UpdatedAt: time.Now(),
// }
//
// m, _ := StructToBSONMap(s)
//
// // m = map[address:0xc000793b60 type:0xc0011a0ed0 updatedAt:2024-07-31 18:58:30.7681785 +0900 JST m=+18.554084601]
// fmt.Println(m)

func StructToBSONMap(data interface{}) bson.M {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	result := bson.M{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("bson")

		if fieldType.Name == "Id" {
			continue
		}

		if tag == "" || tag == "-" {
			continue
		}

		// 處理 bson:"xxx,omitempty"
		tagParts := strings.Split(tag, ",")
		tagKey := tagParts[0]

		// 如果欄位是指標類型且為 nil，則跳過不設定
		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		if tagKey == "" {
			tagKey = fieldType.Name
		}

		result[tagKey] = field.Interface()
	}

	return result
}
