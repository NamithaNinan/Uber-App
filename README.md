# Cmpe273-Assignment3
To implement a trip planner using uber api

##Working

This is a basic simple trip planning application, built as a continuation to Cmpe273-Assignment2. In this application, we use the same database that contains the places entered in Cmpe273-Assignment2.

1) POST http://localhost:8084/trips

Here using the POST command, we try creating finding an optimum path depending on the shortest distance. The user provides the id of the starting location and the id's of the places to be travelled

{

    "starting_from_location_id" : "5629da8f18683bb841ef075d",
    
    "location_ids" : ["5649518d66bae33b8883d654","5649511466bae33b8883d653","5649522c66bae33b8883d655","5649527e66bae33b8883d656"]
    
} 

The application makes use of Uber api to provide the total cost and the time taken in starting from the starting location, making a round trip.

{

  "id": "5666373d66bae316887f04ac",
  
  "Status": "processing",
  
  "Starting_from_location_id": "5629da8f18683bb841ef075d",
  
  "Best_route_location_ids": [
  
    "5649511466bae33b8883d653",
    
    "5649527e66bae33b8883d656",
    
    "5649522c66bae33b8883d655",
    
    "5649518d66bae33b8883d654",
    
    ""
  ],
  
  "Total_uber_costs": 162,
  
  "Total_uber_duration": 12279,
  
  "Total_distance": 110.55000000000001
  
}

 2) GET  http://localhost:8084/trips/5666373d66bae316887f04ac
 
 As each record in the databse is stored with a unique Id, this function only retrieves the record on providing the Id in the link
 
 3) PUT http://localhost:8084/trips/5666373d66bae316887f04ac/request
 
 Starts a call to uber from each destination till the end destination
 
 {
  "Id": "5666373d66bae316887f04ac",
  "Status": "processing",
  "Starting_from_location_id": "5629da8f18683bb841ef075d",
  "Next_destination_location_id": "5649511466bae33b8883d653",
  "Best_route_location_ids": [
    "5649511466bae33b8883d653",
    "5649527e66bae33b8883d656",
    "5649522c66bae33b8883d655",
    "5649518d66bae33b8883d654",
    ""
  ],
  "Total_uber_costs": 162,
  "Total_uber_duration": 12279,
  "Total_distance": 110.55000000000001,
  "Uber_wait_time_eta": 2
}

{
  "Id": "5666373d66bae316887f04ac",
  "Status": "processing",
  "Starting_from_location_id": "5629da8f18683bb841ef075d",
  "Next_destination_location_id": "5649527e66bae33b8883d656",
  "Best_route_location_ids": [
    "5649511466bae33b8883d653",
    "5649527e66bae33b8883d656",
    "5649522c66bae33b8883d655",
    "5649518d66bae33b8883d654",
    ""
  ],
  "Total_uber_costs": 162,
  "Total_uber_duration": 12279,
  "Total_distance": 110.55000000000001,
  "Uber_wait_time_eta": 2
}

{
  "Id": "5666373d66bae316887f04ac",
  "Status": "processing",
  "Starting_from_location_id": "5629da8f18683bb841ef075d",
  "Next_destination_location_id": "5649522c66bae33b8883d655",
  "Best_route_location_ids": [
    "5649511466bae33b8883d653",
    "5649527e66bae33b8883d656",
    "5649522c66bae33b8883d655",
    "5649518d66bae33b8883d654",
    ""
  ],
  "Total_uber_costs": 162,
  "Total_uber_duration": 12279,
  "Total_distance": 110.55000000000001,
  "Uber_wait_time_eta": 2
}

{
  "Id": "5666373d66bae316887f04ac",
  "Status": "completed",
  "Starting_from_location_id": "5629da8f18683bb841ef075d",
  "Next_destination_location_id": "5649518d66bae33b8883d654",
  "Best_route_location_ids": [
    "5649511466bae33b8883d653",
    "5649527e66bae33b8883d656",
    "5649522c66bae33b8883d655",
    "5649518d66bae33b8883d654",
    ""
  ],
  "Total_uber_costs": 162,
  "Total_uber_duration": 12279,
  "Total_distance": 110.55000000000001,
  "Uber_wait_time_eta": 2
}





