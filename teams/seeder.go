package teams

import (
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"league/db"
	"league/models"

	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var userCollection *mongo.Collection = db.GetCollection(db.MongoClient, "users")
var adminEmail string
var adminUser models.User


func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	adminEmail = os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		log.Fatal("Error loading super admin details")
	}

	err = syncAdmin()
	if err != nil {
		log.Fatal("Error loading admin data")
	}
	if empty := isCollectionEmpty(teamCollection); empty {

		SeedTeams()
	}

}

func syncAdmin() error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err := userCollection.FindOne(ctx, bson.M{"email": strings.ToLower(adminEmail)}).Decode(&adminUser)
	if err != nil {
		// Handle error
		if err == mongo.ErrNoDocuments {
			// If no user found with the specified email
			return fmt.Errorf("user with the email %v is not found", adminEmail)
		}
		// If other error occurred
		return fmt.Errorf("failed to fetch user: %v", err)
	}
	return nil
}

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

func generateTeams() []interface{} {
	names := []string{"Arsenal", "Aston Villa", "Brentford", "Brighton & Hove Albion", "Burnley",
		"Chelsea", "Crystal Palace", "Everton", "Leeds United", "Leicester City",
		"Liverpool", "Manchester City", "Manchester United", "Newcastle United", "Norwich City",
		"Southampton", "Tottenham Hotspur", "Watford", "West Ham United", "Wolverhampton Wanderers"}

	stadiums := []string{"Emirates Stadium", "Villa Park", "Brentford Community Stadium", "Amex Stadium",
		"Turf Moor", "Stamford Bridge", "Selhurst Park", "Goodison Park", "Elland Road", "King Power Stadium",
		"Anfield", "Etihad Stadium", "Old Trafford", "St James' Park", "Carrow Road", "St Mary's Stadium",
		"Tottenham Hotspur Stadium", "Vicarage Road", "London Stadium", "Molineux Stadium"}

	sponsors := []string{"Adidas", "Nike", "Chevrolet", "Samsung", "Puma", "Audi", "Coca-Cola", "Amazon", "Toyota",
		"Visa", "Mastercard", "Microsoft", "Apple", "Google", "Facebook", "McDonald's", "Uber", "Tesla", "BMW", "Mercedes-Benz"}

	states := []string{"London", "Manchester", "Birmingham", "Liverpool", "Leeds", "Sheffield", "Bristol", "Newcastle upon Tyne",
		"Nottingham", "Leicester", "Sunderland", "Brighton", "Coventry", "Hull", "Stoke-on-Trent", "Wolverhampton",
		"Derby", "Swansea", "Southampton", "Aberdeen"}

	rand.Seed(time.Now().UnixNano())

	teams := make([]interface{}, len(names))
	for i := 0; i < len(names); i++ {
		team := models.Team{
			ID:          primitive.NewObjectID(),
			Name:        names[i],
			State:       states[i],
			Country:     "England",
			FoundedYear: rand.Intn(150) + 1871,
			Stadium:     stadiums[i],
			Sponsor:     sponsors[i],
			CreatedBy:   adminUser.Id,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		teams[i] = team
	}

	return teams
}

func SeedTeams() {
	teams := generateTeams()
	_, err := teamCollection.InsertMany(context.Background(), teams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Teams seeded successfully!")
}
