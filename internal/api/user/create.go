package user

import (
	"net/http"
)

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//
	//var req model.User
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	http.Error(w, "invalid JSON", http.StatusBadRequest)
	//	return
	//}
	//
	////createdUser, err := i.userService.Create(ctx, &req)
	////if err != nil {
	////	http.Error(w, "failed to create user", http.StatusInternalServerError)
	////	return
	////}
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(createdUser)
}
