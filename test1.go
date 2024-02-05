package main

import (
    "fmt"
    "github.com/go-pdf/fpdf"
)

func main() {
    // Create a new PDF document
    pdf := fpdf.New("P", "mm", "A4", "")

    // Add a page
    pdf.AddPage()

    // Add clinic information
    addClinicInfo(pdf)

    // Add patient information and prescription
    addPatientInfo(pdf)

    // Output the PDF to a file
    err := pdf.OutputFileAndClose("prescription.pdf")
    if err != nil {
        fmt.Println("Error:", err)
    }
}

func addClinicInfo(pdf *fpdf.Fpdf) {
    // Set font for clinic info
    pdf.SetFont("Arial", "B", 14)

	pdf.Image("logo.jpg", 10, 10, 30, 0, false, "", 0, "")

    // Clinic Name
    pdf.CellFormat(0, 10, "Clinic Name", "", 0, "C", false, 0, "")
    pdf.Ln(5)

    // Clinic Address
    pdf.SetFont("Arial", "", 12)
    pdf.CellFormat(0, 10, "123 Main Street", "", 0, "C", false, 0, "")
    pdf.Ln(5)
    pdf.CellFormat(0, 10, "City, Country", "", 0, "C", false, 0, "")
    pdf.Ln(10)
}

func addPatientInfo(pdf *fpdf.Fpdf) {
    // Set font for patient info
    pdf.SetFont("Arial", "B", 12)

    // Patient Name
    pdf.CellFormat(0, 10, "Patient Name: Jane Doe", "", 0, "L", false, 0, "")
    pdf.Ln(5)

    // Patient Details
    pdf.SetFont("Arial", "", 12)
    pdf.CellFormat(0, 10, "Age: 30", "", 0, "L", false, 0, "")
    pdf.Ln(5)
    pdf.CellFormat(0, 10, "Gender: Female", "", 0, "L", false, 0, "")
    pdf.Ln(10)

    // Prescription
    pdf.SetFont("Arial", "B", 12)
    pdf.CellFormat(0, 10, "Prescription", "", 0, "L", false, 0, "")
    pdf.Ln(5)

    // Prescription Details
    pdf.SetFont("Arial", "", 12)
    pdf.MultiCell(0, 10, "1. Diagnosis: Fever\n2. Medication: Paracetamol 500mg\n3. Dosage: 1 tablet every 6 hours", "", "L", false)
}
