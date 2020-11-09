package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

/* Test /article endpoint first check post method then get method retrive info and check with already post data, return error if not same, then check /article
	get method with pagination and compare with expected output
Test /article/<id> check output as same as expected
Test article/search end point */

// Test articleHandle function - Endpoint /articles - get and put method
func TestHandleArticles(t *testing.T) {

	// check post article method
	var addArticle = []byte(`{"title":"Create Some Title","subtitle":"Create Some Subtitle","content":"Create long content comes here"}`)
	var addM map[string]interface{}
	json.Unmarshal(addArticle, &addM)

	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(addArticle))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleArticles)
	handler.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusOK, rr.Code)

	var createM map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &createM)

	id := createM["InsertedID"].(string)
	fmt.Println(id)

	checkArticle(t, id, addM)

	// Test Get Article Method

	getReq, err := http.NewRequest("GET", "/articles?limit=1&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Content-Type", "application/json")
	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(handleArticles)
	getHandler.ServeHTTP(getRR, getReq)

	checkResponseCode(t, http.StatusOK, getRR.Code)

	expectedM := addM

	var getM []map[string]interface{}
	json.Unmarshal(getRR.Body.Bytes(), &getM)

	fmt.Println("get Rquest", getRR.Body.String())

	if getM[0]["title"] != expectedM["title"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["title"], getM[0]["title"])
	}
	if getM[0]["subtitle"] != expectedM["subtitle"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["subtitle"], getM[0]["subtitle"])
	}
	if getM[0]["content"] != expectedM["content"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["content"], getM[0]["content"])
	}

}

// Test getArticle function endpoint -/article/<id>  - get method
func TestGetArticle(t *testing.T) {

	var expected = []byte(`{"title":"Some Title","subtitle":"Some Subtitle","content":"long content comes here"}`)
	var expectedM map[string]interface{}
	json.Unmarshal(expected, &expectedM)

	id := "5fa977d57f3135488b270da1"
	checkArticle(t, id, expectedM)

}

// Test searchArticle function - Endpoint /articles/search - get method
func TestSearchArticle(t *testing.T) {

	var expected = []byte(`{"title":"Pfizer’s Covid Vaccine Prevents 90% of Infections in Study","subtitle":"The Covid-19 vaccine being developed by Pfizer Inc. and BioNTech SE prevented more than 90% of infections.","content":"The Covid-19 vaccine being developed by Pfizer Inc. and BioNTech SE prevented more than 90% of infections in a study of tens of thousands of volunteers, the most encouraging scientific advance so far in the battle against the coronavirus. Eight months into the worst pandemic in a century, the preliminary results pave the way for the companies to seek an emergency-use authorization from regulators."
	
	}`)
	var expectedM map[string]interface{}
	json.Unmarshal(expected, &expectedM)

	getReq, err := http.NewRequest("GET", "/articles/search?q=Covid Vaccine", nil)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Content-Type", "application/json")
	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(searchArticles)
	getHandler.ServeHTTP(getRR, getReq)

	checkResponseCode(t, http.StatusOK, getRR.Code)

	var getM []map[string]interface{}
	json.Unmarshal(getRR.Body.Bytes(), &getM)

	fmt.Println("get Rquest", getRR.Body.String())

	if getM[0]["title"] != expectedM["title"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["title"], getM[0]["title"])
	}
	if getM[0]["subtitle"] != expectedM["subtitle"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["subtitle"], getM[0]["subtitle"])
	}
	if getM[0]["content"] != expectedM["content"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["content"], getM[0]["content"])
	}

}

// Helper function to check request code is same as expected
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Check get method article is same as expected
func checkArticle(t *testing.T, id string, expectedM map[string]interface{}) {
	getReq, _ := http.NewRequest("GET", "/articles/"+id, nil)
	getReq.Header.Set("Content-Type", "application/json")
	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(getArticle)
	getHandler.ServeHTTP(getRR, getReq)

	checkResponseCode(t, http.StatusOK, getRR.Code)
	var addedM map[string]interface{}
	json.Unmarshal(getRR.Body.Bytes(), &addedM)

	fmt.Println(getRR.Body.String())
	if addedM["title"] != expectedM["title"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["title"], addedM["title"])
	}
	if addedM["subtitle"] != expectedM["subtitle"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["subtitle"], addedM["subtitle"])
	}
	if addedM["content"] != expectedM["content"] {
		t.Errorf("Expected the title to remain the same (%v). Got %v", expectedM["content"], addedM["content"])
	}
}
