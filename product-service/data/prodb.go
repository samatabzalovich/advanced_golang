package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: Product{},
	}
}

type Models struct {
	LogEntry Product
}

type Product struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	Price     float32   `bson:"price" json:"price"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *Product) Insert(entry Product) error {
	collection := client.Database("logs").Collection("products")

	_, err := collection.InsertOne(context.TODO(), Product{
		Name:      entry.Name,
		Data:      entry.Data,
		Price:     entry.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (l *Product) All() ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("products")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*Product

	for cursor.Next(ctx) {
		var item Product

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *Product) GetOne(id string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("products")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry Product
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *Product) Get(maxPrice float64, minPrice float64) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("products")

	filter := bson.D{}
	if minPrice != 0 || maxPrice != 0 {
		filter = bson.D{
			{"$and",
				bson.A{
					bson.D{{"price", bson.D{{"$gt", minPrice}}}},
					bson.D{{"price", bson.D{{"$lte", maxPrice}}}},
				},
			},
		}
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var Products []*Product

	for cursor.Next(ctx) {
		var item Product

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			Products = append(Products, &item)
		}
	}

	return Products, nil

}

func (l *Product) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("products")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *Product) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("products")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"price", l.Price},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
