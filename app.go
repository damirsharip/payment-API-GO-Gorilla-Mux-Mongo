package main

import(
	"github.com/damirsharip/gorilla/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
)


var dbClient *mongo.Client
var db *mongo.Database

func AllPaymentsEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var payments []models.Payment

	collection := db.Collection("payments")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var payment models.Payment
		cursor.Decode(&payment)
		payments = append(payments, payment)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(payments)
}

func FindPaymentEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	params := mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	var payment models.Payment
	collection := db.Collection("payments")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err := collection.FindOne(ctx, models.Payment{ID: id}).Decode(&payment)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "error": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func CreatePaymentEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var payment models.Payment
	_ = json.NewDecoder(r.Body).Decode(&payment)
	collection := db.Collection("payments")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, payment)
	json.NewEncoder(w).Encode(result)
}

func UpdatePaymentEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	var payment models.Payment
	_ = json.NewDecoder(r.Body).Decode(&payment)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := db.Collection("payments")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"name", payment.Name}, 
			{"Typeof", payment.Typeof},
			{"comment", payment.Comment},
			{"price", payment.Price},
			{"time", time.Now},
			{"category", payment.Category}}}}
	collection.FindOneAndUpdate(ctx, filter, update) 
}
func DeletePaymentEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := db.Collection("payments")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	filter := bson.D{{"_id", id}}
	collection.FindOneAndDelete(ctx, filter); 
	}


func main(){
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().ApplyURI("mongodb+srv://user:user@cluster0.1iesq.mongodb.net")
	dbClient, _ = mongo.Connect(ctx, clientOptions)
	db = dbClient.Database("development_db")

	r := mux.NewRouter()

	r.HandleFunc("/api/payments", AllPaymentsEndPoint).Methods("GET")
	r.HandleFunc("/api/payments", CreatePaymentEndPoint).Methods("POST")
	r.HandleFunc("/api/payments/{id}", UpdatePaymentEndPoint).Methods("PUT")
	r.HandleFunc("/api/payments/{id}", DeletePaymentEndPoint).Methods("DELETE")
	r.HandleFunc("/api/payments/{id}", FindPaymentEndpoint).Methods("GET")

	fmt.Println("Backend service is up and running on port 8080")
	http.ListenAndServe(":8080", r)
}
