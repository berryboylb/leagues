package fixtures

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	
	"league/models"

	"context"
	"fmt"
	"log"
	"time"
)



func isCollectionEmpty(collection *mongo.Collection) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Println("Failed to count documents:", err)
		return false
	}

	return count == 0
}

func generateCompetitions() []interface{} {
	competitionNames := []string{
		"UEFA Champions League",
		"UEFA Europa League",
		"FA Cup",
		"EFL Cup (Carabao Cup)",
		"FIFA Club World Cup",
		"Community Shield",
	}
	competitionTypes := []string{
		"European",
		"European",
		"Domestic",
		"Domestic",
		"International",
		"Domestic",
	}

	competitions := make([]interface{}, 0)
	for i, competitionName := range competitionNames {
		competition := models.Competition{
			ID:        primitive.NewObjectID(),
			Name:      competitionName,
			Type:      competitionTypes[i],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		competitions = append(competitions, competition)
	}

	return competitions
}

func SeedComps() {
	_, err := competitionCollection.InsertMany(context.Background(), generateCompetitions())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("competitions seeded successfully!")
}

func init() {
	if empty := isCollectionEmpty(competitionCollection); empty {
		SeedComps()
	}

}
