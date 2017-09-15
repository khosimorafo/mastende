package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/configor"
	"github.com/khosimorafo/mastende/db"
	"github.com/khosimorafo/mastende/gql"
	"github.com/khosimorafo/mastende/mastende"
	"github.com/khosimorafo/mastende/tenants"
	"github.com/pkg/errors"
	"github.com/khosimorafo/mastende/invoices"
	"github.com/khosimorafo/mastende/payments"
)

var TenantList []tenants.Tenant
var app *db.App
var m *mastende.Mastende

type MastendeQL struct {
	Router *mux.Router
}

func (a *MastendeQL) Initialize() {

	rand.Seed(time.Now().UnixNano())

	configor.Load(&db.Config, "config.yml")
	app = db.DB(db.Config.DB.DbUrl, db.Config.DB.DbName)
	app.Database.DropDatabase()

	m = mastende.New(app)

	a.Router = mux.NewRouter()
	a.Routes()
}

func (a *MastendeQL) Run(port string) {

	log.Fatal(http.ListenAndServe(port, a.Router))
}

// define custom GraphQL ObjectType `todoType` for our Golang struct `Todo`
// Note that
// - the fields in our todoType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var tenantType = graphql.NewObject(gql.TenantTypeConfig())
var invoiceType = graphql.NewObject(gql.InvoiceTypeConfig())
var paymentType = graphql.NewObject(gql.PaymentTypeConfig())

var resultType = graphql.NewObject(gql.ResultTypeConfig())

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name:   "RootMutation",
	Fields: mutationFields,
})

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func (a *MastendeQL) Routes() {

	a.Router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	//http.Handle("/", a.Router)
}

func Main() {

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	// Display some basic instructions
	fmt.Println("Now server is running on port 8080")
	/*
		fmt.Println("Get single todo: curl -g 'http://localhost:8080/graphql?query={todo(id:\"b\"){id,text,done}}'")
		fmt.Println("Persist new todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:\"My+new+todo\"){id,text,done}}'")
		fmt.Println("Update todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:\"a\",done:true){id,text,done}}'")
		fmt.Println("Load todo list: curl -g 'http://localhost:8080/graphql?query={todoList{id,text,done}}'")
	*/

	fmt.Println("Get single tenants: curl -g 'http://localhost:8080/graphql?query={tenants(id:\"5TXvJj6VlpRbpThYfMhmBPq2k\"){id,name,zaid}}'")
	fmt.Println("Load tenants list: curl -g 'http://localhost:8080/graphql?query={tenantList{_id, name,zaid}}'")
	fmt.Println("Persist new Tenant: curl -g 'http://localhost:8080/graphql?query=mutation+_{createTenant(name:\"Khosi+Morafo\", zaid:\"7704215267089\"){id,name,zaid}}'")

	fmt.Println("Access the web app via browser at 'http://localhost:8080'")

	http.ListenAndServe(":8080", nil)
}



// root query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/graphql?query={lastTodo{id,text,done}}'
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		/*
		   curl -g 'http://localhost:8080/graphql?query={todo(id:"b"){id,text,done}}'
		*/
		"singleTenant": &graphql.Field{
			Type:        tenantType,
			Description: "Get single tenant",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				_, isOK := params.Args["id"].(string)
				if !isOK {

					return tenants.Tenant{}, errors.New("id string invalid")
				}

				if err := m.GetTenant(params.Args); err != nil {

					return nil, err
				} else {

					// return the new Tenant object that saved to DB
					// Note here that
					// - we are returning a `Tenant` struct instance here
					// - we previously specified the return Type to be `tenantType`
					return m.Tenant, nil
				}
			},
		},
		/*
		   curl -g 'http://localhost:8080/graphql?query={tenantList{_id, name,zaid}}'
		*/
		"tenantList": &graphql.Field{
			Type:        graphql.NewList(tenantType),
			Description: "List of tenants",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {



				if err := m.TenantList(params.Args); err != nil {

					return nil, err
				}
				return m.Tenants, nil
			},
		},

		/*
		   curl -g 'http://localhost:8080/graphql?query={tenantList{_id, name,zaid}}'
		*/
		"tenantInvoiceList": &graphql.Field{
			Type:        graphql.NewList(invoiceType),
			Description: "List of invoices",
			Args: graphql.FieldConfigArgument{
				"tenantid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				_, isOK := params.Args["tenantid"].(string)
				if !isOK {

					return []invoices.Invoice{}, errors.New("id string invalid")
				}

				if err := m.InvoiceListByTenant(params.Args); err != nil {

					return nil, err
				}
				return m.Invoices, nil
			},
		},

		/*
		   curl -g 'http://localhost:8080/graphql?query={tenantList{_id, name,zaid}}'
		*/
		"outstandingInvoiceList": &graphql.Field{
			Type:        graphql.NewList(invoiceType),
			Description: "List of outstanding invoices",
			Args: graphql.FieldConfigArgument{
				"periodname": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				_, isOK := params.Args["periodname"].(string)
				if !isOK {

					return []invoices.Invoice{}, errors.New("query string invalid")
				}

				if err := m.OustandingInvoiceList(params.Args); err != nil {

					return nil, err
				}
				return m.Invoices, nil
			},
		},

		/*
		   curl -g 'http://localhost:8080/graphql?query={tenantPaymentList{tenantid}}'
		*/
		"tenantPaymentList": &graphql.Field{
			Type:        graphql.NewList(paymentType),
			Description: "List of tenant payments",
			Args: graphql.FieldConfigArgument{
				"tenantid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				_, isOK := params.Args["tenantid"].(string)
				if !isOK {

					return []payments.Payment{}, errors.New("id string invalid")
				}

				if err := m.PaymentListByTenant(params.Args); err != nil {

					return nil, err
				}
				return m.Payments, nil
			},
		},
	},
})


var mutationFields = graphql.Fields{
	/*
		curl -g 'http://localhost:8080/graphql?query=mutation+
		_{createTenant(name:"", zaid:"", moveindate:""){id,name,zaid}}'

	*/
	"createTenant": &graphql.Field{
		Type:        tenantType, // the return type for this field
		Description: "Persist new tenants",
		Args:        gql.TenantFieldArguments(),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			// perform mutation operation here
			// for e.g. create a Tenant and save to DB.
			if err := m.CreateTenant(params.Args, false); err != nil {

				return nil, err
			} else {

				// return the new Tenant object that saved to DB
				// Note here that
				// - we are returning a `Tenant` struct instance here
				// - we previously specified the return Type to be `tenantType`
				return m.Tenant, nil
			}

		},
	},

	/*
		curl -g 'http://localhost:8080/graphql?query=mutation+
		_{updateTenant(id:"a",name:true){id,name,zaid}}'
	*/
	"updateTenant": &graphql.Field{
		Type:        tenantType, // the return type for this field
		Description: "Update existing tenant record",
		Args:        gql.TenantFieldArguments(),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			// perform mutation operation here
			// for e.g. create a Tenant and save to DB.
			if err := m.UpdateTenant(params.Args); err != nil {

				return nil, err
			} else {

				// return the new Tenant object that saved to DB
				// Note here that
				// - we are returning a `Tenant` struct instance here
				// - we previously specified the return Type to be `tenantType`
				return m.Tenant, nil
			}

		},
	},

	/*
		curl -g 'http://localhost:8080/graphql?query=mutation+
		_{setTenantStatus(id:"a",status:Active){id,name,zaid}}'
	*/

	"setTenantStatus": &graphql.Field{
		Type:        tenantType, // the return type for this field
		Description: "Set tenants status",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"status": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},

		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// marshall and cast the argument value
			//id, _ := params.Args["id"].(string)
			//status, _ := params.Args["status"].(string)

			affectedTenant := tenants.Tenant{}

			// Return affected todo
			return affectedTenant, nil
		},
	},

	/*
		curl -g 'http://localhost:8080/graphql?query=mutation+
			_{createInvoice(tenantid:"", lineitems:{}, date:"", duedate:""){id,total}}'
	*/
	"createInvoice": &graphql.Field{
		Type:        invoiceType, // the return type for this field
		Description: "Persist new invoice",
		Args:        gql.InvoiceFieldArguments(),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			// perform mutation operation here
			// for e.g. create a Invoice and save to DB.
			if err := m.CreateInvoice(params.Args); err != nil {

				return nil, err
			} else {

				// return the new Tenant object that saved to DB
				// Note here that
				// - we are returning a `Tenant` struct instance here
				// - we previously specified the return Type to be `tenantType`
				return m.Invoice, nil
			}

		},
	},

	/*
		curl -g 'http://localhost:8080/graphql?query=mutation+
			_{createInvoice(tenantid:"", lineitems:{}, date:"", duedate:""){id,total}}'
	*/
	"createMonthlyInvoice": &graphql.Field{
		Type:        invoiceType, // the return type for this field
		Description: "Create new monthly tenant invoice",
		Args:        gql.InvoiceFieldArguments(),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			// perform mutation operation here
			// for e.g. create a Invoice and save to DB.

			input := map[string]interface{}{

				"tenantid": params.Args["tenantid"].(string),
				"date":  params.Args["date"].(string),
				"daydueby":  db.Config.Invoices.NettDue,
				"lastdiscountday":  db.Config.Invoices.LastDiscountDay,
			}

			if err := m.CreateMonthlyTenantInvoice(input); err != nil {

				return nil, err
			} else {

				// return the new Tenant object that saved to DB
				// Note here that
				// - we are returning a `Tenant` struct instance here
				// - we previously specified the return Type to be `tenantType`
				return m.Invoice, nil
			}

		},
	},

	/*
		curl -g 'http://localhost:8080/graphql?query=mutation+
			_{makePayment(invoiceid:"", tenantid:"", amount:"", date:"", mode:""){id}}'
	*/
	"makePayment": &graphql.Field{
		Type:        resultType, // the return type for this field
		Description: "Make payment on invoice",
		Args:        gql.PaymentFieldArguments(),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			//log.Println("Params is ", params.Args)

			// perform mutation operation here
			// for e.g. create a payment and save to DB.
			if err := m.MakeInvoicePayment(params.Args); err != nil {

				log.Println("Error is ", err)
				return nil, err
			} else {

				// Note here that
				// - we are returning a `Result` struct instance here
				// - we previously specified the return Type to be `paymentType`
				return m.Result, nil
			}
		},
	},

	/*
		curl -g 'http://localhost:8080/graphql?query=mutation+
			_{paymentExtensionRequest(invoiceid:"", datetopay:""){result, message}}'
	*/
	"paymentExtension": &graphql.Field{
		Type:        resultType, // the return type for this field
		Description: "Make payment on invoice",
		Args:        gql.PaymentExtensionFieldArguments(),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			// perform mutation operation here
			if err := m.PaymentExtensionRequest(params.Args); err != nil {

				log.Println("Error is ", err)
				return nil, err
			} else {

				// return the new Payment object that saved to DB
				// Note here that
				// - we are returning a `Result` struct instance here
				// - we previously specified the return Type to be `paymentType`
				return m.Result, nil
			}
		},
	},
}


/*

3 pole, 5, 150 5 table, stage

*/