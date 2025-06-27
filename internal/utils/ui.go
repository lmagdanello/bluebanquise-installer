package utils

import "fmt"

// ShowCompletionMessage displays the completion message
func ShowCompletionMessage(userName, userHome string) {
	fmt.Println()
	fmt.Println("Bootstrap done.")
	fmt.Printf("You can now login as %s user via 'su - %s'\n", userName, userName)
	fmt.Println()
	fmt.Println("To use BlueBanquise, remember to set Ansible environment variable:")
	fmt.Printf("ANSIBLE_CONFIG=$HOME/bluebanquise/ansible.cfg\n")
	fmt.Println()
	fmt.Println("You can find documentation at http://bluebanquise.com/documentation/")
	fmt.Println("You can ask for help or rise issues at https://github.com/bluebanquise/bluebanquise/")
	fmt.Println()
	fmt.Println("Thank you for using BlueBanquise :)")
	fmt.Println("Have fun!")
	fmt.Println()
}
