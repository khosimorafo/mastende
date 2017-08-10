package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/antonholmquist/jason"
)

var a MastendeQL
var tenantid string
var invoiceid string

func TestMain(m *testing.M) {

	a = MastendeQL{}

	a.Initialize()

	code := m.Run()

	os.Exit(code)
}

func TestCreateTenant(t *testing.T) {

	payload := fmt.Sprintf(`query=mutation+_{
		createTenant(name:"Khosi Morafo",zaid:"7704215267089", moveindate:"2017-08-01")
		{id,name,zaid}}`)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("POST", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	b, _ := json.Marshal(result)
	v, _ := jason.NewObjectFromBytes(b)
	data, _ := v.GetObject("data")
	tenant, _ := data.GetObject("createTenant")

	_id, _ := tenant.GetString("id")

	tenantid = _id

	t.Log(response.Body)

	t.Log(tenantid)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestCreateInvoice(t *testing.T) {

	payload := fmt.Sprintf(`query=mutation+_{
		createMonthlyInvoice(tenantid:"%s", date:"2017-08-01", duedate:"2017-08-05")
		{id, status, total, balance}}`, tenantid)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("POST", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	b, _ := json.Marshal(result)
	v, _ := jason.NewObjectFromBytes(b)
	data, _ := v.GetObject("data")
	tenant, _ := data.GetObject("createInvoice")

	_id, _ := tenant.GetString("id")

	invoiceid = _id

	t.Log(response.Body)

	t.Log(invoiceid)

	//tenantid = response_.Body.

	checkResponseCode(t, http.StatusOK, response.Code)
}

func checkResponseCode(t *testing.T, expected, actual int) {

	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
