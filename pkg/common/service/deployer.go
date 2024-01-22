package service

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type Middleware interface {
	BeforeRun(ctx context.Context, g *errgroup.Group, service Service)
}

type deploymentConfiguration struct {
	service     Service
	middlewares []Middleware
}

type Deployer struct {
	configs []deploymentConfiguration
}

func NewDeployer() *Deployer {
	return &Deployer{}
}

func (deployer *Deployer) Add(service Service, middlewares ...Middleware) *Deployer {
	deployer.configs = append(deployer.configs, deploymentConfiguration{
		service:     service,
		middlewares: middlewares,
	})

	return deployer
}

func (deployer *Deployer) Deploy(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, cfg := range deployer.configs {
		cfg := cfg
		g.Go(func() error {
			for _, mw := range cfg.middlewares {
				mw.BeforeRun(ctx, g, cfg.service)
			}
			err := cfg.service.Run(ctx)
			if err != nil {
				return fmt.Errorf("failed to run (service='%s' id='%s'): %w",
					cfg.service.Name(),
					cfg.service.ID(),
					err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
