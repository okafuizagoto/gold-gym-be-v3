package goldgym

type GetOneStock struct {
	StockID         string  `db:"stock_id" json:"stock_id"`
	StockCode       string  `db:"stock_code" json:"stock_code"`
	StockName       string  `db:"stock_name" json:"stock_name"`
	StockPack       string  `db:"stock_pack" json:"stock_pack"`
	StockQTY        int     `db:"stock_qty" json:"stock_qty"`
	StockPrice      float32 `db:"stock_price" json:"stock_price"`
	StockLastUpdate string  `db:"stock_last_update" json:"stock_last_update"`
	StockUpdateBy   string  `db:"stock_update_by" json:"stock_update_by"`
}

type InsertStock struct {
	StockID       string  `db:"stock_id" json:"stock_id"`
	StockCode     string  `db:"stock_code" json:"stock_code"`
	StockName     string  `db:"stock_name" json:"stock_name"`
	StockPack     string  `db:"stock_pack" json:"stock_pack"`
	StockQTY      int     `db:"stock_qty" json:"stock_qty"`
	StockPrice    float32 `db:"stock_price" json:"stock_price"`
	StockUpdateBy string  `db:"stock_update_by" json:"stock_update_by"`
}

type InsertStockData struct {
	StockData InsertStock `json:"data"`
}
