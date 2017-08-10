package main

import (
	"testing"
	"os"
	"net/http"
	"net/http/httptest"
	"fmt"
	"encoding/json"
	"github.com/antonholmquist/jason"
)

var a MastendeQL
var tenantid string
var invoiceid string

func TestMain(m *testing.M) {

	a = MastendeQL{}

	a.Initialize()

	//a.Run(":8080")

	code := m.Run()

	os.Exit(code)
}

func TestsRun(t *testing.T){

	//testCreateTenant(t)

}

func TestCreateTenant(t *testing.T) {

	payload_ := fmt.Sprintf(`query=mutation+_{
		createTenant(name:"Khosi Morafo",zaid:"7704215267089", moveindate:"2017-08-01")
		{id,name,zaid}}`)

	resource_ := fmt.Sprintf("/graphql?%v", payload_)

	//t.Log(payload_)

	//t.Log(resource_)

	req_, _ := http.NewRequest("POST", resource_, nil)

	response := executeRequest(req_)

	var result_ten map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result_ten)

	//fmt.Print(result_ten)

	b, _ := json.Marshal(result_ten)
	v, _ := jason.NewObjectFromBytes(b)
	data, _ := v.GetObject("data")
	tenant, _ := data.GetObject("createTenant")

	_id, _ := tenant.GetString("id")

	tenantid = _id

	t.Log(response.Body)

	t.Log(tenantid)

	//tenantid = response_.Body.

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestCreateInvoice(t *testing.T) {

	payload_ := fmt.Sprintf(`query=mutation+_{
		createInvoice(tenantid:"%s",date:"2017-08-05", duedate:"2017-08-05")
		{id, status, total, balance}}`, tenantid)

	resource_ := fmt.Sprintf("/graphql?%v", payload_)

	req_, _ := http.NewRequest("POST", resource_, nil)

	response := executeRequest(req_)

	var result_ten map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result_ten)

	//fmt.Print(result_ten)

	b, _ := json.Marshal(result_ten)
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

func TestAddInvoiceLineItem(t *testing.T) {

	payload_ := fmt.Sprintf(`query=mutation+_{
		createLineItem(invoiceid:"%s",date:"2017-08-05", duedate:"2017-08-05")
		{id, status, total, balance}}`, invoiceid)

	resource_ := fmt.Sprintf("/graphql?%v", payload_)

	req_, _ := http.NewRequest("POST", resource_, nil)

	response := executeRequest(req_)

	var result_ten map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result_ten)

	//fmt.Print(result_ten)

	b, _ := json.Marshal(result_ten)
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

