package client

import (
	"bufio"
	"fmt"
	"github.com/stlesnik/goph_keeper/internal/config"
	"golang.org/x/term"
	"os"
	"strings"
)

// App represents the CLI application
type App struct {
	client *Client
}

// NewApp creates a new CLI application
func NewApp(config config.ClientConfig) *App {
	return &App{
		client: NewClient(config),
	}
}

// Run starts the CLI application
func (a *App) Run() error {
	fmt.Println("-------------------GophKeeper Client-------------------")

	for {
		fmt.Println("\n---Available commands:---")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. View Profile")
		fmt.Println("4. List Data")
		fmt.Println("5. Get Data Item")
		fmt.Println("6. Create Data Item")
		fmt.Println("7. Update Data Item")
		fmt.Println("8. Delete Data Item")
		fmt.Println("9. Change Password")
		fmt.Println("0. Exit")

		choice := a.getInput("Choose an option (0-9): ")

		switch choice {
		case "1":
			a.handleRegister()
		case "2":
			a.handleLogin()
		case "3":
			a.handleProfile()
		case "4":
			a.handleListData()
		case "5":
			a.handleGetData()
		case "6":
			a.handleCreateData()
		case "7":
			a.handleUpdateData()
		case "8":
			a.handleDeleteData()
		case "9":
			a.handleChangePassword()
		case "0":
			fmt.Println("Goodbye!")
			return nil
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// handleRegister handles user registration
func (a *App) handleRegister() {
	fmt.Println("\n-----------User Registration-----------")

	username := a.getInput("Username: ")
	email := a.getInput("Email: ")
	password, err := a.getPassword("Password: ")
	if err != nil {
		fmt.Printf("Password error:%v\n", err)
		return
	}
	if err = a.client.Register(username, email, password); err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return
	}

	fmt.Println("Registration successful!")
}

// handleLogin handles user login
func (a *App) handleLogin() {
	fmt.Println("\n-----------User Login-----------")

	email := a.getInput("Email: ")
	password, err := a.getPassword("Password: ")
	if err != nil {
		fmt.Printf("Password error:%v\n", err)
		return
	}
	if err := a.client.Login(email, password); err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}

	fmt.Println("Login successful!")
}

// handleProfile displays user profile
func (a *App) handleProfile() {
	fmt.Println("\n-----------User Profile-----------")

	profile, err := a.client.GetProfile()
	if err != nil {
		fmt.Printf("Failed to get profile: %v\n", err)
		return
	}

	fmt.Printf("Username: %s\n", profile.Username)
	fmt.Printf("Email: %s\n", profile.Email)
	fmt.Printf("Created: %s\n", profile.CreatedAt)
}

// handleListData lists all data items
func (a *App) handleListData() {
	fmt.Println("\n-----------Your Data Items-----------")

	items, err := a.client.GetAllData()
	if err != nil {
		fmt.Printf("Failed to get data: %v\n", err)
		return
	}

	if len(items) == 0 {
		fmt.Println("No data items found.")
		return
	}

	for i, item := range items {
		fmt.Printf("%d. [%s] %s (ID: %s)\n", i+1, item.Type, item.Title, item.ID)
	}
}

// handleGetData gets specific data item
func (a *App) handleGetData() {
	fmt.Println("\n-----------Get Data Item-----------")

	id := a.getInput("Enter data item ID: ")
	if id == "" {
		fmt.Println("ID cannot be empty")
		return
	}

	item, err := a.client.GetDataByID(id)
	if err != nil {
		fmt.Printf("Failed to get data item: %v\n", err)
		return
	}

	fmt.Printf("\n---Data Item Details:---\n")
	fmt.Printf("ID: %s\n", item.ID)
	fmt.Printf("Type: %s\n", item.Type)
	fmt.Printf("Title: %s\n", item.Title)
	fmt.Printf("Data: %s\n", item.Data)
	if item.Metadata != "" {
		fmt.Printf("Metadata: %s\n", item.Metadata)
	}
	fmt.Printf("Created: %s\n", item.CreatedAt)
	fmt.Printf("Updated: %s\n", item.UpdatedAt)
}

// handleCreateData creates new data item
func (a *App) handleCreateData() {
	fmt.Println("\n-----------Create Data Item-----------")

	fmt.Println("Available types: password, text, binary, card")
	dataType := a.getInput("Type: ")
	title := a.getInput("Title: ")
	data := a.getInput("Data: ")
	metadata := a.getInput("Metadata (optional): ")

	if err := a.client.CreateData(dataType, title, data, metadata); err != nil {
		fmt.Printf("Failed to create data item: %v\n", err)
		return
	}

	fmt.Println("Data item created successfully!")
}

// handleUpdateData updates existing data item
func (a *App) handleUpdateData() {
	fmt.Println("\n-----------Update Data Item-----------")

	id := a.getInput("Enter data item ID: ")
	if id == "" {
		fmt.Println("ID cannot be empty")
		return
	}

	fmt.Println("Available types: password, text, binary, card")
	dataType := a.getInput("Type: ")
	title := a.getInput("Title: ")
	data := a.getInput("Data: ")
	metadata := a.getInput("Metadata (optional): ")

	if err := a.client.UpdateData(id, dataType, title, data, metadata); err != nil {
		fmt.Printf("Failed to update data item: %v\n", err)
		return
	}

	fmt.Println("Data item updated successfully!")
}

// handleDeleteData deletes data item
func (a *App) handleDeleteData() {
	fmt.Println("\n-----------Delete Data Item-----------")

	id := a.getInput("Enter data item ID: ")
	if id == "" {
		fmt.Println("ID cannot be empty")
		return
	}

	confirm := a.getInput("Are you sure? (yes/no): ")
	if strings.ToLower(confirm) != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	if err := a.client.DeleteData(id); err != nil {
		fmt.Printf("Failed to delete data item: %v\n", err)
		return
	}

	fmt.Println("Data item deleted successfully!")
}

// handleChangePassword changes user password
func (a *App) handleChangePassword() {
	fmt.Println("\n-----------Change Password-----------")

	currentPassword, currErr := a.getPassword("Current password: ")
	if currErr != nil {
		fmt.Printf("Current password error:%v\n", currErr)
		return
	}
	newPassword, newErr := a.getPassword("New password: ")
	if newErr != nil {
		fmt.Printf("New password error:%v\n", newErr)
		return
	}
	if err := a.client.ChangePassword(currentPassword, newPassword); err != nil {
		fmt.Printf("Failed to change password: %v\n", err)
		return
	}

	fmt.Println("Password changed successfully!")
}

// getInput gets user input from stdin
func (a *App) getInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// getPassword gets password input (hidden)
func (a *App) getPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password:", err)
		return "", err
	}
	return strings.TrimSpace(string(password)), nil
}
