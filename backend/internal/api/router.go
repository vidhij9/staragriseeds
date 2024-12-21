package api

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"
	"backend/internal/service"
	"fmt"

	"github.com/gorilla/mux"
)

func SetupRouter(services *service.Services) *mux.Router {
	r := mux.NewRouter()

	farmerHandler := handlers.NewFarmerHandler(services.Farmer)
	cceHandler := handlers.NewCCEHandler(services.CCE)
	ticketHandler := handlers.NewTicketHandler(services.Ticket, services.Farmer)

	fmt.Println("Inside setuprouter")

	// GET
	// Farmer routes
	r.HandleFunc("/farmers/{id}", farmerHandler.GetFarmer).Methods("GET")
	r.HandleFunc("/farmers", farmerHandler.GetFarmers).Methods("GET")
	r.HandleFunc("/farmer/contact/{contact}", farmerHandler.GetFarmerByContact).Methods("GET")

	// CCE routes
	r.HandleFunc("/cces/{id}", cceHandler.GetCCE).Methods("GET")
	r.HandleFunc("/cces", cceHandler.GetCCEs).Methods("GET")

	// Ticket routes
	r.HandleFunc("/tickets/{id}", ticketHandler.GetTicket).Methods("GET")
	r.HandleFunc("/tickets", ticketHandler.GetTickets).Methods("GET")
	r.HandleFunc("/tickets/farmer/{contact}", ticketHandler.GetTicketsByFarmer).Methods("GET")
	r.HandleFunc("/tickets/cce/{id}", ticketHandler.GetTicketsByCCE).Methods("GET")
	r.HandleFunc("/tickets/cce/{id}/status/{status}", ticketHandler.GetTicketsByCCEAndStatus).Methods("GET")

	// POST
	// Farmer routes
	r.HandleFunc("/farmers", middleware.AuthMiddleware(farmerHandler.CreateFarmer)).Methods("POST")
	// CCE routes
	r.HandleFunc("/cces", middleware.AuthMiddleware(cceHandler.CreateCCE)).Methods("POST")
	// Ticket routes
	r.HandleFunc("/tickets", middleware.AuthMiddleware(ticketHandler.CreateTicket)).Methods("POST")

	// PUT
	// Farmer routes
	r.HandleFunc("/farmers/{id}", middleware.AuthMiddleware(farmerHandler.UpdateFarmer)).Methods("PUT")
	// CCE routes
	r.HandleFunc("/cces/{id}", middleware.AuthMiddleware(cceHandler.UpdateCCE)).Methods("PUT")
	// Ticket routes
	r.HandleFunc("/tickets/{id}", middleware.AuthMiddleware(ticketHandler.UpdateTicket)).Methods("PUT")

	// DELETE
	// Farmer routes
	r.HandleFunc("/farmers/{id}", middleware.AuthMiddleware(farmerHandler.DeleteFarmer)).Methods("DELETE")
	// CCE routes
	r.HandleFunc("/cces/{id}", middleware.AuthMiddleware(cceHandler.DeleteCCE)).Methods("DELETE")
	// Ticket routes
	r.HandleFunc("/tickets/{id}", middleware.AuthMiddleware(ticketHandler.DeleteTicket)).Methods("DELETE")

	return r
}
