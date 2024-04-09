package mongodb

import (
	"fmt"

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
//			ID       string    `bson:"-"`
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
