package handlers

import (
	"app/views"
	"app/views/pages"
	"app/views/section"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetHome(c echo.Context) error {
	ctx := c.Request().Context()

	services, err := h.Client.Service.Query().All(ctx)
	if err != nil {
		return err
	}

	testimonials, err := h.Client.Testimonial.Query().All(ctx)
	if err != nil {
		return err
	}

	return views.Model(pages.Home{
		Title: "Your Local Aluminium Joinery Specialists",
		Content: `With over 25 years of experience, whether it's residential or commercial, we can help with all things window
		and doors in Auckland. Available to travel further when required.`,
		Services:     services,
		Testimonials: testimonials,
		Config: section.ConfigModel{
			Company:     "Rodney Joinery Services",
			PhoneNumber: "022 351 1657",
		},
		Company: section.Company{
			Name:        "Rodney Joinery Services",
			Email:       "mark@rodneyjoineryservices.co.nz",
			PhoneNumber: "022 351 1657",
		},
	}).Send(c)
}
