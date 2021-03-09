/*
Copyright © 2021 Murphyyi <zy84338719@hotmail.com>

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
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rds     Redis
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rd-tool",
	Short: "A redis tool application",
	Long: `一个帮助你实现某些特定功能的redis操作功能的程序 

For example:
	可以帮助你缓慢删除redis 内部某个db下面的存储数据
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		service := NewService(rds)
		service.Del(rds.Key, rds.Pattern, rds.Types, rds.BatchSize)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotools.yaml)")
	rootCmd.PersistentFlags().BoolVar(&rds.ClusterMode, "cluster", false, "集群类型")
	rootCmd.PersistentFlags().StringVar(&rds.Addr, "addr", "localhost:6379", "redis地址，集群版请用 , 隔开")
	rootCmd.PersistentFlags().IntVar(&rds.DB, "db", 0, "redis db ")
	rootCmd.PersistentFlags().StringVar(&rds.Key, "key", "", "删除key")
	rootCmd.PersistentFlags().StringVar(&rds.Password, "password", "", "redis 密码")
	rootCmd.PersistentFlags().StringVar(&rds.Pattern, "pattern", "", "")
	rootCmd.PersistentFlags().StringVar(&rds.Types, "types", "string", "string set hash list zset")
	rootCmd.PersistentFlags().Int64Var(&rds.BatchSize, "batchsize", 10000, "分片大小")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gotools" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gotools")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
