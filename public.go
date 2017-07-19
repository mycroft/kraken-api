package krakenapi

import (
	"net/url"
	"strconv"
	"strings"
)

/*
URL: https://api.kraken.com/0/public/Time

Result: Server's time

unixtime =  as unix timestamp
rfc1123 = as RFC 1123 time format
Note: This is to aid in approximating the skew time between the server and client.
*/
func (api *KrakenApi) ApiServerTime() (interface{}, error) {
	resp, err := api.Query(URL_PUBLIC_TIME, url.Values{}, false)
	if err != nil {
		return nil, err
	}

	return parse(resp, nil)
}

/*
URL: https://api.kraken.com/0/public/Assets

Input:

info = info to retrieve (optional):
    info = all info (default)
aclass = asset class (optional):
    currency (default)
asset = comma delimited list of assets to get info on (optional.  default = all for given asset class)
Result: array of asset names and their info

<asset_name> = asset name
    altname = alternate name
    aclass = asset class
    decimals = scaling decimal places for record keeping
    display_decimals = scaling decimal places for output display
*/
func (api *KrakenApi) ApiAssets() (map[string]Asset, error) {
	resp, err := api.Query(URL_PUBLIC_ASSETS, url.Values{}, false)
	if err != nil {
		return nil, err
	}

	assets := make(map[string]Asset)

	_, err = parse(resp, &assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

/*
URL: https://api.kraken.com/0/public/AssetPairs

Input:

info = info to retrieve (optional):
    info = all info (default)
    leverage = leverage info
    fees = fees schedule
    margin = margin info
pair = comma delimited list of asset pairs to get info on (optional.  default = all)
Result: array of pair names and their info

<pair_name> = pair name
    altname = alternate pair name
    aclass_base = asset class of base component
    base = asset id of base component
    aclass_quote = asset class of quote component
    quote = asset id of quote component
    lot = volume lot size
    pair_decimals = scaling decimal places for pair
    lot_decimals = scaling decimal places for volume
    lot_multiplier = amount to multiply lot volume by to get currency volume
    leverage_buy = array of leverage amounts available when buying
    leverage_sell = array of leverage amounts available when selling
    fees = fee schedule array in [volume, percent fee] tuples
    fees_maker = maker fee schedule array in [volume, percent fee] tuples (if on maker/taker)
    fee_volume_currency = volume discount currency
    margin_call = margin call level
    margin_stop = stop-out/liquidation margin level
Note: If an asset pair is on a maker/taker fee schedule, the taker side is given in "fees"
       and maker side in "fees_maker". For pairs not on maker/taker, they will only be given in "fees".
*/
func (api *KrakenApi) ApiAssetPairs(info, pair string) (map[string]AssetPair, error) {
	params := url.Values{}
	if pair != "" {
		params.Set("pair", pair)
	}

	if info != "" {
		params.Set("info", info)
	}

	resp, err := api.Query(URL_PUBLIC_ASSET_PAIRS, params, false)
	if err != nil {
		return nil, err
	}

	assets_pairs := make(map[string]AssetPair)

	_, err = parse(resp, &assets_pairs)
	if err != nil {
		return nil, err
	}

	return assets_pairs, nil
}

/*
URL: https://api.kraken.com/0/public/Ticker

Input:

pair = comma delimited list of asset pairs to get info on
Result: array of pair names and their ticker info

<pair_name> = pair name
    a = ask array(<price>, <whole lot volume>, <lot volume>),
    b = bid array(<price>, <whole lot volume>, <lot volume>),
    c = last trade closed array(<price>, <lot volume>),
    v = volume array(<today>, <last 24 hours>),
    p = volume weighted average price array(<today>, <last 24 hours>),
    t = number of trades array(<today>, <last 24 hours>),
    l = low array(<today>, <last 24 hours>),
    h = high array(<today>, <last 24 hours>),
    o = today's opening price
Note: Today's prices start at 00:00:00 UTC
*/
func (api *KrakenApi) ApiTicker(pairs []string) (map[string]Ticker, error) {
	params := url.Values{}
	params.Set("pair", strings.Join(pairs, ","))

	resp, err := api.Query(URL_PUBLIC_TICKER, params, false)
	if err != nil {
		return nil, err
	}

	tickers := make(map[string]Ticker)

	_, err = parse(resp, &tickers)
	if err != nil {
		return nil, err
	}

	return tickers, nil
}

/*
URL: https://api.kraken.com/0/public/OHLC

Input:

pair = asset pair to get OHLC data for
interval = time frame interval in minutes (optional):
	1 (default), 5, 15, 30, 60, 240, 1440, 10080, 21600
since = return committed OHLC data since given id (optional.  exclusive)
Result: array of pair name and OHLC data

<pair_name> = pair name
    array of array entries(<time>, <open>, <high>, <low>, <close>, <vwap>, <volume>, <count>)
last = id to be used as since when polling for new, committed OHLC data
Note: the last entry in the OHLC array is for the current, not-yet-committed frame and will always
      be present, regardless of the value of "since".
*/
func (api *KrakenApi) ApiOHLC(pair string, interval int, since uint64) (float64, []OHLCEntry, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if interval != 0 {
		params.Set("interval", strconv.Itoa(interval))
	}

	if since != 0 {
		params.Set("since", strconv.FormatUint(since, 10))
	}

	resp, err := api.Query(URL_PUBLIC_OHLC, params, false)
	if err != nil {
		return 0, nil, err
	}

	content, err := parse(resp, nil)
	if err != nil {
		return 0, nil, err
	}

	var last float64
	ohlc_data := make([]OHLCEntry, 0)

	// Don't have much choice as returned json contains arrays, not structures...
	for key, value := range content.(map[string]interface{}) {
		if key == "last" {
			last = value.(float64)
		} else {
			for _, subvalue := range value.([]interface{}) {
				values := make([]interface{}, len(subvalue.([]interface{})))
				for i, subsubvalue := range subvalue.([]interface{}) {
					values[i] = subsubvalue
				}

				open, err := strconv.ParseFloat(values[1].(string), 64)
				if err != nil {
					return 0, nil, err
				}

				high, err := strconv.ParseFloat(values[2].(string), 64)
				if err != nil {
					return 0, nil, err
				}

				low, err := strconv.ParseFloat(values[3].(string), 64)
				if err != nil {
					return 0, nil, err
				}

				close, err := strconv.ParseFloat(values[4].(string), 64)
				if err != nil {
					return 0, nil, err
				}

				vwap, err := strconv.ParseFloat(values[5].(string), 64)
				if err != nil {
					return 0, nil, err
				}

				volume, err := strconv.ParseFloat(values[6].(string), 64)
				if err != nil {
					return 0, nil, err
				}

				entry := OHLCEntry{
					Time:   values[0].(float64),
					Open:   open,
					High:   high,
					Low:    low,
					Close:  close,
					VWAP:   vwap,
					Volume: volume,
					Count:  values[7].(float64),
				}

				ohlc_data = append(ohlc_data, entry)
			}
		}
	}

	return last, ohlc_data, nil
}

/*
URL: https://api.kraken.com/0/public/Depth

Input:

pair = asset pair to get market depth for
count = maximum number of asks/bids (optional)
Result: array of pair name and market depth

<pair_name> = pair name
    asks = ask side array of array entries(<price>, <volume>, <timestamp>)
    bids = bid side array of array entries(<price>, <volume>, <timestamp>)
*/
func (api *KrakenApi) ApiDepth(pair string, count int) (map[string]PublicOrderBook, error) { // XXX
	params := url.Values{}
	params.Set("pair", pair)

	if count != 0 {
		params.Set("count", strconv.Itoa(count))
	}

	resp, err := api.Query(URL_PUBLIC_ORDER_BOOK, params, false)
	if err != nil {
		return nil, err
	}

	out := make(map[string]PublicOrderBook)

	_, err = parse(resp, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
Input:

pair = asset pair to get trade data for
since = return trade data since given id (optional.  exclusive)
Result: array of pair name and recent trade data

<pair_name> = pair name
    array of array entries(<price>, <volume>, <time>, <buy/sell>, <market/limit>, <miscellaneous>)
last = id to be used as since when polling for new trade data
*/
func (api *KrakenApi) ApiTrades(pair string, since string) (map[string][]RecentTrade, float64, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if since != "" {
		params.Set("since", since)
	}

	resp, err := api.Query(URL_PUBLIC_RECENT_TRADES, params, false)
	if err != nil {
		return nil, 0, err
	}

	content, err := parse(resp, nil)
	if err != nil {
		return nil, 0, err
	}

	var last float64
	out := make(map[string][]RecentTrade)

	for key, value := range content.(map[string]interface{}) {
		if key == "last" {
			last, err = strconv.ParseFloat(value.(string), 64)
			if err != nil {
				return nil, 0, err
			}
			break
		}

		trades := make([]RecentTrade, 0)

		for _, subvalue := range value.([]interface{}) {
			values := subvalue.([]interface{})

			price, err := strconv.ParseFloat(values[0].(string), 64)
			if err != nil {
				return nil, 0, err
			}

			volume, err := strconv.ParseFloat(values[1].(string), 64)
			if err != nil {
				return nil, 0, err
			}

			trades = append(trades, RecentTrade{
				Price:     price,
				Volume:    volume,
				Time:      values[2].(float64),
				Type:      values[3].(string),
				TradeType: values[4].(string),
				Misc:      values[5].(string),
			})
		}

		out[key] = trades
	}

	return out, last, err
}

/*
Input:

pair = asset pair to get spread data for
since = return spread data since given id (optional.  inclusive)
Result: array of pair name and recent spread data

<pair_name> = pair name
    array of array entries(<time>, <bid>, <ask>)
last = id to be used as since when polling for new spread data
Note: "since" is inclusive so any returned data with the same time as the previous set should overwrite all of the previous set's entries at that time
*/
func (api *KrakenApi) ApiSpread(pair string, since string) (map[string][]Spread, float64, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if since != "" {
		params.Set("since", since)
	}

	resp, err := api.Query(URL_PUBLIC_SPREAD, params, false)
	if err != nil {
		return nil, 0, err
	}

	content, err := parse(resp, nil)
	if err != nil {
		return nil, 0, err
	}

	var last float64
	out := make(map[string][]Spread)

	for key, value := range content.(map[string]interface{}) {
		if key == "last" {
			last = value.(float64)
			continue
		}

		spreads := make([]Spread, 0)

		for _, subvalue := range value.([]interface{}) {
			values := subvalue.([]interface{})

			ask, err := strconv.ParseFloat(values[1].(string), 64)
			if err != nil {
				return nil, 0, err
			}

			bid, err := strconv.ParseFloat(values[2].(string), 64)
			if err != nil {
				return nil, 0, err
			}

			spreads = append(spreads, Spread{
				Time: values[0].(float64),
				Ask:  ask,
				Bid:  bid,
			})
		}

		out[key] = spreads
	}

	return out, last, nil
}
