package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BidsAsks struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}
type Results struct {
	SumBid float64
	SumAsk float64
	MaxBid string
	MinAsk string
	AvgBid float64
	AvgAsk float64
}

func main() {
	res := &Results{}
	bidAndAskAnalyzer(res)
}
func parseBidsAndAsks(bidsAsks *BidsAsks) error {
	urls := "https://api.binance.com/api/v3/depth?symbol=BTCUSDT&limit=20"
	resp, err := http.Get(urls)
	if err != nil {
		log.Println("error is coming from request from binance", err)
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error is coming from ioutil reader", err)
		return err
	}
	json.Unmarshal(b, bidsAsks)
	bidsAsks.Bids = bidsAsks.Bids[:15]
	bidsAsks.Asks = bidsAsks.Asks[:15]
	return nil
}
func bidAndAskAnalyzer(res *Results) {
	start := time.Now()
	bidsAsks := &BidsAsks{}
	err := parseBidsAndAsks(bidsAsks)
	if err != nil {
		log.Println("error from parse bids and asks", err)
		return
	}
	var err1, err2 error
	res.SumAsk, res.AvgAsk, err1 = sumAndAvg(bidsAsks.Asks)
	res.SumBid, res.AvgBid, err2 = sumAndAvg(bidsAsks.Asks)
	if err1 != nil || err2 != nil {
		log.Println("error from sumAndAvg", err)
		return
	}
	sortBidsAndAsksSlice(bidsAsks.Bids)
	sortBidsAndAsksSlice(bidsAsks.Asks)
	res.MaxBid = bidsAsks.Bids[len(bidsAsks.Bids)-1][0]
	res.MinAsk = bidsAsks.Asks[0][0]
	printer(res)
	if time.Now().Sub(start).Seconds() < 1 {
		time.Sleep(1 * time.Second)
	}
	bidAndAskAnalyzer(res)
}
func printer(res *Results) {
	fmt.Println("Sum of Bids: ", res.SumBid)
	fmt.Println("Sum of Asks:", res.SumAsk)
	fmt.Println("Max in Bids:", res.MaxBid)
	fmt.Println("Min in Asks:", res.MinAsk)
	fmt.Println("Average of Bids:", res.AvgBid)
	fmt.Println("Average of Asks:", res.AvgAsk)
	fmt.Println()
}
func sumAndAvg(v [][]string) (float64, float64, error) {
	res := 0.0
	for i := 0; i < len(v); i++ {
		num, err := strconv.ParseFloat(v[i][0], 64)
		if err != nil {
			return 0.0, 0.0, err
		}
		res += num
	}
	return res, res / float64(len(v)), nil
}
func sortBidsAndAsksSlice(v [][]string) {
	for i := 0; i < len(v); i++ {
		for j := i; j < len(v); j++ {
			if v[i][0] > v[j][0] {
				v[i], v[j] = v[j], v[i]
			}
		}
	}
}
