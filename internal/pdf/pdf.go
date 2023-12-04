package pdf

import (
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/thejixer/shop-api/internal/models"
	"github.com/thejixer/shop-api/internal/utils"
)

func genMessage(order *models.OrderDto) string {
	switch order.Status {
	case "created":
		return fmt.Sprintf("you're order has been submited and will be verified by our office soon")
	case "verified":
		return fmt.Sprintf("you order has been verified and will be packaged by our stock soon")
	case "packaged":
		return fmt.Sprintf("your order has been packaged and is in queue to be sent by our office soon")
	case "sent":
		return fmt.Sprintf("your order has been sent and will be delivered soon")
	case "delivered":
		return fmt.Sprintf("your order has been successfully delivered")
	default:
		return ""
	}
}

func GeneratePDF(order *models.OrderDto) (string, error) {

	mkDirErr := os.Mkdir("public", 0750)
	if mkDirErr != nil && !os.IsExist(mkDirErr) {
		return "", mkDirErr
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Invoice ", true)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 40)

	pdf.Cell(120, 10, "")
	pdf.Cell(60, 10, "INVOICE")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 14)
	idString := fmt.Sprintf("#%v", order.Id)
	pdf.Cell(float64(178-3*len(idString)), 20, "")
	pdf.Cell(60, 20, idString)
	pdf.Ln(1)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(12, 20, "")
	pdf.Cell(35, 20, "Name:")
	pdf.SetFont("Arial", "", 16)
	pdf.Cell(20, 20, order.User.Name)
	pdf.Ln(9)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(12, 20, "")
	pdf.Cell(35, 20, "Deliver to:")
	pdf.SetFont("Arial", "", 14)
	pdf.Cell(20, 20, order.Address.RecieverName)
	pdf.Ln(6)

	pdf.Cell(47, 20, "")
	pdf.Cell(20, 20, order.Address.RecieverPhone)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(12, 20, "")
	pdf.Cell(35, 20, "Bank")
	pdf.Cell(20, 20, "really a great bank")
	pdf.Ln(6)
	pdf.Cell(12, 20, "")
	pdf.Cell(35, 20, "Account Name")
	pdf.Cell(20, 20, "Jix shop owner")
	pdf.Ln(6)
	pdf.Cell(12, 20, "")
	pdf.Cell(35, 20, "BSB")
	pdf.Cell(20, 20, "000 000")
	pdf.Ln(6)
	pdf.Cell(12, 20, "")
	pdf.Cell(35, 20, "Account Number")
	pdf.Cell(20, 20, "0000 0000")
	pdf.Ln(12)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(12, 20, "")
	pdf.Cell(80, 20, "Product")
	pdf.Cell(40, 20, "Price")
	pdf.Cell(40, 20, "Quantity")
	pdf.Cell(40, 20, "Total")
	pdf.Ln(12)
	var startOfLine float64 = 88
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(.4)
	pdf.SetLineCapStyle("square")
	pdf.Line(22, startOfLine, 199, startOfLine)

	pdf.SetFont("Arial", "", 12)
	for _, p := range order.Items {
		startOfLine = startOfLine + 12
		pdf.Cell(12, 20, "")
		pdf.Cell(80, 20, p.Product.Title)
		pdf.Cell(40, 20, fmt.Sprintf("$%v", p.Product.Price))
		pdf.Cell(40, 20, fmt.Sprintf("%v", p.Quantity))
		pdf.Cell(40, 20, fmt.Sprintf("$%v", utils.ToFixed(p.Product.Price*float64(p.Quantity), 2)))
		pdf.Line(22, startOfLine, 199, startOfLine)
		pdf.Ln(12)
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(12, 20, "")
	pdf.Cell(80, 20, "Total")
	pdf.Cell(80, 20, "")
	pdf.Cell(40, 20, fmt.Sprintf("$%v", order.TotalPrice))

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(30)
	pdf.Cell(12, 10, "")
	pdf.MultiCell(
		180,
		5,
		genMessage(order),
		"",
		"",
		false,
	)

	fileLoccation := fmt.Sprintf("public/order-%v.pdf", order.Id)
	err := pdf.OutputFileAndClose(fileLoccation)
	if err != nil {
		return "", err
	}

	return fileLoccation, nil
}
