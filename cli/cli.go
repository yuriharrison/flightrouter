package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/thatisuday/commando"
	"github.com/yuriharrison/bexs-test/api"
	"github.com/yuriharrison/bexs-test/flightsdb"
	"github.com/yuriharrison/bexs-test/loader"
	"github.com/yuriharrison/bexs-test/util"
)

func cliBasic(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	fmt.Println("Importing file...")
	file := args["file"].Value
	db := flightsdb.New()
	loader.ImportFlightsFromFile(file, db)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Search for flight (e.g. GRU-APQ): ")
	rawRoute, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	origCode, destCode := util.FormatInputRoute(rawRoute)
	cheapestRoute, err := db.CheapestRoute(origCode, destCode)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Cheapest route available:", util.FormatRouteToString(cheapestRoute))
}

func startAPI(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	fmt.Println("Starting server ...")
	db := flightsdb.New()
	if csvFile, ok := flags["data"].Value.(string); ok && csvFile != "NULL" {
		loader.ImportFlightsFromFile(csvFile, db)
	}
	port := flags["port"].Value.(int)
	api.StartServer(db, port)
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
		SetDescription("Starts the Web API").
		AddFlag("port,p", "port", commando.Int, 8080).
		AddFlag("data,d", "csv file to pre fetch data", commando.String, "NULL").
		SetAction(startAPI)

	commando.Parse(nil)
}
