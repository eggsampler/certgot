package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/eggsampler/certgot/cli"
	"github.com/eggsampler/certgot/log"
	"github.com/eggsampler/certgot/util"
	"gopkg.in/ini.v1"
)

const (
	COMMAND_CERTIFICATES = "certificates"
)

var (
	cmdCertificates = &cli.SubCommand{
		Name:       COMMAND_CERTIFICATES,
		Run:        commandCertificates,
		HelpTopics: []string{TOPIC_MANAGE},
		Usage: cli.SubCommandUsage{
			UsageDescription:    "Display information about certificates you have from Certbot",
			ArgumentDescription: "List certificates managed by Certbot",
			Arguments:           []string{ARG_DOMAIN, ARG_CERT_NAME},
		},
	}

	certSections = []string{"cert", "privkey", "chain", "fullchain"}
)

func commandCertificates(app *cli.App) error {
	renewalDir := filepath.Join(argConfigDir.StringOrDefault(), "renewal")
	renewalConfPattern := filepath.Join(renewalDir, "*.conf")
	renewalFiles, err := filepath.Glob(renewalConfPattern)
	if err != nil {
		log.WithField("path", renewalConfPattern).Error("globbing renewal files")
		return fmt.Errorf("error finding renewal files: %v", err)
	}
	log.WithField("path", renewalConfPattern).WithField("count", len(renewalFiles)).Debug("found renewals files")

	type foundCert struct {
		name     string
		domains  []string
		expiry   time.Time
		certPath string
		keyPath  string
		validStr string
	}

	var foundCerts []foundCert

	for _, f := range renewalFiles {
		ll := log.WithField("renewalfile", f)
		ll.Debug("reading")

		cfg, err := ini.Load(f)
		if err != nil {
			ll.WithError(err).Error("opening renewal file")
			return fmt.Errorf("error loading renewal file %s: %v", f, err)
		}

		skip := false
		for _, v := range certSections {
			if !cfg.Section("").HasKey(v) {
				ll.WithField("section", v).Error("missing required section")
				fmt.Printf("Renewal configuration file %s is missing required section %q. Skipping.\n", f, v)
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		certName := filepath.Base(f)

		if argCertName.IsPresent && argCertName.StringOrDefault() != "" {
			if !strings.EqualFold(certName, argCertName.StringOrDefault()) {
				ll.WithField("wantedcertname", argCertName.StringOrDefault()).
					WithField("foundcertname", certName).Debug("skipping due to cert name mismatch")
				continue
			}
		}

		fc := foundCert{
			name:     certName,
			certPath: cfg.Section("").Key("cert").String(),
			keyPath:  cfg.Section("").Key("privkey").String(),
		}

		cert, err := util.ReadCertificate(fc.certPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		chain, err := util.ReadCertificate(cfg.Section("").Key("chain").String())
		if err != nil {
			fmt.Println(err)
			continue
		}

		revoked, err := util.IsRevoked(cert, chain)
		if err != nil {
			ll.WithError(err).Error("checking ocsp revoked status")
			fmt.Printf("Error checking OCSP revocation status on certificate %s: %v", fc.name, err)
		}

		fc.domains = cert.DNSNames
		fc.expiry = cert.NotAfter

		if strings.Contains(cfg.Section("renewalparams").Key("server").String(), "staging") {
			fc.validStr = "INVALID: TEST_CERT"
		} else if fc.expiry.Before(time.Now()) {
			fc.validStr = "INVALID: EXPIRED"
		} else if revoked {
			fc.validStr = "INVALID: REVOKED"
		} else {
			diff := fc.expiry.Sub(time.Now())
			if diff < 24*time.Hour {
				fc.validStr = fmt.Sprintf("VALID: %.2f hour(s)", diff.Hours())
			} else if diff < 48*time.Hour {
				fc.validStr = "VALID: 1 day"
			} else {
				fc.validStr = fmt.Sprintf("VALID: %d days", int(diff.Hours()/24))
			}
		}

		foundCerts = append(foundCerts, fc)
	}

	fmt.Println(strings.Repeat("- ", 40))

	if len(foundCerts) == 0 {
		fmt.Println("No certificates found")
	} else {
		fmt.Println("Found the following certs:")
		for _, c := range foundCerts {
			fmt.Printf("  Certificate Name: %s\n"+
				"    Domains: %s\n"+
				"    Expiry Date: %s (%s)\n"+
				"    Certificate Path: %s\n"+
				"    Private Key Path: %s\n",
				c.name,
				strings.Join(c.domains, " "),
				c.expiry.String(),
				c.validStr,
				c.certPath,
				c.keyPath)
		}
	}

	fmt.Println(strings.Repeat("- ", 40))

	return nil
}
