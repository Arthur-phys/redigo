//go:build !integration && !e2e
// +build !integration,!e2e

package cache

import (
	"testing"

	"github.com/Arthur-phys/redigo/pkg/redigoerr"
)

func TestSet_Should_Return_NIL_When_Set_Is_Done(t *testing.T) {
	cs := New()
	err := cs.Set("KEY", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
}

func TestGet_Should_Return_Error_When_Key_Not_Present(t *testing.T) {
	cs := New()
	if _, err := cs.Get("KEY"); err == nil {
		t.Errorf("A key is present when it should not be!")
	}
}

func TestGet_Should_Return_Value_And_Nil_When_Key_Is_Present(t *testing.T) {
	cs := New()
	err := cs.Set("KEY", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if _, err := cs.Get("KEY"); err != nil {
		t.Errorf("A key is not present when it should be!")
	}
}

func TestRPush_Should_Create_NewCache_List_In_Cache_When_Not_Present(t *testing.T) {
	cs := New()
	err := cs.RPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("KEYVECTOR"); err != nil || val != "REDIGO" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
}

func TestRPush_Should_Add_To_List_In_Cache_When_Present(t *testing.T) {
	cs := New()
	err := cs.RPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("KEYVECTOR", "NIJI")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("KEYVECTOR"); err != nil || val != "NIJI" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
}

func TestRPop_Should_Return_Error_When_List_Is_Not_Present(t *testing.T) {
	cs := New()
	if _, err := cs.RPop("KEYVECTOR"); err == nil {
		t.Errorf("Expected error but obtained nil! %v", err)
	}
}

func TestRPop_Should_Delete_List_When_Empty(t *testing.T) {
	cs := New()
	err := cs.RPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("KEYVECTOR"); err != nil || val != "REDIGO" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.RPop("KEYVECTOR"); !redigoerr.KeyNotFound(err) {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestRPop_Should_Remove_Elements_In_Succession(t *testing.T) {
	cs := New()
	err := cs.RPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("KEYVECTOR", "NIJI")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.RPop("KEYVECTOR"); err != nil || val != "NIJI" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if val, err := cs.RPop("KEYVECTOR"); err != nil || val != "REDIGO" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.RPop("KEYVECTOR"); !redigoerr.KeyNotFound(err) {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestLPush_Should_Add_To_List_In_Cache_When_Present(t *testing.T) {
	cs := New()
	err := cs.LPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("KEYVECTOR", "NIJI")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.LPop("KEYVECTOR"); err != nil || val != "NIJI" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
}

func TestLPop_Should_Delete_List_When_Empty(t *testing.T) {
	cs := New()
	err := cs.LPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.LPop("KEYVECTOR"); err != nil || val != "REDIGO" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.LPop("KEYVECTOR"); !redigoerr.KeyNotFound(err) {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestLPop_Should_Remove_Elements_In_Succession(t *testing.T) {
	cs := New()
	err := cs.LPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("KEYVECTOR", "NIJI")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if val, err := cs.LPop("KEYVECTOR"); err != nil || val != "NIJI" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if val, err := cs.LPop("KEYVECTOR"); err != nil || val != "REDIGO" {
		t.Errorf("An error occured! %v - %v", err, val)
	}
	if _, err := cs.LPop("KEYVECTOR"); !redigoerr.KeyNotFound(err) {
		t.Errorf("Unexpected error occurred! %v", err)
	}
}

func TestLPop_Should_Return_Error_When_List_Is_Not_Present(t *testing.T) {
	cs := New()
	if _, err := cs.LPop("KEYVECTOR"); err == nil {
		t.Errorf("Expected error but obtained nil! %v", err)
	}
}

func TestLIndex_Should_Return_Element_When_Present_In_List(t *testing.T) {
	cs := New()
	err := cs.LPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("KEYVECTOR", "NIJI")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("KEYVECTOR", "ANUBIS")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if str, err := cs.LIndex("KEYVECTOR", 0); err != nil || str != "ANUBIS" {
		t.Errorf("Unable to retrieve first value (0) from 'key' list! %v", err)
	}
}

func TestLIndex_Should_Return_Error_When_Index_Not_Present_In_List(t *testing.T) {
	cs := New()
	err := cs.LPush("KEYVECTOR", "REDIGO")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.RPush("KEYVECTOR", "NIJI")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	err = cs.LPush("KEYVECTOR", "ANUBIS")
	if err != nil {
		t.Errorf("An error occurred! %v", err)
	}
	if s, err := cs.LIndex("KEYVECTOR", 5); err == nil {
		t.Errorf("Was able to retrieve unexistant value! %v - %s", err, s)
	}
}
