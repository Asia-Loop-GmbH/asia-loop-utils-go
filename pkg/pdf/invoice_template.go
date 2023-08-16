package pdf

import (
	"time"
)

type invoiceTemplateItems struct {
	Name     string
	SKU      string
	Amount   int
	Total    string
	TaxClass string
	Tax      string
}

type invoiceTemplateProps struct {
	StoreName       string
	StoreAddress    string
	StoreTaxNumber  string
	StoreTelephone  string
	StoreEmail      string
	CustomerName    string
	CustomerAddress string
	InvoiceNumber   string
	OrderNumber     string
	Date            time.Time
	Items           []invoiceTemplateItems
	Total           string
	Tax             string
	Tax19           string
	Tax7            string
}

const invoiceTemplate = `
<html>
<head>
    <style>
        div.footer {
            display: block;
            position: running(footer);
            bottom: 0;
        }
        @page {
            size: A4;
            margin: 10% 5% 15% 5%;
            @bottom-center {
                content: element(footer);
                height: 100pt;
            }
        }
    </style>
</head>
<body style="font-size: 10pt; font-family: Arial, Helvetica, sans-serif !important;">
    <table style="width: 100%;">
        <tr>
            <td style="width: 70%;"></td>
            <td style="width: 30%;">
                <span>{{.StoreName}}</span> <br/>
                <span>{{.StoreAddress}}</span>
            </td>
        </tr>
    </table>
    <h1>RECHNUNG</h1>
    <table style="width: 100%;">
        <tr>
            <td style="width: 50%;">
                <span>{{.CustomerName}}</span><br/>
                <span>{{.CustomerAddress}}</span><br/>
            </td>
            <td style="width: 30%;">
                Rechnungsnummer:<br/>
                Rechnungsdatum:<br/>
                Bestellnummer:<br/>
            </td>
            <td style="width: 20%;">
                <span>{{.InvoiceNumber}}</span><br/>
                <span>{{.Date | DateTime}}</span><br/>
                <span>{{.OrderNumber}}</span>
            </td>
        </tr>
    </table>
    <table style="width: 100%; margin-top: 2rem;">
        <thead>
            <tr style="text-align: left; background-color: #23a638; color: white;">
                <th style="width: 30%;">Produkt</th>
                <th style="width: 10%;">Anzahl</th>
                <th style="width: 20%;">Preis</th>
                <th style="width: 20%;">MwSt.</th>
                <th style="width: 20%;"></th>
            </tr>
        </thead>
        <tbody>
            {{range .Items}}
            <tr>
                <td style="border-bottom: 1px solid grey;">
                    <span>{{.Name}}</span><br/>
                    <span>Art.-Nr.: {{.SKU}}</span>
                </td>
                <td style="border-bottom: 1px solid grey;">
                    <span>{{.Amount}}</span>
                </td>
                <td style="border-bottom: 1px solid grey;">
                    <span>{{.Total}}€</span>
                </td>
                <td style="border-bottom: 1px solid grey;">
                    <span>{{.TaxClass}}</span>
                </td>
                <td style="border-bottom: 1px solid grey;">
                    <span>{{.Tax}}€</span>
                </td>
            </tr>
            {{end}}
        </tbody>
    </table>
    <table style="width: 100%; margin-top: 1rem;">
        <tr>
            <td style="width: 50%;"></td>
            <td style="width: 30%;">
                <strong>Gesamt</strong>
            </td>
            <td style="width: 20%;">
                <span>{{.Total}}€</span> (inkl. <span>{{.Tax}}€</span> MwSt.)
            </td>
        </tr>
		<tr>
			<td style="width: 50%;"></td>
            <td style="width: 30%;">
                <strong>MwSt. 19%</strong>
            </td>
            <td style="width: 20%;">
                <span>{{.Tax19}}€</span>
            </td>
		</tr>
		<tr>
			<td style="width: 50%;"></td>
            <td style="width: 30%;">
                <strong>MwSt. 7%</strong>
            </td>
            <td style="width: 20%;">
                <span>{{.Tax7}}€</span>
            </td>
		</tr>
    </table>
    <div class="footer">
        <table style="width: 100%; border-top: 2px solid #23a638;">
            <tr>
                <td style="width: 25%;">
                    <span>{{.StoreName}}</span><br/>
                    <span>{{.StoreAddress}}</span><br/>
                    Steuernummer: <span>{{.StoreTaxNumber}}</span>
                </td>
                <td style="width: 25%;">
                    Tel.: <span>{{.StoreTelephone}}</span><br/>
                    E-Mail: <span>{{.StoreEmail}}</span><br/>
                    Web: www.asialoop.de
                </td>
                <td style="width: 50%;">
                    Stadt- u. Kreissparkasse Erlangen Höchstadt Herzogenaurach<br/>
                    IBAN: DE65 7635 0000 0060 1146 75<br/>
                    BIC: BYLADEM1ERH
                </td>
            </tr>
        </table>
    </div>
</body>
</html>
`
