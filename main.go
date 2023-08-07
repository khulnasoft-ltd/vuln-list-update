package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	githubql "github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"

	"github.com/khulnasoft-ltd/vuln-list-update-update/alma"
	"github.com/khulnasoft-ltd/vuln-list-update-update/alpine"
	alpineunfixed "github.com/khulnasoft-ltd/vuln-list-update-update/alpine-unfixed"
	"github.com/khulnasoft-ltd/vuln-list-update-update/amazon"
	arch_linux "github.com/khulnasoft-ltd/vuln-list-update-update/arch"
	"github.com/khulnasoft-ltd/vuln-list-update-update/chainguard"
	"github.com/khulnasoft-ltd/vuln-list-update-update/cwe"
	"github.com/khulnasoft-ltd/vuln-list-update-update/debian/tracker"
	"github.com/khulnasoft-ltd/vuln-list-update-update/ghsa"
	"github.com/khulnasoft-ltd/vuln-list-update-update/glad"
	"github.com/khulnasoft-ltd/vuln-list-update-update/kevc"
	"github.com/khulnasoft-ltd/vuln-list-update-update/mariner"
	"github.com/khulnasoft-ltd/vuln-list-update-update/nvd"
	oracleoval "github.com/khulnasoft-ltd/vuln-list-update-update/oracle/oval"
	"github.com/khulnasoft-ltd/vuln-list-update-update/osv"
	"github.com/khulnasoft-ltd/vuln-list-update-update/photon"
	redhatoval "github.com/khulnasoft-ltd/vuln-list-update-update/redhat/oval"
	"github.com/khulnasoft-ltd/vuln-list-update-update/redhat/securitydataapi"
	"github.com/khulnasoft-ltd/vuln-list-update-update/rocky"
	susecvrf "github.com/khulnasoft-ltd/vuln-list-update-update/suse/cvrf"
	"github.com/khulnasoft-ltd/vuln-list-update-update/ubuntu"
	"github.com/khulnasoft-ltd/vuln-list-update-update/utils"
	"github.com/khulnasoft-ltd/vuln-list-update-update/wolfi"
)

var (
	target = flag.String("target", "", "update target (nvd, alpine, alpine-unfixed, redhat, redhat-oval, "+
		"debian, ubuntu, amazon, oracle-oval, suse-cvrf, photon, arch-linux, ghsa, glad, cwe, osv, mariner, kevc, wolfi, chainguard)")
	vulnListDir  = flag.String("vuln-list-dir", "", "vuln-list dir")
	targetUri    = flag.String("target-uri", "", "alternative repository URI (only glad)")
	targetBranch = flag.String("target-branch", "", "alternative repository branch (only glad)")
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()
	now := time.Now().UTC()

	if *vulnListDir != "" {
		utils.SetVulnListDir(*vulnListDir)
	}

	switch *target {
	case "nvd":
		if err := nvd.Update(now.Year()); err != nil {
			return xerrors.Errorf("NVD update error: %w", err)
		}
	case "redhat":
		if err := securitydataapi.Update(); err != nil {
			return xerrors.Errorf("Red Hat Security Data API update error: %w", err)
		}
	case "redhat-oval":
		rc := redhatoval.NewConfig()
		if err := rc.Update(); err != nil {
			return xerrors.Errorf("Red Hat OVALv2 update error: %w", err)
		}
	case "debian":
		dc := tracker.NewClient()
		if err := dc.Update(); err != nil {
			return xerrors.Errorf("Debian update error: %w", err)
		}
	case "ubuntu":
		if err := ubuntu.Update(); err != nil {
			return xerrors.Errorf("Ubuntu update error: %w", err)
		}
	case "alpine":
		au := alpine.NewUpdater()
		if err := au.Update(); err != nil {
			return xerrors.Errorf("Alpine update error: %w", err)
		}
	case "alpine-unfixed":
		au := alpineunfixed.NewUpdater()
		if err := au.Update(); err != nil {
			return xerrors.Errorf("Alpine Secfixes Tracker update error: %w", err)
		}
	case "amazon":
		ac := amazon.NewConfig()
		if err := ac.Update(); err != nil {
			return xerrors.Errorf("Amazon Linux update error: %w", err)
		}
	case "oracle-oval":
		oc := oracleoval.NewConfig()
		if err := oc.Update(); err != nil {
			return xerrors.Errorf("Oracle OVAL update error: %w", err)
		}
	case "suse-cvrf":
		sc := susecvrf.NewConfig()
		if err := sc.Update(); err != nil {
			return xerrors.Errorf("SUSE CVRF update error: %w", err)
		}
	case "photon":
		pc := photon.NewConfig()
		if err := pc.Update(); err != nil {
			return xerrors.Errorf("Photon update error: %w", err)
		}
	case "ghsa":
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		httpClient := oauth2.NewClient(context.Background(), src)

		gc := ghsa.NewConfig(githubql.NewClient(httpClient))
		if err := gc.Update(); err != nil {
			return xerrors.Errorf("GitHub Security Advisory update error: %w", err)
		}
	case "glad":
		gu := glad.NewUpdater(*targetUri, *targetBranch)
		if err := gu.Update(); err != nil {
			return xerrors.Errorf("GitLab Advisory Database update error: %w", err)
		}
	case "cwe":
		c := cwe.NewCWEConfig()
		if err := c.Update(); err != nil {
			return xerrors.Errorf("CWE update error: %w", err)
		}
	case "arch-linux":
		al := arch_linux.NewArchLinux()
		if err := al.Update(); err != nil {
			return xerrors.Errorf("Arch Linux update error: %w", err)
		}
	case "alma":
		ac := alma.NewConfig()
		if err := ac.Update(); err != nil {
			return xerrors.Errorf("AlmaLinux update error: %w", err)
		}
	case "rocky":
		rc := rocky.NewConfig()
		if err := rc.Update(); err != nil {
			return xerrors.Errorf("Rocky Linux update error: %w", err)
		}
	case "osv":
		p := osv.NewOsv()
		if err := p.Update(); err != nil {
			return xerrors.Errorf("OSV update error: %w", err)
		}
	case "mariner":
		src := mariner.NewConfig()
		if err := src.Update(); err != nil {
			return xerrors.Errorf("CBL-Mariner Vulnerability Data update error: %w", err)
		}
	case "kevc":
		src := kevc.NewConfig()
		if err := src.Update(); err != nil {
			return xerrors.Errorf("Known Exploited Vulnerability Catalog update error: %w", err)
		}
	case "wolfi":
		wu := wolfi.NewUpdater()
		if err := wu.Update(); err != nil {
			return xerrors.Errorf("Wolfi update error: %w", err)
		}
	case "chainguard":
		cu := chainguard.NewUpdater()
		if err := cu.Update(); err != nil {
			return xerrors.Errorf("Chainguard update error: %w", err)
		}
	default:
		return xerrors.New("unknown target")
	}

	return nil
}
