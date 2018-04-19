package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"encoding/json"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		log.Print("arguments: <dbAddress> <campaignID> <numRequests> \n EG: Flaconi12 10")
		os.Exit(1)
	}

	dbAddr := os.Args[1]
	campaignID := os.Args[2]
	numRequests, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal("wrong numRequests: %v", numRequests)
	}

	s, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatalf("connection to '%s' error: %s", dbAddr, err)
	}

	for i := 0; i < numRequests ; i++{
		start := time.Now()
		coupon, err := GetNextCoupon(s, campaignID)
		if err != nil {
			log.Fatalf("get next coupon error: %v", err)
		} else {
			elapsedTime := time.Now().Sub(start)
			jsonStr, _ := json.Marshal(coupon)
			log.Printf("took %v to get result: %v", elapsedTime, string(jsonStr))
		}
	}
}

type Coupon struct {
	ID         string `json:"id" bson:"_id"`
	Code       string `json:"code" bson:"code"`
	CampaignID string `json:"campaign_id" bson:"campaign_id"`
	Fetched    bool   `json:"fetched" bson:"fetched"`
}


func GetNextCoupon(session *mgo.Session, campaignCode string) (*Coupon, error) {
	query := bson.M{"campaign_id": campaignCode, "fetched": false}
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"fetched": true}},
		ReturnNew: false,
	}
	coupon := &Coupon{}

	ss := session.Copy()
	defer ss.Close()

	collection := ss.DB("campaigndb").C("coupons")
	info, err := collection.Find(query).Limit(1).Apply(change, &coupon)
	if info != nil && (info.Updated == 1 || err == nil) {
		return coupon, nil
	}
	return nil, err
}