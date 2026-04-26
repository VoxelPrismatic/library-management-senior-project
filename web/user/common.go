package user

type FormEntry struct {
	ID    string
	Type  string
	Label string
	Ghost string
}

var (
	EntryFirstName   = FormEntry{"firstName", "text", "First Name", "John"}
	EntryLastName    = FormEntry{"lastName", "text", "Last Name", "Smith"}
	EntryEmailAddr   = FormEntry{"emailAddr", "email", "Email Address", "jsmith27@depaul.edu"}
	EntrySecret      = FormEntry{"secret", "password", "Password", "Password1!"}
	EntrySecretAgain = FormEntry{"secret_again", "password", "Re-type Password", "Password1!"}

	RegisterFormEntries = []FormEntry{
		EntryFirstName,
		EntryLastName,
		EntryEmailAddr,
		EntrySecret,
		EntrySecretAgain,
	}
)
