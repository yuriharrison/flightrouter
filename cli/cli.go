package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/thatisuday/commando"
	"github.com/yuriharrison/flightrouter/api"
	"github.com/yuriharrison/flightrouter/flightsdb"
	"github.com/yuriharrison/flightrouter/loader"
	"github.com/yuriharrison/flightrouter/util"
)

func cliBasic(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	fmt.Println("Importing file...")
	file := args["file"].Value
	db := flightsdb.New()
	loader.ImportFlightsFromFile(file, db)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Search for flight (e.g. GRU-APQ): ")
		rawRoute, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		origCode, destCode := util.FormatInputRoute(rawRoute)
		cheapestRoute, err := db.CheapestRoute(origCode, destCode)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Cheapest route available:", util.FormatRouteToString(cheapestRoute))
	}
}

func startAPI(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	fmt.Println("Starting server ...")
	db := flightsdb.New()
	if csvFile, ok := flags["data"].Value.(string); ok {
		loader.ImportFlightsFromFile(csvFile, db)
		db.SetFile(csvFile)
	}
	api.StartServer(db)
}

// Run configure commando and start application
func Run() {
	commando.
		SetExecutableName("flightrouter").
		SetVersion("1.0.0").
		SetDescription("Flight router engine.")

	commando.
		Register(nil).
		AddArgument("file", "flights csv file", "").
		SetAction(cliBasic)

	commando.
		Register("api").
		SetDescription("Starts the Web API on port 8080").
		AddFlag("data,d", "csv file", commando.String, "./data.csv").
		SetAction(startAPI)

	commando.Parse(nil)
}
