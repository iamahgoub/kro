package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kro-run/kro/api/v1alpha1"
	kroclient "github.com/kro-run/kro/pkg/client"
	"github.com/kro-run/kro/pkg/graph"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type GenerateConfig struct {
	resourceGraphDefinitionFile string
	outputFormat                string
}

var config = &GenerateConfig{}

func init() {
	generateCmd.PersistentFlags().StringVarP(&config.resourceGraphDefinitionFile, "file", "f", "",
		"Path to the ResourceGraphDefinition file")
	generateCmd.PersistentFlags().StringVarP(&config.outputFormat, "format", "o", "yaml", "Output format (yaml|json)")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the CRD from a given ResourceGraphDefinition",
	Long: "Generate the CRD from a given ResourceGraphDefinition." +
		"This command generates a CustomResourceDefinition (CRD) based on the provided ResourceGraphDefinition file.",
}

var generateCRDCmd = &cobra.Command{
	Use:   "crd",
	Short: "Generate a CustomResourceDefinition (CRD)",
	Long: "Generate a CustomResourceDefinition (CRD) from a " +
		"ResourceGraphDefinition file. This command reads the " +
		"ResourceGraphDefinition and outputs the corresponding CRD " +
		"in the specified format.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.resourceGraphDefinitionFile == "" {
			return fmt.Errorf("ResourceGraphDefinition file is required")
		}

		data, err := os.ReadFile(config.resourceGraphDefinitionFile)
		if err != nil {
			return fmt.Errorf("failed to read ResourceGraphDefinition file: %w", err)
		}

		var rgd v1alpha1.ResourceGraphDefinition
		if err = yaml.Unmarshal(data, &rgd); err != nil {
			return fmt.Errorf("failed to unmarshal ResourceGraphDefinition: %w", err)
		}

		if err = generateCRD(&rgd); err != nil {
			return fmt.Errorf("failed to generate CRD: %w", err)
		}

		return nil
	},
}

var generateInstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Generate a ResourceGraphDefinition (RGD) instance",
	Long: "Generate a ResourceGraphDefinition (RGD) instance from a " +
		"ResourceGraphDefinition file. This command reads the " +
		"ResourceGraphDefinition and outputs the corresponding RGD instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.resourceGraphDefinitionFile == "" {
			return fmt.Errorf("ResourceGraphDefinition file is required")
		}

		data, err := os.ReadFile(config.resourceGraphDefinitionFile)
		if err != nil {
			return fmt.Errorf("failed to read ResourceGraphDefinition file: %w", err)
		}

		var rgd v1alpha1.ResourceGraphDefinition
		if err = yaml.Unmarshal(data, &rgd); err != nil {
			return fmt.Errorf("failed to unmarshal ResourceGraphDefinition: %w", err)
		}

		if err = generateInstance(&rgd); err != nil {
			return fmt.Errorf("failed to generate RGD instance: %w", err)
		}

		return nil
	},
}

func generateCRD(rgd *v1alpha1.ResourceGraphDefinition) error {
	rgdGraph, err := createGraphBuilder(rgd)
	if err != nil {
		return fmt.Errorf("failed to setup rgd graph: %w", err)
	}

	crd := rgdGraph.Instance.GetCRD()
	crd.SetAnnotations(map[string]string{"kro.run/version": "dev"})

	b, err := marshalObject(crd, config.outputFormat)
	if err != nil {
		return fmt.Errorf("failed to marshal CRD: %w", err)
	}

	fmt.Println(string(b))

	return nil
}

func generateInstance(rgd *v1alpha1.ResourceGraphDefinition) error {
	rgdGraph, err := createGraphBuilder(rgd)
	if err != nil {
		return fmt.Errorf("failed to create resource graph definition: %w", err)
	}

	emulatedObj := rgdGraph.Instance.GetEmulatedObject()
	emulatedObj.SetAnnotations(map[string]string{"kro.run/version": "dev"})

	delete(emulatedObj.Object, "status")

	b, err := marshalObject(emulatedObj, config.outputFormat)
	if err != nil {
		return fmt.Errorf("failed to marshal CRD: %w", err)
	}

	fmt.Println(string(b))

	return nil
}

func createGraphBuilder(rgd *v1alpha1.ResourceGraphDefinition) (*graph.Graph, error) {
	set, err := kroclient.NewSet(kroclient.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create client set: %w", err)
	}

	restConfig := set.RESTConfig()

	builder, err := graph.NewBuilder(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph builder: %w", err)
	}

	rgdGraph, err := builder.NewResourceGraphDefinition(rgd)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource graph definition: %w", err)
	}

	return rgdGraph, nil
}

func marshalObject(obj interface{}, outputFormat string) ([]byte, error) {
	var b []byte
	var err error

	if obj == nil {
		return nil, fmt.Errorf("nil object provided for marshaling")
	}

	switch outputFormat {
	case "json":
		b, err = json.Marshal(obj)
		if err != nil {
			return nil, err
		}
	case "yaml":
		b, err = yaml.Marshal(obj)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	return b, nil
}

func AddGenerateCommands(rootCmd *cobra.Command) {
	generateCmd.AddCommand(generateCRDCmd)
	generateCmd.AddCommand(generateInstanceCmd)
	rootCmd.AddCommand(generateCmd)
}
