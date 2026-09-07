package service

import (
	"app/src/model"
	"testing"
)

func TestCanStartSync(t *testing.T) {
	if !canStartSync(0, nil) {
		t.Fatal("empty should allow sync")
	}
	if canStartSync(3, nil) {
		t.Fatal("existing cfdis should block")
	}
	if canStartSync(0, &model.SatMasivaJob{Status: "running"}) {
		t.Fatal("running should block")
	}
	if canStartSync(0, &model.SatMasivaJob{Status: "done"}) {
		t.Fatal("done should block")
	}
	if !canStartSync(2, &model.SatMasivaJob{Status: "failed"}) {
		t.Fatal("failed should allow retry")
	}
	if !canStartSync(2, &model.SatMasivaJob{Status: "aborted"}) {
		t.Fatal("aborted should allow retry")
	}
}

func TestParseMonto(t *testing.T) {
	v := parseMonto("1234.50")
	if v == nil || *v != 1234.50 {
		t.Fatalf("got %v", v)
	}
	if parseMonto("") != nil {
		t.Fatal("empty")
	}
}

func TestParseSATTime(t *testing.T) {
	if parseSATTime("2026-03-01T12:00:00") == nil {
		t.Fatal("expected time")
	}
}
