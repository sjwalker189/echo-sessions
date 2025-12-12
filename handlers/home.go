package handlers

import (
	"app/view"
	"app/view/pages"
	"app/view/section"
	"github.com/labstack/echo/v4"
)

func Welcome(c echo.Context) error {
	// services, err := h.Client.Service.Query().All(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// testimonials, err := h.Client.Testimonial.Query().All(ctx)
	// if err != nil {
	// 	return err
	// }

	return view.Model(pages.Home{
		Title: "Your Local Aluminium Joinery Specialists",
		Content: `With over 25 years of experience, whether it's residential or commercial, we can help with all things window
		and doors in Auckland. Available to travel further when required.`,
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
