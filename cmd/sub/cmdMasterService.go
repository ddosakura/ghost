package sub

import (
	"github.com/spf13/cobra"
)

// Cmd
var (
	MasterServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "manage status of service",
		Long:  `Service Start/Stop/Status (recommend starting by system service).`,
	}
)

func init() {
	//MasterServiceCmd.PersistentFlags().StringVarP(
	//	&serviceSign,
	//)
}
