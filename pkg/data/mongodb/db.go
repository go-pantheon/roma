package mongodb

import "go.mongodb.org/mongo-driver/v2/mongo"

type DB struct {
	DB *mongo.Database
}

func NewMongoDB(mdb *mongo.Database) *DB {
	return &DB{DB: mdb}
}
