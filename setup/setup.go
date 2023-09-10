package setup

import (
	"DadJokesAPI/server/database"
	"DadJokesAPI/shared"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

func Run() error {
	if err := setupConfig(); err != nil {
		return fmt.Errorf("failed to setup server config: %w", err)
	}
	if err := createDatabase(); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	if err := setupDatabase(); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	return nil
}

func setupConfig() error {
	_, err := os.Stat("config.toml")
	if err == nil && readInput("Config file already exists? Do you want to overwrite it? (Y/n): ")[0] != byte('Y') {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check for config file: %w", err)
	}

	configFile, err := os.Create("config.toml")
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer configFile.Close()

	lines := []string{
		"[server]",
		"port = 8080",
		"debug = true",

		"[postgres]",
		"port = 12345",
		"db = \"dad-jokes-api\"",
		"user = \"db-username\"",
		"password = \"db-password\"",
		"max_idle_connections = 8",
		"max_open_connections = 64",

		"[rate_limiter]",
		"enabled = false",
		"rate = 1",
		"burst = 10",
		"timeout = 180",
	}
	if err = writeLines(configFile, lines); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func createDatabase() error {
	dbPort, err := shared.GetInt("postgres.port")
	if err != nil {
		return fmt.Errorf("failed to get postgres port: %w", err)
	}
	dbName, err := shared.GetStr("postgres.db")
	if err != nil {
		return fmt.Errorf("failed to get postgres db name: %w", err)
	}
	dbUser, err := shared.GetStr("postgres.user")
	if err != nil {
		return fmt.Errorf("failed to get postgres user: %w", err)
	}
	dbPass, err := shared.GetStr("postgres.password")
	if err != nil {
		return fmt.Errorf("failed to get postgres password: %w", err)
	}

	envFile, err := os.Create("setup/.env")
	if err != nil {
		return fmt.Errorf("failed to create docker .env file: %w", err)
	}
	defer envFile.Close()

	lines := []string{
		"CONTAINER_PORT=5432",
		fmt.Sprintf("HOST_PORT=%d", dbPort),
		fmt.Sprintf("POSTGRES_DB=%s", dbName),
		fmt.Sprintf("POSTGRES_USER=%s", dbUser),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPass),
	}
	if err = writeLines(envFile, lines); err != nil {
		return fmt.Errorf("failed to write docker .env file: %w", err)
	}

	cmd := exec.Command("docker-compose", "-f", "setup/docker-compose.yaml", "up", "-d") // -d to run in detached mode
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run docker-compose: %w", err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	for i := 0; i < 10; i++ {
		inspectResponse, err := cli.ContainerInspect(context.Background(), "dad-jokes-postgres")
		if err != nil {
			return fmt.Errorf("failed to inspect container: %w", err)
		}

		if inspectResponse.State.Running {
			time.Sleep(5 * time.Second)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("failed to start database container")
}

func setupDatabase() error {
	db, err := shared.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	setupSql, err := os.ReadFile("setup/setup.sql")
	if err != nil {
		return fmt.Errorf("failed to open setup.sql: %w", err)
	}

	err = db.Exec(string(setupSql)).Error
	if err != nil {
		return fmt.Errorf("failed to execute setup.sql: %w", err)
	}

	jokeList, err := os.ReadFile("setup/jokes.txt")
	if err != nil {
		return fmt.Errorf("failed to read preset jokes.txt: %w", err)
	}

	jokes := strings.Split(string(jokeList), "\n")
	for _, joke := range jokes {
		if _, err = database.CreateJoke(database.Joke{
			Text: joke,
		}); err != nil {
			return fmt.Errorf("failed to create joke: %w", err)
		}
	}

	return nil
}

func readInput(prompt string) string {
	var input string

	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

func writeLines(file *os.File, lines []string) error {
	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write file line %s: %w", line, err)
		}
	}
	return nil
}
