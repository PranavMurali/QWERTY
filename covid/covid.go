package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"encoding/json"
)


func main() {
    response, err := http.Get("https://coronavirus-19-api.herokuapp.com/countries/india")

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		var r map[string]interface{}
		json.Unmarshal([]byte(responseData), &r)
		fmt.Printf( "Country:  " +"args[1]"+"\n")
		fmt.Printf( "Cases:         %.f\n",r["cases"])
		fmt.Printf( "Deaths:        %.f\n",r["deaths"])
		fmt.Printf( "Recovered:     %.f\n",r["recovered"])
		fmt.Printf( "Active Cases:  %.f\n",r["active"])
		fmt.Printf( "Deaths Today:  %.f\n",r["todayDeaths"])
		fmt.Printf( "Cases Today:   %.f\n",r["todayCases"])		
		fmt.Printf( "Cases per 1 Million:         %.f\n",r["casesPerOneMillion"])
		fmt.Printf( "Deaths per 1 Million:        %.f\n",r["deathsPerOneMillion"])
		fmt.Printf( "Total Tests:                 %.f\n",r["totalTests"])
		fmt.Printf( "Tests per 1 Million:         %.f\n",r["testsPerOneMillion"])
}