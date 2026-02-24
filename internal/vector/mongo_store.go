package vector

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoStore struct {
	coll *mongo.Collection
}

func NewMongoStore(db *mongo.Database, coll string) *MongoStore {
	return &MongoStore{
		coll: db.Collection(coll),
	}
}

func (m *MongoStore) Upsert(ctx context.Context, docs []Document) error {
	var models []mongo.WriteModel

	for _, d := range docs {
		models = append(models, mongo.NewReplaceOneModel().
			SetFilter(bson.M{"_id": d.ID}).
			SetReplacement(bson.M{
				"_id":       d.ID,
				"domain":    d.Domain,
				"text":      d.Text,
				"embedding": d.Vector,
			}).
			SetUpsert(true))
	}

	_, err := m.coll.BulkWrite(ctx, models)
	return err
}

func (m *MongoStore) Search(
	ctx context.Context,
	vec []float32,
	domain string,
	topK int,
) ([]Hit, error) {

	pipeline := mongo.Pipeline{
		{{
			Key: "$vectorSearch", Value: bson.M{
				"index":         "vector_index",
				"path":          "embedding",
				"queryVector":   vec,
				"numCandidates": 100,
				"limit":         topK,
				"filter": bson.M{
					"domain": domain,
				},
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"text":   1,
				"domain": 1,
				"score":  bson.M{"$meta": "vectorSearchScore"},
			},
		}},
	}

	cur, err := m.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []Hit
	for cur.Next(ctx) {
		var r struct {
			ID     string  `bson:"_id"`
			Text   string  `bson:"text"`
			Domain string  `bson:"domain"`
			Score  float64 `bson:"score"`
		}
		if err := cur.Decode(&r); err != nil {
			continue
		}

		out = append(out, Hit{
			ID: r.ID, Text: r.Text, Score: r.Score, Domain: r.Domain,
		})
	}

	return out, nil
}
