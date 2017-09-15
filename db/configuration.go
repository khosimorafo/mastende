package db

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		DbUrl	string `default:"mongodb://mastende:mastende@ds115573.mlab.com:15573/mastende-test"`
		DbName  string `default:"mastende-test"`
		Port    uint   `default:"3306"`
	}

	Invoices struct {

		NettDue     	int   `default:"4"`
		LastDiscountDay int   `default:"24"`
	}
}{}

