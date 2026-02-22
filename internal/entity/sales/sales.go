package goldgym

type SalesHeader struct {
	SaleID           string  `db:"sale_id" json:"sale_id"`
	SaleTransdate    string  `db:"sale_transdate" json:"sale_transdate"`
	SaleTransTime    string  `db:"sale_transtime" json:"sale_transtime"`
	SaleTranstotal   float64 `db:"sale_transtotal" json:"sale_transtotal"`
	SaleTranspayment float64 `db:"sale_transpayment" json:"sale_transpayment"`
	SaleTranschange  float64 `db:"sale_transchange" json:"sale_transchange"`
	SaleSalesperson  string  `db:"sale_salesperson" json:"sale_salesperson"`
}

type SalesDetail struct {
	SaleID         string `db:"sale_id" json:"sale_id"`
	SaleStockID    string `db:"sale_stockid" json:"sale_stockid"`
	SaleStockcode  string `db:"sale_stockcode" json:"sale_stockcode"`
	SaleStockname  string `db:"sale_stockname" json:"sale_stockname"`
	SaleQty        string `db:"sale_qty" json:"sale_qty"`
	SaleSalesprice string `db:"sale_salesprice" json:"sale_salesprice"`
	SalePack      string `db:"sale_pack" json:"sale_pack"`
	SaleLastupdate      string `db:"sale_lastupdate" json:"sale_lastupdate"`
}
