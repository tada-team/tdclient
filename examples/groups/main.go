package main

import (
	"fmt"
	"log"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
	"github.com/tealeg/xlsx/v3"
)

const UNKNOWN = "unknown"

func main() {

	settings := examples.NewSettings()
	settings.RequireToken()
	settings.RequireTeam()
	settings.RequireChat()
	settings.Parse()

	session, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	session.SetToken(settings.Token)

	contacts, err := session.Contacts(settings.TeamUid)
	if err != nil {
		panic(err)
	}
	log.Println("Total team contacts count:", len(contacts))

	groupUid := tdproto.JID(settings.Chat)
	groupMembers, _ := session.GroupMembers(settings.TeamUid, groupUid)

	companies := splitByCompanies(contacts)
	printCompaniesContentReport(companies)

	users := make([]Users, 0)
	for _, v := range companies {
		for _, contact := range v {
			if !checkMembership(contact, groupMembers) {
				continue
			}

			users = append(users, Users{
				IsArchive:   contact.IsArchive,
				DisplayName: contact.DisplayName,
				Company:     contact.CustomFields.Company,
				Department:  contact.CustomFields.Department,
				Title:       contact.Role,
				Email:       contact.ContactEmail,
				Jid:         contact.Jid.String(),
			})
		}
	}
	save(users)
}

func checkMembership(contact tdproto.Contact, membership []tdproto.GroupMembership) bool {
	for _, v := range membership {
		if v.Jid.String() == contact.Jid.String() {
			return true
		}
	}
	return false
}

func updateContactCustomFields(contact tdproto.Contact) tdproto.Contact {
	if contact.CustomFields == nil {
		contact.CustomFields = &tdproto.ContactCustomFields{}
	}
	return contact
}

func splitByCompanies(contacts []tdproto.Contact) map[string][]tdproto.Contact {
	companies := make(map[string][]tdproto.Contact)
	for _, contact := range contacts {
		contact = updateContactCustomFields(contact)

		if contact.CustomFields.Company != "" {
			company := contact.CustomFields.Company
			companies[company] = append(companies[company], contact)
		} else {
			companies[UNKNOWN] = append(companies[UNKNOWN], contact)
		}
	}
	return companies
}

func printCompaniesContentReport(companies map[string][]tdproto.Contact) {
	fmt.Println("Participants count of each company")
	for k, v := range companies {
		log.Println(len(v), k)
	}
}

type Users struct {
	IsArchive   bool
	DisplayName string
	Company     string
	Department  string
	Title       string
	Email       string
	Jid         string
}

type xlsxStatCol struct {
	title string
	fn    func(cell *xlsx.Cell, users Users)
}

func save(users []Users) {
	cols := []xlsxStatCol{
		{"Name", func(cell *xlsx.Cell, users Users) {
			cell.SetString(users.DisplayName)
		}},
		{"Company", func(cell *xlsx.Cell, users Users) {
			cell.SetString(users.Company)
		}},
		{"Department", func(cell *xlsx.Cell, users Users) {
			cell.SetString(users.Department)
		}},
		{"Position", func(cell *xlsx.Cell, users Users) {
			cell.SetString(users.Title)
		}},
		{"Email", func(cell *xlsx.Cell, users Users) {
			cell.SetString(users.Email)
		}},
		{"Archived", func(cell *xlsx.Cell, users Users) {
			cell.SetBool(users.IsArchive)
		}},
		{"UID", func(cell *xlsx.Cell, users Users) {
			cell.SetString(users.Jid)
		}},
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Users")
	if err != nil {
		panic(err)
	}

	row := sheet.AddRow()
	for _, v := range cols {
		cell := row.AddCell()
		cell.SetString(v.title)
	}

	for _, user := range users {
		row := sheet.AddRow()
		for _, v := range cols {
			v.fn(row.AddCell(), user)
		}
	}
	err = file.Save("report.xlsx")
	if err != nil {
		panic(err)
	}
}
