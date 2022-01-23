package grafana

import (
	"context"
	"fmt"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/grabana/alertmanager/email"
	"github.com/K-Phoen/grabana/alertmanager/slack"
	"github.com/go-logr/logr"
)

var ErrInvalidContactPointType = fmt.Errorf("invalid contact point type")
var ErrInvalidRoutingRule = fmt.Errorf("invalid routing rule")

type AlertManager struct {
	logger        logr.Logger
	grabanaClient *grabana.Client
}

func NewAlertManager(logger logr.Logger, grabanaClient *grabana.Client) *AlertManager {
	return &AlertManager{
		logger:        logger,
		grabanaClient: grabanaClient,
	}
}

func (manager *AlertManager) Reset(ctx context.Context) error {
	return nil
}

func (manager *AlertManager) Configure(ctx context.Context, manifest v1alpha1.AlertManager) error {
	var managerOpts []alertmanager.Option

	// contact points
	contactPointsOpts, err := manager.contactPointsOpts(manifest)
	if err != nil {
		return err
	}
	if len(contactPointsOpts) != 0 {
		managerOpts = append(managerOpts, alertmanager.ContactPoints(contactPointsOpts...))
	}

	// routing rules
	routingOpts, err := manager.routingOpts(manifest)
	if err != nil {
		return err
	}
	if len(routingOpts) != 0 {
		managerOpts = append(managerOpts, alertmanager.Routing(routingOpts...))
	}

	// default contact point
	if manifest.Spec.DefaultContactPoint != "" {
		managerOpts = append(managerOpts, alertmanager.DefaultContactPoint(manifest.Spec.DefaultContactPoint))
	}

	return manager.grabanaClient.ConfigureAlertManager(ctx, alertmanager.New(managerOpts...))
}

func (manager *AlertManager) contactPointsOpts(manifest v1alpha1.AlertManager) ([]alertmanager.Contact, error) {
	var opts []alertmanager.Contact

	for _, contactPointSpec := range manifest.Spec.ContactPoints {
		opt, err := manager.contactPointOpt(contactPointSpec)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (manager *AlertManager) contactPointOpt(contactPointSpec v1alpha1.ContactPoint) (alertmanager.Contact, error) {
	var contactPointTypeOpts []alertmanager.ContactPointOption

	for _, contactPointType := range contactPointSpec.Contacts {
		opt, err := manager.contactPointTypeOpt(contactPointType)
		if err != nil {
			return alertmanager.Contact{}, err
		}

		contactPointTypeOpts = append(contactPointTypeOpts, opt)
	}

	return alertmanager.ContactPoint(contactPointSpec.Name, contactPointTypeOpts...), nil
}

func (manager *AlertManager) contactPointTypeOpt(contactPointType v1alpha1.ContactPointType) (alertmanager.ContactPointOption, error) {
	if contactPointType.Email != nil {
		return manager.contactPointTypeEmail(*contactPointType.Email)
	}
	if contactPointType.Slack != nil {
		return manager.contactPointTypeSlack(*contactPointType.Slack)
	}

	return nil, ErrInvalidContactPointType
}

func (manager *AlertManager) contactPointTypeEmail(contactPointType v1alpha1.EmailContactType) (alertmanager.ContactPointOption, error) {
	var opts []email.Option

	if contactPointType.Single {
		opts = append(opts, email.Single())
	}
	if contactPointType.Message != "" {
		opts = append(opts, email.Message(contactPointType.Message))
	}

	return email.To(contactPointType.To, opts...), nil
}

func (manager *AlertManager) contactPointTypeSlack(contactPointType v1alpha1.SlackContactType) (alertmanager.ContactPointOption, error) {
	var opts []slack.Option

	if contactPointType.Title != "" {
		opts = append(opts, slack.Title(contactPointType.Title))
	}
	if contactPointType.Body != "" {
		opts = append(opts, slack.Body(contactPointType.Body))
	}

	return slack.Webhook(contactPointType.Webhook, opts...), nil
}

func (manager *AlertManager) routingOpts(manifest v1alpha1.AlertManager) ([]alertmanager.RoutingPolicy, error) {
	var opts []alertmanager.RoutingPolicy

	for _, routingPolicy := range manifest.Spec.Routing {
		opt, err := manager.routingPolicyOpt(routingPolicy)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (manager *AlertManager) routingPolicyOpt(policySpec v1alpha1.RoutingPolicy) (alertmanager.RoutingPolicy, error) {
	var opts []alertmanager.RoutingPolicyOption

	for _, rule := range policySpec.Rules {
		labelOpts, err := manager.routingLabelRules(rule)
		if err != nil {
			return alertmanager.RoutingPolicy{}, err
		}

		opts = append(opts, labelOpts...)
	}

	return alertmanager.Policy(policySpec.ContactPoint, opts...), nil
}

func (manager *AlertManager) routingLabelRules(rule v1alpha1.LabelsMatchingRule) ([]alertmanager.RoutingPolicyOption, error) {
	var operatorFunc func(tag string, value string) alertmanager.RoutingPolicyOption
	var labels map[string]string
	var opts []alertmanager.RoutingPolicyOption

	if rule.Eq != nil {
		operatorFunc = alertmanager.TagEq
		labels = rule.Eq
	} else if rule.Neq != nil {
		operatorFunc = alertmanager.TagNeq
		labels = rule.Neq
	} else if rule.Matches != nil {
		operatorFunc = alertmanager.TagMatches
		labels = rule.Matches
	} else if rule.NotMatches != nil {
		operatorFunc = alertmanager.TagNotMatches
		labels = rule.NotMatches
	} else {
		return nil, ErrInvalidRoutingRule
	}

	for key, value := range labels {
		opts = append(opts, operatorFunc(key, value))
	}

	return opts, nil
}
