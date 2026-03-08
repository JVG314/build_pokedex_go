package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JVG314/build_pokedex_go/internal/pokeapi"
)

type config struct {
	Next     *string
	Previous *string
	Client   *pokeapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// defined here so we can use it later inside commandHelp to iterate over it
// but initialized inside startRepl() to avoid initialization circular dependency
// as commands is being initialized as  package load time, commandHelp is referenc
// ed in that initialization but commandHelp tries to use commands to iterate over it
// so we define here but initalize later
var commands map[string]cliCommand

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := ""
	if cfg.Next != nil {
		url = *cfg.Next
	}

	res, err := cfg.Client.GetLocationAreas(url)
	if err != nil {
		return err
	}

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous

	return nil
}

func commandMapb(cfg *config, args []string) error {
	url := ""
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	url = *cfg.Previous
	res, err := cfg.Client.GetLocationAreas(url)
	if err != nil {
		return err
	}

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a location area name")
	}
	name := args[0]
	fmt.Printf("Exploring %s area...\n", name)
	res, err := cfg.Client.GetLocationDetails(name)
	if err != nil {
		return err
	}

	for _, location := range res.PokemonEncounters {
		fmt.Printf("- %s\n", location.Pokemon.Name)
	}

	// cfg.Next = res.Next
	// cfg.Previous = res.Previous

	return nil
}

func startRepl() {
	cfg := &config{
		Client: pokeapi.NewClient(10 * time.Second),
	}
	scanner := bufio.NewScanner(os.Stdin)
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Explore the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Explore the Pokemon world backwards",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a particular location",
			callback:    commandExplore,
		},
	}
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		cleaned := cleanInput(input)
		if len(cleaned) == 0 {
			continue
		}
		// fmt.Printf("Your command was: %s\n", cleaned[0])
		commandName := cleaned[0]
		args := cleaned[1:]
		command, exists := commands[commandName]
		if exists {
			err := command.callback(cfg, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}
