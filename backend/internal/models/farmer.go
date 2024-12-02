package models

type Farmer struct {
    ID       string   `json:"id" dynamodbav:"id"`
    Name     string   `json:"name" dynamodbav:"name"`
    Contact  string   `json:"contact" dynamodbav:"contact"`
    State    string   `json:"state" dynamodbav:"state"`
    District string   `json:"district" dynamodbav:"district"`
    Tehsil   string   `json:"tehsil" dynamodbav:"tehsil"`
    Village  string   `json:"village" dynamodbav:"village"`
    Pincode  string   `json:"pincode" dynamodbav:"pincode"`
    Address  string   `json:"address" dynamodbav:"address"`
    Tag      string   `json:"tag" dynamodbav:"tag"`
    Crop     []string `json:"crop" dynamodbav:"crop"`
}
