package fixtures

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"league/db"
	"league/models"

	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var teamCollection *mongo.Collection = db.GetCollection(db.MongoClient, "teams")

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

func fetchTeams() ([]models.Team, error) {
	var teams []models.Team
	cursor, err := teamCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find teams: %w", err)
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &teams); err != nil {
		return nil, fmt.Errorf("failed to decode teams: %w", err)
	}
	return teams, nil
}

func fetchComp() ([]models.Competition, error) {
	var teams []models.Competition
	cursor, err := competitionCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find Competition: %w", err)
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &teams); err != nil {
		return nil, fmt.Errorf("failed to decode Competition: %w", err)
	}
	return teams, nil
}

func generateFixtures(competition []models.Competition, teams []models.Team) []interface{} {
	fixtures := make([]interface{}, 0)
	ref := []string{"Adidas", "Nike", "Chevrolet", "Samsung", "Puma", "Audi", "Coca-Cola", "Amazon", "Toyota",
		"Visa", "Mastercard", "Microsoft", "Apple", "Google", "Facebook", "McDonald's", "Uber", "Tesla", "BMW", "Mercedes-Benz"}
	formation := []string{"4-4-2", "4-3-3", "3-4-3", "4-3-2", "3-4-2", "3-3-4", "2-3-4", "3-2-4", "2-4-3", "2-3-3"}
	for i := 0; i < 299; i++ {
		hash, _ := generateRandomString(10)
		awayLineUp := make([]string, 0)
		homeLineUp := make([]string, 0)

		awaySubs := make([]string, 0)
		homeSubs := make([]string, 0)
		for i := 0; i < 11; i++ {
			awayLineUp = append(awayLineUp, ref[rand.Intn(len(ref))])
			homeLineUp = append(homeLineUp, ref[rand.Intn(len(ref))])
		}

		for i := 0; i < 5; i++ {
			awaySubs = append(awayLineUp, ref[rand.Intn(len(ref))])
			homeSubs = append(homeLineUp, ref[rand.Intn(len(ref))])
		}
		fixtures = append(fixtures, models.Fixture{
			HomeTeamID:    teams[rand.Intn(len(teams))].ID,
			AwayTeamID:    teams[rand.Intn(len(teams))].ID,
			CompetitionID: competition[rand.Intn(len(competition))].ID,
			Status:        []models.Status{models.Ongoing, models.Completed, models.Pending}[rand.Intn(3)],
			Date:          time.Now(),
			Stadium: []string{"Emirates Stadium", "Villa Park", "Brentford Community Stadium", "Amex Stadium",
				"Turf Moor", "Stamford Bridge", "Selhurst Park", "Goodison Park", "Elland Road", "King Power Stadium",
				"Anfield", "Etihad Stadium", "Old Trafford", "St James' Park", "Carrow Road", "St Mary's Stadium",
				"Tottenham Hotspur Stadium", "Vicarage Road", "London Stadium", "Molineux Stadium"}[rand.Intn(20)],
			Referee:    fmt.Sprint("%v %v", ref[rand.Intn(len(ref))], ref[rand.Intn(len(ref))]),
			UniqueLink: hash,
			Away: models.Details{
				Substitutes:    awaySubs,
				Lineup:         awayLineUp,
				Formation:      formation[rand.Intn(10)],
				Goals:          rand.Intn(10),
				GoalScorers:    make([]string, 0),
				Shots:          rand.Intn(10),
				ShotsOnTarget:  rand.Intn(10),
				Possession:     rand.Float64()*(45.0-1.0) + 1.0,
				Passes:         rand.Intn(10),
				PassesAccuracy: rand.Intn(10),
				Fouls:          rand.Intn(10),
				YellowCards:    rand.Intn(10),
				RedCards:       rand.Intn(10),
				OffSides:       rand.Intn(10),
				Corners:        rand.Intn(10),
			},
			Home: models.Details{
				Substitutes:    homeSubs,
				Lineup:         homeLineUp,
				Formation:      formation[rand.Intn(10)],
				Goals:          rand.Intn(10),
				GoalScorers:    make([]string, 0),
				Shots:          rand.Intn(10),
				ShotsOnTarget:  rand.Intn(10),
				Possession:     rand.Float64()*(45.0-1.0) + 1.0,
				Passes:         rand.Intn(10),
				PassesAccuracy: rand.Intn(10),
				Fouls:          rand.Intn(10),
				YellowCards:    rand.Intn(10),
				RedCards:       rand.Intn(10),
				OffSides:       rand.Intn(10),
				Corners:        rand.Intn(10),
			},
		},
		)
	}
	return fixtures
}

func init() {
	if empty := isCollectionEmpty(competitionCollection); empty {
		SeedComps()
	}

	fmt.Println("starting fixtures seeding")
	if fixtureEmpty := isCollectionEmpty(fixtureCollection); fixtureEmpty {
		teamsCount, err := teamCollection.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			fmt.Printf("Error counting documents in teams collection: %v\n", err)
			return
		}

		fmt.Println("number of teams", teamsCount)

		if teamsCount == 0 {
			fmt.Println("No teams found, waiting for teams collection to be seeded...")
			// Implement logic to wait until teams are seeded (e.g., using a loop with a sleep)
			for {
				time.Sleep(1 * time.Minute) // Adjust sleep duration as needed
				teamsCount, err = teamCollection.CountDocuments(context.Background(), bson.M{})
				if err != nil {
					fmt.Printf("Error counting documents in teams collection: %v\n", err)
					return
				}
				if teamsCount > 0 {
					break
				}
			}
			fmt.Println("Teams collection has been seeded.")
		}

		competitionCount, err := competitionCollection.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			fmt.Printf("Error counting documents in competition collection: %v\n", err)
			return
		}

		if competitionCount == 0 {
			fmt.Println("No competition found, waiting for competition collection to be seeded...")
			// Implement logic to wait until teams are seeded (e.g., using a loop with a sleep)
			for {
				time.Sleep(1 * time.Minute) // Adjust sleep duration as needed
				competitionCount, err = competitionCollection.CountDocuments(context.Background(), bson.M{})
				if err != nil {
					fmt.Printf("Error counting documents in competition collection: %v\n", err)
					return
				}
				if competitionCount > 0 {
					break
				}
			}
			fmt.Println("competition collection has been seeded.")
		}

		

		comps, err := fetchComp()
		if err != nil {
			fmt.Printf("Error counting documents in competition collection: %v\n", err)
			return
		}

		teams, err := fetchTeams()
		if err != nil {
			fmt.Printf("Error counting documents in teams collection: %v\n", err)
			return
		}

		fixtures := generateFixtures(comps, teams)
		_, err = fixtureCollection.InsertMany(context.Background(), fixtures)
		if err != nil {
			fmt.Printf("Error inserting fixtures: %v\n", err)
			return
		}
		fmt.Println("Fixtures collection has been seeded.")

	}


}
