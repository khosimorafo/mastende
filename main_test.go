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
		createTenant(name:"Khosi Morafo Http",zaid:"7704215267089", moveindate:"2017-08-01")
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

func TestUpdateTenant(t *testing.T) {

	payload := fmt.Sprintf(`query=mutation+_{
		updateTenant(id:"%s", name:"Khosi Morafo Updated",zaid:"8004215267089", site:"Mganka")
		{id, name, zaid, site}}`, tenantid)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("PUT", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log(response.Body)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestTenantList(t *testing.T) {

	payload := fmt.Sprintf(`query={tenantList {id,name,zaid} }`)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	//b, _ := json.Marshal(result)
	//v, _ := jason.NewObjectFromBytes(b)
	//data, _ := v.GetObject("data")
	//tenant, _ := data.GetObject("singleTenant")

	t.Log("--------------------------------------------------------")
	t.Log("Tenant list is : ", response.Body)
	t.Log("--------------------------------------------------------")

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
	tenant, _ := data.GetObject("createMonthlyInvoice")

	_id, _ := tenant.GetString("id")

	invoiceid = _id

	t.Log("--------------------------------------------------------")
	t.Log(response.Body)
	t.Log("--------------------------------------------------------")


	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreatePaymentExtension(t *testing.T){

	payload := fmt.Sprintf(`query=mutation+_{
		paymentExtension(invoiceid: "%s", paybydate: "%s", requestdate:"%s", requestby:"%s", requestmode:"%s")
		{result, message}}`, invoiceid, "2017-08-14", "2017-08-14", "Khosi", "Test")

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("POST", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log("--------------------------------------------------------")
	t.Log("Payment extension result : ", response.Body)
	t.Log("--------------------------------------------------------")


	checkResponseCode(t, http.StatusOK, response.Code)
}

/* Test before payment*/
func TestOutstandingInvoiceList(t *testing.T){

	payload := fmt.Sprintf(`query={outstandingInvoiceList
	 (periodname:"%s")
	 {id, periodname, total, balance, status} }`, "August-2017")

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log("--------------------------------------------------------")
	t.Log("Outstanding Invoice list is : ", response.Body)
	t.Log("--------------------------------------------------------")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateInvoicePayment(t *testing.T) {

	payload := fmt.Sprintf(`query=mutation+_{
		makePayment(invoiceid: "%s", tenantid:"%s", description:"Httpayment", date:"2017-08-10", amount:330.0, mode: "Cash")
		{result, message}}`, invoiceid, tenantid)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("POST", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log("--------------------------------------------------------")
	t.Log("Invoice payment result : ", response.Body)
	t.Log("--------------------------------------------------------")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestTenantInvoiceList(t *testing.T){

	payload := fmt.Sprintf(`query={tenantInvoiceList
	 (tenantid: "%s")
	 {id, periodname, total, balance, status} }`, tenantid)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log("--------------------------------------------------------")
	t.Log("Tenant Invoice list is : ", response.Body)
	t.Log("--------------------------------------------------------")

	checkResponseCode(t, http.StatusOK, response.Code)
}

/* Test after payment*/
func TestOutstandingInvoiceList1(t *testing.T){

	payload := fmt.Sprintf(`query={outstandingInvoiceList
	 (periodname:"%s")
	 {id, periodname, total, balance, status} }`, "August-2017")

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log("--------------------------------------------------------")
	t.Log("Outstanding Invoice list is : ", response.Body)
	t.Log("--------------------------------------------------------")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestTenantPaymentList(t *testing.T){

	payload := fmt.Sprintf(`query={tenantPaymentList
	 (tenantid: "%s")
	 {id, tenantid, invoiceid, amount, date} }`, tenantid)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log("--------------------------------------------------------")
	t.Log("Tenant Payment list is : ", response.Body)
	t.Log("--------------------------------------------------------")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestReadTenant(t *testing.T) {

	payload := fmt.Sprintf(`query={
		singleTenant(id:"%s")
		{	id, name, zaid, site, room, mobile, outstanding, overdue
				invoices{ id, balance, total, duedate, date, status }
				payments{ id, amount, date, mode, description }

		}}`, tenantid)

	resource := fmt.Sprintf("/graphql?%v", payload)

	request, _ := http.NewRequest("GET", resource, nil)

	response := executeRequest(request)

	var result map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &result)

	t.Log("Retrieved Tenant is : ", response.Body)

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
