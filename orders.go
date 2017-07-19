package krakenapi

import (
	"net/url"
	"strconv"
)

/*
Input:

pair = asset pair
type = type of order (buy/sell)
ordertype = order type:
    market
    limit (price = limit price)
    stop-loss (price = stop loss price)
    take-profit (price = take profit price)
    stop-loss-profit (price = stop loss price, price2 = take profit price)
    stop-loss-profit-limit (price = stop loss price, price2 = take profit price)
    stop-loss-limit (price = stop loss trigger price, price2 = triggered limit price)
    take-profit-limit (price = take profit trigger price, price2 = triggered limit price)
    trailing-stop (price = trailing stop offset)
    trailing-stop-limit (price = trailing stop offset, price2 = triggered limit offset)
    stop-loss-and-limit (price = stop loss price, price2 = limit price)
    settle-position
price = price (optional.  dependent upon ordertype)
price2 = secondary price (optional.  dependent upon ordertype)
volume = order volume in lots

leverage = amount of leverage desired (optional.  default = none)
oflags = comma delimited list of order flags (optional):
    viqc = volume in quote currency (not available for leveraged orders)
    fcib = prefer fee in base currency
    fciq = prefer fee in quote currency
    nompp = no market price protection
    post = post only order (available when ordertype = limit)
starttm = scheduled start time (optional):
    0 = now (default)x
    +<n> = schedule start time <n> seconds from now
    <n> = unix timestamp of start time
expiretm = expiration time (optional):
    0 = no expiration (default)
    +<n> = expire <n> seconds from now
    <n> = unix timestamp of expiration time
userref = user reference id.  32-bit signed number.  (optional)
validate = validate inputs only.  do not submit order (optional)

optional closing order to add to system when order gets filled:
    close[ordertype] = order type
    close[price] = price
    close[price2] = secondary price
Result:

descr = order description info
    order = order description
    close = conditional close order description (if conditional close set)
txid = array of transaction ids for order (if order was added successfully)
*/
func (api *KrakenApi) ApiAddOrder(pair, bstype, ordertype string, price, price2, volume float64, oflags string) (*OrderResult, error) {
	params := url.Values{}
	params.Set("pair", pair)
	params.Set("type", bstype)
	params.Set("ordertype", ordertype)
	if ordertype != "market" {
		params.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	}

	params.Set("volume", strconv.FormatFloat(volume, 'f', -1, 64))
	// params.Set("validate", "1") /* Used  for debug */

	if price2 != 0 {
		params.Set("price2", strconv.FormatFloat(price, 'f', -1, 64))
	}

	if oflags != "" {
		params.Set("oflags", oflags)
	}

	resp, err := api.Query(URL_PRIVATE_ADD_ORDER, params, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, &OrderResult{})
	if err != nil {
		return nil, err
	}

	return content.(*OrderResult), nil
}

/*
Input:

txid = transaction id
Result:

count = number of orders canceled
pending = if set, order(s) is/are pending cancellation
*/
func (api *KrakenApi) ApiCancelOrder(txid string) (*CancelResult, error) {
	params := url.Values{}
	params.Set("txid", txid)

	resp, err := api.Query(URL_PRIVATE_CANCEL_ORDER, params, true)
	if err != nil {
		return nil, err
	}

	content, err := parse(resp, &CancelResult{})
	if err != nil {
		return nil, err
	}

	return content.(*CancelResult), err
}
