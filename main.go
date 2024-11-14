package main

import (
	"os"

	"go.pberndt.com/iwebd/dlnad"
	"go.pberndt.com/iwebd/ftpd"
	"go.pberndt.com/iwebd/httpd"
	"go.pberndt.com/iwebd/util"
	"go.pberndt.com/iwebd/webdavd"

	"github.com/spf13/cobra"
)

const longDoc = `iwebd combines a bunch of means to share files over the local network

dlna functionality is provided by dms by anacrolix (BSD 3),
ftp functionality is provided by goftp by yob et al (MIT),
icons are taken from the Yaru theme from Ubuntu (CC BY-SA 4.0),
and this entire project benefits a lot from the huge stdlib Go brings,
rest is by Phillip Berndt <phillip.berndt@googlemail.com>.

iwebd is available under the terms of the GNU General Public License version 3.

`

func main() {
	// I'm deviating from standard Cobra use here for simplicity. All
	// commands are defined in this file, with implementations in modules
	// of their own.

	rootCmd := &cobra.Command{
		Use:   "iwebd",
		Short: "iwebd combines a bunch of means to share files over the local network",
		Long:  longDoc,
		Args:  cobra.NoArgs,
	}
	rootCmd.PersistentFlags().String("addr", "", "Local address:port to bind to")
	rootCmd.PersistentFlags().Bool("bonjour", false, "Announce service via zeroconf")

	// httpd -- http server
	httpd := &cobra.Command{
		Use:   "http",
		Short: "Spawn a http server",
		Args:  cobra.NoArgs,
		Run:   httpd.Serve,
	}
	httpd.Flags().Var(&util.UserPass{}, "auth", "Credentials to use")
	httpd.Flags().Bool("tls", false, "Serve https rather than http")
	httpd.Flags().Bool("read-only", false, "Deny clients write-access")
	httpd.Flags().Bool("live-reload", false, "Enable/embed live-reload to HTML pages")
	rootCmd.AddCommand(httpd)

	// ftpd -- ftp server
	ftpd := &cobra.Command{
		Use:   "ftp",
		Short: "Spawn an ftp server",
		Args:  cobra.NoArgs,
		Run:   ftpd.Serve,
	}
	ftpd.Flags().Var(&util.UserPass{}, "auth", "Credentials to use")
	ftpd.Flags().Bool("tls", false, "Serve https rather than http")
	ftpd.Flags().Bool("read-only", false, "Deny clients write-access")
	rootCmd.AddCommand(ftpd)

	// dlnad -- dlna server
	dlnad := &cobra.Command{
		Use:   "dlna",
		Short: "Spawn a dlna server",
		Args:  cobra.NoArgs,
		Run:   dlnad.Serve,
	}
	dlnad.Flags().Bool("no-ffmpeg", false, "Disable probing/transcoding using ffmpeg")
	rootCmd.AddCommand(dlnad)

	// webdavd -- webdav server
	webdavd := &cobra.Command{
		Use:   "webdav",
		Short: "Spawn a webdav server",
		Args:  cobra.NoArgs,
		Run:   webdavd.Serve,
	}
	webdavd.Flags().Var(&util.UserPass{}, "auth", "Credentials to use")
	webdavd.Flags().Bool("tls", false, "Serve https rather than http")
	webdavd.Flags().Bool("read-only", false, "Deny clients write-access")
	rootCmd.AddCommand(webdavd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
