package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iamneuron/students-check-api/internal/storage"
	"github.com/iamneuron/students-check-api/internal/types"
	"github.com/iamneuron/students-check-api/internal/utils/responce"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("welcome to students api"))

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			responce.WriteJson(w, http.StatusBadRequest, responce.GeneralError(err))
			return

		}
		if err != nil {
			responce.WriteJson(w, http.StatusBadRequest, responce.GeneralError(err))
		}

		//validate request

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			responce.WriteJson(w, http.StatusBadRequest, responce.ValidattionError(validateErrs))
			return
		}
		lastID, err := storage.CreateStudent(
			student.Name,
			student.Email,
			int(student.Age),
		)

		slog.Info("User Creaetd successfuly")
		if err != nil {
			responce.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		responce.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastID})
	}
}
