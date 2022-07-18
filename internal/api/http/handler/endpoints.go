package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kosdirus/andintern/internal/model"
	"net/http"
	"strconv"
)

var cars = []model.Car{
	{1, "audi", 40000},
	{2, "bmw", 55555},
	{3, "volvo", 65015},
	{4, "maserati", 150450},
	{5, "mercedes", 123000},
}

func getCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var car model.Car
	idStr, ok := values["id"]
	if !ok {
		brands, ok := values["brand"]
		if !ok {
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("%v", cars)))
			return
		}
		brand := brands[0]
		var present bool
		for _, v := range cars {
			if v.Brand == brand {
				car = v
				present = true
				break
			}
		}
		if !present {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("no car in database with brand: %s", brand)))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("requested data for brand: %v, car: %v.\n Values:%v", brand[0], car, values)))
		return
	}
	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("wrong id:\"%s\" query parameter: should be integer", idStr)))
		return
	}
	var present bool
	for _, v := range cars {
		if v.Id == id {
			car = v
			present = true
			break
		}
	}
	if !present {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("no car in database with id: %d", id)))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("requested data for id: %d, car: %v.\n Values:%v", id, car, values)))

}

func createCar(w http.ResponseWriter, r *http.Request) {
	car := model.Car{}
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("wrong body parameters"))
		return
	}

	for _, v := range cars {
		if car.Brand == v.Brand {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("provided car brand already exists: '%s'", car.Brand)))
			return
		}
	}
	cars = append(cars, car)
	w.Write([]byte(fmt.Sprintf("created car: %v", car)))
}

func updateCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	idStr, ok := values["id"]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("missing id in URL query"))
		return
	}
	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("wrong 'id' query parameter: should be integer"))
		return
	}

	car := model.Car{}
	err = json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("wrong body parameters"))
		return
	}

	for i := range cars {
		if id == cars[i].Id {
			w.WriteHeader(200)
			if car.Brand != "" {
				cars[i].Brand = car.Brand
			}
			if car.Price != 0 {
				cars[i].Price = car.Price
			}
			w.Write([]byte(fmt.Sprintf("updated car by id: %d, car: %v", id, cars[i])))
			return
		}
	}

	w.WriteHeader(400)
	w.Write([]byte(fmt.Sprintf("no car with such id: %d", id)))
}

func deleteCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if idStr, ok := values["id"]; ok {
		id, err := strconv.Atoi(idStr[0])
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("wrong 'id' query parameter: should be integer"))
			return
		}
		carsLen := len(cars)
		for i, v := range cars {
			if v.Id == id {
				copy(cars[i:], cars[i+1:])
				cars = cars[:len(cars)-1]
				break
			}
		}
		if carsLen == len(cars) {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("no car in database with id:%d", id)))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("successfully deleted car with id:%d", id)))
		return
	}

	if brand, ok := values["brand"]; ok {
		carsLen := len(cars)
		for i, v := range cars {
			if v.Brand == brand[0] {
				copy(cars[i:], cars[i+1:])
				cars = cars[:len(cars)-1]
				break
			}
		}
		if carsLen == len(cars) {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("no car in database with brand:%s", brand)))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("successfully deleted car with id:%s", brand)))
		return
	}

	if priceLowerThanStr, ok := values["priceLowerThan"]; ok {
		priceLowerThan, err := strconv.Atoi(priceLowerThanStr[0])
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("wrong 'priceLowerThan' query parameter: should be integer"))
			return
		}
		carsLen := len(cars)
		for i := len(cars) - 1; i >= 0; i-- {
			if cars[i].Price <= priceLowerThan {
				cars = append(cars[:i], cars[i+1:]...)

			}
		}
		if carsLen == len(cars) {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("no car in database with price lower than:%d", priceLowerThan)))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("successfully deleted cars with price lower than:%d", priceLowerThan)))
		return
	}

	w.WriteHeader(400)
	w.Write([]byte(fmt.Sprintf("none of required fields provided (id/brand/pricelowerthan")))
}
