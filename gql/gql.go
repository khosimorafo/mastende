package gql

import (
	"github.com/graphql-go/graphql"
	"time"
	"math/rand"
)

func init() {

	rand.Seed(time.Now().UnixNano())
}

func ResultTypeConfig() graphql.ObjectConfig {

	return graphql.ObjectConfig{
		Name: "Result",
		Fields: graphql.Fields{
			"result": &graphql.Field{
				Type: graphql.String,
			},
			"message": &graphql.Field{
				Type: graphql.String,
			},
		},
	}
}

func TenantTypeConfig() graphql.ObjectConfig {

	return graphql.ObjectConfig{
		Name: "Tenant",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"zaid": &graphql.Field{
				Type: graphql.String,
			},
			"mobile": &graphql.Field{
				Type: graphql.String,
			},
			"telephone": &graphql.Field{
				Type: graphql.String,
			},
			"site": &graphql.Field{
				Type: graphql.String,
			},
			"room": &graphql.Field{
				Type: graphql.String,
			},
			"gender": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"invoices": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(InvoiceTypeConfig())),
			},
			"payments": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(PaymentTypeConfig())),
			},
			"outstanding": &graphql.Field{
				Type: graphql.Float,
			},
			"overdue": &graphql.Field{
				Type: graphql.Float,
			},
		},
	}
}

func InvoiceTypeConfig() graphql.ObjectConfig {

	return graphql.ObjectConfig{
		Name: "Invoice",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"tenantid": &graphql.Field{
				Type: graphql.String,
			},
			"tenantname": &graphql.Field{
				Type: graphql.String,
			},
			"number": &graphql.Field{
				Type: graphql.String,
			},
			"reference": &graphql.Field{
				Type: graphql.String,
			},
			"total": &graphql.Field{
				Type: graphql.Float,
			},
			"balance": &graphql.Field{
				Type: graphql.Float,
			},
			"date": &graphql.Field{
				Type: graphql.String,
			},
			"duedate": &graphql.Field{
				Type: graphql.String,
			},
			"periodindex": &graphql.Field{
				Type: graphql.Int,
			},
			"periodname": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"lineitems": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(PaymentTypeConfig())),
			},
		},
	}
}

func PaymentTypeConfig() graphql.ObjectConfig {

	return graphql.ObjectConfig{

		Name: "Payment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"tenantid": &graphql.Field{
				Type: graphql.String,
			},
			"invoiceid": &graphql.Field{
				Type: graphql.String,
			},
			"number": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"amount": &graphql.Field{
				Type: graphql.Float,
			},
			"date": &graphql.Field{
				Type: graphql.String,
			},
			"mode": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
		},
	}
}

func TenantFieldArguments() graphql.FieldConfigArgument {

	return graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"zaid": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"moveindate": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"telephone": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"mobile": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"site": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"room": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"gender": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"status": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"moveoutdate": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"imageurl": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}
}

func InvoiceFieldArguments() graphql.FieldConfigArgument {

	return graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"tenantid": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"tenantname": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"number": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"reference": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"total": &graphql.ArgumentConfig{
			Type: graphql.Float,
		},
		"balance": &graphql.ArgumentConfig{
			Type: graphql.Float,
		},
		"lineitems": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"date": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"duedate": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"periodindex": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"periodname": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"status": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}
}

func PaymentFieldArguments() graphql.FieldConfigArgument {

	return graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"tenantid": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"invoiceid": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"number": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"description": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"amount": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"date": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"mode": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"status": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}
}

func ItemFieldArguments() graphql.FieldConfigArgument {

	return graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"rate": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"status": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}
}

func LineItemFieldArguments() graphql.FieldConfigArgument {

	return graphql.FieldConfigArgument{

		"name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"rate": &graphql.ArgumentConfig{
			Type: graphql.Float,
		},
		"quantity": &graphql.ArgumentConfig{
			Type: graphql.Float,
		},
		"total": &graphql.ArgumentConfig{
			Type: graphql.Float,
		},
		"discount": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}
}

func PaymentExtensionFieldArguments() graphql.FieldConfigArgument {

	return graphql.FieldConfigArgument{

		"invoiceid": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"paybydate": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"requestdate": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"requestby": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"requestmode": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	}
}



