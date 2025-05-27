package http

import (
	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/repository"
	"Effective-Mobile/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.PeopleService
}

func NewHandler(s *service.PeopleService) *Handler { return &Handler{svc: s} }

func (h *Handler) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	people := v1.Group("/people")
	{
		people.POST("", h.create)
		people.GET("", h.list)
		people.GET("/:id", h.get)
		people.PUT("/:id", h.update)
		people.DELETE("/:id", h.delete)
	}
}

// @Summary      Create person
// @Description  Добавляет нового человека и мгновенно обогащает возрастом, полом, национальностью
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        person  body      http.CreatePersonRequest  true  "ФИО"
// @Success      201     {object}  map[string]int64          "{"id":1}"
// @Failure      400     {object}  map[string]interface{}    "ошибка валидации"
// @Failure      500     {object}  map[string]string         "внутренняя ошибка"
// @Router       /api/v1/people [post]
func (h *Handler) create(c *gin.Context) {
	var req CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	if err := validate.Struct(&req); err != nil {
		badValidation(c, err)
		return
	}

	p := domain.Person{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	id, err := h.svc.Add(c.Request.Context(), &p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary      List people
// @Description  Вывести список людей с фильтрацией и пагинацией
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "максимум элементов (default 10)"
// @Param        offset     query     int     false  "смещение (default 0)"
// @Param        name       query     string  false  "фильтр по имени, ILIKE"
// @Param        surname    query     string  false  "фильтр по фамилии, ILIKE"
// @Success      200  {object}  map[string]interface{} "total & items"
// @Router       /api/v1/people [get]
func (h *Handler) list(c *gin.Context) {
	lim, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	off, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	f := repository.ListFilter{
		Name:    c.Query("name"),
		Surname: c.Query("surname"),
		Limit:   lim,
		Offset:  off,
	}
	people, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "items": people})
}

// @Summary      Get person
// @Description  Получить человека по ID
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  domain.Person
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/people/{id} [get]
func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	p, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// @Summary      Update person
// @Description  Полностью обновить ФИО (и заново обогатить данные)
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        id      path      int                     true  "ID"
// @Param        person  body      http.UpdatePersonRequest  true  "Новое ФИО"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/people/{id} [put]
func (h *Handler) update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	if err := validate.Struct(&req); err != nil {
		badValidation(c, err)
		return
	}

	p := domain.Person{
		ID:         id,
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	if err := h.svc.Update(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary      Delete person
// @Description  Удалить человека по ID
// @Tags         people
// @Param        id   path      int  true  "ID"
// @Success      204
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/people/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
