package utils

const BookingEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Bokningsbekräftelse - Ren & Flytt</title>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 20px; }
		.header { background-color: #f4f4f4; padding: 20px; border-radius: 5px; text-align: center; margin-bottom: 30px; }
		.booking-table { width: 100%; border-collapse: collapse; margin: 20px 0; }
		.booking-table th, .booking-table td { border: 1px solid #ddd; padding: 12px; text-align: left; }
		.booking-table th { background-color: #f2f2f2; font-weight: bold; }
		.section-header { background-color: #e8f4fd; font-weight: bold; }
		.services-list { margin: 0; padding-left: 20px; }
		.footer { margin-top: 30px; padding: 20px; background-color: #f4f4f4; border-radius: 5px; text-align: center; }
	</style>
</head>
<body>
	<div class="header">
		<h1>Bokningsbekräftelse</h1>
		<h2>Välkommen till Ren & Flytt!</h2>
		<p>Tack för att du valde våra tjänster. Nedan följer dina bokningsuppgifter:</p>
	</div>

	<table class="booking-table">
		<tr class="section-header">
			<td colspan="2"><strong>Kontaktinformation</strong></td>
		</tr>
		<tr>
			<th>Namn</th>
			<td>{{.Contact.Name}}</td>
		</tr>
		<tr>
			<th>Personnummer</th>
			<td>{{.Contact.SSN}}</td>
		</tr>
		<tr>
			<th>Email</th>
			<td>{{.Contact.Email}}</td>
		</tr>
		<tr>
			<th>Telefonnummer</th>
			<td>{{.Contact.Phone}}</td>
		</tr>
		<tr>
			<th>RUT-avdrag</th>
			<td>{{if .Contact.Rutavdrag}}Yes{{else}}No{{end}}</td>
		</tr>
		{{if .Contact.Message}}
		<tr>
			<th>Meddelande</th>
			<td>{{.Contact.Message}}</td>
		</tr>
		{{end}}

		<tr class="section-header">
			<td colspan="2"><strong>Services Requested</strong></td>
		</tr>
		<tr>
			<th>Valda Tjänster</th>
			<td>
				<ul class="services-list">
					{{range .Services}}
					<li>{{.}}</li>
					{{end}}
				</ul>
			</td>
		</tr>

		<tr class="section-header">
			<td colspan="2"><strong>Datum</strong></td>
		</tr>
		<tr>
			<th>Flyttdatum</th>
			<td>{{.MovingDate.Format "2006-01-02"}}</td>
		</tr>
		<tr>
			<th>Flexibelt datum</th>
			<td>{{if .FlexibleDate}}Yes{{else}}No{{end}}</td>
		</tr>
		<tr>
			<th>Städdatum</th>
			<td>{{.CleaningDate.Format "2006-01-02"}}</td>
		</tr>

		<tr class="section-header">
			<td colspan="2"><strong>Nuvarande adress</strong></td>
		</tr>
		<tr>
			<th>Adress</th>
			<td>{{.CurrentAddress.Address}}</td>
		</tr>
		<tr>
			<th>Typ av boende</th>
			<td>{{.CurrentAddress.ResidenceType}}</td>
		</tr>
		<tr>
			<th>Boyta</th>
			<td>{{.CurrentAddress.LivingSpace}} m²</td>
		</tr>
		<tr>
			<th>Våning</th>
			<td>{{.CurrentAddress.Floor}}</td>
		</tr>
		<tr>
			<th>Tillgång</th>
			<td>{{.CurrentAddress.Accessibility}}</td>
		</tr>

		<tr class="section-header">
			<td colspan="2"><strong>Flyttadress</strong></td>
		</tr>
		<tr>
			<th>Adress</th>
			<td>{{.NewAddress.Address}}</td>
		</tr>
		<tr>
			<th>Typ av boende</th>
			<td>{{.NewAddress.ResidenceType}}</td>
		</tr>
		<tr>
			<th>Boyta</th>
			<td>{{.NewAddress.LivingSpace}} m²</td>
		</tr>
		<tr>
			<th>Våning</th>
			<td>{{.NewAddress.Floor}}</td>
		</tr>
		<tr>
			<th>Tilgång</th>
			<td>{{.NewAddress.Accessibility}}</td>
		</tr>
	</table>

	<div class="footer">
		<p><strong>Vad händer härnäst?</strong></p>
		<p>Vi kommer att granska din bokningsförfrågan och kontakta dig inom 24 timmar för att bekräfta uppgifterna och ge en offert.</p>
		<p>Om du har fler frågor, vänligen kontakta oss:</p>
		<p>📧 Email: <a href="mailto:info@renochflytt.se">info@renochflytt.se</a></p>
		<p>🌐 Webbplats: <a href="https://renochflytt.se">www.renochflytt.se</a></p>
		<p><em>Tack för att du väljer Ren & Flytt!</em></p>
	</div>
</body>
</html>
`