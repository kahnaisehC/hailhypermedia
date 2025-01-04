package contacts

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type ContactError struct {
	Message string
}

func (c ContactError) Error() string {
	return c.Message
}

type Contact struct {
	Id    int
	Name  string
	Email string
	Phone string
}

// (C)reate
func verifyContact(name, email, phone string, id int) (Contact, error) {
	contact := Contact{
		Id:    id,
		Name:  name,
		Email: email,
		Phone: phone,
	}

	contacts, err := GetContacts()
	if err != nil {
		return contact, err
	}

	// check Id
	if contact.Id != -1 {
		found := false
		for _, v := range contacts {
			if v.Id == contact.Id {
				found = true
				break
			}
		}
		if !found {
			return contact, ContactError{
				Message: fmt.Sprintf("No contact with id %d was found", contact.Id),
			}
		}
	} else {
		if len(contacts) == 0 {
			contact.Id = 1
		} else {
			contact.Id = contacts[len(contacts)-1].Id + 1
		}
	}
	// check name
	if len(name) >= 16 {
		return contact, ContactError{
			Message: "Name is more than 16 characters long",
		}
	}
	// check email
	emailIsCorrect := false
	atCount := 0
	for emailIndex := 0; emailIndex < len(email); emailIndex++ {
		if (emailIndex == 0 || emailIndex == len(email)-1) && email[emailIndex] == '@' {
			emailIsCorrect = false
			break
		}
		if email[emailIndex] == '@' {
			atCount++
			emailIsCorrect = atCount == 1
		}
	}
	if !emailIsCorrect {
		return contact, ContactError{
			Message: "Email format is Incorrect",
		}
	}

	// check phone
	if len(phone) >= 16 {
		return contact, ContactError{
			Message: "Phone is more than 16 characters long",
		}
	}
	return contact, nil
}

func CreateContact(newContact Contact) error {
	newContact, err := verifyContact(newContact.Name, newContact.Email, newContact.Phone, newContact.Id)
	if err != nil {
		return err
	}
	file, err := os.OpenFile("contacts.csv", os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(
		[]string{
			strconv.Itoa(newContact.Id),
			newContact.Name,
			newContact.Email,
			newContact.Phone,
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(newContact.Email)

	return nil
}

// (R)ead

func GetContact(id int) (Contact, error) {
	var contact Contact
	f, err := os.Open("contacts.csv")
	if err != nil {
		return contact, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		return contact, err
	}
	for _, csvContact := range data {
		csvId, err := strconv.Atoi(csvContact[0])
		if err != nil {
			return contact, err
		}
		if csvId == id {
			contact.Id = csvId
			contact.Name = csvContact[1]
			contact.Email = csvContact[2]
			contact.Phone = csvContact[3]
			return contact, nil
		}
	}

	return contact, ContactError{
		Message: "contact id not found",
	}
}

func GetContacts() ([]Contact, error) {
	var contacts []Contact
	f, err := os.Open("contacts.csv")
	if err != nil {
		f, err = os.Create("contacts.csv")
		if err != nil {
			return []Contact{}, err
		}
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		return []Contact{}, err
	}
	for _, csvContact := range data {
		contactId, err := strconv.Atoi(csvContact[0])
		if err != nil {
			fmt.Printf("cannot convert %v to integer\n", csvContact[0])
			continue
		}
		contactName := csvContact[1]
		contactEmail := csvContact[2]
		contactPhone := csvContact[3]

		contact := Contact{
			contactId,
			contactName,
			contactEmail,
			contactPhone,
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func FilterContacts(reg string, contacts []Contact) []Contact {
	if reg == "" {
		return contacts
	}
	var filteredContacts []Contact

	for _, v := range contacts {
		name := v.Name
		regIndex := 0
		for nameIndex := 0; nameIndex < len(name); nameIndex++ {
			if name[nameIndex] == reg[regIndex] {
				regIndex++
			}
			if regIndex == len(reg) {
				filteredContacts = append(filteredContacts, v)
				break
			}
		}

	}
	return filteredContacts
}

// (U)pdate
func UpdateContact(id int, name, email, phone string) error {
	if id < 1 {
		return ContactError{
			Message: "Invalid Id to update",
		}
	}
	updatedContact, err := verifyContact(name, email, phone, id)
	if err != nil {
		return err
	}

	contacts, err := GetContacts()
	if err != nil {
		return err
	}
	file, err := os.OpenFile("contacts.csv", os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, v := range contacts {
		if v.Id == updatedContact.Id {
			err = writer.Write(
				[]string{
					strconv.Itoa(updatedContact.Id),
					updatedContact.Name,
					updatedContact.Email,
					updatedContact.Phone,
				},
			)
			if err != nil {
				return err
			}
			continue
		}

		err = writer.Write(
			[]string{
				strconv.Itoa(v.Id),
				v.Name,
				v.Email,
				v.Phone,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// (D)elete
func DeleteContact(id int) error {
	contacts, err := GetContacts()
	if err != nil {
		return err
	}
	err = os.Remove("contacts.csv")
	if err != nil {
		return err
	}
	file, err := os.OpenFile("contacts.csv", os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, v := range contacts {
		if v.Id == id {
			continue
		}

		err = writer.Write(
			[]string{
				strconv.Itoa(v.Id),
				v.Name,
				v.Email,
				v.Phone,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
