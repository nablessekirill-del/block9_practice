package handler

import (
	"block9_practice/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

var count int = 0

type ProductHandler struct {
	store  *domain.Store
	logger *logrus.Logger
}

func New(store *domain.Store, l *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		store:  store,
		logger: l,
	}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	/*
		b, err := json.MarshalIndent(list, "", "	")
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(b); err != nil {
			fmt.Println("Failed to write http responce:", err)
			return
		}*/
	query := r.URL.Query().Get("search")
	list := h.store.List(query)

	writeJSON(w, http.StatusOK, list)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Не удалось конвертировать строку в инт")
		return
	}

	c, ok := h.store.Get((id))
	if ok == false {
		writeError(w, http.StatusNotFound, "Нет такого продукта")
		return
	}
	writeJSON(w, http.StatusFound, c)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if count < 5 {
		var product domain.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			writeError(w, http.StatusBadRequest, "Не получилось прочитать данные из запроса")
			return
		}
		if err := product.ValidateForCreate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			h.logger.Error(err.Error())
			return
		}
		product1 := h.store.Create(product)
		h.logger.Info("Добавлен продукт:", product1.Name, "Цена:", product1.Price)

		writeJSON(w, http.StatusCreated, product1)
		count++
		h.logger.Debug("Количество продуктов:", count)
	} else {
		h.logger.Warn("Превышен лимит товаров Т_Т")
	}

}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var tempPr domain.Product
	if err := json.NewDecoder(r.Body).Decode(&tempPr); err != nil {
		writeError(w, http.StatusBadRequest, "Не получилось прочитать данные из запроса")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Не удалось конвертировать строку в инт")
		return
	}

	prod, ok := h.store.Update(id, tempPr)
	if ok == false {
		writeError(w, http.StatusNotFound, "Не найден продукт")
		return
	}
	writeJSON(w, http.StatusAccepted, prod)

}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Не удалось конвертировать строку в инт")
		return
	}
	h.store.Delete(id)
	writeJSON(w, http.StatusNoContent, "")
	count--
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {

	writeJSON(w, status, map[string]string{"error": msg})

}
