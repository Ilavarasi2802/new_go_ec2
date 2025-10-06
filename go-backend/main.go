package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Employee
type Employee struct {
	ID   int    `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

// Department
type Department struct {
	Name  string `bson:"name" json:"name"`
	EmpID int    `bson:"emp_id" json:"emp_id"`
}

// Developer
type Developer struct {
	Language string `bson:"language" json:"language"`
	EmpID    int    `bson:"emp_id" json:"emp_id"`
}

// Tester
type Tester struct {
	Language string `bson:"language" json:"language"`
	EmpID    int    `bson:"emp_id" json:"emp_id"`
}

// Response row
type EmployeeFull struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Language   string `json:"language"`
}

// Globals
var client *mongo.Client
var empCol, deptCol, devCol, testCol *mongo.Collection

// Get next emp id
func getNextEmpID(ctx context.Context) (int, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	var last Employee
	err := empCol.FindOne(ctx, bson.D{}, opts).Decode(&last)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 1, nil
		}
		return 0, err
	}
	return last.ID + 1, nil
}

// Safe extract string from first element
func safeGetFirstStringField(row bson.M, arrName, fieldName string) string {
	a, ok := row[arrName]
	if !ok {
		return ""
	}
	arr, ok := a.(bson.A)
	if !ok {
		if ia, ok2 := a.([]interface{}); ok2 && len(ia) > 0 {
			if m, ok3 := ia[0].(bson.M); ok3 {
				if s, ok4 := m[fieldName].(string); ok4 {
					return s
				}
			}
		}
		return ""
	}
	if len(arr) == 0 {
		return ""
	}
	if m, ok := arr[0].(bson.M); ok {
		if s, ok2 := m[fieldName].(string); ok2 {
			return s
		}
	}
	return ""
}

// Convert any numeric type to int
func intFromAny(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	case float64:
		return int(t)
	default:
		return 0
	}
}

// ---------- Handlers ----------

// Create employee
func createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Name       string `json:"name"`
		Department string `json:"department"`
		Language   string `json:"language"`
		Role       string `json:"role"` // "developer" or "tester"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Department == "" || req.Role == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	newID, err := getNextEmpID(ctx)
	if err != nil {
		http.Error(w, "error generating emp id", http.StatusInternalServerError)
		log.Println("getNextEmpID err:", err)
		return
	}

	// Insert into employees
	_, err = empCol.InsertOne(ctx, Employee{ID: newID, Name: req.Name})
	if err != nil {
		http.Error(w, "error inserting employee", http.StatusInternalServerError)
		log.Println("insert emp err:", err)
		return
	}

	// Insert into department
	_, err = deptCol.InsertOne(ctx, Department{Name: req.Department, EmpID: newID})
	if err != nil {
		http.Error(w, "error inserting department", http.StatusInternalServerError)
		log.Println("insert dept err:", err)
		return
	}

	// Insert into developer/tester
	if req.Role == "developer" {
		_, err = devCol.InsertOne(ctx, Developer{Language: req.Language, EmpID: newID})
		if err != nil {
			http.Error(w, "error inserting developer", http.StatusInternalServerError)
			log.Println("insert dev err:", err)
			return
		}
	} else if req.Role == "tester" {
		_, err = testCol.InsertOne(ctx, Tester{Language: req.Language, EmpID: newID})
		if err != nil {
			http.Error(w, "error inserting tester", http.StatusInternalServerError)
			log.Println("insert tester err:", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "inserted", "id": newID})
}

// Get all employees with lookup
func getEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "departments"},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "department_info"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "developers"},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "developer_info"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "testers"},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "emp_id"},
			{Key: "as", Value: "tester_info"},
		}}},
	}

	cursor, err := empCol.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, "aggregation error", http.StatusInternalServerError)
		log.Println("aggregate err:", err)
		return
	}

	var rows []bson.M
	if err := cursor.All(ctx, &rows); err != nil {
		http.Error(w, "cursor read error", http.StatusInternalServerError)
		return
	}

	var out []EmployeeFull
	for _, r := range rows {
		lang := safeGetFirstStringField(r, "developer_info", "language")
		if lang == "" {
			lang = safeGetFirstStringField(r, "tester_info", "language")
		}
		dept := safeGetFirstStringField(r, "department_info", "name")
		id := intFromAny(r["id"])
		name, _ := r["name"].(string)
		out = append(out, EmployeeFull{
			ID:         id,
			Name:       name,
			Department: dept,
			Language:   lang,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// ---------- Main ----------
func main() {
	uri := "mongodb+srv://ilavarasinataraj_db_user:RmkY5omt7xbLViA9@cluster0.d497idh.mongodb.net/?retryWrites=true&w=majority"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("connect:", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("ping:", err)
	}
	fmt.Println("Connected to MongoDB!")

	db := client.Database("new_company")
	empCol = db.Collection("employees")
	deptCol = db.Collection("departments")
	devCol = db.Collection("developers")
	testCol = db.Collection("testers")

	r := mux.NewRouter()
	r.HandleFunc("/employee", createEmployeeHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/employees", getEmployeesHandler).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	fmt.Println("Backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
