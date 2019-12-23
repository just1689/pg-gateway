package query

import "testing"

func TestBuildQueryFromURLForSimpleQuery(t *testing.T) {
	url := "people?age=eq.25"
	query, err := BuildQueryFromURL(url)
	if err != nil {
		t.Error("unexpected error returned for good url")
		return
	}
	if query.Entity != "people" {
		t.Log("Got", query.Entity, "and not", "people")
		t.Fail()
	}

	if len(query.Comparisons) != 1 {
		t.Error("Expected one query not", len(query.Comparisons))
		return
	}
	if query.Comparisons[0].Field != "age" {
		t.Error("Expected age not", query.Comparisons[0].Field)
		return
	}
	if query.Comparisons[0].Comparator != "eq" {
		t.Error("Expected eq not", query.Comparisons[0].Comparator)
		return
	}
	if query.Comparisons[0].Value != "25" {
		t.Error("Expected 25 not", query.Comparisons[0].Value)
		return
	}

}

func TestBuildQueryFromURLForMultiQuery(t *testing.T) {
	url := "people?age=gte.25&active=is.true"
	query, err := BuildQueryFromURL(url)
	if err != nil {
		t.Error("unexpected error returned for good url")
		return
	}
	if query.Entity != "people" {
		t.Log("Got", query.Entity, "and not", "people")
		t.Fail()
	}

	if len(query.Comparisons) != 2 {
		t.Error("Expected 2 query not", len(query.Comparisons))
		return
	}
	if query.Comparisons[0].Field != "age" {
		t.Error("Expected age not", query.Comparisons[0].Field)
		return
	}
	if query.Comparisons[0].Comparator != "gte" {
		t.Error("Expected eq not", query.Comparisons[0].Comparator)
		return
	}
	if query.Comparisons[0].Value != "25" {
		t.Error("Expected 25 not", query.Comparisons[0].Value)
		return
	}

	if query.Comparisons[1].Field != "active" {
		t.Error("Expected active not", query.Comparisons[1].Field)
		return
	}
	if query.Comparisons[1].Comparator != "is" {
		t.Error("Expected is not", query.Comparisons[1].Comparator)
		return
	}
	if query.Comparisons[1].Value != "true" {
		t.Error("Expected true not", query.Comparisons[1].Value)
		return
	}

}

func TestBuildQueryFromURLForMultiQuery3(t *testing.T) {
	url := "people?age=gte.25&active=is.true&name=eq.Justin"
	query, err := BuildQueryFromURL(url)
	if err != nil {
		t.Error("unexpected error returned for good url")
		return
	}
	if query.Entity != "people" {
		t.Log("Got", query.Entity, "and not", "people")
		t.Fail()
	}

	if len(query.Comparisons) != 3 {
		t.Error("Expected 3 query not", len(query.Comparisons))
		return
	}
	if query.Comparisons[0].Field != "age" {
		t.Error("Expected age not", query.Comparisons[0].Field)
		return
	}
	if query.Comparisons[0].Comparator != "gte" {
		t.Error("Expected eq not", query.Comparisons[0].Comparator)
		return
	}
	if query.Comparisons[0].Value != "25" {
		t.Error("Expected 25 not", query.Comparisons[0].Value)
		return
	}

	if query.Comparisons[1].Field != "active" {
		t.Error("Expected active not", query.Comparisons[1].Field)
		return
	}
	if query.Comparisons[1].Comparator != "is" {
		t.Error("Expected is not", query.Comparisons[1].Comparator)
		return
	}
	if query.Comparisons[1].Value != "true" {
		t.Error("Expected true not", query.Comparisons[1].Value)
		return
	}

	if query.Comparisons[2].Field != "name" {
		t.Error("Expected name not", query.Comparisons[2].Field)
		return
	}
	if query.Comparisons[2].Comparator != "eq" {
		t.Error("Expected eq not", query.Comparisons[2].Comparator)
		return
	}
	if query.Comparisons[2].Value != "Justin" {
		t.Error("Expected Justin not", query.Comparisons[2].Value)
		return
	}

}

func TestBuildQueryFromURLForLimit(t *testing.T) {
	url := "people?age=gte.25&active=is.true&name=eq.Justin&limit=10"
	query, err := BuildQueryFromURL(url)
	if err != nil {
		t.Error("unexpected error returned for good url")
		return
	}
	if query.Limit != 10 {
		t.Error("did not get limit of 10, instead got", query.Limit)
		return
	}
}

func TestBuildQueryFromURLForLimitFirst(t *testing.T) {
	url := "people?limit=10"
	query, err := BuildQueryFromURL(url)
	if err != nil {
		t.Error("unexpected error returned for good url")
		return
	}
	if query.Limit != 10 {
		t.Error("did not get limit for first param of 10, instead got", query.Limit)
		return
	}
}

func TestBuildQueryFromURLForSelect1(t *testing.T) {
	url := "people?age=gte.25&active=is.true&name=eq.Justin&limit=10&select=id"
	query, err := BuildQueryFromURL(url)
	if err != nil {
		t.Error("unexpected error returned for good url")
		return
	}
	if len(query.Select) != 1 {
		t.Error("did not find a select! should be one")
		return
	}
}

func TestBuildQueryFromURLForSelect2(t *testing.T) {
	url := "people?age=gte.25&active=is.true&name=eq.Justin&limit=10&select=id,name"
	query, err := BuildQueryFromURL(url)
	if err != nil {
		t.Error("unexpected error returned for good url")
		return
	}
	if len(query.Select) != 2 {
		t.Error("did not find 2 select! should be 2")
		return
	}

	if query.Select[0] != "id" || query.Select[1] != "name" {
		t.Error("expected selects to be id and name, not", query.Select[0], query.Select[1])
	}
}

func TestBuildQueryFromURL(t *testing.T) {
	u := "http://localhost:8080/friends?id=100"
	_, err := BuildQueryFromURL(u)
	if err == nil {
		t.Error("expected bad url to return result in error for BuildQueryFromURL")
		return
	}
}
