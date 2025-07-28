package riverdata

import (
	"context"
	"fmt"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/schemas"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Update() {
	init_rivers()
	fmt.Println("Clearing Database...")
	_, err := db.RIVER_DEFINITIONS.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	fmt.Println("...Database Cleared!")

	rivers := init_rivers()

	for _, river := range rivers {
		fmt.Printf("inserting river: %s\n", river.RiverName)
		river_bson, err := bson.Marshal(river)
		if err != nil {
			fmt.Printf("unable to marshal %s", river.RiverName)
		}
		db.RIVER_DEFINITIONS.InsertOne(context.TODO(), river_bson)
	}
}

func init_rivers() []schemas.River {
	output := []schemas.River{}

	output = append(output, coloradoRiver)
	output = append(output, sanGabrielRiver)

	for _, river := range output {
		err := validate(river)
		if err != nil {
			panic(err)
		}
	}
	return output
}
