package handlers

import "net/http"

type Interface interface {
	Get(w http.ResponseWriter, r *http.Request)
	CreateRecord(w http.ResponseWriter, r *http.Request)
	UpdateRecord(w http.ResponseWriter, r *http.Request)
	DeleteRecord(w http.ResponseWriter, r *http.Request)
	UnhandledMethod(w http.ResponseWriter, r *http.Request)
}

type ControlDataInsertInterface interface {
	CreateRecord(w http.ResponseWriter, r *http.Request)
	CreateRecordFromOtherEnv(w http.ResponseWriter, r *http.Request)
	//CreateRecordInternal(w http.ResponseWriter, r *http.Request)
}
