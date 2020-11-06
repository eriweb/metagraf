/*
Copyright 2018-2020 The metaGraf Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	log "k8s.io/klog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
)

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
	createCmd.PersistentFlags().BoolVar(&Output, "output", false, "also output objects")
	createCmd.PersistentFlags().StringVarP(&Format, "format","o","json", "specify json or yaml, json id default")
	createCmd.PersistentFlags().StringVar(&Version, "version", "", "Override version in metaGraf specification.")
	createCmd.PersistentFlags().BoolVar(&Dryrun, "dryrun", false, "do not create objects, only output")
	createCmd.AddCommand(createConfigMapCmd)
	createCmd.AddCommand(createDeploymentCmd)
	createCmd.AddCommand(createDeploymentConfigCmd)
	createCmd.AddCommand(createDotCmd)
	createCmd.AddCommand(createRefCmd)
	createCmd.AddCommand(createSecretCmd)
	createCmd.AddCommand(createRouteCmd)
	createDeploymentCmd.Flags().StringVarP(&Namespace, "namespace", "n","", "namespace to work on, if not supplied it will use current working namespace")
	createDeploymentCmd.Flags().StringVar(&OName, "name", "", "Overrides name of deployment.")
	createDeploymentCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createDeploymentCmd.Flags().StringVar(&CVfile, "cvfile","", "File with component configuration values. (key=value pairs)")
	createDeploymentCmd.Flags().BoolVar(&BaseEnvs, "baseenv", false, "Hydrate deploymentconfig with baseimage environment variables")
	createDeploymentCmd.Flags().BoolVar(&Defaults, "defaults", false, "Populate Environment variables with default values from metaGraf")
	createDeploymentCmd.Flags().StringVarP(&ImageNS,"imagens", "i", "", "Image Namespace, used to override default namespace")
	createDeploymentCmd.Flags().StringVarP(&Registry,"registry", "r",viper.GetString("registry"), "Specify container registry host")
	createDeploymentCmd.Flags().StringVarP(&Tag,"tag", "t", "latest", "specify custom tag")
	createDeploymentConfigCmd.Flags().StringVarP(&Namespace, "namespace", "n","", "namespace to work on, if not supplied it will use current working namespace")
	createDeploymentConfigCmd.Flags().StringVar(&OName, "name", "", "Overrides name of deployment.")
	createDeploymentConfigCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createDeploymentConfigCmd.Flags().StringVar(&CVfile, "cvfile","", "File with component configuration values. (key=value pairs)")
	createDeploymentConfigCmd.Flags().BoolVar(&BaseEnvs, "baseenv", false, "Hydrate deploymentconfig with baseimage environment variables")
	createDeploymentConfigCmd.Flags().BoolVar(&Defaults, "defaults", false, "Populate Environment variables with default values from metaGraf")
	createDeploymentConfigCmd.Flags().StringVarP(&ImageNS,"imagens", "i", "", "Image Namespace, used to override default namespace")
	createDeploymentConfigCmd.Flags().StringVarP(&Registry,"registry", "r",viper.GetString("registry"), "Specify container registry host")
	createDeploymentConfigCmd.Flags().StringVarP(&Tag,"tag", "t", "latest", "specify custom tag")
	createSecretCmd.Flags().StringVarP(&Namespace, "namespace", "n", "","namespace to work on, if not supplied it will use current working namespace")
	createSecretCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createSecretCmd.Flags().BoolVarP(&CreateGlobals, "globals", "g", false, "Override default behavior and force creation of global secrets. Will not overwrite existing ones.")
	createConfigMapCmd.Flags().StringVarP(&Namespace, "namespace", "n", "","namespace to work on, if not supplied it will use current working namespace")
	createConfigMapCmd.Flags().StringVar(&OName, "name", "", "Overrides name of application used to prefix configmaps.")
	createConfigMapCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createConfigMapCmd.Flags().StringVar(&CVfile, "cvfile","", "File with component configuration values. (key=value pairs)")
	createRouteCmd.Flags().StringVarP(&Namespace, "namespace", "n", "","namespace to work on, if not supplied it will use current working namespace")
	createRouteCmd.Flags().StringVar(&OName, "name", "", "Overrides name of application.")
	createRouteCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createRouteCmd.Flags().StringVarP(&Context,"context", "c","/","Application context root. (\"/<context>\")")

}


var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create operations",
	Long:  MGBanner + ` create `,
}

var createConfigMapCmd = &cobra.Command{
	Use:   "configmap <metagraf>",
	Short: "create ConfigMaps from metaGraf file",
	Long:  MGBanner + `create ConfigMap`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
			log.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				log.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}
		FlagPassingHack()

		mg := metagraf.Parse(args[0])

		modules.Variables = mg.GetProperties()
		OverrideProperties(modules.Variables)
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}

		modules.GenConfigMaps(&mg)
	},
}

var createDeploymentCmd = &cobra.Command{
	Use:   "deployment <metagraf>",
	Short: "create Deployment from metaGraf file",
	Long:  MGBanner + `create Deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
			log.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				log.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}

		mg := metagraf.Parse(args[0])
		FlagPassingHack()

		modules.Variables = mg.GetProperties()
		OverrideProperties(modules.Variables)

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}
		modules.GenDeployment(&mg, Namespace)
	},
}

var createDeploymentConfigCmd = &cobra.Command{
	Use:   "deploymentconfig <metagraf>",
	Short: "create DeploymentConfig from metaGraf file",
	Long:  MGBanner + `create DeploymentConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
			log.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				log.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}

		mg := metagraf.Parse(args[0])
		FlagPassingHack()

		modules.Variables = mg.GetProperties()
		OverrideProperties(modules.Variables)
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}

		modules.GenDeploymentConfig(&mg)
	},
}



var createDotCmd = &cobra.Command{
	Use:   "dot <collection directory>",
	Short: "create Graphviz service graph from collectio of metaGraf's",
	Long:  MGBanner + `create dot`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(StrMissingCollection)
			return
		}
		FlagPassingHack()
		modules.GenDotFromPath(args[0])
	},
}

var createSecretCmd = &cobra.Command{
	Use:   "secret <metaGraf>",
	Short: "create Secrets from metaGraf specification",
	Long:  MGBanner + `create Secret`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
			log.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				log.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}
		FlagPassingHack()
		mg := metagraf.Parse(args[0])

		modules.GenSecrets(&mg)
	},
}

var createRouteCmd = &cobra.Command{
	Use:   "route <metaGraf>",
	Short: "create Route from metaGraf specification",
	Long:  MGBanner + `create route`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
			log.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				log.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}
		FlagPassingHack()
		mg := metagraf.Parse(args[0])

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}

		modules.GenRoute(&mg)
	},
}