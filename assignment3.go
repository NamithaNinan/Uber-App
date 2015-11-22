package main
import (
    "encoding/json"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/binary"
    "fmt"
    "bytes"
    "strings"
    "time"
    "gopkg.in/mgo.v2"
    "os"
    //"os/exec"
    "gopkg.in/mgo.v2/bson"
    //"requests"
)


var msg,loc1,loc2,loc3,loc4 outputMongo
var msg2 outputMongo2
var msg3 outputMongo3
var collection *mgo.Collection
var t Tripinput
var x Tripoutput
var etaoutput Tripoutputeta 
var price Uberprice
var route1,route2,route3,route4 int
var route [5][5]float64
var min,dist float64
var mindist float64=0
var Lat [5]float64
var Long [5]float64
var b1 []byte
var b2 []byte
var b []byte
var optimum [4]float64
var dur[5][5] float64
var timex float64=0
var cost[5][5] float64
var tcost float64=0
var prod[5][5] string
var prodid string
var stat[5][5] string
var s Statusfind
var si StatusInput

type Uberprice struct {
  Prices []struct {
    CurrencyCode    string  `json:"currency_code"`
    DisplayName     string  `json:"display_name"`
    Distance        float64 `json:"distance"`
    Duration        float64     `json:"duration"`
    Estimate        string  `json:"estimate"`
    HighEstimate    float64     `json:"high_estimate"`
    LowEstimate     float64     `json:"low_estimate"`
    ProductID       string  `json:"product_id"`
    SurgeMultiplier float64     `json:"surge_multiplier"`
  } `json:"prices"`
}

type Tripinput struct {
    Starting_from_location_id string
    Location_ids [5]string
}

type Tripoutput struct {
     Id string
     Status string
     Starting_from_location_id string
     Best_route_location_ids [5]string
     Total_uber_costs float64
     Total_uber_duration float64
     Total_distance float64
}

type Tripoutputeta struct {
     Id bson.ObjectId
     Status string
     Starting_from_location_id string
     Next_destination_location_id string
     Best_route_location_ids [5]string
     Total_uber_costs float64
     Total_uber_duration float64
     Total_distance float64
     Uber_wait_time_eta float64
}

type outputMongo struct {
    ID  bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Status  string `json:"name"`
    Address     string `json:"address"`
    City        string `json:"city"`
    State string `json:"state"`
    Zip   string `json:"zip"`
    Coordinates struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    } `json:"coordinates"`
       
}

type outputMongo2 struct{
    ID  bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Status string
     Starting_from_location_id string
     Best_route_location_ids [5]string
     Total_uber_costs float64
     Total_uber_duration float64
     Total_distance float64
}

type outputMongo3 struct{
    ID  bson.ObjectId `json:"id" bson:"_id,omitempty"`
    obj_id string
    Status string
     Starting_from_location_id string
     Next_destination_location_id string
     Best_route_location_ids [5]string
     Total_uber_costs float64
     Total_uber_duration float64
     Total_distance float64
     Uber_wait_time_eta float64
}

type Statusfind struct {
    Driver          interface{} `json:"driver"`
    Eta             int         `json:"eta"`
    Location        interface{} `json:"location"`
    RequestID       string      `json:"request_id"`
    Status          string      `json:"status"`
    SurgeMultiplier int         `json:"surge_multiplier"`
    Vehicle         interface{} `json:"vehicle"`
}


type StatusInput struct {
    product_id string
    start_latitude float64
    start_longitude float64
    end_latitude float64
    end_longitude float64
}

func connectdb(){
    uri := "mongodb://dbuser:dbuser@ds053090.mongolab.com:53090/location"
    if uri == "" {
    fmt.Println("no connection string provided")
    os.Exit(1)
  }
 
  sess, err := mgo.Dial(uri)
  if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
  }
  sess.SetSafe(&mgo.Safe{})
  
  collection = sess.DB("location").C("address")
}

func prettyprint(b []byte) ([]byte, error) {
    var out bytes.Buffer
    err := json.Indent(&out, b, "", "  ")
    return out.Bytes(), err
}

func statuscalc(product string,first int,second int,startlong float64,startlati float64,endlong float64, endlati float64) {
    si.product_id=product
    si.start_latitude=startlong
    si.start_longitude=startlati
    si.end_latitude=endlong
    si.end_longitude=endlati
    
    /*b1,err:=json.Marshal(si)   
    buf:=new(bytes.Buffer)
    err=binary.Write(buf,binary.BigEndian,&b1)
*/
    b1,_=json.Marshal(si) 
    bx1, _ := prettyprint(b1)  
    url1:=fmt.Sprintf("https://sandbox-api.uber.com/v1/requests")
    timeout := time.Duration(15 * time.Second)
        client := http.Client{Timeout: timeout}
        req,_:=http.NewRequest("POST",url1,bytes.NewBuffer(bx1))
        req.Header.Set("Authorization","Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3QiLCJoaXN0b3J5Il0sInN1YiI6IjNmNjY5ODZiLWY1YTQtNDc2ZS1hZDM3LWE4NjMzODU2MTliYyIsImlzcyI6InViZXItdXMxIiwianRpIjoiNzdlMWE0Y2ItNDhmMi00ODJkLTg1NGUtM2RlMjk5YzQ5Njc3IiwiZXhwIjoxNDUwNjQ1NDQwLCJpYXQiOjE0NDgwNTM0MzksInVhY3QiOiJ6MHRVRkcyTTFhb0lKSWczWnpRTHRhdGFVRThSOVEiLCJuYmYiOjE0NDgwNTMzNDksImF1ZCI6ImZLaElaWFpOQU9HRjdwamNmaGdWV1c1UmxPWTVpbm9xIn0.BtIHZDmKw3uQY509ISQczs7wa8syN0x8irmajVaPVito5yF7HOCP2ziXCxhY0CZ-wrQ8a7DaGNMAfXiPj7RxMyT0kQqTmwXa1fjWGaRkhuKT3MBGhaQaaxOgEWiPbSSZFEFgjtHpb9cY1l5VoD516-sYuTC4-g1hTXWQBUrtG5qn7B69ABUNMAPukEHe2Ho31Nga-JljY7AqnHzY1Z7EhAMYZncy038c9_XxBHLAeiGZJ-91ubI20l8fuIEe_vjUuvKCp25JtQEt53hRSSa4arnBtC756Ff_5Trrhi6zWCOYq2cvWoTD6cFM0yx6_SShBlm-brcLaCXprJO-pF0vbw")
        req.Header.Set("Content-Type","application/json")
        res,err:=client.Do(req)
        if err!=nil{
        log.Fatal("there is an error in uber api ",err)
        }
        decoder:=json.NewDecoder(res.Body)
        
        err=decoder.Decode(&s)
        stat[first][second]=s.Status
        fmt.Println(s.Status)
}

func distcalc(first int,second int,startlong float64,startlati float64,endlong float64, endlati float64) (float64){
        
        url := fmt.Sprintf("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=%f&start_longitude=%f&end_latitude=%f&end_longitude=%f",startlong,startlati,endlong,endlati)
        timeout := time.Duration(5 * time.Second)
        client := http.Client{Timeout: timeout}
        req,_:=http.NewRequest("GET",url,nil)
        req.Header.Set("Authorization","Token 1U9BtzwI5Luh1Ebc5ZTlUeA-h8n6Xj4V9BQt4rIV")
        res,err:=client.Do(req)
        if err!=nil{
        log.Fatal("there is an error in uber api")
        }
        decoder:=json.NewDecoder(res.Body)
        err=decoder.Decode(&price)
        dist=price.Prices[0].Distance
        dur[first][second]=price.Prices[0].Duration
        cost[first][second]=price.Prices[0].LowEstimate
        prod[first][second]=price.Prices[0].ProductID
        statuscalc(prod[first][second],first,second,startlong,startlati,endlong,endlati)
        //fmt.Println(prod[first][second],startlati,startlong,endlati,endlong)
        if err!=nil{
        log.Fatal("there is an error in getting price api", err)
        }
        return(dist)

}
func planner(rw http.ResponseWriter, req *http.Request) {
    
    if(req.Method=="POST"){//------------------------> POST---------------------------------->
        body, err := ioutil.ReadAll(req.Body)
        err = json.Unmarshal(body, &t)
        if err != nil {
          panic("error in unmarshalling")
        }
        start:=t.Starting_from_location_id
        x.Starting_from_location_id=t.Starting_from_location_id
        sid := bson.ObjectIdHex(start)
        err = collection.FindId(sid).One(&msg)
        if err != nil {
        fmt.Printf("Error in searching! %v\n", err)
        os.Exit(1)
       }

        st1:=t.Location_ids[0]
        st2:=t.Location_ids[1]
        st3:=t.Location_ids[2]
        st4:=t.Location_ids[3]
        id1 := bson.ObjectIdHex(st1)
        id2 := bson.ObjectIdHex(st2)
        id3 := bson.ObjectIdHex(st3)
        id4 := bson.ObjectIdHex(st4)
        err = collection.FindId(id1).One(&loc1)
        err = collection.FindId(id2).One(&loc2)
        err = collection.FindId(id3).One(&loc3)
        err = collection.FindId(id4).One(&loc4)
        Lat[0]=msg.Coordinates.Latitude
        Lat[1]=loc1.Coordinates.Latitude//for point t.Location_ids[0]
        Lat[2]=loc2.Coordinates.Latitude//for point t.Location_ids[1]
        Lat[3]=loc3.Coordinates.Latitude//for point t.Location_ids[2]
        Lat[4]=loc4.Coordinates.Latitude//for point t.Location_ids[3]
        Long[0]=msg.Coordinates.Longitude
        Long[1]=loc1.Coordinates.Longitude
        Long[2]=loc2.Coordinates.Longitude
        Long[3]=loc3.Coordinates.Longitude
        Long[4]=loc4.Coordinates.Longitude
        
        if err != nil {
          fmt.Printf("Error in searching! %v\n", err)
          os.Exit(1)
        }
        route[0][1]=distcalc(0,1,Long[0],Lat[0],Long[1],Lat[1])
        route[0][2]=distcalc(0,2,Long[0],Lat[0],Long[2],Lat[2])
        route[0][3]=distcalc(0,3,Long[0],Lat[0],Long[3],Lat[3])
        route[0][4]=distcalc(0,4,Long[0],Lat[0],Long[4],Lat[4])
        route[1][1]=distcalc(1,1,Long[1],Lat[1],Long[1],Lat[1])
        route[1][2]=distcalc(1,2,Long[1],Lat[1],Long[2],Lat[2])
        route[1][3]=distcalc(1,3,Long[1],Lat[1],Long[3],Lat[3])
        route[1][4]=distcalc(1,4,Long[1],Lat[1],Long[4],Lat[4])
        route[2][1]=distcalc(2,1,Long[2],Lat[2],Long[1],Lat[1])
        route[2][2]=distcalc(2,2,Long[2],Lat[2],Long[2],Lat[2])
        route[2][3]=distcalc(2,3,Long[2],Lat[2],Long[3],Lat[3])
        route[2][4]=distcalc(2,4,Long[2],Lat[2],Long[4],Lat[4])
        route[3][1]=distcalc(3,1,Long[3],Lat[3],Long[1],Lat[1])
        route[3][2]=distcalc(3,2,Long[3],Lat[3],Long[2],Lat[2])
        route[3][3]=distcalc(3,3,Long[3],Lat[3],Long[3],Lat[3])
        route[3][4]=distcalc(3,4,Long[3],Lat[3],Long[4],Lat[4])
        route[4][1]=distcalc(4,1,Long[4],Lat[3],Long[1],Lat[1])
        route[4][2]=distcalc(4,2,Long[4],Lat[3],Long[2],Lat[2])
        route[4][3]=distcalc(4,3,Long[4],Lat[3],Long[3],Lat[3])
        route[4][4]=distcalc(4,4,Long[4],Lat[3],Long[4],Lat[4])
        min=route[0][1]
        route1=1
        for i:=1;i<5;i++{
            if(route[0][i]<min){
                    min=route[0][i]
                    route1=i
        //            fmt.Println(route1)
            }          
        } 
        //fmt.Println(min,route1)  
        timex=timex+dur[0][route1]
        tcost=tcost+cost[0][route1]
        mindist=mindist+min
        min=999
        route2=1
        x.Best_route_location_ids[0]=t.Location_ids[route1-1]
        //fmt.Println(x.Best_route_location_ids[0])
        for i:=1;i<5;i++{
            if(route[route1][i]<min&&route1!=i){
                    min=route[route1][i]
                    route2=i
        //            fmt.Println(route2)
            }                      
        }   
        //fmt.Println(min,route2) 
        timex=timex+dur[route1][route2]
        tcost=tcost+cost[route1][route2]
        mindist=mindist+min
        x.Best_route_location_ids[1]=t.Location_ids[route2-1]
        //fmt.Println(x.Best_route_location_ids[1])
        min=999
        route3=1
        for i:=1;i<5;i++{
            if(route[route2][i]<min&&route2!=i&&i!=route1){
                    min=route[route2][i]
                    route3=i
          //          fmt.Println(route3)
            }                      
        }
        //fmt.Println(min,route3)  
        mindist=mindist+min
        timex=timex+dur[route2][route3]
        tcost=tcost+cost[route2][route3]
        route4=4
        x.Best_route_location_ids[2]=t.Location_ids[route3-1]
        //fmt.Println(x.Best_route_location_ids[2])
        min=999
        for i:=1;i<5;i++{
            if(route[route3][i]<min&&route3!=i&&i!=route1&&i!=route2){
                    min=route[0][i]
                    route4=i
          //          fmt.Println(route4)
            }                     
        }
        //fmt.Println(min,route4)  
        mindist=mindist+min
        timex=timex+dur[route3][route4]+dur[0][route4]
        tcost=tcost+cost[0][route4]+cost[route3][route4]
        x.Best_route_location_ids[3]=t.Location_ids[route4-1]
        //fmt.Println(x.Best_route_location_ids[3])
        mindist=mindist+route[0][route4]
        x.Total_distance=mindist
        x.Total_uber_duration=timex
        x.Total_uber_costs=tcost
        x.Status="processing"
        //fmt.Println(x.Total_distance,x.Total_uber_duration,x.Total_uber_costs)
        msg2=outputMongo2{bson.NewObjectId(),x.Status,x.Starting_from_location_id,x.Best_route_location_ids,x.Total_uber_costs,x.Total_uber_duration,x.Total_distance}
        //fmt.Println(msg2)
        b,_=json.Marshal(msg2)   
        bx, _ := prettyprint(b)
        n:=binary.Size(bx)
        s := string(bx[:n])
        fmt.Fprintf(rw,s)
        
        err = collection.Insert(msg2)
         if err != nil {
        fmt.Printf("Can't insert document: %v\n", err)
        }
    }else if(req.Method=="GET"){//------------------------> GET---------------------------------->
        s1:=req.URL.Path[1:]
        st1:=string(s1[6:])
        oid := bson.ObjectIdHex(st1)
  
        err := collection.FindId(oid).One(&msg2)
        if err != nil {
        fmt.Printf("Error in searching! %v\n", err)
        os.Exit(1)
         }
       
        b,_=json.Marshal(msg2)   
        bx, _ := prettyprint(b)
        n:=binary.Size(bx)
        s := string(bx[:n])
        fmt.Fprintf(rw,s)
         

    }else if(req.Method=="PUT"){//-------------------------> PUT --------------------------------->
        s1:=req.URL.Path[1:]
        st1:=string(s1[6:])
        req_id:=strings.Split(st1,"/")
        oid:=bson.ObjectIdHex(req_id[0])
        err := collection.FindId(oid).One(&msg2)
        if err != nil {
        fmt.Printf("Error in searching! %v\n", err)
        os.Exit(1)
        }
        i:=0
        for i=0;i<4;i++{
            
            etaoutput.Id=msg2.ID
            
            etaoutput.Starting_from_location_id=msg2.Starting_from_location_id
            if(i==3){
                etaoutput.Status="completed"
            }else{
                etaoutput.Status=msg2.Status    
            }
            etaoutput.Next_destination_location_id=msg2.Best_route_location_ids[i]
            etaoutput.Best_route_location_ids=msg2.Best_route_location_ids
            etaoutput.Total_uber_costs=msg2.Total_uber_costs
            etaoutput.Total_uber_duration=msg2.Total_uber_duration
            etaoutput.Total_distance=msg2.Total_distance
            etaoutput.Uber_wait_time_eta=2
            b,_=json.Marshal(etaoutput)   
            bx, _ := prettyprint(b)
            n:=binary.Size(bx)
            s := string(bx[:n])
            fmt.Fprintf(rw,s)
            msg3=outputMongo3{bson.NewObjectId(),req_id[0],etaoutput.Status,etaoutput.Starting_from_location_id,etaoutput.Next_destination_location_id,etaoutput.Best_route_location_ids,etaoutput.Total_uber_costs,etaoutput.Total_uber_duration,etaoutput.Total_distance,etaoutput.Uber_wait_time_eta}
            err = collection.Insert(msg3)
            if err != nil {
            fmt.Printf("Can't insert document: %v\n", err)
            }
        }
        
    }
}

func main() {
   
    http.HandleFunc("/trips", planner)
    http.HandleFunc("/trips/", planner)
    connectdb()
    log.Fatal(http.ListenAndServe(":8084", nil))
}
/*product_id 3ab64887-4842-4c8e-9780-ccecd3a0391d
client_id 1cb4zEWGGF-ui8ZcfwLfewSpFRjhuaff
code=lf7S2rMow95WiEBF5nGv0WixQTRTCY
{
    "last_authenticated": 1448011546,
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3QiLCJoaXN0b3J5Il0sInN1YiI6IjNmNjY5ODZiLWY1YTQtNDc2ZS1hZDM3LWE4NjMzODU2MTliYyIsImlzcyI6InViZXItdXMxIiwianRpIjoiYTE0ODkzYTQtNjdiMy00NGU0LTk1NGUtYzY5MjVjNDhjNWM0IiwiZXhwIjoxNDUwNjAzNzE4LCJpYXQiOjE0NDgwMTE3MTcsInVhY3QiOiJtaUszM2hvRkNrZ2NlemVpclZZTUFmT1M4Q3ZZdkciLCJuYmYiOjE0NDgwMTE2MjcsImF1ZCI6IjFjYjR6RVdHR0YtdWk4WmNmd0xmZXdTcEZSamh1YWZmIn0.TrNs-yJdmFnHia2dpsOF6mHGcbHet7WRREgww3MzU_PT3yaXUOce298K9ijrssEeHYG10W-oO-KxFJHsexN0xRtb6RXLz12QFD_glauUK6WE4RPG2CFgFVNvreyVuhMVA1ClCOePZ4oGMezi6mbpTvR_h0V40BNJykmSKK6YxEyQEthLW12rMetgTi1oWskFqezOyytieIgyf83kMUb7OL1nb04zse_wDXxnId9I0W0Lz3x6pYMB23JMKZqav-HnB4n5EpCQJ-ZpTBdmVbVU3huRa-kXQLgCgImm8o_HEaUNypl_LjsfdbCGtHj5vMrQE3zJnLZR8hyUR946xwUgPA",
    "expires_in": 2592000,
    "token_type": "Bearer",
    "scope": "profile request history",
    "refresh_token": "2y16WHgMliI2XH0q5ZDA0xn7J8JJZr"
}
 eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3QiLCJoaXN0b3J5Il0sInN1YiI6IjNmNjY5ODZiLWY1YTQtNDc2ZS1hZDM3LWE4NjMzODU2MTliYyIsImlzcyI6InViZXItdXMxIiwianRpIjoiNzdlMWE0Y2ItNDhmMi00ODJkLTg1NGUtM2RlMjk5YzQ5Njc3IiwiZXhwIjoxNDUwNjQ1NDQwLCJpYXQiOjE0NDgwNTM0MzksInVhY3QiOiJ6MHRVRkcyTTFhb0lKSWczWnpRTHRhdGFVRThSOVEiLCJuYmYiOjE0NDgwNTMzNDksImF1ZCI6ImZLaElaWFpOQU9HRjdwamNmaGdWV1c1UmxPWTVpbm9xIn0.BtIHZDmKw3uQY509ISQczs7wa8syN0x8irmajVaPVito5yF7HOCP2ziXCxhY0CZ-wrQ8a7DaGNMAfXiPj7RxMyT0kQqTmwXa1fjWGaRkhuKT3MBGhaQaaxOgEWiPbSSZFEFgjtHpb9cY1l5VoD516-sYuTC4-g1hTXWQBUrtG5qn7B69ABUNMAPukEHe2Ho31Nga-JljY7AqnHzY1Z7EhAMYZncy038c9_XxBHLAeiGZJ-91ubI20l8fuIEe_vjUuvKCp25JtQEt53hRSSa4arnBtC756Ff_5Trrhi6zWCOYq2cvWoTD6cFM0yx6_SShBlm-brcLaCXprJO-pF0vbw
INPUT :
{
    "starting_from_location_id" : "5629da8f18683bb841ef075d",
    "location_ids" : ["5649518d66bae33b8883d654","5649511466bae33b8883d653","5649522c66bae33b8883d655","5649527e66bae33b8883d656"]
} 

*/