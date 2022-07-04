package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mysql/app/entity"
	"mysql/app/service"
	"mysql/http"
	apphttp "mysql/http"
	appsql "mysql/sql"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
)

type City struct {
	Id         int64
	Name       string
	Population int
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/go test?parseTime=true")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	if err := run(ctx, db); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}

// Run starts the main application
func run(ctx context.Context, db *sql.DB) error {

	sqlCityService := appsql.NewCityService(db)

	HTTPServerAPI := apphttp.NewServerAPI()

	HTTPServerAPI.Addr = ":8080"
	HTTPServerAPI.CityService = sqlCityService

	if err := HTTPServerAPI.Open(); err != nil {
		return err
	}

	if HTTPServerAPI.UseTLS() {
		go func() {
			log.Fatal(http.ListenAndServeTLSRedirect(""))
		}()
	}

	return nil
}

func testSql() {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/go-test?parseTime=true")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	ctx := context.Background()

	sqlCityService := appsql.NewCityService(db)

	if err := printTable(ctx, db, sqlCityService); err != nil {
		panic(err)
	}

	fmt.Println("--------------------------")

	prato := entity.City{
		Name:       "Prato",
		Population: 200000,
	}

	if err := sqlCityService.CreateCity(ctx, &prato); err != nil {
		panic(err)
	}

	/*if err := sqlCityService.DeleteCity(ctx, 4); err != nil {
		panic(err)
	}*/

	newPop := 0
	cup := service.CityUpdate{Population: &newPop}

	if err := sqlCityService.UpdateCity(ctx, 8, cup); err != nil {
		panic(err)
	}

	popLte := 400000
	popGte := 600000
	popEq := 200000

	if cities, err := sqlCityService.FindCityByPopulationLte(ctx, popLte); err != nil {
		panic(err)
	} else {
		fmt.Println("Città filtrate con popolazione < ", popLte)
		for _, city := range cities {
			fmt.Printf("%v\n", city)
		}
		fmt.Println("--------------------------")
	}

	if cities, err := sqlCityService.FindCityByPopulationGte(ctx, popGte); err != nil {
		panic(err)
	} else {
		fmt.Println("Città filtrate con popolazione > ", popGte)
		for _, city := range cities {
			fmt.Printf("%v\n", city)
		}
		fmt.Println("--------------------------")
	}

	if cities, err := sqlCityService.FindCityByPopulation(ctx, popEq); err != nil {
		panic(err)
	} else {
		fmt.Println("Città filtrate con popolazione = ", popEq)
		for _, city := range cities {
			fmt.Printf("%v\n", city)
		}
		fmt.Println("--------------------------")
	}

	if err := printTable(ctx, db, sqlCityService); err != nil {
		panic(err)
	}
}

/*func createCity(db *sql.DB, city *City) error {
	sql := "INSERT INTO cities(name, population) VALUES (?,?)"
	res, err := db.Exec(sql, city.Name, city.Population)

	if err != nil {
		return err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return err
	}

	city.Id = lastId
	return nil
}

func deleteCityById(db *sql.DB, Id int64) error {
	sql := "DELETE FROM cities WHERE id = ?"
	_, err := db.Exec(sql, Id)

	if err != nil {
		return err
	}

	return nil
}

func findCityById(db *sql.DB, Id int64) (*City, error) {
	rs := db.QueryRow("SELECT * FROM `cities` WHERE id = ?;", Id)
	var praga City
	if err := rs.Scan(&praga.Id, &praga.Name, &praga.Population); err != nil {
		return nil, err
	}
	return &praga, nil
}

func updateCityById(db *sql.DB, newCity City, Id int64) error {
	sql := "UPDATE cities SET Name = ?, Population = ? WHERE id = ?;"
	_, err := db.Exec(sql, newCity.Name, newCity.Population, Id)
	if err != nil {
		return err
	}

	return nil
}

func updateCityByName(db *sql.DB, name string, newPopulation int) error {
	sql := "UPDATE cities SET Population = ? WHERE name = ?;"
	_, err := db.Exec(sql, newPopulation, name)
	if err != nil {
		return err
	}

	return nil
}*/

func printTable(ctx context.Context, db *sql.DB, cityService service.CityService) error {

	cities, err := cityService.FindCities(ctx, service.CityFilter{})
	if err != nil {
		return err
	}
	for _, city := range cities {
		fmt.Printf("%v\n", city)
	}

	return nil
}
