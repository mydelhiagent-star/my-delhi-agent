package services

import (
    "context"
    "myapp/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// Assume you have a GCPService with these methods:
// GenerateSignedURL(objectName string, expiry time.Duration) (string, error)
// PublicFileURL(objectName string) string

type PropertyService struct {
    PropertyCollection *mongo.Collection

}

func (s *PropertyService) CreateProperty(property models.Property) (primitive.ObjectID, error) {
    property.ID = primitive.NewObjectID()
    _, err := s.PropertyCollection.InsertOne(context.Background(), property)
    return property.ID, err
}

func (s *PropertyService) GetPropertyByID(id primitive.ObjectID) (*models.Property, error) {
    var property models.Property
    err := s.PropertyCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&property)
    if err != nil {
        return nil, err
    }
    return &property, nil
}

func (s *PropertyService) UpdateProperty(id primitive.ObjectID, updates models.Property) error {
    update := bson.M{"$set": updates}
    _, err := s.PropertyCollection.UpdateByID(context.Background(), id, update)
    return err
}

func (s *PropertyService) DeleteProperty(id primitive.ObjectID) error {
    _, err := s.PropertyCollection.DeleteOne(context.Background(), bson.M{"_id": id})
    return err
}

func (s *PropertyService) GetAllProperties(ctx context.Context) ([]models.Property, error) {
    cur, err := s.PropertyCollection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)

    var properties []models.Property
    if err := cur.All(ctx, &properties); err != nil {
        return nil, err
    }
    return properties, nil
}



