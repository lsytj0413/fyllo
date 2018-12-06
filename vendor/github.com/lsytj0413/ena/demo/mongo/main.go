package main

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// User ...
type User struct {
	Id        bson.ObjectId `bson:"_id"`
	Name      string        `bson:"name"`
	Parent    string        `bson:"parent"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}

func (u User) String() string {
	return fmt.Sprintf("User(Id=%s,Name=%s,Parent=%s,CreatedAt=%s,UpdatedAt=%s)", string(u.Id), u.Name, u.Parent, u.CreatedAt.String(), u.UpdatedAt.String())
}

const url string = "mongodb://anquanz:anquanz@127.0.0.1:27017/anfu?maxPoolSize=4096&minPoolSize=120"

var c *mgo.Collection
var c2 *mgo.Collection
var session *mgo.Session

func init() {
	var err error
	session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	c = session.DB("anfu").C("AssetCategory")
	c2 = session.DB("anfu").C("vuls")
}

func find() {
	session := session.Copy()
	defer session.Close()

	fmt.Println("find01: -------------------")
	var users []User
	c.Find(bson.M{"name": "newparent01"}).All(&users)
	for _, value := range users {
		fmt.Println(value)
	}

	fmt.Println("find02: -------------------")
	users = make([]User, 0)
	c.Find(bson.M{"$or": []bson.M{bson.M{"name": "newparent01"}, bson.M{"name": "newparent02"}}}).Select(bson.M{"_id": 1, "name": 1}).All(&users)
	for _, value := range users {
		fmt.Println(value)
	}

	fmt.Println("find03: -------------------")
	users = make([]User, 0)
	c.Find(bson.M{"name": bson.RegEx{Pattern: "Parent", Options: "i"}}).All(&users)
	for _, value := range users {
		fmt.Println(value)
	}
}

func agg() {
	session := session.Copy()
	defer session.Close()

	resp := []bson.M{}
	pipe := c.Pipe([]bson.M{
		{"$match": bson.M{"name": bson.M{"$exists": true}}},
		{"$group": bson.M{
			"_id": "$parent",
			"category": bson.M{
				"$push": bson.M{
					"_id":    "$_id",
					"name":   "$name",
					"parent": "$parent",
				},
			},
		}},
	})
	err := pipe.All(&resp)
	if err != nil {
		panic(err)
	}
	for _, v := range resp {
		fmt.Println(v)
	}

	fmt.Println("\n\n\nagg vuls: --------------------------")
	resp = []bson.M{}
	pipe = c2.Pipe([]bson.M{
		{"$match": bson.M{"detail.id": bson.M{"$exists": true}}},
		{"$lookup": bson.M{
			"from":         "sites",
			"localField":   "site",
			"foreignField": "_id",
			"as":           "site",
		}},
		{"$unwind": "$site"},
		{"$facet": bson.M{
			"new": []bson.M{
				{
					"$match": bson.M{
						"fixStatus": "new",
					},
				},
				{
					"$group": bson.M{
						"_id": "$site.user",
						"count": bson.M{
							"$sum": 1,
						},
					},
				},
				{
					"$project": bson.M{
						"user":  "$_id",
						"count": "$count",
						"name":  "new",
					},
				},
			},
			"normal": []bson.M{
				{"$match": bson.M{
					"fixStatus": "normal",
				}},
				{
					"$group": bson.M{
						"_id": "$site.user",
						"count": bson.M{
							"$sum": 1,
						},
					},
				},
				{
					"$project": bson.M{
						"user":  "$_id",
						"count": "$count",
						"name":  "normal",
					},
				},
			},
		}},
		{"$project": bson.M{
			"status": bson.M{
				"$concatArrays": []string{"$new", "$normal"},
			},
		}},
		{"$unwind": "$status"},
		{
			"$group": bson.M{
				"_id": "$status.user",
				"user": bson.M{
					"$first": "$status.user",
				},
				"status": bson.M{
					"$push": bson.M{
						"name":  "$status.name",
						"count": "$status.count",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{
			"$project": bson.M{
				"_id":    "$_id",
				"user":   "$user._id",
				"name":   "$user.realname",
				"status": "$status",
			},
		},
	})
	err = pipe.All(&resp)
	if err != nil {
		panic(err)
	}
	for _, v := range resp {
		fmt.Println(v)
	}
}

func main() {
	find()

	fmt.Println("agg: -----------------------")
	agg()
	fmt.Println("done")
}
