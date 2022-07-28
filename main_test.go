package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var a App

var (
	uuidTest string
	uuidDrawingTest string
)

func TestMain(m *testing.M) {
	a.Initialize("./app_test.toml")
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM canvas")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS canvas
(
  id varchar(64) NOT NULL,
  drawing varchar(1000) NOT NULL,
  creationdate datetime NOT NULL,
  PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/canvasList", nil)
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

func TestGetNonExistentCanvas(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/canvas/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Canvas not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Canvas not found'. Got '%s'", m["error"])
	}
}

func TestCreateCanvas(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"drawing": "api testing create"}`)
	req, _ := http.NewRequest("POST", "/canvas", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	uuidTest = m["id"].(string)

	if m["drawing"] != "api testing create" {
		t.Errorf("Expected canvas drawing to be 'api testing create'. Got '%v'", m["drawing"])
	}

}

func TestGetCanvas(t *testing.T) {
	clearTable()
	addCanvas(1)

	req, _ := http.NewRequest("GET", "/canvas/"+uuidTest, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCanvas(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO canvas(id,drawing, creationdate) VALUES(?, ?,?)", uuidTest,"Inserting Drawing Test ", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func TestUpdateCanvas(t *testing.T) {

	clearTable()
	addCanvas(1)

	req, _ := http.NewRequest("GET", "/canvas/"+uuidTest, nil)
	response := executeRequest(req)
	var originalCanvas map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCanvas)

	var jsonStr = []byte(`{"drawing":"drawing updated canvas"}`)
	req, _ = http.NewRequest("PUT", "/canvas/"+uuidTest, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalCanvas["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalCanvas["id"], m["id"])
	}

	if m["drawing"] == originalCanvas["drawing"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalCanvas["drawing"], m["drawing"], m["drawing"])
	}
}

func TestDeleteCanvas(t *testing.T) {
	clearTable()
	addCanvas(1)

	req, _ := http.NewRequest("GET", "/canvas/"+uuidTest, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/canvas/"+uuidTest, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/canvas/"+uuidTest, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestCreateCanvasRequest(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`[
		{
		   "RectangleAt" : [14,0],
		   "Width": 7,
		   "Height":6,
		   "Outline": "none",
		   "Fill": "."
		}, 
		{
		   "RectangleAt" : [0,3],
		   "Width": 8,
		   "Height":4,
		   "Outline": "o",
		   "Fill": "none"
		 }
	]`)
	req, _ := http.NewRequest("POST", "/canvasCreateRequest", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var (
		m map[string]interface{}
		returnSplit []string
	)
	json.Unmarshal(response.Body.Bytes(), &m)

	uuidDrawingTest = m["id"].(string)
	drawing:="              .......\n              .......\n              .......\noooooooo      .......\no      o      .......\no      o      .......\noooooooo             "
	drawingSplit :=strings.Split(drawing, "\n")

	for _, v := range m["drawing"].([]interface{}) {
		returnSplit = append(returnSplit, v.(string))
	}
	if  !reflect.DeepEqual(returnSplit,drawingSplit) {
		t.Errorf("Expected canvas drawing to be '%v'. Got '%v'", drawingSplit,returnSplit)
	}

}

func TestGetCanvasResponse(t *testing.T) {

	req, _ := http.NewRequest("GET", "/canvasResponse/"+uuidDrawingTest, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}