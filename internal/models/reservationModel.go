package models

import "time"

type Reservation struct {
	ID                            int       `json:"reservation_id"`
	USERID                        int       `json:"user_id"`
	TABLENUMBER                   int       `json:"table_number"`
	RESERVATIONPERSONNAME         string    `json:"reservation_person_name"`
	RESERVATIONPERSONEMAIL        string    `json:"reservation_person_email"`
	RESERVATIONPERSONMOBILENUMBER string    `json:"reservation_person_mobile_number"`
	RESERVATIONTIME               string    `json:"reservation_time"`
	RESERVATIONDATE               string    `json:"reservation_date"`
	NUMBEROFGUESTS                int       `json:"number_of_guests"`
	SPECIALREQUESTS               *string   `json:"special_requests"`
	STATUS                        string    `json:"status"`
	DISHITEMS                     []int     `json:"dish_items"`
	DISHDETAILS                   []Dish    `json:"dish_details"`
	CREATEDAT                     time.Time `json:"created_at"`
	PAYEMENTID                    string    `json:"payment_id"`
	SIGNATURE                     string    `json:"signature"`
	PAYMENTSTATUS                 string    `json:"payment_status"`
	INVOICEURL                    string    `json:"invoice_url"`
	UPDATEDAT                     time.Time `json:"updated_at"`
}
