package magma

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Proposal struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Specs       []Spec `json:"specs"`
}

type Spec struct {
	Index                         string      `json:"index"`
	Name                          string      `json:"name"`
	Enabled                       bool        `json:"enabled"`
	ReliabilityThreshold          uint32      `json:"reliability_threshold"`
	DataReliabilityEnabled        bool        `json:"data_reliability_enabled"`
	BlockDistanceForFinalizedData uint64      `json:"block_distance_for_finalized_data"`
	BlocksInFinalizationProof     uint8       `json:"blocks_in_finalization_proof"`
	AverageBlockTime              string      `json:"average_block_time"`
	AllowedBlockLagForQosSync     string      `json:"allowed_block_lag_for_qos_sync"`
	MinStakeProvider              MinStake    `json:"min_stake_provider"`
	MinStakeClient                MinStake    `json:"min_stake_client"`
	APIs                          APIDataList `json:"apis"`
}

type MinStake struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type APIData struct {
	Name          string             `json:"name"`
	BlockParsing  BlockParsingData   `json:"block_parsing"`
	ComputeUnits  string             `json:"compute_units"`
	Enabled       bool               `json:"enabled"`
	ApiInterfaces []ApiInterfaceData `json:"api_interfaces"`
}

type BlockParsingData struct {
	ParserArg  []string `json:"parser_arg"`
	ParserFunc string   `json:"parser_func"`
}

type ApiInterfaceData struct {
	Category          CategoryData `json:"category"`
	Interface         string       `json:"interface"`
	Type              string       `json:"type"`
	ExtraComputeUnits string       `json:"extra_compute_units"`
}

type CategoryData struct {
	Deterministic bool `json:"deterministic"`
	Local         bool `json:"local"`
	Subscription  bool `json:"subscription"`
	Stateful      int  `json:"stateful"`
}

type APIDataList struct {
	Apis []APIData `json:"apis"`
}

// LOGIC:
func GenerateSpec(fileName string, chainNameFlag string, chainIdxFlag string) error {

	// Check if fileName has ".txt" extension, and add it if not
	if !strings.HasSuffix(fileName, ".txt") {
		fileName += ".txt"
	}

	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// Read the file contents into memory as a byte slice
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Convert the byte slice to a string and split it into lines
	fileContent := string(fileBytes)
	lines := strings.Split(fileContent, "\n")

	data := APIDataList{
		Apis: make([]APIData, 0),
	}

	// Iterate through the lines
	for _, line := range lines {
		fmt.Println(line)
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}
		newData := APIData{
			Name: line,
			BlockParsing: BlockParsingData{
				ParserArg:  []string{"latest"},
				ParserFunc: "DEFAULT",
			},
			ComputeUnits: "10",
			Enabled:      true,
			ApiInterfaces: []ApiInterfaceData{
				{
					Category: CategoryData{
						Deterministic: false,
						Local:         false,
						Subscription:  false,
						Stateful:      0,
					},
					Interface:         "jsonrpc",
					Type:              "POST",
					ExtraComputeUnits: "0",
				},
			},
		}
		data.Apis = append(data.Apis, newData)

	}

	// Write the JSON data to a file
	err = WriteJSONFile("output.json", data, chainNameFlag, chainIdxFlag)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return nil
	}
	fmt.Println("JSON file written successfully.")

	return nil
}

func WriteJSONFile(fileName string, data APIDataList, chainNameFlag string, chainIdxFlag string) error {
	// Write the JSON data to a file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	dataWithHeader := Proposal{
		Title:       "Add Specs: " + chainNameFlag,
		Description: "Adding new specification support for relaying " + chainNameFlag + " data on Lava",
		Specs: []Spec{
			{Index: chainIdxFlag,
				Name:                          chainNameFlag,
				Enabled:                       true,
				ReliabilityThreshold:          268435455,
				DataReliabilityEnabled:        true,
				BlockDistanceForFinalizedData: 64,
				BlocksInFinalizationProof:     3,
				AverageBlockTime:              "13000",
				AllowedBlockLagForQosSync:     "2",
				MinStakeProvider: MinStake{
					Denom:  "ulava",
					Amount: "50000000000",
				},
				MinStakeClient: MinStake{
					Denom:  "ulava",
					Amount: "50000000000",
				},
				APIs: data,
			},
		},
	}
	// Marshal the header into JSON format
	jsonData, err := json.MarshalIndent(dataWithHeader, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}