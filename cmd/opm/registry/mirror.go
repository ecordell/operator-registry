package registry

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/operator-framework/operator-registry/pkg/mirror"
)

func MirrorCmd() *cobra.Command {
	o := mirror.DefaultImageIndexMirrorerOptions()
	cmd := &cobra.Command{
		Hidden: true,
		Use:    "mirror",
		Short:  "mirror an operator-registry catalog",
		Long:   `mirror an operator-registry catalog image from one registry to another`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if debug, _ := cmd.Flags().GetBool("debug"); debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			dest := args[1]

			mirrorer, err := mirror.NewIndexImageMirror(o, mirror.WithSource(src), mirror.WithDest(dest))
			if err != nil {
				return err
			}
			_, err = mirrorer.Mirror()
			if err != nil {
				return err
			}
			return nil
		},
	}
	flags := cmd.Flags()

	cmd.Flags().Bool("debug", false, "Enable debug logging.")
	flags.StringVar(&o.ManifestDir, "--to-manifests",  "manifests", "Local path to store manifests.")

	return cmd
}
