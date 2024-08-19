package main

import (
	"fmt"
	"log"
	"time"
	"os/exec"
	"os"
	"strings"

	. "github.com/compliance-framework/assessment-runtime/provider"
	"github.com/google/uuid"
)

type PrivateerExampleProvider struct {
	message string
}

func (p *PrivateerExampleProvider) Evaluate(input *EvaluateInput) (*EvaluateResult, error) {
	// Extract the passed-in YAML into code here.
	yamlString, ok := input.Configuration["yaml"]
	if !ok {
		return nil, fmt.Errorf("yaml parameter is missing")
	}
	log.Printf("yamlString: %s", yamlString)

	// Take that YAML string and create a config file /raid-wireframe-config.yml, eg
	// Define the path where the YAML file should be created
	configFilePath := "/raid-wireframe-config.yml"

	// Create the config file
	file, err := os.Create(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the YAML string to the file
	_, err = file.WriteString(yamlString)
	if err != nil {
		return nil, fmt.Errorf("failed to write to file: %v", err)
	}

	log.Printf("Successfully created config file at %s", configFilePath)

	// There can be an array of subjects if needed, but here we have only one
	subjects := make([]*Subject, 0)
	subject_id := fmt.Sprintf("Subject identifier: %s", "Privateer example Subject ID placeholder")
	subjects = append(subjects, &Subject{
		Id:    subject_id,
		Type:  SubjectType_INVENTORY_ITEM,
		Title: "Privateer Example Subject",
		Props: map[string]string{
			"id": subject_id,
		},
	})

	// Return the result with subjects and additional props if necessary
	return &EvaluateResult{
		Subjects: subjects,
	}, nil
}

func (p PrivateerExampleProvider) Execute(input *ExecuteInput) (*ExecuteResult, error) {
	start_time := time.Now().Format(time.RFC3339)

	var obs *Observation
	var fndngs *Finding

	observations := []*Observation{}
	findings := []*Finding{}

	obs_id := uuid.New().String()

	// Work out whether we're compliant. Assume false first.
	compliant := false

	// Run the command and capture the output
	cmd := exec.Command("bash", "-c", "privateer -c /raid-wireframe-config.yml sally | grep -c ERROR")
	outputBytes, _ := cmd.Output()
	// TODO: handle err
	// Convert the output to a string and trim whitespace
	outputStr := strings.TrimSpace(string(outputBytes))
	// Check the output value
	if outputStr == "0" {
		compliant = true
	}

	if (!compliant) {
		// observation and finding
		obs = &Observation{
			Id:               obs_id,
			Title:            "Privateer Example Observation",
			Description:      "Description of the observation that did not succeed",
			Collected:        time.Now().Format(time.RFC3339),
			Expires:          time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
			Links:            []*Link{},
			Props:            []*Property{},
			RelevantEvidence: []*Evidence{
				{
					Description: "Raid ran and failed",
				},
			},
			Remarks:          "A remark",
		}
		fndngs = &Finding{
			Id:                  uuid.New().String(),
			Title:               "Raid finding title",
			Description:         "Raid finding description",
			Remarks:             "Some relevant remarks",
			RelatedObservations: []string{obs_id},
		}
		observations = append(observations, obs)
		findings = append(findings, fndngs)
	} else {
		// observation only
		obs = &Observation{
			Id:          obs_id,
			Title:       "Privateer Example Observation",
			Description: "Description of the observation that succeeded",
			Collected:   time.Now().Format(time.RFC3339),
			Expires:     time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
			Links:       []*Link{},
			Props: []*Property{},
			RelevantEvidence: []*Evidence{
				{
					Description: "Raid ran and succeeded",
				},
			},
			Remarks: "All OK.",
		}
		observations = append(observations, obs)
	}

	// Log that the check has successfully run
	logEntry := &LogEntry{
		Title:       "Privateer example log entry title",
		Description: "Privateer example log entry description",
		Start:       start_time,
		End:         time.Now().Format(time.RFC3339),
	}

	// Return the result
	return &ExecuteResult{
		Status:       ExecutionStatus_SUCCESS,
		Observations: observations,
		Findings:     findings,
		Logs:         []*LogEntry{logEntry},
	}, nil
}

func main() {
	Register(&PrivateerExampleProvider{
		message: "PrivateerExample completed",
	})
}
