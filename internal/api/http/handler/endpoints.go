package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kosdirus/andintern/internal/api"
	httpv1 "github.com/kosdirus/andintern/internal/api/http/handler/v1"
	"github.com/kosdirus/andintern/internal/database/dataprovider"
	"github.com/kosdirus/andintern/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (srv Server) getCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carFilter, err := api.ParseCarFilter(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	dbCars, err := srv.andintern.GetCars(ctx, carFilter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if len(dbCars) == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("no car in database with %v", carFilter))
		return
	}

	cars := make([]*httpv1.Car, 0, len(dbCars))
	for _, car := range dbCars {
		cars = append(cars, httpv1.ToCar(car))
	}

	api.RespondDataOK(ctx, w, api.RangeItemsResponse{
		Cars: cars,
	})
}

func (srv *Server) createCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carToCreate, err := api.ParseCreateCarRequest(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	} else if carToCreate.Brand == "" || carToCreate.Price == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("body parameters should not be empty"))
		return
	}

	createdCar, err := srv.andintern.CreateCar(ctx, carToCreate)
	if err != nil && strings.Contains(err.Error(), "SQLSTATE 23505") {
		api.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("provided car brand already exists: '%s'",
			carToCreate.Brand))
		return
	} else if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	api.RespondDataOK(ctx, w, api.RangeItemsResponse{Cars: createdCar})
}

func (srv *Server) updateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carToUpdate, err := api.ParseUpdateCarRequest(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	} else if carToUpdate.Price.Set && carToUpdate.Price.Value == 0 {
		api.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("price can't be zero"))
		return
	}

	if carToUpdate.Brand.Set {
		car, err := srv.andintern.GetCar(ctx, dataprovider.NewCarFilter().ByBrand(carToUpdate.Brand.Value))
		if err != nil {
			api.RespondError(ctx, w, http.StatusBadRequest, err)
			return
		} else if car != nil {
			api.RespondError(ctx, w, http.StatusBadRequest,
				fmt.Errorf("provided car brand already exists: '%s'", car.Brand))
			return
		}
	}

	car, err := srv.andintern.GetCar(ctx, dataprovider.NewCarFilter().ByID(carToUpdate.Id))
	if err == nil && car == nil {
		api.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("not found car with such id: '%d'",
			carToUpdate.Id))
		return
	}

	updatedCar, err := srv.andintern.UpdateCar(ctx, carToUpdate)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	api.RespondDataOK(ctx, w, updatedCar)
}

func (srv *Server) deleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carFilter, err := api.ParseCarFilter(r)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	err = srv.andintern.DeleteCar(ctx, carFilter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	api.RespondDataOK(ctx, w, fmt.Sprintf("successfully deleted car with query parameter: %s", carFilter))
}

//
//
//
func (srv Server) getCarArchive(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("in database no car with such id: %d", id)))
			return
		}
		bytes, err := json.Marshal(get)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error while marshaling json: %s", err.Error())))
			return
		}
		w.Header().Set("Content-Type", "application/json")
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
		w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return

}

func (srv *Server) updateCarArchive(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("not found car with such id: %d", id)))
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
			w.WriteHeader(500)
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

func (srv *Server) deleteCarArchive(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("in database no car with such id: %d", id)))
			return
		}

		_, err = srv.andintern.DB().Exec("DELETE FROM andintern.cars "+
			"WHERE id = $1;", id)
		if err != nil {
			w.WriteHeader(500)
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
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("in database no car with such brand: %s", brand)))
			return
		}

		_, err := srv.andintern.DB().Exec("DELETE FROM andintern.cars "+
			"WHERE brand = $1;", brand)
		if err != nil {
			w.WriteHeader(500)
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
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("no cars in database with price lower than %d", priceLowerThan)))
			return
		}

		_, err = srv.andintern.DB().Exec("DELETE FROM andintern.cars "+
			"WHERE price<$1", priceLowerThan)
		if err != nil {
			w.WriteHeader(500)
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
