package mongoDB

import (
	"context"
	"searchandfind/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostModel struct {
	DB *mongo.Client
}

func (m *PostModel) Insert(p *models.Post) (string, error) {
	collection := m.DB.Database("searchandfind").Collection("posts")
	insertResult, err := collection.InsertOne(context.TODO(), p)
	if err != nil {
		return "", err
	}
	id := insertResult.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (m *PostModel) Get(f string) (*models.Post, error) {
	var post *models.Post
	// filter := bson.D{{"_id", "ObjectId(\"" + f + "\")"}}
	// fmt.Println("ObjectId(\"" + f + "\")")
	objID, _ := primitive.ObjectIDFromHex(f)
	collection := m.DB.Database("searchandfind").Collection("posts")

	err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&post)
	// fmt.Println(post.ID, "ok")
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	var posts []*models.Post
	collection := m.DB.Database("searchandfind").Collection("posts")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var post *models.Post

		err := cur.Decode(&post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())
	return posts, nil
}
