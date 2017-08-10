package invoices

type LineItem struct {

	ItemID      		string  	`json:"itemid"`
	InvoiceID     		string  	`json:"invoiceid"`
	Name        		string  	`json:"name,omitempty"`
	Description	 	string  	`json:"description,omitempty"`
	Rate        		float64 	`json:"rate"`
	Quantity    		float64   	`json:"quantity"`
}
