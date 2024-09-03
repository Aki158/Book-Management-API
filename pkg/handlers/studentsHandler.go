package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Aki158/School-API/pkg/db"
	"github.com/Aki158/School-API/pkg/structs"
)

func StudentsHandler(mydb *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Content-Typeヘッダーをapplication/jsonに設定する
		w.Header().Set("Content-Type", "application/json")

		var response structs.Response
		// クエリパラメータを解析する
		query := r.URL.Query()

		facilitatorIdStr := query.Get("facilitator_id")
		facilitatorId, err := strconv.Atoi(facilitatorIdStr)
		// facilitator_idに数字以外の文字列または1未満の値が入っている場合は、400 Bad Request を返す
		if err != nil || facilitatorId < 1 {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Invalid facilitator_id parameter: '%s', returning 400 Bad Request", facilitatorIdStr)
			return
		}

		pageStr := query.Get("page")
		page := 1

		// pageに数字以外の文字列または1未満の値が入っている場合は、400 Bad Request を返す
		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("Invalid page parameter: '%s', returning 400 Bad Request", pageStr)
				return
			}
		}

		limitStr := query.Get("limit")
		limit := 5
	
		// limitに数字以外の文字列または1未満の値が入っている場合は、400 Bad Request を返す
		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit < 1 {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("Invalid limit parameter: '%s', returning 400 Bad Request", limitStr)
				return
			}
		}

		sort := query.Get("sort")

		// sortにidまたはname、loginId以外の値が入っている場合は、400 Bad Request を返す
		if sort == "" {
			sort = "id"
		} else if !(sort == "id" || sort == "name" || sort == "loginId") {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Invalid sort parameter: '%s', returning 400 Bad Request", sort)
			return
		}

		order := query.Get("order")

		// orderにascまたはdesc以外の値が入っている場合は、400 Bad Request を返す
		if order == "" {
			order = "asc"
		} else if !(order == "asc" || order == "desc") {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Invalid order parameter: '%s', returning 400 Bad Request", order)
			return
		}

		name_like := query.Get("name_like")
		loginId_like := query.Get("loginId_like")
		
		response.StudentArr = mydb.Read(facilitatorId, page, limit, sort, order, name_like, loginId_like)
		response.TotalCount = len(response.StudentArr)
		// リクエストに該当する生徒が存在しない場合は、404 Not Found を返す
		if response.TotalCount == 0 {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("Request failed: resource not found, returning 404 Not Found")
			return
		}

		// マップをJSONにエンコードしてレスポンスとして送信する
		json.NewEncoder(w).Encode(response)
	}
}
