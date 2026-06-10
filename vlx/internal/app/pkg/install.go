package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rvelte/vlx/internal/core/cmd"
	"github.com/rvelte/vlx/internal/core/printer"
)

func Install(ctx context.Context, s *Scheme) error {
	p := printer.New()

	if err := runHook(ctx, p, "pre", s.PreHook); err != nil {
		return err
	}

	if len(s.Repos) > 0 {
		p.Header(fmt.Sprintf("Adding %d zypper repo(s)...", len(s.Repos)))
		for _, repo := range s.Repos {
			err := cmd.New("sudo", "zypper", "addrepo", repo.URL, repo.Alias).
				Stdout(os.Stdout).
				Stderr(os.Stderr).
				Stdin(os.Stdin).
				Run(ctx)
			if err != nil {
				p.Warn(fmt.Sprintf("Repo %s may already exist, continuing...", repo.Alias))
			}
		}
	}

	if len(s.Zypper) > 0 {
		p.Info("Refreshing zypper repos...")
		err := cmd.New("sudo", "zypper", "refresh").
			Stdout(os.Stdout).
			Stderr(os.Stderr).
			Stdin(os.Stdin).
			Run(ctx)
		if err != nil {
			return fmt.Errorf("zypper refresh: %w", err)
		}

		p.Header(fmt.Sprintf("Installing %d zypper package(s)...", len(s.Zypper)))
		args := append([]string{"install", "-y"}, s.Zypper...)
		err = cmd.New("sudo", append([]string{"zypper"}, args...)...).
			Stdout(os.Stdout).
			Stderr(os.Stderr).
			Stdin(os.Stdin).
			Run(ctx)
		if err != nil {
			return fmt.Errorf("zypper install: %w", err)
		}
	}

	if len(s.Flatpak) > 0 {
		p.Info("Refreshing flatpak remotes...")
		err := cmd.New("flatpak", "update", "--appstream", "-y").
			Stdout(os.Stdout).
			Stderr(os.Stderr).
			Stdin(os.Stdin).
			Run(ctx)
		if err != nil {
			return fmt.Errorf("flatpak refresh: %w", err)
		}

		p.Header(fmt.Sprintf("Installing %d flatpak package(s)...", len(s.Flatpak)))
		for _, ref := range s.Flatpak {
			err := cmd.New("flatpak", "install", "-y", "flathub", ref).
				Stdout(os.Stdout).
				Stderr(os.Stderr).
				Stdin(os.Stdin).
				Run(ctx)
			if err != nil {
				return fmt.Errorf("flatpak install %s: %w", ref, err)
			}
		}
	}

	if err := runHook(ctx, p, "post", s.PostHook); err != nil {
		return err
	}

	return nil
}

func runHook(ctx context.Context, p *printer.Printer, label string, steps []string) error {
	if len(steps) == 0 {
		return nil
	}
	p.Info(fmt.Sprintf("Running %s hook...", label))
	script := strings.Join(steps, " && ")
	err := cmd.New("sh", "-c", script).
		Stdout(os.Stdout).
		Stderr(os.Stderr).
		Stdin(os.Stdin).
		Run(ctx)
	if err != nil {
		return fmt.Errorf("%s hook: %w", label, err)
	}
	return nil
}
