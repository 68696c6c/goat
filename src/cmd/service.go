package cmd

import "github.com/spf13/cobra"

// Goat uses the Cobra framework.
type Service interface {
	RootCMD() *cobra.Command
	//GenMigrationCMD(db *gorm.DB) *cobra.Command
	//GenMigrateCMD(db *gorm.DB) *cobra.Command
	//GenServerCMD() *cobra.Command
}

type Config struct {
}

type ServiceCobra struct {
	config Config
}

func NewServiceCobra(c Config) ServiceCobra {
	return ServiceCobra{
		config: c,
	}
}

func (s ServiceCobra) RootCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "goat",
		Short: "Root command for Goat CLI",
	}
}
