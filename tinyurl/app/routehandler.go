package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"tinyurl/app/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) *appError {
	log.Info("Inside index")
	return indexTmpl.Execute(w, r, nil)
}

func HandleTinyUrlPost(w http.ResponseWriter, r *http.Request) *appError {

	formData, err := formHandler(r)

	if err != nil {
		return errorTmpl.Execute(w, r, model.ErrorMessage{Message: err.Error()})
	}
	log.Debug("Form data:", formData)

	cacheAvailable, cacheResult := getCachedURL(formData.LongURL)

	if cacheAvailable {
		return detailTmpl.Execute(w, r, model.TinyUrl{fmt.Sprintf("%s:%s/%s", appConfig.HostName, appConfig.Port, cacheResult.ShortString), formData.LongURL, cacheResult.CreatedAt})
	}

	shortString := GenerateShortString()
	log.Info("Generated new url", shortString)

	result, err := SelectTinyURL(formData.LongURL)

	if err != nil {
		log.Warn("Error while selecting tinyurl from db")
	} else {
		return detailTmpl.Execute(w, r, model.TinyUrl{fmt.Sprintf("%s:%s/%s", appConfig.HostName, appConfig.Port, result.ShortURL), result.LongURL, result.CreatedAt})
	}

	currentTime := time.Now()
	err = InsertTinyURLs(shortString, formData.LongURL, currentTime)
	if err != nil {
		log.Fatal(err.Error())
		return errorTmpl.Execute(w, r, nil)
	}

	setCacheURL(formData.LongURL, shortString, currentTime)
	http.Redirect(w, r, fmt.Sprintf("/tinyurl/%s/detail", shortString), http.StatusFound)
	return nil
}

func HandleDetail(w http.ResponseWriter, r *http.Request) *appError {

	var tinyURL model.TinyUrl
	shortURL := mux.Vars(r)["short_url"]
	result, err := SelectTinyURL(shortURL)

	if err != nil {
		log.Warn("Error while selecting tinyurl from db")
	}

	if result == tinyURL {
		log.Info("result:", result)
		return notFoundTmpl.Execute(w, r, nil)
	}

	return detailTmpl.Execute(w, r, model.TinyUrl{fmt.Sprintf("%s:%s/%s", appConfig.HostName, appConfig.Port, result.ShortURL), result.LongURL, result.CreatedAt})
}

func HandleTinyUrl(w http.ResponseWriter, r *http.Request) *appError {

	shortURL := mux.Vars(r)["short_url"]
	result, err := SelectTinyURL(shortURL)

	if err != nil {
		log.Warn("Error while selecting tinyurl from db")
	}

	if result.LongURL == "" {
		log.Info("result:", result)
		return notFoundTmpl.Execute(w, r, nil)
	}

	currentTime := time.Now()
	err = InsertTinyURLVisits(shortURL, currentTime)
	if err != nil {
		log.Fatal(err.Error())
		return errorTmpl.Execute(w, r, nil)
	}

	http.Redirect(w, r, result.LongURL, http.StatusFound)
	return nil

}

func HandleList(w http.ResponseWriter, r *http.Request) *appError {
	results, err := Conn.Query("SELECT * FROM tiny_urls")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer results.Close()

	type listRow struct {
		model.TinyUrl
		StatsURL string
	}

	//var tinyURLs []model.TinyUrl
	var listData []listRow
	for results.Next() {
		var tinyURL model.TinyUrl
		err = results.Scan(&tinyURL.ShortURL, &tinyURL.LongURL, &tinyURL.CreatedAt)
		if err != nil {
			log.Fatal(err.Error())
		}

		listData = append(listData, listRow{
			model.TinyUrl{
				fmt.Sprintf("%s:%s/%s", appConfig.HostName, appConfig.Port, tinyURL.ShortURL),
				tinyURL.LongURL,
				tinyURL.CreatedAt,
			},
			fmt.Sprintf("%s:%s/tinyurl/%s/stats", appConfig.HostName, appConfig.Port, tinyURL.ShortURL),
		})
	}

	return listTmpl.Execute(w, r, listData)
}

func HandleStats(w http.ResponseWriter, r *http.Request) *appError {
	shortURL := mux.Vars(r)["short_url"]
	count, err := CountTinyURLVisited(shortURL)

	if err != nil {
		log.Fatal(err.Error())
		notFoundTmpl.Execute(w, r, nil)
	}

	result, err := SelectTinyURL(shortURL)

	if err != nil {
		log.Warn("Error while selecting tinyurl from db")
	}

	statsData := struct {
		ShortURL, LongURL string
		Count             int
	}{
		ShortURL: shortURL,
		LongURL:  result.LongURL,
		Count:    count,
	}

	return statsTmpl.Execute(w, r, statsData)

}

func HandleRquestCount(w http.ResponseWriter, r *http.Request) *appError {
	shortURL := mux.Vars(r)["short_url"]

	requestsData := GetRequestCountPerDay(shortURL, 7)

	statsData := struct {
		RequestsData []model.RequestCount
	}{
		RequestsData: requestsData,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statsData)
	return nil
}

func formHandler(r *http.Request) (*model.Form, error) {

	long_url, err := url.ParseRequestURI(r.FormValue("long_url"))

	if err != nil {
		log.Warn("Invalid long_url")
		return nil, err
	}

	formData := &model.Form{
		LongURL: long_url.String(),
	}

	return formData, nil
}
