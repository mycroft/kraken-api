package krakenapi

import (
	"net/url"
	"strconv"
)

/*
URL: https://api.kraken.com/0/private/TradeBalance

Input:

aclass = asset class (optional):
    currency (default)
asset = base asset used to determine balance (default = ZUSD)
Result: array of trade balance info

eb = equivalent balance (combined balance of all currencies)
tb = trade balance (combined balance of all equity currencies)
m = margin amount of open positions
n = unrealized net profit/loss of open positions
c = cost basis of open positions
v = current floating valuation of open positions
e = equity = trade balance + unrealized net profit/loss
mf = free margin = equity - initial margin (maximum margin available to open new positions)
ml = margin level = (equity / initial margin) * 100
Note: Rates used for the floating valuation is the midpoint of the best bid and ask prices
*/
func (api *KrakenApi) ApiBalance() (map[string]float64, error) {
	resp, err := api.Query(URL_PRIVATE_BALANCE, url.Values{}, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, nil)
	if err != nil {
		return nil, err
	}

	balance := make(map[string]float64)

	for v, r := range content.(map[string]interface{}) {
		value, err := strconv.ParseFloat(r.(string), 64)
		if err != nil {
			return balance, err
		}

		balance[v] = value
	}

	return balance, nil
}

/*
Input:

aclass = asset class (optional):
    currency (default)
asset = base asset used to determine balance (default = ZUSD)
Result: array of trade balance info

eb = equivalent balance (combined balance of all currencies)
tb = trade balance (combined balance of all equity currencies)
m = margin amount of open positions
n = unrealized net profit/loss of open positions
c = cost basis of open positions
v = current floating valuation of open positions
e = equity = trade balance + unrealized net profit/loss
mf = free margin = equity - initial margin (maximum margin available to open new positions)
ml = margin level = (equity / initial margin) * 100
Note: Rates used for the floating valuation is the midpoint of the best bid and ask prices
*/
func (api *KrakenApi) ApiTradeBalance(asset string) (*TradeBalance, error) {
	params := url.Values{}
	params.Set("asset", asset)

	resp, err := api.Query(URL_PRIVATE_TRADE_BALANCE, params, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, &TradeBalance{})
	if err != nil {
		return nil, err
	}

	return content.(*TradeBalance), nil
}

/*
URL: https://api.kraken.com/0/private/OpenOrders

Input:

trades = whether or not to include trades in output (optional.  default = false)
userref = restrict results to given user reference id (optional)
Result: array of order info in open array with txid as the key

refid = Referral order transaction id that created this order
userref = user reference id
status = status of order:
    pending = order pending book entry
    open = open order
    closed = closed order
    canceled = order canceled
    expired = order expired
opentm = unix timestamp of when order was placed
starttm = unix timestamp of order start time (or 0 if not set)
expiretm = unix timestamp of order end time (or 0 if not set)
descr = order description info
    pair = asset pair
    type = type of order (buy/sell)
    ordertype = order type (See Add standard order)
    price = primary price
    price2 = secondary price
    leverage = amount of leverage
    order = order description
    close = conditional close order description (if conditional close set)
vol = volume of order (base currency unless viqc set in oflags)
vol_exec = volume executed (base currency unless viqc set in oflags)
cost = total cost (quote currency unless unless viqc set in oflags)
fee = total fee (quote currency)
price = average price (quote currency unless viqc set in oflags)
stopprice = stop price (quote currency, for trailing stops)
limitprice = triggered limit price (quote currency, when limit based order type triggered)
misc = comma delimited list of miscellaneous info
    stopped = triggered by stop price
    touched = triggered by touch price
    liquidated = liquidation
    partial = partial fill
oflags = comma delimited list of order flags
    viqc = volume in quote currency
    fcib = prefer fee in base currency (default if selling)
    fciq = prefer fee in quote currency (default if buying)
    nompp = no market price protection
trades = array of trade ids related to order (if trades info requested and data available)
Note: Unless otherwise stated, costs, fees, prices, and volumes are in the asset pair's scale,
      not the currency's scale. For example, if the asset pair uses a lot size that has a scale
      of 8, the volume will use a scale of 8, even if the currency it represents only has a scale
      of 2. Similarly, if the asset pair's pricing scale is 5, the scale will remain as 5, even
      if the underlying currency has a scale of 8.
*/
func (api *KrakenApi) ApiOpenOrders(trades bool, userref string) (*OpenOrders, error) {
	params := url.Values{}
	if trades {
		params.Set("trades", "true")
	}
	if userref != "" {
		params.Set("userref", userref)
	}

	resp, err := api.Query(URL_PRIVATE_OPEN_ORDERS, params, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, &OpenOrders{})
	if err != nil {
		return nil, err
	}

	return content.(*OpenOrders), nil
}

/*
URL: https://api.kraken.com/0/private/ClosedOrders

Input:

trades = whether or not to include trades in output (optional.  default = false)
userref = restrict results to given user reference id (optional)
start = starting unix timestamp or order tx id of results (optional.  exclusive)
end = ending unix timestamp or order tx id of results (optional.  inclusive)
ofs = result offset
closetime = which time to use (optional)
    open
    close
    both (default)
Result: array of order info

closed = array of order info.  See Get open orders.  Additional fields:
    closetm = unix timestamp of when order was closed
    reason = additional info on status (if any)
count = amount of available order info matching criteria
Note: Times given by order tx ids are more accurate than unix timestamps.
      If an order tx id is given for the time, the order's open time is used
*/
func (api *KrakenApi) ApiClosedOrders(trades bool, userref, start, end string, ofs int, closetime string) (*ClosedOrders, error) {
	params := url.Values{}
	if trades {
		params.Set("trades", "true")
	}

	if userref != "" {
		params.Set("userref", userref)
	}

	if start != "" {
		params.Set("start", start)
	}

	if end != "" {
		params.Set("end", end)
	}

	if ofs != 0 {
		params.Set("ofs", strconv.Itoa(ofs))
	}

	if closetime != "" {
		params.Set("closetime", closetime)
	}

	resp, err := api.Query(URL_PRIVATE_CLOSED_ORDERS, params, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, &ClosedOrders{})
	if err != nil {
		return nil, err
	}

	return content.(*ClosedOrders), nil
}

/*
URL: https://api.kraken.com/0/private/QueryOrders

Input:

trades = whether or not to include trades in output (optional.  default = false)
userref = restrict results to given user reference id (optional)
txid = comma delimited list of transaction ids to query info about (20 maximum)
Result: associative array of orders info

<order_txid> = order info.  See Get open orders/Get closed orders
*/
func (api *KrakenApi) ApiQueryOrders(trades bool, userref string, txid string) (*QueryOrder, error) {
	params := url.Values{}
	if trades {
		params.Set("trades", "true")
	}
	if userref != "" {
		params.Set("userref", userref)
	}
	params.Set("txid", txid)

	resp, err := api.Query(URL_PRIVATE_QUERY_ORDERS, params, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, &QueryOrder{})
	if err != nil {
		return nil, err
	}

	return content.(*QueryOrder), nil
}

/*
type = type of trade (optional)
    all = all types (default)
    any position = any position (open or closed)
    closed position = positions that have been closed
    closing position = any trade closing all or part of a position
    no position = non-positional trades
trades = whether or not to include trades related to position in output (optional.  default = false)
start = starting unix timestamp or trade tx id of results (optional. exclusive)
end = ending unix timestamp or trade tx id of results (optional. inclusive)
ofs = result offset
*/
func (api *KrakenApi) ApiTradesHistory(trade_type string, incl_trades bool, start, end string, ofs int) (map[string]Trade, error) {
	params := url.Values{}
	if trade_type != "" {
		params.Set("type", trade_type)
	}

	if incl_trades {
		params.Set("trades", "true")
	}

	if start != "" {
		params.Set("start", start)
	}

	if end != "" {
		params.Set("end", end)
	}

	if ofs != 0 {
		params.Set("ofs", strconv.Itoa(ofs))
	}

	resp, err := api.Query(URL_PRIVATE_TRADES_HISTORY, params, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, &TradeHistoryResult{})
	if err != nil {
		return nil, err
	}

	return content.(*TradeHistoryResult).Trades, nil
}

/*
URL: https://api.kraken.com/0/private/QueryTrades

Input:

txid = comma delimited list of transaction ids to query info about (20 maximum)
trades = whether or not to include trades related to position in output (optional.  default = false)

Result: associative array of trades info

<trade_txid> = trade info.  See Get trades history

*/
func (api *KrakenApi) ApiQueryTrades(txid string, incl_trades bool) (map[string]Trade, error) {
	params := url.Values{}
	params.Set("txid", txid)

	if incl_trades {
		params.Set("trades", "true")
	}

	resp, err := api.Query(URL_PRIVATE_QUERY_TRADES, params, true)
	if err != nil {
		return nil, err
	}

	out := make(map[string]Trade)

	_, err = parse(resp, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
URL: https://api.kraken.com/0/private/OpenPositions

Input:

txid = comma delimited list of transaction ids to restrict output to
docalcs = whether or not to include profit/loss calculations (optional.  default = false)
Result: associative array of open position info

<position_txid> = open position info
    ordertxid = order responsible for execution of trade
    pair = asset pair
    time = unix timestamp of trade
    type = type of order used to open position (buy/sell)
    ordertype = order type used to open position
    cost = opening cost of position (quote currency unless viqc set in oflags)
    fee = opening fee of position (quote currency)
    vol = position volume (base currency unless viqc set in oflags)
    vol_closed = position volume closed (base currency unless viqc set in oflags)
    margin = initial margin (quote currency)
    value = current value of remaining position (if docalcs requested.  quote currency)
    net = unrealized profit/loss of remaining position (if docalcs requested.  quote currency, quote currency scale)
    misc = comma delimited list of miscellaneous info
    oflags = comma delimited list of order flags
        viqc = volume in quote currency
*/
func (api *KrakenApi) ApiOpenPositions(txid string, docalcs bool) (map[string]OpenPosition, error) {
	params := url.Values{}
	if txid != "" {
		params.Set("txid", txid)
	}

	if docalcs {
		params.Set("docalcs", "true")
	}

	resp, err := api.Query(URL_PRIVATE_OPEN_POSITIONS, params, true)
	if err != nil {
		return nil, err
	}

	out := make(map[string]OpenPosition)

	_, err = parse(resp, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
URL: https://api.kraken.com/0/private/Ledgers

Input:

aclass = asset class (optional):
    currency (default)
asset = comma delimited list of assets to restrict output to (optional.  default = all)
type = type of ledger to retrieve (optional):
    all (default)
    deposit
    withdrawal
    trade
    margin
start = starting unix timestamp or ledger id of results (optional.  exclusive)
end = ending unix timestamp or ledger id of results (optional.  inclusive)
ofs = result offset
Result: associative array of ledgers info

<ledger_id> = ledger info
    refid = reference id
    time = unx timestamp of ledger
    type = type of ledger entry
    aclass = asset class
    asset = asset
    amount = transaction amount
    fee = transaction fee
    balance = resulting balance
Note: Times given by ledger ids are more accurate than unix timestamps.
*/
func (api *KrakenApi) ApiLedgers(asset, ledger_type, start, end string, ofs int) (map[string]Ledger, error) {
	params := url.Values{}
	if asset != "" {
		params.Set("asset", asset)
	}

	if ledger_type != "" {
		params.Set("type", ledger_type)
	}

	if start != "" {
		params.Set("start", start)
	}

	if end != "" {
		params.Set("end", end)
	}

	if ofs != 0 {
		params.Set("ofs", strconv.Itoa(ofs))
	}

	resp, err := api.Query(URL_PRIVATE_LEDGERS, params, true)
	if err != nil {
		return nil, err
	}

	out := new(LedgerResponse)

	_, err = parse(resp, &out)
	if err != nil {
		return nil, err
	}

	return out.Ledger, nil
}

/*
Query ledgers
URL: https://api.kraken.com/0/private/QueryLedgers

Input:

id = comma delimited list of ledger ids to query info about (20 maximum)
Result: associative array of ledgers info

<ledger_id> = ledger info.  See Get ledgers info
*/
func (api *KrakenApi) ApiQueryLedgers(id string) (map[string]Ledger, error) {
	params := url.Values{}
	params.Set("id", id)

	resp, err := api.Query(URL_PRIVATE_QUERY_LEDGERS, params, true)
	if err != nil {
		return nil, err
	}

	out := make(map[string]Ledger)

	_, err = parse(resp, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
Get trade volume
URL: https://api.kraken.com/0/private/TradeVolume

Input:

pair = comma delimited list of asset pairs to get fee info on (optional)
fee-info = whether or not to include fee info in results (optional)
Result: associative array

currency = volume currency
volume = current discount volume
fees = array of asset pairs and fee tier info (if requested)
    fee = current fee in percent
    minfee = minimum fee for pair (if not fixed fee)
    maxfee = maximum fee for pair (if not fixed fee)
    nextfee = next tier's fee for pair (if not fixed fee.  nil if at lowest fee tier)
    nextvolume = volume level of next tier (if not fixed fee.  nil if at lowest fee tier)
    tiervolume = volume level of current tier (if not fixed fee.  nil if at lowest fee tier)
fees_maker = array of asset pairs and maker fee tier info (if requested) for any pairs on maker/taker schedule
    fee = current fee in percent
    minfee = minimum fee for pair (if not fixed fee)
    maxfee = maximum fee for pair (if not fixed fee)
    nextfee = next tier's fee for pair (if not fixed fee.  nil if at lowest fee tier)
    nextvolume = volume level of next tier (if not fixed fee.  nil if at lowest fee tier)
    tiervolume = volume level of current tier (if not fixed fee.  nil if at lowest fee tier)
Note: If an asset pair is on a maker/taker fee schedule, the taker side is given in "fees" and maker side in "fees_maker". For pairs not on maker/taker, they will only be given in "fees".
*/
func (api *KrakenApi) ApiTradeVolume(pair string, feeinfo bool) (*TradeVolume, error) {
	params := url.Values{}
	params.Set("pair", pair)
	if feeinfo {
		params.Set("fee-info", "true")
	}

	resp, err := api.Query(URL_PRIVATE_TRADE_VOLUME, params, true)
	if err != nil {
		return nil, err
	}

	out := new(TradeVolume)

	_, err = parse(resp, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
