package main

import (
	"bytes"
	//"database/sql"
	"encoding/json"
	//"fmt"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/http/httptest"
	"os"
	//"reflect"
	"testing"
	"time"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Router = mux.NewRouter()
	code := m.Run()
	os.Exit(code)
}

func (a *App) assertJSON(actual []byte, data interface{}, t *testing.T) {
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("the expected json: %s is different from actual %s", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		//Println("Expected response code:", expected, ". Got:", actual)
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func Test_Budget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	app := &App{DB: db}
	app.Router = a.Router

	a.Router.HandleFunc("/budgetcategories", app.budgetCategory_createNew).Methods("POST")
	a.Router.HandleFunc("/budgetcategories/{id:[0-9]+}", app.budgetCategory_getOne).Methods("GET")
	a.Router.HandleFunc("/budgetcategories", app.budgetCategory_getAll).Methods("GET")
	a.Router.HandleFunc("/budgetcategories/{id:[0-9]+}", app.budgetCategory_update).Methods("PUT")
	a.Router.HandleFunc("/budgetcategories/{id:[0-9]+}", app.budgetCategory_delete).Methods("DELETE")

	a.Router.HandleFunc("/lawfirmactions", app.lawFirmActions_createNew).Methods("POST")
	a.Router.HandleFunc("/lawfirmactions/{id:[0-9]+}", app.lawFirmActions_getOne).Methods("GET")
	a.Router.HandleFunc("/lawfirmactions", app.lawFirmActions_getAll).Methods("GET")
	a.Router.HandleFunc("/lawfirmactions/{id:[0-9]+}", app.lawFirmActions_update).Methods("PUT")
	a.Router.HandleFunc("/lawfirmactions/{id:[0-9]+}", app.lawFirmActions_delete).Methods("DELETE")

	a.Router.HandleFunc("/lawfirmcost", app.lawFirmCost_createNew).Methods("POST")
	a.Router.HandleFunc("/lawfirmcost/{id:[0-9]+}", app.lawFirmCost_getOne).Methods("GET")
	a.Router.HandleFunc("/lawfirmcost", app.lawFirmCost_getAll).Methods("GET")
	a.Router.HandleFunc("/lawfirmcost/{id:[0-9]+}", app.lawFirmCost_update).Methods("PUT")
	a.Router.HandleFunc("/lawfirmcost/{id:[0-9]+}", app.lawFirmCost_delete).Methods("DELETE")

	a.Router.HandleFunc("/lawfirmcostusptofees", app.lawFirmCostUsptoFees_createNew).Methods("POST")
	a.Router.HandleFunc("/lawfirmcostusptofees/{id:[0-9]+}", app.lawFirmCostUsptoFees_getOne).Methods("GET")
	a.Router.HandleFunc("/lawfirmcostusptofees", app.lawFirmCostUsptoFees_getAll).Methods("GET")
	a.Router.HandleFunc("/lawfirmcostusptofees/{id:[0-9]+}", app.lawFirmCostUsptoFees_update).Methods("PUT")
	a.Router.HandleFunc("/lawfirmcostusptofees/{id:[0-9]+}", app.lawFirmCostUsptoFees_delete).Methods("DELETE")

	Convey("Test_Law_Firm_Actions", t, func() {

		Convey("Law_Firm_Actions_Create_New", func() {
			payload := []byte(`{"user_company_id": 3, "name":"TestName"}`)
			req, err := http.NewRequest("POST", "/lawfirmactions", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			mock.ExpectExec("INSERT INTO law_firm_action_list").
				WithArgs([]byte("3"), "TestName").
				WillReturnResult(sqlmock.NewResult(1, 1))

			response := executeRequest(req)
			checkResponseCode(t, http.StatusCreated, response.Code)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Actions_Get_One", func() {
			req, err := http.NewRequest("GET", "/lawfirmactions/1", nil)
			So(err, ShouldBeNil)

			rows := sqlmock.NewRows([]string{"user_company_id", "name"}).AddRow(1, "TestName")

			mock.ExpectQuery("SELECT user_company_id, name FROM law_firm_action_list WHERE (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			data := lawFirmActions{ID: 1, UserCompanyId: NullInt64{Int64: 1, Valid: true}, Name: "TestName"}

			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Actions_Get_All", func() {
			req, err := http.NewRequest("GET", "/lawfirmactions", nil)
			So(err, ShouldBeNil)

			rows := sqlmock.NewRows([]string{"id", "user_company_id", "name"}).
				AddRow(1, 1, "Category1").
				AddRow(2, nil, "Category2")

			mock.ExpectQuery("SELECT id, user_company_id, name FROM law_firm_action_list (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			data := []budgetCategory{
				budgetCategory{ID: 1, UserCompanyId: NullInt64{Int64: 1, Valid: true}, Name: "Category1"},
				budgetCategory{ID: 2, UserCompanyId: NullInt64{Int64: 0, Valid: false}, Name: "Category2"},
			}

			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Actions_Update", func() {
			lfa := lawFirmActions{}
			lfa.ID = 1
			lfa.Name = "New_TestName"

			mock.ExpectExec("UPDATE law_firm_action_list").
				WithArgs([]byte("2"), lfa.Name, lfa.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			payload := []byte(`{"user_company_id":2, "name":"New_TestName"}`)
			req, err := http.NewRequest("PUT", "/lawfirmactions/1", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := lawFirmActions{ID: 1, UserCompanyId: NullInt64{Int64: 2, Valid: true}, Name: "New_TestName"}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Actions_Delete", func() {
			req, err := http.NewRequest("DELETE", "/lawfirmactions/1", nil)
			So(err, ShouldBeNil)

			mock.ExpectExec("DELETE FROM law_firm_action_list").
				WillReturnResult(sqlmock.NewResult(1, 1))
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := map[string]string{"result": "success"}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})
	})
	//-------------------------------------------------------------------------------------------------
	Convey("Test_Budget_Categories", t, func() {

		Convey("Budget_Categories_Create_New", func() {
			payload := []byte(`{"user_company_id": 3, "name":"TestName"}`)
			req, err := http.NewRequest("POST", "/budgetcategories", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			mock.ExpectExec("INSERT INTO budget_categories_list").
				WithArgs([]byte("3"), "TestName").
				WillReturnResult(sqlmock.NewResult(1, 1))

			response := executeRequest(req)
			checkResponseCode(t, http.StatusCreated, response.Code)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Budget_Categories_Get_One", func() {
			req, err := http.NewRequest("GET", "/budgetcategories/1", nil)
			So(err, ShouldBeNil)

			rows := sqlmock.NewRows([]string{"user_company_id", "name"}).AddRow(1, "hello")

			mock.ExpectQuery("SELECT user_company_id, name FROM budget_categories_list WHERE (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			data := budgetCategory{ID: 1, UserCompanyId: NullInt64{Int64: 1, Valid: true}, Name: "hello"}

			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Budget_Categories_Get_All", func() {
			req, err := http.NewRequest("GET", "/budgetcategories", nil)
			So(err, ShouldBeNil)

			rows := sqlmock.NewRows([]string{"id", "user_company_id", "name"}).
				AddRow(1, 1, "Category1").
				AddRow(2, nil, "Category2")

			mock.ExpectQuery("SELECT id, user_company_id, name FROM budget_categories_list (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			data := []budgetCategory{
				budgetCategory{ID: 1, UserCompanyId: NullInt64{Int64: 1, Valid: true}, Name: "Category1"},
				budgetCategory{ID: 2, UserCompanyId: NullInt64{Int64: 0, Valid: false}, Name: "Category2"},
			}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Budget_Categories_Update", func() {
			bc := budgetCategory{}
			bc.ID = 1
			bc.Name = "New_TestName"

			mock.ExpectExec("UPDATE budget_categories_list").
				WithArgs([]byte("2"), bc.Name, bc.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			payload := []byte(`{"user_company_id":2, "name":"New_TestName"}`)
			req, err := http.NewRequest("PUT", "/budgetcategories/1", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := budgetCategory{ID: 1, UserCompanyId: NullInt64{Int64: 2, Valid: true}, Name: "New_TestName"}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Budget_Categories_Delete", func() {
			req, err := http.NewRequest("DELETE", "/budgetcategories/1", nil)
			So(err, ShouldBeNil)

			mock.ExpectExec("DELETE FROM budget_categories_list").
				WillReturnResult(sqlmock.NewResult(1, 1))
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := map[string]string{"result": "success"}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})
	})

	Convey("Test_Law_Firm_Cost", t, func() {

		Convey("Law_Firm_Cost_Create_New", func() {
			payload := []byte(`{"user_id":2, "user_company_id": 3, "customer_number_id":4, "valid_from":"2018-06-12T21:31:15Z"}`)
			req, err := http.NewRequest("POST", "/lawfirmcost", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			mock.ExpectExec("INSERT INTO law_firm_cost").
				WithArgs([]byte("2"), []byte("3"), []byte("4"), "2018-06-12T21:31:15Z").
				WillReturnResult(sqlmock.NewResult(1, 1))

			response := executeRequest(req)
			checkResponseCode(t, http.StatusCreated, response.Code)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Cost_Get_One", func() {
			req, err := http.NewRequest("GET", "/lawfirmcost/1", nil)
			So(err, ShouldBeNil)
			myTime, err := time.Parse("2006-01-02T15:04:05Z", "2018-06-12T21:31:15Z")
			rows := sqlmock.NewRows([]string{"user_id", "user_company_id", "customer_number_id", "valid_from"}).AddRow(2, 3, 4, myTime)

			mock.ExpectQuery("SELECT user_id, user_company_id, customer_number_id, valid_from FROM law_firm_cost WHERE (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := lawFirmCost{ID: 1, UserId: NullInt64{Int64: 2, Valid: true}, UserCompanyId: NullInt64{Int64: 3, Valid: true},
				CustomerNumberId: NullInt64{Int64: 4, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}}

			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Cost_Get_All", func() {
			req, err := http.NewRequest("GET", "/lawfirmcost", nil)
			So(err, ShouldBeNil)
			myTime, err := time.Parse("2006-01-02T15:04:05Z", "2018-06-12T21:31:15Z")
			rows := sqlmock.NewRows([]string{"id", "user_id", "user_company_id", "customer_number_id", "valid_from"}).
				AddRow(1, 2, 3, 4, myTime).
				AddRow(2, 5, 6, 7, myTime)

			mock.ExpectQuery("SELECT id, user_id, user_company_id, customer_number_id, valid_from FROM law_firm_cost (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			data := []lawFirmCost{
				lawFirmCost{ID: 1, UserId: NullInt64{Int64: 2, Valid: true}, UserCompanyId: NullInt64{Int64: 3, Valid: true},
					CustomerNumberId: NullInt64{Int64: 4, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}},
				lawFirmCost{ID: 2, UserId: NullInt64{Int64: 5, Valid: true}, UserCompanyId: NullInt64{Int64: 6, Valid: true},
					CustomerNumberId: NullInt64{Int64: 7, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}},
			}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Cost_Update", func() {
			lfc := lawFirmCost{}
			lfc.ID = 1
			myTime, err := time.Parse("2006-01-02T15:04:05Z", "2018-06-12T21:31:15Z")
			mock.ExpectExec("UPDATE law_firm_cost").
				WithArgs([]byte("5"), []byte("6"), []byte("7"), "2018-06-12T21:31:15Z", lfc.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			payload := []byte(`{"user_id":5, "user_company_id":6, "customer_number_id":7, "valid_from":"2018-06-12T21:31:15Z"}`)
			req, err := http.NewRequest("PUT", "/lawfirmcost/1", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := lawFirmCost{ID: 1, UserId: NullInt64{Int64: 5, Valid: true}, UserCompanyId: NullInt64{Int64: 6, Valid: true},
				CustomerNumberId: NullInt64{Int64: 7, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_Cost_Delete", func() {
			req, err := http.NewRequest("DELETE", "/lawfirmcost/1", nil)
			So(err, ShouldBeNil)

			mock.ExpectExec("DELETE FROM law_firm_cost").
				WillReturnResult(sqlmock.NewResult(1, 1))
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := map[string]string{"result": "success"}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})
	})
	//-------------------------------------------------------------------------------------------------

	Convey("Test_Law_Firm_CostUsptoFees", t, func() {

		Convey("Law_Firm_CostUsptoFees_Create_New", func() {
			payload := []byte(`{"user_id":2, "user_company_id": 3, "valid_from":"2018-06-12T21:31:15Z"}`)
			req, err := http.NewRequest("POST", "/lawfirmcostusptofees", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			mock.ExpectExec("INSERT INTO law_firm_cost_uspto_fees").
				WithArgs([]byte("2"), []byte("3"), "2018-06-12T21:31:15Z").
				WillReturnResult(sqlmock.NewResult(1, 1))

			response := executeRequest(req)
			checkResponseCode(t, http.StatusCreated, response.Code)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_CostUsptoFees_Get_One", func() {
			req, err := http.NewRequest("GET", "/lawfirmcostusptofees/1", nil)
			So(err, ShouldBeNil)
			myTime, err := time.Parse("2006-01-02T15:04:05Z", "2018-06-12T21:31:15Z")

			rows := sqlmock.NewRows([]string{"user_id", "user_company_id", "valid_from"}).AddRow(2, 3, myTime)
			mock.ExpectQuery("SELECT user_id, user_company_id, valid_from FROM law_firm_cost_uspto_fees WHERE (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			data := lawFirmCostUsptoFees{ID: 1, UserId: NullInt64{Int64: 2, Valid: true},
				UserCompanyId: NullInt64{Int64: 3, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}}

			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_CostUsptoFees_Get_All", func() {
			req, err := http.NewRequest("GET", "/lawfirmcostusptofees", nil)
			So(err, ShouldBeNil)
			myTime, err := time.Parse("2006-01-02T15:04:05Z", "2018-06-12T21:31:15Z")

			rows := sqlmock.NewRows([]string{"id", "user_id", "user_company_id", "valid_from"}).
				AddRow(1, 3, 4, myTime).
				AddRow(2, 5, 6, myTime)

			mock.ExpectQuery("SELECT id, user_id, user_company_id, valid_from FROM law_firm_cost_uspto_fees (.+)").WillReturnRows(rows)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			data := []lawFirmCostUsptoFees{
				lawFirmCostUsptoFees{ID: 1, UserId: NullInt64{Int64: 3, Valid: true},
					UserCompanyId: NullInt64{Int64: 4, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}},
				lawFirmCostUsptoFees{ID: 2, UserId: NullInt64{Int64: 5, Valid: true},
					UserCompanyId: NullInt64{Int64: 6, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}},
			}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_CostUsptoFees_Update", func() {
			lfcuf := lawFirmCostUsptoFees{}
			lfcuf.ID = 1
			myTime, err := time.Parse("2006-01-02T15:04:05Z", "2018-06-12T21:31:15Z")

			mock.ExpectExec("UPDATE law_firm_cost_uspto_fees").
				WithArgs([]byte("2"), []byte("3"), "2018-06-12T21:31:15Z", lfcuf.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			payload := []byte(`{"user_id":2, "user_company_id":3, "valid_from":"2018-06-12T21:31:15Z"}`)
			req, err := http.NewRequest("PUT", "/lawfirmcostusptofees/1", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)

			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := lawFirmCostUsptoFees{ID: 1, UserId: NullInt64{Int64: 2, Valid: true},
				UserCompanyId: NullInt64{Int64: 3, Valid: true}, ValidFrom: NullTime{Time: myTime, Valid: true}}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

		Convey("Law_Firm_CostUsptoFees_Delete", func() {
			req, err := http.NewRequest("DELETE", "/lawfirmcostusptofees/1", nil)
			So(err, ShouldBeNil)

			mock.ExpectExec("DELETE FROM law_firm_cost_uspto_fees").
				WillReturnResult(sqlmock.NewResult(1, 1))
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			data := map[string]string{"result": "success"}
			app.assertJSON(response.Body.Bytes(), data, t)

			err = mock.ExpectationsWereMet()
			So(err, ShouldBeNil)
		})

	})

}
