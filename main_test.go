package main_test

import (
	"."
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(tableCalcFunction); err != nil {
		log.Fatal(err)
	}

	a.DB.Exec(tableCalcTrigger)
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS shoes (
    id                          SERIAL,
    name                        TEXT NOT NULL,
    "trueToSizeData"              int[] NOT NULL DEFAULT '{}',
    "trueToSizeCalculation"       numeric(14,13) NOT NULL DEFAULT 0.00,
    CONSTRAINT shoes_pkey PRIMARY KEY (id)
)`

const tableCalcFunction = `CREATE OR REPLACE FUNCTION CalculationTrueToSize()
RETURNS TRIGGER AS $$
DECLARE
    arrLen integer;
    arrSum integer;
BEGIN
    SELECT cardinality("trueToSizeData") INTO arrLen
    FROM shoes WHERE id=new.id;

    SELECT SUM(UNNEST(t)) INTO arrSum
    FROM (SELECT UNNEST("trueToSizeData") FROM shoes WHERE id=new.id) t;

    UPDATE shoes SET "trueToSizeCalculation" = (arrSum::NUMERIC/arrLen::NUMERIC) WHERE id=new.id;

    RETURN new;
END;
$$ LANGUAGE plpgsql;`

const tableCalcTrigger = `CREATE TRIGGER update_trueToSizeCalculation
AFTER UPDATE OF "trueToSizeData" ON shoes
FOR EACH ROW EXECUTE PROCEDURE CalculationTrueToSize();`

func clearTable() {
	a.DB.Exec("DELETE FROM shoes")
	a.DB.Exec("ALTER SEQUENCE shoes_id_seq RESTART WITH 1")
	a.DB.Exec("DROP FUNCTION IF EXISTS CalculationTrueToSize()")
	// a.DB.Exec("DROP TRIGGER IF EXISTS update_trueToSizeCalculation ON shoes")
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/shoes", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentShoe(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/shoes/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Shoe not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Shoe not found'. Got '%s'", m["error"])
	}
}

func TestCreateShoe(t *testing.T) {
	clearTable()

	payload := []byte(`{"shoe":"AJ 1 Mid Cool Blue","trueToSizeVal":5}`)

	req, _ := http.NewRequest("POST", "/shoes", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["shoe"] != "AJ 1 Mid Cool Blue" {
		t.Errorf("Expected shoe name to be 'AJ 1 Mid Cool Blue'. Got '%v'", m["shoe"])
	}

	if m["trueToSizeCalculation"] != 5.0 {
		t.Errorf("Expected shoe trueToSizeCalculation to be '5.0'. Got '%v'", m["trueToSizeCalculation"])
	}
}

func TestGetShoe(t *testing.T) {
	clearTable()
	addShoes(1)

	req, _ := http.NewRequest("GET", "/shoes/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteShoe(t *testing.T) {
	clearTable()
	addShoes(5)

	req, _ := http.NewRequest("DELETE", "/shoes/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var count int

	err := a.DB.QueryRow("SELECT COUNT(*) FROM shoes").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count != 4 {
		t.Errorf("Expected shoe record count to be '4'. Got '%v'", count)
	}
}

func addShoes(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO shoes(name, \"trueToSizeData\") VALUES($1, $2)", "Shoe "+strconv.Itoa(i), "{1,2,3,4,5}")
	}
}
