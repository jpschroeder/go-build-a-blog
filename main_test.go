package main

import "testing"

func TestMain(t *testing.T) {
	db, err := initDb(":memory:")
	if err != nil {
		t.Errorf("Cannot open in-memory database: %s", err)
	}

	if DefaultBlogExistsQuery(db) {
		t.Errorf("Default blog exists")
	}

	err = AddDefaultBlogCommand(db, "testkey", "testtitle")
	if err != nil {
		t.Errorf("Error adding default blog: %s", err)
	}

	if !DefaultBlogExistsQuery(db) {
		t.Errorf("Default blog not created")
	}

	blog, err := DefaultBlogQuery(db)
	if err != nil {
		t.Errorf("Error getting default blog: %s", err)
	}

	if blog.Title != "testtitle" {
		t.Errorf("Blog title mismatch: %s", blog.Title)
	}
}
