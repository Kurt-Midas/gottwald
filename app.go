package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./landing")))
	http.HandleFunc("/tlf", rustIsBest)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./bower_components"))))
	fmt.Println("Listening...")
	http.ListenAndServe(":8081", nil)
}

type orderInfo struct {
	Title   string
	LpItems []lpItemInfo
}

type lpItemInfo struct {
	ID            int         `json:"id"`
	Quantity      int         `json:"quantity"`
	LpCost        int         `json:"lpCost"`
	IskCost       float64     `json:"iskCost"`
	RequiredItems []inputItem `json:"requiredItems"`
	Item          typeDetails `json:"item"`
	Ratio         float64     `json:"ratio"`
	Wanted        int         `json:"wanted"`
}

type inputItem struct {
	Quantity int         `json:"quantity"`
	Item     typeDetails `json:"item"`
}

type typeDetails struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func parseLpJsonToStruct(js string) lpItemInfo {
	// fmt.Printf("jsonifying :: %v\n", js)
	obj := new(lpItemInfo)
	jsonErr := json.Unmarshal([]byte(js), obj)
	if jsonErr != nil {
		fmt.Printf("Got json error %v at string %v\n\n", jsonErr, js)
		return lpItemInfo{}
	}
	return *obj
}

/*type lpItemInfo struct {
	itemInfo
	LpCost  int
	IskCost float32
	Inputs  []itemInfo
	Ratio   float32
}

type itemInfo struct {
	TypeID   string
	Quantity int
	Price    float32
}*/

/*priceMap := map[int]float64 {
	15614 : 100000.0,
	15612 : 50000.0,
	3829 : 55000.00,
	15615 : 500000.00,
	184 : 24.10,
	587 : 499000.00,
	31924 : 5324999.39,
	3554 : 1685.99,
	9899 : 10999999.99,
	20171 : 100000.00,
	21924 : 589.00,
	17815 : 450891.30,
	264 : 639.95,
	31982 : 39998.46,
	181 : 14.30,
	28336 : 650.00,
	17812 : 11249993.50,
	377 : 30000.00,
	31928 : 11879959.36,
	33332 : 19998.51,
	31990 : 56000.00,
	11283 : 3100.00
}*/

func ratioAndStuff(item lpItemInfo, priceMap map[int]float64) lpItemInfo {
	var revenue = priceMap[item.Item.ID] * float64(item.Quantity)
	item.Item.Price = priceMap[item.Item.ID]
	var cost = 0.0
	cost += item.IskCost
	for i := range item.RequiredItems {
		cost += priceMap[item.RequiredItems[i].Item.ID] * float64(item.RequiredItems[i].Quantity)
		item.RequiredItems[i].Item.Price = priceMap[item.RequiredItems[i].Item.ID]
	}
	var profit = revenue - cost
	item.Ratio = profit / float64(item.LpCost)
	// fmt.Printf("Ratio :: %v\n", item.Ratio)
	return item
}

func rustIsBest(w http.ResponseWriter, r *http.Request) {
	var data = orderInfo{}
	data.Title = "templated stuff"
	list := make([]lpItemInfo, 10)
	list[0] = parseLpJsonToStruct(`{"iskCost": 125000, "lpCost_str": "125", "iskCost_str": "125000", "requiredItems": [{ "item": { "id_str": "264", "href": "https://crest-tq.eveonline.com/inventory/types/264/", "id": 264, "name": "Cap Booster 50" }, "quantity_str": "20", "quantity": 20 }], "lpCost": 125, "item": { "id_str": "33332", "href": "https://crest-tq.eveonline.com/inventory/types/33332/", "id": 33332, "name": "Navy Cap Booster 50" }, "id_str": "15803", "quantity_str": "20", "id": 15803, "quantity": 20 }`)
	list[1] = parseLpJsonToStruct(`{"iskCost": 250000, "lpCost_str": "250", "iskCost_str": "250000", "requiredItems": [{"item": {"id_str": "3554", "href": "https://crest-tq.eveonline.com/inventory/types/3554/", "id": 3554, "name": "Cap Booster 100"}, "quantity_str": "20", "quantity": 20}], "lpCost": 250, "item": {"id_str": "31982", "href": "https://crest-tq.eveonline.com/inventory/types/31982/", "id": 31982, "name": "Navy Cap Booster 100"}, "id_str": "14693", "quantity_str": "20", "id": 14693, "quantity": 20}`)
	list[2] = parseLpJsonToStruct(`{"iskCost": 375000, "lpCost_str": "375", "iskCost_str": "375000", "requiredItems": [{"item": {"id_str": "11283", "href": "https://crest-tq.eveonline.com/inventory/types/11283/", "id": 11283, "name": "Cap Booster 150"}, "quantity_str": "20", "quantity": 20}], "lpCost": 375, "item": {"id_str": "31990", "href": "https://crest-tq.eveonline.com/inventory/types/31990/", "id": 31990, "name": "Navy Cap Booster 150"}, "id_str": "14694", "quantity_str": "20", "id": 14694, "quantity": 20}`)
	list[3] = parseLpJsonToStruct(`{"iskCost": 1200000, "lpCost_str": "1200", "iskCost_str": "1200000", "requiredItems": [{"item": {"id_str": "181", "href": "https://crest-tq.eveonline.com/inventory/types/181/", "id": 181, "name": "Depleted Uranium S"}, "quantity_str": "5000", "quantity": 5000}], "lpCost": 1200, "item": {"id_str": "28336", "href": "https://crest-tq.eveonline.com/inventory/types/28336/", "id": 28336, "name": "Republic Fleet Depleted Uranium S"}, "id_str": "4623", "quantity_str": "5000", "id": 4623, "quantity": 5000}`)
	list[4] = parseLpJsonToStruct(`{"iskCost": 1200000, "lpCost_str": "1200", "iskCost_str": "1200000", "requiredItems": [{"item": {"id_str": "184", "href": "https://crest-tq.eveonline.com/inventory/types/184/", "id": 184, "name": "Phased Plasma S"}, "quantity_str": "5000", "quantity": 5000}], "lpCost": 1200, "item": {"id_str": "21924", "href": "https://crest-tq.eveonline.com/inventory/types/21924/", "id": 21924, "name": "Republic Fleet Phased Plasma S"}, "id_str": "3774", "quantity_str": "5000", "id": 3774, "quantity": 5000}`)
	list[5] = parseLpJsonToStruct(`{"iskCost": 0, "lpCost_str": "10000", "iskCost_str": "0", "requiredItems": [{"item": {"id_str": "587", "href": "https://crest-tq.eveonline.com/inventory/types/587/", "id": 587, "name": "Rifter"}, "quantity_str": "1", "quantity": 1}, {"item": {"id_str": "17815", "href": "https://crest-tq.eveonline.com/inventory/types/17815/", "id": 17815, "name": "Minmatar UUA Nexus Chip"}, "quantity_str": "1", "quantity": 1}], "lpCost": 10000, "item": {"id_str": "17812", "href": "https://crest-tq.eveonline.com/inventory/types/17812/", "id": 17812, "name": "Republic Fleet Firetail"}, "id_str": "14809", "quantity_str": "1", "id": 14809, "quantity": 1}`)
	list[6] = parseLpJsonToStruct(`{"iskCost": 500000, "lpCost_str": "1000", "iskCost_str": "500000", "requiredItems": [{"item": {"id_str": "377", "href": "https://crest-tq.eveonline.com/inventory/types/377/", "id": 377, "name": "Small Shield Extender I"}, "quantity_str": "1", "quantity": 1}, {"item": {"id_str": "15614", "href": "https://crest-tq.eveonline.com/inventory/types/15614/", "id": 15614, "name": "Imperial Navy Colonel Insignia I"}, "quantity_str": "4", "quantity": 4}, {"item": {"id_str": "15612", "href": "https://crest-tq.eveonline.com/inventory/types/15612/", "id": 15612, "name": "Imperial Navy Captain Insignia I"}, "quantity_str": "2", "quantity": 2}], "lpCost": 1000, "item": {"id_str": "31924", "href": "https://crest-tq.eveonline.com/inventory/types/31924/", "id": 31924, "name": "Republic Fleet Small Shield Extender"}, "id_str": "14680", "quantity_str": "1", "id": 14680, "quantity": 1}`)
	list[7] = parseLpJsonToStruct(`{"iskCost": 2000000, "lpCost_str": "3000", "iskCost_str": "2000000", "requiredItems": [{"item": {"id_str": "3829", "href": "https://crest-tq.eveonline.com/inventory/types/3829/", "id": 3829, "name": "Medium Shield Extender I"}, "quantity_str": "1", "quantity": 1}, {"item": {"id_str": "15615", "href": "https://crest-tq.eveonline.com/inventory/types/15615/", "id": 15615, "name": "Imperial Navy General Insignia I"}, "quantity_str": "2", "quantity": 2}], "lpCost": 3000, "item": {"id_str": "31928", "href": "https://crest-tq.eveonline.com/inventory/types/31928/", "id": 31928, "name": "Republic Fleet Medium Shield Extender"}, "id_str": "14681", "quantity_str": "1", "id": 14681, "quantity": 1}`)
	list[8] = parseLpJsonToStruct(`{"iskCost": 5250000, "lpCost_str": "5250", "iskCost_str": "5250000", "requiredItems": [], "lpCost": 5250, "item": {"id_str": "9899", "href": "https://crest-tq.eveonline.com/inventory/types/9899/", "id": 9899, "name": "Ocular Filter - Basic"}, "id_str": "3433", "quantity_str": "1", "id": 3433, "quantity": 1}`)
	list[9] = parseLpJsonToStruct(`{"iskCost": 250000, "lpCost_str": "250", "iskCost_str": "250000", "requiredItems": [], "lpCost": 250, "item": {"id_str": "20171", "href": "https://crest-tq.eveonline.com/inventory/types/20171/", "id": 20171, "name": "Datacore - Hydromagnetic Physics"}, "id_str": "15753", "quantity_str": "5", "id": 15753, "quantity": 5}`)

	// itemMap := make(map[int]string)
	// for _, item := range list {
	// 	itemMap[item.Item.ID] = item.Item.Name
	// 	for _, in := range item.RequiredItems {
	// 		itemMap[in.Item.ID] = in.Item.Name
	// 	}
	// }
	priceMap := map[int]float64{
		15614: 100000.0, 15612: 50000.0, 3829: 55000.00, 15615: 500000.00,
		184: 24.10, 587: 499000.00, 31924: 5324999.39, 3554: 1685.99, 9899: 10999999.99,
		20171: 100000.00, 21924: 589.00, 17815: 450891.30, 264: 639.95,
		31982: 39998.46, 181: 14.30, 28336: 650.00, 17812: 11249993.50,
		377: 30000.00, 31928: 11879959.36, 33332: 19998.51, 31990: 56000.00, 11283: 3100.00}
	for i := range list {
		list[i] = ratioAndStuff(list[i], priceMap)
	}
	data.LpItems = list
	fmt.Printf("Item Map :: %v\n", priceMap)
	serveOrderPage(w, r, data)
}

func serveOrderPage(w http.ResponseWriter, r *http.Request, data orderInfo) {
	t, err := template.New("index.html").Delims("[[", "]]").ParseFiles("./order/index.html")
	// t, err := template.ParseFiles("./order/index.html")
	if err != nil {
		panic(err)
	}
	fmt.Printf("serveOrderPage %v\n", t.Name())
	// t.ExecuteTemplate(w, data.Title, data)
	t.Execute(w, data)
}
