package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kosdirus/andintern/internal/model"
	"net/http"
	"strconv"
	"strings"
)

var cars = []model.Car{
	{1, "audi", 40000},
	{2, "bmw", 55555},
	{3, "volvo", 65015},
	{4, "maserati", 150450},
	{5, "mercedes", 123000},
}

func (srv Server) getCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if idStr, ok := values["id"]; ok {
		id, err := strconv.Atoi(idStr[0])
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("wrong 'id' query parameter: should be integer"))
			return
		}

		var get model.Car
		srv.andintern.DB().Get(&get, "SELECT * FROM andintern.cars "+
			"WHERE id=$1", id)
		if get.Id == 0 {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("in database no car with such id: %d", id)))
			return
		}
		bytes, err := json.Marshal(get)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while marshaling json: %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
		return
	}

	if brands, ok := values["brand"]; ok {
		var get model.Car
		brand := brands[0]
		srv.andintern.DB().Get(&get, "SELECT * FROM andintern.cars "+
			"WHERE brand=$1", brand)
		if get.Id == 0 {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("in database no car with such brand: %s", brand)))
			return
		}
		bytes, err := json.Marshal(get)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while marshaling json: %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
		return
	}

	var get []model.Car
	srv.andintern.DB().Select(&get, "SELECT * FROM andintern.cars")
	bytes, err := json.Marshal(get)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error while marshaling json: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return

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

func (srv *Server) createCar(w http.ResponseWriter, r *http.Request) {
	car := model.Car{}
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("wrong body parameters"))
		return
	}

	_, err = srv.andintern.DB().Exec("insert into andintern.cars (brand, price)"+
		" values ($1, $2);", /*srv.cfg.DB.SchemaName,*/ car.Brand, car.Price)

	if err != nil && strings.Contains(err.Error(), "SQLSTATE 23505") {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("provided car brand already exists: '%s'", car.Brand)))
		return
	} else if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("error while inserting %v: %v", car, err.Error())))
		return
	}

	w.Write([]byte(fmt.Sprintf("created car with brand: %s, price: %d", car.Brand, car.Price)))
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

func (srv *Server) updateCar(w http.ResponseWriter, r *http.Request) {
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

	var res model.Car
	srv.andintern.DB().Get(&res, "SELECT * FROM andintern.cars "+
		"WHERE id=$1", id)
	if res.Id == 0 {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("in database no car with such id: %d", id)))
		return
	}

	car := model.Car{}
	err = json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("wrong body parameters"))
		return
	}

	if car.Brand != "" && car.Price != 0 {
		// update brand and price
		_, err := srv.andintern.DB().Exec("update andintern.cars "+
			"set brand = $1, price = $2 where id = $3;", car.Brand, car.Price, id)
		if err != nil && strings.Contains(err.Error(), "SQLSTATE 23505") {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("provided car brand already exists: '%s'", car.Brand)))
			return
		} else if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("error while updating car with id %d: %v", id, err.Error())))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("updated car with id: %d, brand: %s, price: %d", id, car.Brand, car.Price)))
	} else if car.Brand != "" {
		// update brand
		_, err = srv.andintern.DB().Exec("update andintern.cars "+
			"set brand = $1 where id = $2;", car.Brand, id)
		if err != nil && strings.Contains(err.Error(), "SQLSTATE 23505") {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("provided car brand already exists: '%s'", car.Brand)))
			return
		} else if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("error while updating car with id %d: %v", id, err.Error())))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("updated car with id: %d, brand: %s", id, car.Brand)))
	} else if car.Price != 0 {
		// update price
		_, err = srv.andintern.DB().Exec("update andintern.cars "+
			"set price = $1 where id = $2;", car.Price, id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("error while updating car with id %d: %v", id, err.Error())))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("updated car with id: %d, price: %d", id, car.Price)))
	}
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

func (srv *Server) deleteCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if idStr, ok := values["id"]; ok {
		id, err := strconv.Atoi(idStr[0])
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("wrong 'id' query parameter: should be integer"))
			return
		}

		var get model.Car
		srv.andintern.DB().Get(&get, "SELECT * FROM andintern.cars "+
			"WHERE id=$1", id)
		if get.Id == 0 {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("in database no car with such id: %d", id)))
			return
		}

		_, err = srv.andintern.DB().Exec("DELETE FROM andintern.cars "+
			"WHERE id = $1;", id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("error while deleting car with id %d: %v", id, err.Error())))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("successfully deleted car with id: %d", id)))
		return
	}

	if brands, ok := values["brand"]; ok {
		var get model.Car
		brand := brands[0]
		srv.andintern.DB().Get(&get, "SELECT * FROM andintern.cars "+
			"WHERE brand=$1", brand)
		if get.Brand == "" {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("in database no car with such brand: %s", brand)))
			return
		}

		_, err := srv.andintern.DB().Exec("DELETE FROM andintern.cars "+
			"WHERE brand = $1;", brand)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("error while deleting car with brand '%s': %v", brand, err.Error())))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("successfully deleted car with brand: %s", brand)))
		return
	}

	if priceLowerThanStr, ok := values["priceLowerThan"]; ok {
		priceLowerThan, err := strconv.Atoi(priceLowerThanStr[0])
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("wrong 'priceLowerThan' query parameter: should be integer"))
			return
		}

		var get model.Car
		err = srv.andintern.DB().Get(&get, "SELECT * FROM andintern.cars "+
			"WHERE price<$1 LIMIT 1", priceLowerThan)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("no cars in database with price lower than %d", priceLowerThan)))
			return
		}

		_, err = srv.andintern.DB().Exec("DELETE FROM andintern.cars "+
			"WHERE price<$1", priceLowerThan)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("error while deleting car with price lower than '%d': %v",
				priceLowerThan, err.Error())))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("successfully deleted cars with price lower than: %d", priceLowerThan)))
		return
	}

	w.WriteHeader(400)
	w.Write([]byte(fmt.Sprintf("none of required fields provided (id/brand/pricelowerthan")))
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
