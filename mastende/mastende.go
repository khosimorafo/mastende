package mastende

import (
	"github.com/mitchellh/mapstructure"
	"github.com/khosimorafo/mastende/tenants"
	"errors"
	"github.com/khosimorafo/mastende/db"
	"github.com/khosimorafo/mastende/items"
	"github.com/khosimorafo/mastende/invoices"
	"github.com/khosimorafo/mastende/utils"
	"github.com/khosimorafo/mastende/payments"
	"github.com/khosimorafo/mastende/periods"
	"fmt"
)

/**
*
* Create new tenants record.
*/
func New(app *db.App) *Mastende {

	m := &Mastende{App:app}
	return m
}

type Mastende struct {

	App *db.App

	Tenant *tenants.Tenant
	Item *items.Item
	Invoice *invoices.Invoice
	Payment *payments.Payment

	Tenants []tenants.Tenant
	Items []items.Item
	Invoices []invoices.Invoice
	Payments []payments.Payment
}

/********************************************************************************/
/*********************************Tenant*****************************************/
/********************************************************************************/

func (mastende *Mastende) CreateTenant(input map[string]interface{}, with_invoice bool) error {

	mastende.Tenant = tenants.New(mastende.App)

	// 1. Marshal input into mastende._tenant.
	if err := mapstructure.Decode(input, mastende.Tenant); err != nil{

		return errors.New("Data marshalling failure.")
	}
	// 2. Check if the data is valid.
	if err := mastende.Tenant.Validate(); err != nil{

		return errors.New("Failed to validate submitted _tenant data.")
	}
	// #. Persist _tenant into the database.
	if err := mastende.Tenant.Persist(); err != nil {

		error_str := fmt.Sprintf(err.Error())
		return errors.New(error_str)
	}

	if with_invoice{

		createNewTenantInvoice(mastende)
	}

	return nil
}

func (mastende *Mastende) GetTenant(input map[string]interface{}) error {

	id := input["id"].(string)
	t, err := tenants.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else { mastende.Tenant = t }

	return nil
}

func (mastende *Mastende) UpdateTenant(input map[string]interface{}) error {

	id := input["id"].(string)
	// Retrieve the tenants
	t, err := tenants.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else { mastende.Tenant = t }

	// 1. Marshal input into mastende._tenant.
	if err := mapstructure.Decode(input, mastende.Tenant); err != nil{

		return errors.New("Data marshalling failure.")
	}

	// #. Update tenants into the database.
	if err := mastende.Tenant.Update(); err != nil {

		return errors.New("Error updating tenants.")
	}

	return nil
}

func (mastende *Mastende) DeleteTenant(input map[string]interface{}) error {

	id := input["id"].(string)
	t, err := tenants.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else {

		mastende.Tenant = t
		mastende.Tenant.Delete()
	}

	return nil
}

func (mastende *Mastende) TenantList(input map[string]interface{}) error {



	return nil
}

func createNewTenantInvoice(mastende *Mastende) error {

	i_chan := make(chan invoices.Invoice)

	go func() {
		i, _ := mastende.Tenant.MonthlyInvoice()
		i_chan <- *i
	}()

	invoice := <- i_chan

	//m := structs.Map(invoice)

	input := map[string]interface{}{
		"tenantid":       	invoice.TenantID,
		"tenantname":       	invoice.TenantName,
		"number":		invoice.Number,
		"reference":		invoice.Reference,
		"lineitems":		invoice.LineItems,
		"date":       		invoice.Date,
		"duedate": 		invoice.DueDate,
		"periodindex":		invoice.PeriodIndex,
		"periodname":		invoice.PeriodName,
		"status":		invoice.Status,
	}

	if err := mastende.CreateInvoice(input); err != nil {

		return errors.New("Error creating new tenant invoice.")
	}
	return nil
}


/********************************************************************************/
/***********************************Item*****************************************/
/********************************************************************************/

func (mastende *Mastende) CreateItem(input map[string]interface{}) error {

	mastende.Item = items.New(mastende.App)

	// 1. Marshal input into mastende._item.
	if err := mapstructure.Decode(input, mastende.Item); err != nil{

		return errors.New("Data marshalling failure.")
	}
	// 2. Check if the data is valid.
	if err := mastende.Item.Validate(); err != nil{

		return errors.New("Failed to validate submitted _item data.")
	}
	// #. Persist _item into the database.
	if err := mastende.Item.Persist(); err != nil {

		return errors.New("Error creating items.")
	}

	return nil
}

func (mastende *Mastende) GetItem(input map[string]interface{}) error {

	id := input["id"].(string)
	t, err := items.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else { mastende.Item = t }

	return nil
}

func (mastende *Mastende) UpdateItem(input map[string]interface{}) error {

	id := input["id"].(string)
	// Retrieve the items
	t, err := items.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else { mastende.Item = t }

	// 1. Marshal input into mastende._item.
	if err := mapstructure.Decode(input, mastende.Item); err != nil{

		return errors.New("Data marshalling failure.")
	}

	// 2. Check if the data is valid.
	if err := mastende.Item.Validate(); err != nil{

		return errors.New("Failed to validate submitted items data.")
	}
	// #. Update items into the database.
	if err := mastende.Item.Update(); err != nil {

		return errors.New("Error updating items.")
	}

	return nil
}

func (mastende *Mastende) DeleteItem(input map[string]interface{}) error {

	id := input["id"].(string)
	t, err := items.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else {

		mastende.Item = t
		mastende.Item.Delete()
	}

	return nil
}


/********************************************************************************/
/********************************Invoice*****************************************/
/********************************************************************************/

func (mastende *Mastende) CreateInvoice(input map[string]interface{}) error {

	mastende.Invoice = invoices.New(mastende.App)

	var _date	string

	_date 		= input["date"].(string)

	period, err := periods.NewInstanceWithDate(mastende.App, _date)

	if err != nil {

		mastende.Invoice = &invoices.Invoice{}
		error_str := fmt.Sprintf("Could not derive period. %s", err.Error())
		return errors.New(error_str)
	}

	mastende.Invoice.PeriodIndex 	= period.Index
	mastende.Invoice.PeriodName 	= period.Name

	// 1. Marshal input into mastende._invoice.
	if err := mapstructure.Decode(input, mastende.Invoice); err != nil{

		mastende.Invoice = &invoices.Invoice{}
		return errors.New("Data marshalling failure. ")
	}
	// 2. Check if the data is valid.
	if err := mastende.Invoice.Validate(); err != nil{

		mastende.Invoice = &invoices.Invoice{}
		return errors.New("Failed to validate submitted _invoice data. ")
	}
	// #. Persist _invoice into the database.
	if err := mastende.Invoice.Persist(); err != nil {

		mastende.Invoice = &invoices.Invoice{}
		return errors.New("Error creating invoices. ")
	}

	return nil
}

func (mastende *Mastende) GetInvoice(input map[string]interface{}) error {

	id := input["id"].(string)
	t, err := invoices.NewInstanceWithId(mastende.App, id)

	if err != nil {

		mastende.Invoice = &invoices.Invoice{}
		return err
	} else {

		mastende.Invoice = t
	}

	return nil
}

func (mastende *Mastende) UpdateInvoice(input map[string]interface{}) error {

	id := input["id"].(string)
	// Retrieve the invoices
	t, err := invoices.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else { mastende.Invoice = t }

	// 1. Marshal input into mastende._invoice.
	if err := mapstructure.Decode(input, mastende.Invoice); err != nil{

		mastende.Invoice = &invoices.Invoice{}
		return errors.New("Data marshalling failure.")
	}

	// 2. Check if the data is valid.
	if err := mastende.Invoice.Validate(); err != nil{

		mastende.Invoice = &invoices.Invoice{}
		return errors.New("Failed to validate submitted invoices data.")
	}
	// #. Update invoices into the database.
	if err := mastende.Invoice.Update(); err != nil {

		mastende.Invoice = &invoices.Invoice{}
		return errors.New("Error updating invoices.")
	}

	return nil
}

func (mastende *Mastende) AddInvoiceLineItem(input map[string]interface{}) error {

	id := input["id"].(string)
	invoice, err := invoices.NewInstanceWithId(mastende.App, id)

	// Check if invoice valid
	if err != nil {

		return err
	} else {

		lineitem := input["lineitem"].(map[string]interface{})

		if err := invoice.AddLineItem(lineitem); err != nil{

			error_str := fmt.Sprintf(err.Error())
			return errors.New(error_str)
		}
	}

	return nil
}

func (mastende *Mastende) DeleteInvoice(input map[string]interface{}) error {

	id := input["id"].(string)
	t, err := invoices.NewInstanceWithId(mastende.App, id)

	if err != nil {

		return err
	} else {

		mastende.Invoice = t
		mastende.Invoice.Delete()
	}

	return nil
}

func (mastende *Mastende) MakeInvoicePayment(input map[string]interface{}) error {

	var (
		_invoiceid      	string
		_tenantid       	string
		_number		    	string
		_description    	string
		_date           	string
		_payment_amount 	float64 = 0.0
		_mode		       	string

	)

	_invoiceid 		= input["invoiceid"].(string)
	_tenantid 		= input["tenantid"].(string)
	_number 		= input["number"].(string)
	_description 		= input["description"].(string)
	_payment_amount 	= input["amount"].(float64)
	_date 			= input["date"].(string)
	_mode 			= input["mode"].(string)

	// 1. Check if invoice exists
	invoice, err := invoices.NewInstanceWithId(mastende.App, _invoiceid)
	if err != nil {

		return errors.New("Error retrieving invoice. ")
	}

	// 2. Check if amount is valid
	if invoice.Balance < _payment_amount {

		return errors.New("Payment amount bigger than invoice balance. ")
	}

	// 3. Check if date is valid
	_, _, err = utils.DateFormatter( _date)
	if err != nil {

		return errors.New("Date is invalid. ")
	}

	mastende.Payment = payments.New(mastende.App)

	mastende.Payment.InvoiceID 	= _invoiceid
	mastende.Payment.TenantID 	= _tenantid
	mastende.Payment.Number		= _number
	mastende.Payment.Date		= _date
	mastende.Payment.Mode		= _mode
	mastende.Payment.Description	= _description

	// #. Persist _tenant into the database.
	if err := mastende.Payment.Persist(); err != nil {

		mastende.Payment = &payments.Payment{}
		return errors.New("Error creating payment. ")
	}

	return nil
}

func (mastende *Mastende) DiscountInvoice(input map[string]interface{}) error {

	var (
		_invoiceid       string
		_discount_amount float64 = 0.0
	)

	_invoiceid 		= input["invoiceid"].(string)
	_discount_amount 	= input["amount"].(float64)

	// 1. Check if invoice exists
	invoice, err := invoices.NewInstanceWithId(mastende.App, _invoiceid)
	if err != nil {

		error_str := fmt.Sprintf("Error retrieving invoice. %s", err.Error())
		return errors.New(error_str)
	}

	// 2. Check if amount is valid
	if invoice.Balance < _discount_amount {

		return errors.New("Discount amount bigger than invoice balance. ")
	}

	if err := invoice.ApplyDiscount(input); err != nil{

		return errors.New("Error applying discount. ")
	}

	return nil
}


/********************************************************************************/
/********************************Item*****************************************/
/********************************************************************************/

