package main

import (
	"fmt"
	"net/http"
	"log"
	"math"
	"encoding/json"
	"github.com/gorilla/mux"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
	"time"
)


type  Items struct{
    ShortDescription string `json:"shortDescription"` 
	Price string `json:"price"`
}
type Receipt struct{
	Retailer string `json:"retailer"` 
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items[] Items `json:"items"`
	Total string `json:"total"`
}
type  ID struct{
    ID string `json:"id"` 
	
}
type  Points struct{
     Points int `json:"points"`
}

/*
	store: An in-memory key-val store
*/
var store = make(map[string]int)

/*
	add(*): add key as ID and points as val in the store
*/
func add(k string, b int) bool {
	if k == "" {
		return false
	}
	store[k] = b
	return true;
}

/*
	lookUp(*): lookUp key as ID and returns points as val 
*/
func lookUp(k string) int {
	_, ok := store[k]
	if ok {
		b := store[k]
		return b
	}
	return -1
}

/*
	getID(*): Returns Points by ID specification
*/

func getID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	var payload Receipt
	err = json.Unmarshal(body, &payload)

	if err != nil {
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	hash := sha256.New()
	hash.Write(body)
	hashID := hex.EncodeToString(hash.Sum(nil))
	
    points := getPoints(payload.Retailer,payload.Total,payload.Items,payload.PurchaseDate,payload.PurchaseTime)

	add(hashID,points)
	id := ID{
        ID: "",
    }
	id.ID = hashID
	json.NewEncoder(w).Encode(id)

}

/*
	getPointsByID(*): Returns Points by ID specification
*/
func getPointsByID(w http.ResponseWriter, r *http.Request){
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	w.Header().Set("Content-Type", "text/plain")
	points := Points{
        Points: -1,
    }

	points.Points = lookUp(id)
	
	json.NewEncoder(w).Encode(points)
}

/*
		getPoints(*): Returns Points
*/
func getPoints(Retailer string, total string, Items []Items, PurchaseDate string,  PurchaseTime string) int {
	total_int, err := strconv.ParseFloat(total, 64) 
    if err != nil {
        fmt.Println(err)
        return 0
    }

	/*
		One point for every alphanumeric character in the retailer name.
	*/
	sum := countAlphanumeric(Retailer)
	
	/*
		50 points if the total is a round dollar amount with no cents.
	*/
	if math.Mod((total_int*10.0),10)==0.0{
			sum+=50
	}

	/*
		25 points if the total is a multiple of 0.25.
	*/
	if math.Mod(total_int,0.25)==0.0{
		sum+=25
	}

	/*
		5 points for every two items on the receipt.
	*/
	sum=sum+len(Items)/2*5

	/*
		If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	*/
	for i := 0; i < len(Items); i++ {

			if (len(strings.TrimSpace(Items[i].ShortDescription))%3==0){
				price_int, err := strconv.ParseFloat(Items[i].Price, 64) 
				if err != nil {
				fmt.Println(err)
				return 0
				}
				sum=sum+int(math.Ceil(price_int*0.2))
			}
	}

	/*
		6 points if the day in the purchase date is odd.
	*/
		if (isDayOdd(PurchaseDate)){
			sum+=6
		}

	/*
		10 points if the time of purchase is after 2:00pm and before 4:00pm.
	*/
		if(isTimeRange(PurchaseTime)){
			sum+=10
		}
		
		return sum
}

/*
		countAlphanumeric(*): Returns the count of Alphanumeric values.
*/
func countAlphanumeric(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

/*
	isDayOdd(*): Returns bool value for the day is odd or not
*/
 func isDayOdd(dateStr string) bool {

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false
	}

	dayOfYear := date.Day()
	isOdd := dayOfYear%2 != 0

	if isOdd {
		return true
	} 
	return false
 }

/*
	isTimeRange(*): Returns bool value for the condition of  time of purchase is after 2:00pm and before 4:00pm.
*/
 func isTimeRange (timeStr string) bool {

	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		fmt.Println("Error parsing time - ", err)
		return false
	}

	milliseconds := t.Hour()*3600000 + t.Minute()*60000 + t.Second()*1000

	if(milliseconds>50400000 && milliseconds<57600000){
	    return true
	}
	return false
 }

func handleRequests(){
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/receipts/process",getID)
	myRouter.HandleFunc("/receipts/{id}/points",getPointsByID)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	
    handleRequests()
}